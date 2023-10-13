package search

import (
	"fmt"
	"html/template"
	"os"
	"regexp"

	dir "github.com/oworkshop/ow_wiki/app/directory"
)

type RegList struct {
	List    []*regexp.Regexp
	Percent int //0-100
}

func NewRegList(percent int) *RegList {
	return &RegList{
		List:    []*regexp.Regexp{},
		Percent: percent,
	}
}

func (r *RegList) Add(reg ...*regexp.Regexp) *RegList {
	r.List = append(r.List, reg...)
	return r
}

func (r *RegList) Len() int {
	return len(r.List)
}

func (r *RegList) Match(s string) bool {
	num := 0
	for _, v := range r.List {
		if v.MatchString(s) {
			num += 1
		}
	}
	return num*100 >= r.Percent*len(r.List)
}

func (r *RegList) Highlight(s string) string {
	partList := []*Part{NewPart(s, false)}
	for _, v := range r.List {
		partList = HighlightPartList(partList, v)
	}
	var result string
	for _, v := range partList {
		result += v.String()
	}
	return result
}

type RegMap struct {
	Map   map[int][]*RegList
	Level int
}

func (r *RegMap) AddLevel() {
	r.Level += 1
}

func (r *RegMap) Add(l *RegList) {
	if r.Map == nil {
		r.Map = make(map[int][]*RegList)
	}
	r.Map[r.Level] = append(r.Map[r.Level], l)
}

func (r *RegMap) Get(level int) []*RegList {
	if r.Map == nil {
		return nil
	}
	return r.Map[level]
}

/*
优先级高的方法能搜到结果，就不再搜优先级低的
*/
func (r *RegMap) Search(filelist []*File) ([]Result, error) {
	var resultList []Result
	for level := 1; level <= r.Level; level++ {
		for _, regList := range r.Get(level) {
			regNewLine := regexp.MustCompile(`\r|\n`)
			for _, v := range filelist {
				contentBytes, err := os.ReadFile(v.Fullpath)
				if err != nil {
					return nil, fmt.Errorf("get content of %v error: %v", v.Fullpath, err)
				}
				content := string(regNewLine.ReplaceAll(contentBytes, []byte{' '}))
				if regList.Match(content) {
					content = regList.Highlight(content)
					resultList = append(resultList, Result{
						Urlpath: v.GetUrlpath(),
						Text:    template.HTML(content),
					})
				}
			}
		}
		if len(resultList) > 0 {
			break
		}
	}
	return resultList, nil
}

/*
匹配优先级：
1，将整个keyword当作一个正则表达式
2，如果有空白字符，按空白字符分割，形成正则数组，按匹配比例降低优先级

按字符分割，形成正则数组，按匹配比例降低优先级。
！！！不能这么搞，一旦字符较多，会造成极大的性能压力
*/
func ParseKeyword(keyword string) *RegMap {
	m := &RegMap{
		Map:   make(map[int][]*RegList),
		Level: 1,
	}
	//将整个keyword当作一个正则表达式
	m.Add(NewRegList(100).Add(regexp.MustCompile(keyword)))

	regBlank := regexp.MustCompile(`\s+`)
	if regBlank.MatchString(keyword) {
		var list []*regexp.Regexp
		//按空白字符分割
		for _, v := range regBlank.Split(keyword, -1) {
			if len(v) > 0 {
				list = append(list, regexp.MustCompile(v))
			}
		}
		//按匹配比例降低优先级
		if len(list) > 0 {
			m.AddLevel()
			m.Add(NewRegList(100).Add(list...))
			m.AddLevel()
			m.Add(NewRegList(80).Add(list...))
			m.AddLevel()
			m.Add(NewRegList(50).Add(list...))
			m.AddLevel()
			m.Add(NewRegList(30).Add(list...))
			m.AddLevel()
			m.Add(NewRegList(1).Add(list...))
		}
	}

	return m
}

func SearchRegexp(dl *dir.DirList, keyword string) ([]Result, error) {
	var filelist []*File
	var errMsg string
	for _, v := range dl.List() {
		l, err := GetFileList(v, v.Path)
		if err != nil {
			errMsg += fmt.Sprintf("%v\n", err)
			continue
		}
		filelist = append(filelist, l...)
	}
	resultList, err := ParseKeyword(keyword).Search(filelist)
	if err != nil {
		errMsg += fmt.Sprintf("%v\n", err)
		return nil, fmt.Errorf(errMsg)
	}
	return resultList, nil
}

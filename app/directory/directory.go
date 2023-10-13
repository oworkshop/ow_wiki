package dir

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/oworkshop/ow_wiki/app/utils"
)

// Dir of wiki documents
type Dir struct {
	ID          string `json:"ID"`
	Name        string `json:"Name"`
	Path        string `json:"Path"`
	UseChildren bool   `json:"UseChildren"`
	IsPublic    bool   `json:"IsPublic"`
}

func (d *Dir) validate() error {
	if d.Path == "" {
		return fmt.Errorf("path of directory is empty")
	}
	return nil
}

func (d *Dir) clone(newName, newPath string) *Dir {
	new := *d
	if newName != "" {
		new.Name = newName
	}
	if newPath != "" {
		new.Path = newPath
	}
	return &new
}

func (d *Dir) Compute() ([]*Dir, error) {
	var result []*Dir

	if err := d.validate(); err != nil {
		return result, err
	}

	if !filepath.IsAbs(d.Path) {
		d.Path = filepath.Join(utils.Basepath(), d.Path)
	}
	if d.Name == "" {
		d.Name = filepath.Base(d.Path)
	}
	d.ID = d.Name

	if !d.UseChildren {
		result = append(result, d)
	} else {
		entryList, err := utils.GetSortedEntryList(d.Path)
		if err != nil {
			return result, err
		}
		for _, v := range entryList {
			if utils.IsDir(d.Path, v) {
				result = append(result,
					d.clone(v.Name(), filepath.Join(d.Path, v.Name())))
			}
		}
	}

	return result, nil
}

func (d *Dir) UpdateID(id string) *Dir {
	d.ID = id
	return d
}

type DirList struct {
	dirList []*Dir
	dirMap  map[string]*Dir
}

func (dl *DirList) Add(d *Dir) error {
	list, err := d.Compute()
	if err != nil {
		return err
	}
	reg := regexp.MustCompile(`^\w+$`)
	for k, v := range list {
		if _, ok := dl.Get(v.ID); ok || !reg.MatchString(v.ID) {
			v.UpdateID(fmt.Sprint(len(dl.dirList)))
		}
		dl.dirMap[v.ID] = list[k]
		dl.dirList = append(dl.dirList, list[k])
	}
	return nil
}

func (dl *DirList) Get(id string) (*Dir, bool) {
	d, ok := dl.dirMap[id]
	return d, ok
}

func (dl *DirList) List() []*Dir {
	return dl.dirList
}

func NewDirList(originalList []Dir) *DirList {
	dl := &DirList{
		dirMap:  make(map[string]*Dir),
		dirList: []*Dir{},
	}
	for _, v := range originalList {
		dl.Add(v.clone("", ""))
	}
	return dl
}

func SplitUrlpath(urlpath string) (string, string) {
	var idStr, subpath string
	urlpath = strings.TrimPrefix(urlpath, "/")
	pos := strings.Index(urlpath, "/")
	if pos < 0 {
		idStr, subpath = urlpath, ""
	} else {
		idStr, subpath = urlpath[0:pos], urlpath[pos+1:]
	}
	return idStr, subpath
}

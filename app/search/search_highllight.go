package search

import "regexp"

const (
	HightlightPreTag  = `<span class="highlight">`
	HightlightPostTag = `</span>`
)

type Part struct {
	Content   string
	Highlight bool
}

func (p *Part) String() string {
	if p.Highlight {
		return HightlightPreTag + p.Content + HightlightPostTag
	}
	return p.Content
}

func NewPart(content string, highlight bool) *Part {
	return &Part{
		Content:   content,
		Highlight: highlight,
	}
}

func HighlightPart(p *Part, reg *regexp.Regexp) []*Part {
	if p.Highlight {
		return []*Part{p}
	}
	indexList := reg.FindAllStringIndex(p.Content, -1)
	var result []*Part
	pos := 0
	for _, v := range indexList {
		if v[0] > pos {
			result = append(result, NewPart(p.Content[pos:v[0]], false))
			pos = v[1] + 1
		}
		result = append(result, NewPart(p.Content[v[0]:v[1]], true))
	}
	if len(p.Content)+1 > pos {
		result = append(result, NewPart(p.Content[pos:], false))
	}
	return result
}

func HighlightPartList(partList []*Part, reg *regexp.Regexp) []*Part {
	var result []*Part
	for k, v := range partList {
		if v.Highlight || !reg.MatchString(v.Content) {
			result = append(result, v)
			continue
		}
		var rest []*Part
		if len(partList) > k {
			rest = partList[k+1:]
		}
		result = append(result, HighlightPart(v, reg)...)
		result = append(result, HighlightPartList(rest, reg)...)
	}
	return result
}

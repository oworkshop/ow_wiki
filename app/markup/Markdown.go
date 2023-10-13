package markup

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
)

const (
	topContainer = "<div id='top-container' class='top-container'></div>\n"
)

// return TOC, body
func ConvertMd2Html(content string) (string, string) {
	var toc TOC
	var body []string
	var piece string

	if strings.Contains(content, "\r") {
		content = convertCRLF2LF(content)
	}

	regHead := regexp.MustCompile(`^(#+) +(.*)$`)

	for _, line := range strings.Split(content, "\n") {
		if regHead.MatchString(line) {
			if piece != "" {
				html := markdown.ToHTML([]byte(piece), nil, nil)
				body = append(body, string(html))
				piece = ""
			}

			m := regHead.FindStringSubmatch(line)
			s := toc.NewSection(len(m[1]))
			toc.Add(m[2], s, false)
			body = append(body,
				fmt.Sprintf(`<h%v id="sec-%v"><span class="section-number-%v">%v</span> %v</h%v>`,
					len(m[1])+1, s.String(), len(m[1])+1, s.String(), m[2], len(m[1])+1))
		} else {
			piece += line + "\n"
		}
	}

	if piece != "" {
		html := markdown.ToHTML([]byte(piece), nil, nil)
		body = append(body, string(html))
		piece = ""
	}

	body = alignTableCenter(body)
	body = removeEmptyThead(body)

	return toc.HTML(), topContainer + strings.Join(body, "\n")
}

func alignTableCenter(body []string) []string {
	for i, v := range body {
		v = strings.ReplaceAll(v, "<table>",
			`<div class="div-table" style="text-align: center;">
<table style="margin: auto">`)
		body[i] = strings.ReplaceAll(v, "</table>", "</table>\n</div>")
	}
	return body
}

/*
<thead>
<tr>
<th></th>
<th></th>
</tr>
</thead>
*/
func removeEmptyThead(body []string) []string {
	for i, v := range body {
		if !strings.Contains(v, "thead") {
			continue
		}
		var newV, thead []string
		isEmpty := true
		inThead := false
		for _, line := range strings.Split(v, "\n") {
			if line == "<thead>" {
				thead = append(thead, line)
				inThead = true
				continue
			}
			if line == "</thead>" {
				thead = append(thead, line)
				inThead = false
				if !isEmpty {
					newV = append(newV, thead...)
				}
				thead = nil
				isEmpty = true
				continue
			}
			if inThead {
				if strings.HasPrefix(line, "<th>") &&
					strings.HasSuffix(line, "</th>") &&
					line != "<th></th>" {
					isEmpty = false
				}
				thead = append(thead, line)
				continue
			}
			newV = append(newV, line)
		}
		body[i] = strings.Join(newV, "\n")
	}
	return body
}

func convertCRLF2LF(content string) string {
	var b strings.Builder
	var crlf string
	for _, v := range content {
		if v == '\r' || v == '\n' {
			crlf += string(v)
		} else {
			if len(crlf) > 0 {
				crlf = strings.ReplaceAll(crlf, "\r\n", "\n")
				crlf = strings.ReplaceAll(crlf, "\r", "\n")
				b.WriteString(crlf)
				crlf = ""
			}
			b.WriteRune(v)
		}
	}
	return b.String()
}

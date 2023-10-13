package cmd

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"github.com/oworkshop/ow_wiki/app/utils"
)

func genIndexLink(name, urlpath string, isDir bool) string {
	tCard := `<div class="index-card"><div class="index-card-body">%v</div></div>`
	showName := utils.TemplateExecuteReplacer.Replace(name)
	if isDir {
		link := fmt.Sprintf(`<a class="index-link" href="%v?lang=%v">%v</a> `,
			html.EscapeString(urlpath), mainConfig.DefaultLanguage, showName)
		return fmt.Sprintf(tCard, link)
	}
	if name == "favicon.ico" {
		return ""
	}
	link := fmt.Sprintf(`<a class="index-link" target="_blank" href="%v?lang=%v">%v</a> `,
		html.EscapeString(urlpath), mainConfig.DefaultLanguage, showName)
	if strings.HasSuffix(name, ".md") {
		link += fmt.Sprintf(`<a class="index-a-extra" target="_blank" href="%v?print&toc=n&lang=de&toc=y&lang=en&readme=相同参数，比如2个toc，则前面的有效">print</a>`,
			html.EscapeString(urlpath))
	}
	return fmt.Sprintf(tCard, link)
}

func getIndexContent(realpath, urlpath string) (string, error) {
	indexfile := filepath.Join(realpath, "index.md")
	_, err := os.Stat(indexfile)
	if err == nil {
		contentBytes, err := os.ReadFile(indexfile)
		if err != nil {
			return "", err
		}
		return string(contentBytes), nil
	}

	content := ``
	entryList, err := utils.GetSortedEntryList(realpath)
	if err != nil {
		return "", err
	}
	if len(entryList) == 0 {
		return "", nil
	}

	addHeader := false
	for _, v1 := range entryList {
		if utils.IsDir(realpath, v1) {
			continue
		}
		if !addHeader {
			content += "# .\n"
			addHeader = true
		}
		content += genIndexLink(v1.Name(), filepath.Join(urlpath, v1.Name()), false)
	}

	for _, v1 := range entryList {
		if !utils.IsDir(realpath, v1) {
			continue
		}
		el, err := utils.GetSortedEntryList(filepath.Join(realpath, v1.Name()))
		if err != nil {
			return "", err
		}
		if len(el) == 0 {
			continue
		}
		content += fmt.Sprintf("\n# %v <a class=\"index-a-extra\" href=\"/%v\">open</a>\n",
			v1.Name(), html.EscapeString(v1.Name()))
		for _, v2 := range el {
			content += genIndexLink(v2.Name(), filepath.Join(urlpath, v1.Name(), v2.Name()),
				utils.IsDir(filepath.Join(realpath, v1.Name()), v2))
		}
	}

	return content, nil
}

func getRootIndexContent() string {
	content := ``
	for _, ddir := range dirList.List() {
		entryList, err := utils.GetSortedEntryList(ddir.Path)
		if err != nil {
			content += fmt.Sprintf("`````\n%v\n`````\n\n", err)
			continue
		}
		if len(entryList) == 0 {
			continue
		}
		content += fmt.Sprintf("# %v <a class=\"index-a-extra\" href=\"/%v\">open</a>\n",
			ddir.Name, ddir.ID)
		addHeader := false
		for _, v1 := range entryList {
			if utils.IsDir(ddir.Path, v1) {
				continue
			}
			if !addHeader {
				content += "## .\n"
				addHeader = true
			}
			content += genIndexLink(v1.Name(), filepath.Join(ddir.ID, v1.Name()), false)
		}
		for _, v1 := range entryList {
			if !utils.IsDir(ddir.Path, v1) {
				continue
			}
			el, err := utils.GetSortedEntryList(filepath.Join(ddir.Path, v1.Name()))
			if err != nil {
				content += fmt.Sprintf("`````\n%v\n`````\n\n", err)
				continue
			}
			if len(el) == 0 {
				continue
			}
			content += fmt.Sprintf("\n## %v <a class=\"index-a-extra\" href=\"/%v\">open</a>\n",
				v1.Name(), html.EscapeString(filepath.Join(ddir.ID, v1.Name())))
			for _, v2 := range el {
				content += genIndexLink(v2.Name(), filepath.Join(ddir.ID, v1.Name(), v2.Name()),
					utils.IsDir(filepath.Join(ddir.Path, v1.Name()), v2))
			}
		}
		content += "\n\n"
	}
	return content
}

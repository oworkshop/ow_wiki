package cmd

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	dir "github.com/oworkshop/ow_wiki/app/directory"
	"github.com/oworkshop/ow_wiki/app/markup"
	"github.com/oworkshop/ow_wiki/app/search"
	"github.com/oworkshop/ow_wiki/app/utils"
)

func NewHandler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	urlpath := r.URL.Path
	query := r.URL.Query()

	if urlpath == "/favicon.ico" {
		http.ServeFile(w, r, filepath.Join(utils.Basepath(), "favicon.ico"))
		return
	}

	if urlpath == "/search" {
		handleSearch(w, query)
		return
	}

	lang := "en"
	if _, found := query["lang"]; found {
		lang = query["lang"][0]
	}

	dirID, subpath := dir.SplitUrlpath(urlpath)
	if dirID == "" && subpath == "" {
		handleContent(w, lang, "index", getRootIndexContent())
		return
	}

	d, ok := dirList.Get(dirID)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	realpath := filepath.Join(d.Path, subpath)

	fi, err := os.Stat(realpath)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	if fi.IsDir() {
		content, err := getIndexContent(realpath, urlpath)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		handleContent(w, lang, filepath.Base(realpath), content)
		return
	}

	if _, found := query["print"]; found {
		withTOC := false
		if _, found := query["toc"]; found {
			if query["toc"][0] == "y" {
				withTOC = true
			}
		}
		handlePrintPage(w, realpath, lang, filepath.Base(realpath), withTOC)
		return
	}

	if strings.HasSuffix(fi.Name(), ".md") {
		contentBytes, err := os.ReadFile(realpath)
		if err != nil {
			fmt.Fprintf(w, "read file %v error: %v", realpath, err)
			return
		}
		handleContent(w, lang, filepath.Base(realpath), string(contentBytes))
	} else {
		http.ServeFile(w, r, realpath)
	}
}

func handleSearch(w http.ResponseWriter, query url.Values) {
	t, err := template.ParseFiles(
		filepath.Join(utils.Basepath(), "template/html/search.html"))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	_, found := query["keyword"]
	if !found {
		io.WriteString(w, "no keyword found")
		return
	}
	resultList, err := search.SearchRegexp(dirList, query["keyword"][0])
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	t.Execute(w, resultList)
}

func handleContent(w http.ResponseWriter, lang, title, content string) {
	toc, body := markup.ConvertMd2Html(content)

	tpFile := filepath.Join(utils.Basepath(), "template/html/markdown.html")
	tp, err := template.ParseFiles(tpFile)
	if err != nil {
		fmt.Fprintf(w, "parse template %v error: %v", tpFile, err)
		return
	}
	c := struct {
		Lang  string
		Title string
		TOC   template.HTML
		Body  template.HTML
	}{
		Lang:  lang,
		Title: title,
		TOC:   template.HTML(toc),
		Body:  template.HTML(body),
	}
	tp.Execute(w, c)
}

func handlePrintPage(w http.ResponseWriter, mdfile, lang, title string, withTOC bool) {
	contentBytes, err := os.ReadFile(mdfile)
	if err != nil {
		fmt.Fprintf(w, "open file %v error: %v", mdfile, err)
		return
	}
	toc, body := markup.ConvertMd2Html(string(contentBytes))
	if !withTOC {
		toc = ""
	}

	tpFile := filepath.Join(utils.Basepath(), "template/html/print.html")
	tp, err := template.ParseFiles(tpFile)
	if err != nil {
		fmt.Fprintf(w, "parse template %v error: %v", tpFile, err)
		return
	}

	c := struct {
		Lang  string
		Title string
		TOC   template.HTML
		Body  template.HTML
	}{
		Lang:  lang,
		Title: title,
		TOC:   template.HTML(toc),
		Body:  template.HTML(body),
	}
	tp.Execute(w, c)
}

func checkAccessAllowed(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	urlpath := r.URL.Path

	dirID, _ := dir.SplitUrlpath(urlpath)
	if d, ok := dirList.Get(dirID); ok && d.IsPublic {
		next(w, r)
		return
	}

	realip := utils.GetIPFromRequest(r).String()
	for _, v := range allowList {
		if v.MatchString(realip) {
			next(w, r)
			return
		}
	}
	io.WriteString(w, "Access denied!\nYour ip is "+realip+"\n")
}

package search

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	dir "github.com/oworkshop/ow_wiki/app/directory"
	"github.com/oworkshop/ow_wiki/app/utils"
)

type File struct {
	DDir     *dir.Dir
	Fullpath string
}

func (f *File) GetUrlpath() string {
	urlpath, err := filepath.Rel(f.DDir.Path, f.Fullpath)
	if err != nil {
		fmt.Printf("get urlpath of %v error: %v\n", f.Fullpath, err)
		return ""
	}
	return filepath.Join(f.DDir.ID, urlpath)
}

type Result struct {
	Urlpath string
	Text    template.HTML
}

func GetFileList(ddir *dir.Dir, directory string) ([]*File, error) {
	entryList, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	var result []*File
	for _, v := range entryList {
		if utils.IsDir(directory, v) {
			l, err := GetFileList(ddir, filepath.Join(directory, v.Name()))
			if err != nil {
				return nil, err
			}
			result = append(result, l...)
		} else {
			if strings.HasSuffix(v.Name(), ".md") {
				result = append(result, &File{
					DDir:     ddir,
					Fullpath: filepath.Join(directory, v.Name()),
				})
			}
		}
	}
	return result, nil
}

func GetFileID(f string) string {
	h := sha256.New()
	h.Write([]byte(f))
	return hex.EncodeToString(h.Sum(nil))
}

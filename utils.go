package poutil

import (
	"io/ioutil"
	"path"
)

/*
func ScanLocaleDirectory(localeDir string) {
	langs, err := GetLocaleLanguages(localeDir)
	if err != nil {
		panic(err)
	}
	info := map[string]string{}
	for _, lang := range *langs {
		info[lang]
	}
}
*/

func GetLocaleLanguages(localeDir string) (*[]string, error) {
	var langs = []string{}
	fileInfos, err := ioutil.ReadDir(localeDir)
	if err != nil {
		return nil, err
	}
	for _, fi := range fileInfos {
		if fi.IsDir() {
			langs = append(langs, fi.Name())
		}
	}
	return &langs, nil
}

func BuildLCMessageFilePath(localeDir, lang, domain string) string {
	return path.Join(localeDir, lang, "LC_MESSAGES", domain) + ".po"
}

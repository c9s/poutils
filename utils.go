package poutil

import "io/ioutil"

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

package main

import (
	"flag"
	"fmt"
	"github.com/c9s/poutil"
	"github.com/c9s/poutil/pofile"
	"io/ioutil"
	"os"
	"path"
)

var domainOpt = flag.String("domain", "", "domain")
var localeOpt = flag.String("locale", "locale", "locale directory")
var jsonDirOpt = flag.String("json-dir", "public/locale", "json output directory")

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}
	return true
}

func Die(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

func ParallelConvertPoFileToJsonFile(poFile, jsonFile string, c chan bool) bool {
	dict, err := pofile.ParseFile(poFile)
	if err != nil {
		Die(err)
	}

	fmt.Println("Writing JSON", jsonFile)
	jsonOutput := dict.JSONString()
	err = ioutil.WriteFile(jsonFile, []byte(jsonOutput), 0666)
	if err != nil {
		Die("Can not write json file", jsonFile)
	}
	c <- true
	return true
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("Usage: ")
		fmt.Println("   po2json [file] > app.json")
		fmt.Println("   po2json --domain app")
		fmt.Println("   po2json --domain app --locale path/to/locale_dir --json-dir public/locale")
		os.Exit(0)
	}

	// if domain is specified, we should get the po files from locale directory.
	if *domainOpt != "" {
		var domain = *domainOpt
		var localeDir = *localeOpt
		var jsonDir = *jsonDirOpt

		langs, err := poutil.GetLocaleLanguages(localeDir)
		if err != nil {
			Die(err)
		}

		os.MkdirAll(jsonDir, 0777)

		var c = make(chan bool)
		var routineCounter = len(*langs)

		for _, lang := range *langs {
			var poFile = poutil.BuildLCMessageFilePath(localeDir, lang, domain)
			var jsonFile = path.Join(jsonDir, lang) + ".json"

			if FileExists(poFile) {
				fmt.Println("Start Processing", poFile)
				go ParallelConvertPoFileToJsonFile(poFile, jsonFile, c)
			}
		}
		for routineCounter > 0 {
			<-c
			routineCounter--
		}
		fmt.Println("Done")
	} else {
		files := os.Args[1:]

		mainDict := pofile.NewDictionary()

		for _, filename := range files {
			if !FileExists(filename) {
				Die(filename, "does not exist.")
			}

			err := mainDict.MergeFile(filename)
			if err != nil {
				Die("PO File Parsing Error", err)
			}
		}
		fmt.Println(mainDict.JSONString())
	}
}

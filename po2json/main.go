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

func ConvertPoToJson(poFile, jsonFile string, c chan bool) bool {
	dict, err := pofile.ParseFile(poFile)
	if err != nil {
		fmt.Println("PO File Parsing Error", err)
		os.Exit(1)
	}

	fmt.Println("Writing JSON", jsonFile)
	jsonOutput := dict.JSONString()
	err = ioutil.WriteFile(jsonFile, []byte(jsonOutput), 0666)
	if err != nil {
		fmt.Println("Can not write json file", jsonFile)
		os.Exit(1)
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
			fmt.Println(err)
			os.Exit(1)
		}

		os.MkdirAll(jsonDir, 0777)

		var c = make(chan bool)
		var routineCounter = len(*langs)

		for _, lang := range *langs {
			var poFile = path.Join(localeDir, lang, "LC_MESSAGES", domain) + ".po"
			var jsonFile = path.Join(jsonDir, lang) + ".json"

			if FileExists(poFile) {
				fmt.Println("Start Processing", poFile)
				go ConvertPoToJson(poFile, jsonFile, c)
			}
		}
		for routineCounter > 0 {
			<-c
			routineCounter--
		}
		fmt.Println("Done")
	} else {
		filename := os.Args[1]
		if _, err := os.Stat(filename); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dict, err := pofile.ParseFile(filename)
		if err != nil {
			fmt.Println("PO File Parsing Error", err)
			os.Exit(1)
		}
		fmt.Println(dict)
	}
}

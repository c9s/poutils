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
var outputDirOpt = flag.String("output-dir", "output", "CSV output directory.")

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

func ParallelConvertPoFileToCSVFile(poFile, csvFile string, c chan bool) bool {
	dict, err := pofile.ParseFile(poFile)
	if err != nil {
		Die(err)
	}

	fmt.Println("Writing CSV", csvFile)
	if err = dict.WriteCSVFile(csvFile); err != nil {
		Die("Can not write json file", csvFile)
	}
	c <- true
	return true
}

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("Usage: ")
		fmt.Println("   po2csv [file] > app.json")
		fmt.Println("   po2csv --domain app")
		fmt.Println("   po2csv --domain app --locale path/to/locale_dir --output-dir csv")
		os.Exit(0)
	}

	// if domain is specified, we should get the po files from locale directory.
	if *domainOpt != "" {
		var domain = *domainOpt
		var localeDir = *localeOpt
		var outputDir = *outputDirOpt

		langs, err := poutil.GetLocaleLanguages(localeDir)
		if err != nil {
			Die(err)
		}

		os.MkdirAll(outputDir, 0777)

		var c = make(chan bool)
		var routineCounter = len(*langs)

		for _, lang := range *langs {
			var poFile = poutil.BuildLCMessageFilePath(localeDir, lang, domain)
			var csvFile = path.Join(outputDir, lang) + ".csv"
			if FileExists(poFile) {
				fmt.Println("Start Processing", poFile)
				go ParallelConvertPoFileToCSVFile(poFile, csvFile, c)
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
		fmt.Println(mainDict.CSVString())
	}
}

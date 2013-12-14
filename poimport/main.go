package main

import (
	"flag"
	"fmt"
	"github.com/c9s/poutil/pofile"
	"os"
	// "path"
)

/*
poimport imports messages from csv/json to an existing po file.


poimport --csv file.csv --to locale/en/LC_MESSAGES/jifty.po
poimport --json file.json --to locale/en/LC_MESSAGES/jifty.po
*/

var csvFileOpt = flag.String("csv", "", "csv file")
var jsonFileOpt = flag.String("json", "", "json file")
var toOpt = flag.String("to", "", "po file")

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

func main() {
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Println("Usage: ")
		fmt.Println("   poimport --csv en.csv --to en.po")
		fmt.Println("   poimport --json en.json --to en.po")
		os.Exit(0)
	}

	if *toOpt == "" {
		Die("--to [file] option is required.")
	}

	if *csvFileOpt == "" || *jsonFileOpt == "" {
		Die("--csv [file] or --json [file] option is required.")
	}

	file := pofile.NewPOFile()
	if err := file.LoadFile(*toOpt); err != nil {
		Die(err)
	}

}

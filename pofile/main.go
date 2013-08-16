package pofile

import (
	"regexp"
)

const (
	STATE_COMMENT  = iota
	STATE_MSGID    // waiting for msgid
	STATE_MSGSTR   // waiting for msgstr
	STATE_COMPLETE // complete state, waiting for comment or msgid
)

var commentRegExp = regexp.MustCompile("^\\s*#")
var emptyLineRegExp = regexp.MustCompile("^\\s*$")
var msgIdRegExp = regexp.MustCompile("^msgid\\s+\"(.*)\"")
var msgStrRegExp = regexp.MustCompile("^msgstr\\s+\"(.*)\"")
var stringRegExp = regexp.MustCompile("\"(.*)\"")

func ParseFiles(files []string) (*Dictionary, error) {
	mainDict := NewDictionary()
	for _, filename := range files {
		err := mainDict.MergeFile(filename)
		if err != nil {
			return nil, err
		}
	}
	return mainDict, nil
}

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

func ParseMessagesFromFile(filename string) (*Dictionary, error) {
	dict := NewDictionary()
	err := dict.ParseAndLoadFromFile(filename)
	return dict, err
}

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

package pofile

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
)

// import "fmt"
// import "bufio"

/*

To create a empty dictinoary:

	dict := pofile.NewDictionary()
	dict := pofile.Dictionary{}

To add message:

	dict.AddMessage(msgId, msgStr)

To remove message:

	dict.RemoveMessage(msgId)

*/

type Dictionary map[string]string

func NewDictionary() *Dictionary {
	return &Dictionary{}
}

func (self Dictionary) AddMessage(msgId string, msgStr string) {
	self[msgId] = msgStr
}

func (self Dictionary) HasMessage(msgId string) bool {
	_, ok := self[msgId]
	return ok
}

func (self Dictionary) RemoveMessage(msgId string) {
	delete(self, msgId)
}

func (self Dictionary) Merge(dict *Dictionary) {
	for key, val := range *dict {
		self[key] = val
	}
}

func (self Dictionary) MergeFile(filename string) error {
	newDict, err := ParseFile(filename)
	if err != nil {
		return err
	}
	self.Merge(newDict)
	return nil
}

func (self Dictionary) CSVString() string {
	var buf = bytes.NewBufferString("")
	var writer = csv.NewWriter(buf)
	for key, val := range self {
		writer.Write([]string{key, val})
	}
	writer.Flush()
	return buf.String()
}

func (self Dictionary) JSONString() string {
	jsonBytes, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func (self Dictionary) String() string {
	return self.JSONString()
}

func (self Dictionary) WriteJSONFile(filepath string) error {
	var output = self.JSONString()
	return ioutil.WriteFile(filepath, []byte(output), 0666)
}

func (self Dictionary) WriteCSVFile(filepath string) error {
	var output = self.CSVString()
	return ioutil.WriteFile(filepath, []byte(output), 0666)
}

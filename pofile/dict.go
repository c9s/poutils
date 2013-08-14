package pofile

import "encoding/json"
import "encoding/csv"

// import "fmt"
// import "bufio"
import "bytes"

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
	return new(Dictionary)
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
	for key, value := range *dict {
		self[key] = value
	}
}

func (self Dictionary) MergeFile(filename string) error {
	dict, err := ParseFile(filename)
	if err != nil {
		return err
	}
	self.Merge(dict)
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

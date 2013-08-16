package pofile

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// import "fmt"
// import "bufio"

/*

To create a empty dictionary:

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

func (self Dictionary) ParseAndLoadFromFile(filename string) error {
	// process(filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return self.ParseAndLoad(string(bytes))
}

func (self Dictionary) ParseAndLoad(content string) error {
	lines := strings.Split(content, "\n")
	lastMsgId := []string{}
	lastMsgStr := []string{}
	lastComments := []string{}

	state := STATE_COMPLETE

	for _, line := range lines {
		if len(line) == 0 || emptyLineRegExp.MatchString(line) { // skip empty lines
			if state == STATE_MSGSTR {
				self.AddMessage(strings.Join(lastMsgId, ""), strings.Join(lastMsgStr, ""))
				lastMsgId = []string{}
				lastMsgStr = []string{}
				lastComments = []string{}
				state = STATE_COMPLETE
			}
			continue
		}

		if line[0] == '#' || commentRegExp.MatchString(line) {
			lastComments = append(lastComments, line)
			state = STATE_COMMENT
			continue
		}

		if strings.HasPrefix(line, "msgid") || msgIdRegExp.MatchString(line) {

			state = STATE_MSGID
			msgId := msgIdRegExp.FindStringSubmatch(line)[1]
			lastMsgId = append(lastMsgId, msgId)

		} else if strings.HasPrefix(line, "msgstr") || msgStrRegExp.MatchString(line) {
			state = STATE_MSGSTR
			msgStr := msgStrRegExp.FindStringSubmatch(line)[1]
			lastMsgStr = append(lastMsgStr, msgStr)
		} else if stringRegExp.MatchString(line) {
			var str = stringRegExp.FindStringSubmatch(line)[1]
			if state == STATE_MSGID {
				lastMsgId = append(lastMsgId, str)
			} else if state == STATE_MSGSTR {
				lastMsgStr = append(lastMsgStr, str)
			}
		}
	}
	return nil
}

func (self Dictionary) Merge(dict *Dictionary) {
	for key, val := range *dict {
		self[key] = val
	}
}

func (self Dictionary) MergeFile(filename string) error {
	newDict := NewDictionary()
	if err := newDict.ParseAndLoadFromFile(filename); err != nil {
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

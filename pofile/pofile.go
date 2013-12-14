package pofile

import (
	"fmt"
	"io/ioutil"
	"strings"
)

/*
A po file contains multiple messages
*/
type POFile struct {
	// A po file contains multiple messages
	Messages []Message

	// an ID map for checking message existence.
	IdMap map[string]bool
}

/*
A message struct contains multiple lines of MsgId, MsgStr and comment.
*/
type Message struct {
	MsgIds   []string
	MsgStrs  []string
	Comments []string
}

/*
create a new POFile object with empty message, empty id map
*/
func NewPOFile() *POFile {
	return &POFile{Messages: []Message{}, IdMap: map[string]bool{}}
}

/*
Load message from a po file
*/
func (self *POFile) LoadFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return self.ParseAndLoad(string(bytes))
}

/*
Get the count of messages
*/
func (self *POFile) Length() int {
	return len(self.Messages)
}

/*
Output the message to a string, supports Stringer interface.
*/
func (self *POFile) String() string {
	var output string = ""
	for _, msg := range self.Messages {

		output += strings.Join(msg.Comments, "\n") + "\n"

		output += "msgid "
		for _, str := range msg.MsgIds {
			output += "\"" + str + "\"\n"
		}

		output += "msgstr "
		for _, str := range msg.MsgStrs {
			output += "\"" + str + "\"\n"
		}
		output += "\n"
	}
	return output
}

/*
Write current POFile struct to a file.
*/
func (self *POFile) WriteFile(filename string) error {
	output := self.String()
	return ioutil.WriteFile(filename, []byte(output), 0666)
}

/*
Import a dictionary structure (map[string]string) to the current POFile object.
*/
func (self *POFile) ImportDictionary(dict *Dictionary, override bool) {
	for msgId, msgStr := range *dict {
		if _, ok := self.IdMap[msgId]; ok && !override {
			continue
		}
		self.IdMap[msgId] = true
		self.Messages = append(self.Messages, Message{
			MsgIds:   []string{msgId},
			MsgStrs:  []string{msgStr},
			Comments: []string{},
		})
	}
}

/*
Parse & load by givn po message content
*/
func (self *POFile) ParseAndLoad(content string) error {
	lines := strings.Split(content, "\n")
	ids := []string{}
	strs := []string{}
	comments := []string{}

	state := STATE_COMPLETE

	for linenr, line := range lines {
		if len(line) == 0 || EmptyLineRegExp.MatchString(line) { // skip empty lines
			if state == STATE_MSGSTR {
				var msgid = strings.Join(ids, "")

				if _, ok := self.IdMap[msgid]; ok {
					fmt.Println("Duplicate message", msgid)
					continue
				}

				self.IdMap[msgid] = true
				self.Messages = append(self.Messages, Message{MsgIds: ids, MsgStrs: strs, Comments: comments})

				// reset all stacks
				ids = []string{}
				strs = []string{}
				comments = []string{}
				state = STATE_COMPLETE
			}
			continue
		}

		if line[0] == '#' || CommentRegExp.MatchString(line) {
			comments = append(comments, line)
			state = STATE_COMMENT
			continue
		}

		if strings.HasPrefix(line, "msgid") || MsgIdRegExp.MatchString(line) {
			if state == STATE_MSGID {
				panic(fmt.Sprintf("Duplicate msgid statement at line %d", linenr))
			}

			state = STATE_MSGID
			msgId := MsgIdRegExp.FindStringSubmatch(line)[1]
			ids = append(ids, msgId)

		} else if strings.HasPrefix(line, "msgstr") || MsgStrRegExp.MatchString(line) {

			if state == STATE_MSGSTR {
				panic(fmt.Sprintf("Duplicate msgstr statement at line %d", linenr))
			}

			state = STATE_MSGSTR
			msgStr := MsgStrRegExp.FindStringSubmatch(line)[1]
			strs = append(strs, msgStr)
		} else if StringRegExp.MatchString(line) {
			var str = StringRegExp.FindStringSubmatch(line)[1]
			if state == STATE_MSGID {
				ids = append(ids, str)
			} else if state == STATE_MSGSTR {
				strs = append(strs, str)
			}
		}
	}
	return nil
}

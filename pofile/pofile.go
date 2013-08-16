package pofile

import (
	"io/ioutil"
	"strings"
)

type POFile struct {
	// read-only dictionary
	Dictionary Dictionary
	Messages   []Message
}

type Message struct {
	MsgId     string
	MsgString string
	Comments  []string
}

func (self Message) AppendComment(comment string) {
	self.Comments = append(self.Comments, comment)
}

func (self *POFile) ParseFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	_ = bytes
	return nil
}

func (self POFile) ParseAndLoad(content string) error {
	lines := strings.Split(content, "\n")
	lastMsgId := []string{}
	lastMsgStr := []string{}
	lastComments := []string{}

	dict := NewDictionary()

	state := STATE_COMPLETE

	for _, line := range lines {
		if len(line) == 0 || EmptyLineRegExp.MatchString(line) { // skip empty lines
			if state == STATE_MSGSTR {
				dict.AddMessage(strings.Join(lastMsgId, ""), strings.Join(lastMsgStr, ""))
				lastMsgId = []string{}
				lastMsgStr = []string{}
				lastComments = []string{}
				state = STATE_COMPLETE
			}
			continue
		}

		if line[0] == '#' || CommentRegExp.MatchString(line) {
			lastComments = append(lastComments, line)
			state = STATE_COMMENT
			continue
		}

		if strings.HasPrefix(line, "msgid") || MsgIdRegExp.MatchString(line) {

			state = STATE_MSGID
			msgId := MsgIdRegExp.FindStringSubmatch(line)[1]
			lastMsgId = append(lastMsgId, msgId)

		} else if strings.HasPrefix(line, "msgstr") || MsgStrRegExp.MatchString(line) {
			state = STATE_MSGSTR
			msgStr := MsgStrRegExp.FindStringSubmatch(line)[1]
			lastMsgStr = append(lastMsgStr, msgStr)
		} else if StringRegExp.MatchString(line) {
			var str = StringRegExp.FindStringSubmatch(line)[1]
			if state == STATE_MSGID {
				lastMsgId = append(lastMsgId, str)
			} else if state == STATE_MSGSTR {
				lastMsgStr = append(lastMsgStr, str)
			}
		}
	}
	return nil
}

package pofile

import (
	"io/ioutil"
	"strings"
)

type POFile struct {
	// read-only dictionary
	Dictionary Dictionary
	Comments   map[string]string
}

func NewPOFile() *POFile {
	return &POFile{Dictionary: Dictionary{}, Comments: map[string]string{}}
}

func (self POFile) LoadFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return self.ParseAndLoad(string(bytes))
}

func (self POFile) Length() int {
	return len(self.Dictionary)
}

func (self POFile) ParseAndLoad(content string) error {
	lines := strings.Split(content, "\n")
	lastMsgId := []string{}
	lastMsgStr := []string{}
	lastComments := []string{}

	state := STATE_COMPLETE

	for _, line := range lines {
		if len(line) == 0 || EmptyLineRegExp.MatchString(line) { // skip empty lines
			if state == STATE_MSGSTR {
				msgId := strings.Join(lastMsgId, "")

				// map assignment is faster.
				self.Dictionary[msgId] = strings.Join(lastMsgStr, "")
				self.Comments[msgId] = strings.Join(lastComments, "")

				// reset all stacks
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

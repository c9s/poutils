package pofile

import (
	"io/ioutil"
	"regexp"
	"strings"
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

func ParseMessages(content string) (*Dictionary, error) {
	lines := strings.Split(content, "\n")
	lastMsgId := []string{}
	lastMsgStr := []string{}
	lastComments := []string{}

	dictionary := Dictionary{}

	state := STATE_COMPLETE

	for _, line := range lines {
		if len(line) == 0 || emptyLineRegExp.MatchString(line) { // skip empty lines
			if state == STATE_MSGSTR {
				dictionary.AddMessage(strings.Join(lastMsgId, ""), strings.Join(lastMsgStr, ""))
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
	return &dictionary, nil
}

func ParseMessagesFromFile(filename string) (*Dictionary, error) {
	// process(filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseMessages(string(bytes))
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

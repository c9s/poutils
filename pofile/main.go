package pofile

import (
	"io/ioutil"
	"regexp"
	"strings"
)

const (
	STATE_MSGID    = iota // waiting for msgid
	STATE_MSGSTR          // waiting for msgstr
	STATE_COMPLETE        // complete state
)

func ParseContent(content string) (*Dictionary, error) {
	lines := strings.Split(content, "\n")
	lastMsgId := []string{}
	lastMsgStr := []string{}
	lastComments = []string{}

	dictionary := Dictionary{}

	state := STATE_MSGID

	commentRegExp := regexp.MustCompile("^\\s*#")
	emptyLineRegExp := regexp.MustCompile("^\\s*$")
	msgIdRegExp := regexp.MustCompile("^msgid\\s+\"(.*)\"")
	msgStrRegExp := regexp.MustCompile("^msgstr\\s+\"(.*)\"")
	stringRegExp := regexp.MustCompile("\"(.*)\"")

	for _, line := range lines {
		if len(line) == 0 || emptyLineRegExp.MatchString(line) { // skip empty lines
			continue
		}

		if line[0] == '#' ||
			commentRegExp.MatchString(line) {
			lastComments = append(lastComments, line)
			continue
		}

		if strings.HasPrefix(line, "msgid") || msgIdRegExp.MatchString(line) {
			if len(lastMsgId) > 0 && len(lastMsgStr) > 0 {
				// push to the dictionary
				dictionary.AddMessage(strings.Join(lastMsgId, ""), strings.Join(lastMsgStr, ""))
				lastMsgId = []string{}
				lastMsgStr = []string{}
				lastComments = []string{}
			}

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

func ParseFile(filename string) (*Dictionary, error) {
	// process(filename)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseContent(string(bytes))
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

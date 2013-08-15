package pofile

import "io/ioutil"

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

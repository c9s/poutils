package pofile

import "encoding/json"

type Dictionary map[string]string

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

func (self Dictionary) JSONString() string {
	jsonBytes, err := json.MarshalIndent(self, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func (self Dictionary) Merge(dict *Dictionary) {
	for key, value := range *dict {
		self[key] = value
	}
}

package pofile

import (
	"testing"
)

func TestDictionary(t *testing.T) {
	dict := Dictionary{}
	dict.AddMessage("en", "English")

	if !dict.HasMessage("en") {
		t.Fatal("msgid should be defined")
	}

	if val, ok := dict["en"]; ok {
		if val != "English" {
			t.Fatal("Wrong msgstr")
		}
	} else {
		t.Fatal("msgid not found")
	}
	dict.RemoveMessage("en")
	dict.AddMessage("zh_TW", "繁體中文")

	json := dict.JSONString()
	t.Log(json)
	if json == "" {
		t.Fatal("Can not encoding json")
	}

	csv := dict.CSVString()
	t.Log(csv)
	_ = csv
}

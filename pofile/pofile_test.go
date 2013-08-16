package pofile

import (
	"os"
	"testing"
)

func TestPoFile(t *testing.T) {
	pofile := NewPOFile()
	if err := pofile.LoadFile("locale/en/LC_MESSAGES/jifty.po"); err != nil {
		t.Fatal(err)
	}
	if pofile.Length() == 0 {
		t.Log("Length Error")
	}
	if err := pofile.WriteFile("test_output.po"); err != nil {
		t.Fatal(err)
		return
	}
	os.Remove("test_output.po")
}

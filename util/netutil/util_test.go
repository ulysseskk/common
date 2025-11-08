package netutil

import (
	"bytes"
	"os"
	"testing"
)

func Test_parseSocktab(t *testing.T) {
	data, err := os.ReadFile("tcp.out")
	if err != nil {
		t.Fatal(err)
	}
	buffer := bytes.NewBuffer(data)
	result, err := parseSocktab(buffer, func(entry *SockTabEntry) bool {
		return entry.State == Listen
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("no result")
	}

}

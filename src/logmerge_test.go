package main

import (
	"bufio"
	"strings"
	"testing"
)

func Test_readLine(t *testing.T) {
	reader := strings.NewReader("foo")
	scanner := bufio.NewScanner(reader)

	var line string
	var eof bool

	// First read
	line, eof = readLine(scanner)
	if eof {
		t.Error("Unexpected EOF")
	}
	if line != "foo" {
		t.Errorf("Wrong result. Expected %s, got %s\n", "foo", line)
	}

	// Second read
	line, eof = readLine(scanner)
	if !eof {
		t.Error("EOF expected")
	}
	if line != "" {
		t.Errorf("Should be empty. Got %s", line)
	}
}

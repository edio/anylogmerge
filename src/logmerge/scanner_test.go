package logmerge

import (
	"strings"
	"testing"
)

func Test_readLine_oneline(t *testing.T) {
	reader := strings.NewReader("foo")
	scanner := newScannerDefault(reader)

	// First read
	scanner.readLine()
	if scanner.eof {
		t.Error("Unexpected EOF")
	}
	if scanner.Line() != "foo" {
		t.Errorf("Wrong result. Expected '%s', got '%s'\n", "foo", scanner.Line())
	}

	// Second read
	scanner.readLine()
	if !scanner.eof {
		t.Error("EOF expected")
	}
	if scanner.Line() != "" {
		t.Errorf("Should be empty. Got '%s'\n", scanner.Line())
	}
}

func Test_readLine(t *testing.T) {
	reader := strings.NewReader("foo\nbar\n")
	scanner := newScanner(reader, func(line string) string { return string(line[0]) })

	// First read
	scanner.readLine()
	if scanner.eof {
		t.Error("Unexpected EOF")
	}
	if scanner.Line() != "foo" {
		t.Errorf("Wrong result. Expected '%s', got '%s'\n", "foo", scanner.Line())
	}
	if scanner.SortKey() != "f" {
		t.Errorf("Wrong result. Expected '%s', got '%s'\n", "f", scanner.Line())
	}

	// Second read
	scanner.readLine()
	if scanner.eof {
		t.Error("Unexpected EOF")
	}
	if scanner.Line() != "bar" {
		t.Errorf("Wrong result. Expected '%s', got '%s'\n", "bar", scanner.Line())
	}
	if scanner.SortKey() != "b" {
		t.Errorf("Wrong result. Expected '%s', got '%s'\n", "b", scanner.Line())
	}
}

func Test_readLine_empty(t *testing.T) {
	reader := strings.NewReader("\n")
	scanner := newScannerDefault(reader)

	// First read
	scanner.readLine()
	if scanner.eof {
		t.Error("Unexpected EOF")
	}
	if scanner.Line() != "" {
		t.Errorf("Wrong result. Expected empty, got '%s'\n", scanner.Line())
	}

	// Second read
	scanner.readLine()
	if !scanner.eof {
		t.Error("EOF Expected")
	}
}

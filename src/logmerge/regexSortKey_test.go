package logmerge

import (
	"testing"
)

func Test_RegexSortKey(t *testing.T) {
	sortKeyFunc := RegexSortKey("a([a-z])")

	key := sortKeyFunc("abacad")

	if key != "b" {
		t.Errorf("Expected '%s', got '%s'", "b", key)
	}

	key = sortKeyFunc("bacad")

	if key != "c" {
		t.Errorf("Expected '%s', got '%s'", "c", key)
	}
}

func Test_RegexSortKey_nomatch(t *testing.T) {
	sortKeyFunc := RegexSortKey("^a([a-z])$")

	key := sortKeyFunc("abacad")

	if key != "" {
		t.Errorf("Expected empty key, got '%s'", "b", key)
	}
}

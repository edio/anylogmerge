package logmerge

import (
	"testing"
)

func Test_parseSlice(t *testing.T) {
	parsed := parseSlice("10:42")

	if parsed[0] != 10 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 0, 10, parsed[0])
	}
	if parsed[1] != 42 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 1, 42, parsed[1])
	}
}

func Test_parseSlice_startOnly(t *testing.T) {
	parsed := parseSlice("10:")

	if parsed[0] != 10 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 0, 10, parsed[0])
	}
	if parsed[1] != -1 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 1, -1, parsed[1])
	}
}

func Test_parseSlice_endOnly(t *testing.T) {
	parsed := parseSlice(":10")

	if parsed[0] != 0 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 0, 0, parsed[0])
	}
	if parsed[1] != 10 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 1, 10, parsed[1])
	}
}

func Test_parseSlice_noSlice(t *testing.T) {
	parsed := parseSlice(":")

	if parsed[0] != 0 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 0, 0, parsed[0])
	}
	if parsed[1] != -1 {
		t.Errorf("slice[%d] is wrong. Expected %d, got %d", 1, -1, parsed[1])
	}
}

func Test_parseColumnSpec(t *testing.T) {
	parsed := parseColumnSpec(":10,20:42,50:")

	if len(parsed) != 3 {
		t.Errorf("result length is wrong. Expected %d, got %d", 3, len(parsed))
	}
	if parsed[0] != [2]int{0, 10} {
		t.Errorf("slice[0] is wrong: %s", parsed[0])
	}
	if parsed[1] != [2]int{20, 42} {
		t.Errorf("slice[0] is wrong: %s", parsed[1])
	}
	if parsed[2] != [2]int{50, -1} {
		t.Errorf("slice[0] is wrong: %s", parsed[2])
	}
}

func Test_ColumnSortKey(t *testing.T) {
	key := ColumnSortKey(":3,4:9,38:")("The quick brown fox jumps over a lazy dog")
	if key != "Thequickdog" {
		t.Errorf("Key is incorret. Expected '%s'. Got '%s'", "Thequickdog", key)
	}
}

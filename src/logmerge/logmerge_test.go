package logmerge

import (
	"testing"
)

func Test_remove(t *testing.T) {
	slice := []*SortableScanner{
		&SortableScanner{sortKey: "1"},
		&SortableScanner{sortKey: "2"},
		&SortableScanner{sortKey: "3"},
		&SortableScanner{sortKey: "4"},
	}

	s := remove(slice, 2)

	if len(s) != 3 {
		t.Error("Length should have been decreased")
	}
	if s[2].sortKey == "3" {
		t.Error("New element expected on removed position")
	}
}

func Test_remove_first(t *testing.T) {
	slice := []*SortableScanner{
		&SortableScanner{sortKey: "1"},
		&SortableScanner{sortKey: "2"},
		&SortableScanner{sortKey: "3"},
		&SortableScanner{sortKey: "4"},
	}

	s := remove(slice, 0)

	if len(s) != 3 {
		t.Error("Length should have been decreased")
	}
	if s[0].sortKey == "1" {
		t.Error("New element expected on removed position")
	}
}

func Test_remove_last(t *testing.T) {
	slice := []*SortableScanner{
		&SortableScanner{sortKey: "1"},
		&SortableScanner{sortKey: "2"},
		&SortableScanner{sortKey: "3"},
		&SortableScanner{sortKey: "4"},
	}

	s := remove(slice, 3)

	if len(s) != 3 {
		t.Error("Length should have been decreased")
	}
}

func Test_remove_one(t *testing.T) {
	slice := []*SortableScanner{
		&SortableScanner{sortKey: "1"},
	}

	s := remove(slice, 0)

	if len(s) != 0 {
		t.Error("Length should have been decreased")
	}
}

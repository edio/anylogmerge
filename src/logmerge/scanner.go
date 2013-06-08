package logmerge

import (
	"bufio"
	"io"
)

/* Extracts sort key from the whole line */
type SortKeyFunc func(line string) (key string)

/* Default sort key is the whole line itself */
func DefaultSortKey(line string) string {
	return line
}

type SortableScanner struct {
	scanner      *bufio.Scanner
	keyExtractor SortKeyFunc
	sortKey      string
	eof          bool
}

func (scanner *SortableScanner) readLine() {
	if scanner.scanner.Scan() {
		scanner.sortKey = scanner.keyExtractor(scanner.scanner.Text())
	} else if err := scanner.scanner.Err(); err == nil {
		scanner.eof = true
	} else {
		panic(err)
	}
}

/* Line last read */
func (scanner *SortableScanner) Line() (line string) {
	if !scanner.eof {
		line = scanner.scanner.Text()
	}
	return
}

/* Sort key for current line */
func (scanner *SortableScanner) SortKey() (key string) {
	if !scanner.eof {
		key = scanner.sortKey
	}
	return
}

func newScannerDefault(reader io.Reader) *SortableScanner {
	return newScanner(reader, DefaultSortKey)
}

func newScanner(reader io.Reader, sortKeyFunc SortKeyFunc) *SortableScanner {
	sortable := new(SortableScanner)
	sortable.scanner = bufio.NewScanner(reader)
	sortable.keyExtractor = sortKeyFunc
	return sortable
}

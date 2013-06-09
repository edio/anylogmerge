package logmerge

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

/* Print key regex when matched for first time */
var LogRegexMatch bool

var regexLogged bool

/* Extracts sort key from the whole line */
type SortKeyFunc func(line string) (key string)

/* Default sort key is the whole line itself */
func DefaultSortKey(line string) string {
	return line
}

func RegexSortKey(regex string) SortKeyFunc {
	r := regexp.MustCompile(regex)
	return func(line string) string {
		matches := r.FindStringSubmatch(line)
		if len(matches) > 1 {
			key := strings.Join(matches[1:], "")
			if LogRegexMatch && !regexLogged {
				Logger.Printf("Matched key example '%s' : '%s'", regex, key)
				regexLogged = true
			}
			return key
		} else {
			return ""
		}
	}
}

type sortableScanner struct {
	scanner      *bufio.Scanner
	keyExtractor SortKeyFunc
	sortKey      string
	eof          bool
}

func (scanner *sortableScanner) readLine() {
	if scanner.scanner.Scan() {
		scanner.sortKey = scanner.keyExtractor(scanner.scanner.Text())
	} else if err := scanner.scanner.Err(); err == nil {
		scanner.eof = true
	} else {
		panic(err)
	}
}

/* Line last read */
func (scanner *sortableScanner) Line() (line string) {
	if !scanner.eof {
		line = scanner.scanner.Text()
	}
	return
}

func newScannerDefault(reader io.Reader) *sortableScanner {
	return newScanner(reader, DefaultSortKey)
}

func newScanner(reader io.Reader, sortKeyFunc SortKeyFunc) *sortableScanner {
	sortable := new(sortableScanner)
	sortable.scanner = bufio.NewScanner(reader)
	sortable.keyExtractor = sortKeyFunc
	return sortable
}

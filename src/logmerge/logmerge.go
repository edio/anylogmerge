package logmerge

import (
	"fmt"
	"io"
	"log"
)

/* Logger to print progress, etc */
var Logger *log.Logger

/* Print progress every n written lines */
var LogProgress int

var linesWritten int

/* Get index of the smallest string */
type MinFunc func(keys []string) (index int)

func LexicographicOrder(keys []string) (index int) {
	if len(keys) > 1 {
		initial := keys[0]
		for i, k := range keys[1:] {
			if k < initial {
				index = i + 1
			}
		}
	}
	return
}

type Merger struct {
	output  io.Writer
	order   MinFunc
	keyFunc SortKeyFunc
}

func (m *Merger) Merge(input []io.Reader) {
	scanners := make([]*sortableScanner, len(input))
	for index, reader := range input {
		scanners[index] = newScanner(reader, m.keyFunc)
	}

	for i, s := range scanners {
		s.readLine()
		if s.eof {
			scanners = remove(scanners, i)
		}
	}

	for len(scanners) > 0 {
		min := m.order(keys(scanners))
		minScanner := scanners[min]

		// push current minscanner to output and re-read line
		fmt.Fprintln(m.output, minScanner.Line())
		minScanner.readLine()
		if LogProgress > 0 {
			linesWritten += 1
			if linesWritten%LogProgress == 0 {
				Logger.Println("Merged", linesWritten, "lines")
			}
		}
		if minScanner.eof {
			scanners = remove(scanners, min)
		}
	}

	if LogProgress > 0 {
		Logger.Println("Total", linesWritten, "lines")
	}
}

func keys(scanners []*sortableScanner) []string {
	keys := make([]string, len(scanners))
	for i, s := range scanners {
		keys[i] = s.sortKey
	}
	return keys
}

func remove(slice []*sortableScanner, index int) []*sortableScanner {
	if index == len(slice)-1 {
		slice[len(slice)-1], slice = nil, slice[:len(slice)-1]
	} else {
		slice[len(slice)-1], slice[index], slice = nil, slice[len(slice)-1], slice[:len(slice)-1]
	}
	return slice
}

func NewMerger(order MinFunc, keyFunc SortKeyFunc, out io.Writer) *Merger {
	m := new(Merger)
	m.output = out
	m.order = order
	m.keyFunc = keyFunc
	return m
}

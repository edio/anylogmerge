package logmerge

import (
	"strconv"
	"strings"
)

// Extracts sort key from line by slicing the line into pieces.
// columnSpec string received as parameter is a comma-separated list of slices.
// Example: :20,21:30
// means that key consists of 2 string slices, one starts at the beginning of the line and
// lasts till 20th column, another starts at 21th column and lasts till 30th column.
// Slices may overlap.
// Slices that do not fit in string are ignored.
// For example for string "foobar" and columnSpec "2:4,5:9,10:" key will be "obr"
func ColumnSortKey(columnSpec string) SortKeyFunc {
	columns := parseColumnSpec(columnSpec)
	return func(line string) (key string) {
		for _, column := range columns {
			start := column[0]
			end := column[1]

			if start >= len(line) {
				continue
			}

			if end < 0 {
				key += line[start:]
			} else {
				key += line[start:min(end, len(line))]
			}

		}
		return
	}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func parseColumnSpec(columnSpec string) [][2]int {
	slices := strings.Split(columnSpec, ",")

	parsed := make([][2]int, len(slices))
	for index, slice := range slices {
		parsed[index] = parseSlice(slice)
	}
	return parsed
}

func parseSlice(slice string) (parsed [2]int) {
	indexes := strings.Split(slice, ":")

	if len(indexes) != 2 {
		panic("Invalid columnSpec " + slice)
	}

	for i, indexRaw := range indexes {
		if len(indexRaw) > 0 {
			index, err := strconv.Atoi(indexRaw)
			if err != nil {
				panic(err)
			}
			parsed[i] = index
		} else {
			// for i == 0 will return 0, for i == 1 will return -1
			// This makes sense, because ":x" is effectively "0:x"
			// but "x:" is "x:len(str)" which is unknown
			parsed[i] = -1 * i
		}
	}
	return
}

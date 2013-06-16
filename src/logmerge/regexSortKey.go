package logmerge

import (
	"regexp"
	"strings"
)

/* Print key regex when matched for first time */
var LogRegexMatch bool

var regexLogged bool

/* Extract sort key by match line against regular expression with capturing groups */
func RegexSortKey(regex string) SortKeyFunc {
	r := regexp.MustCompile(regex)
	return func(line string) string {
		matches := r.FindStringSubmatch(line)
		if len(matches) > 1 {
			key := strings.Join(matches[1:], "")
			logMatch(regex, key)
			return key
		} else {
			return ""
		}
	}
}

func logMatch(regex, key string) {
	if LogRegexMatch && !regexLogged {
		Logger.Printf("Matched key example '%s' : '%s'", regex, key)
		regexLogged = true
	}
}

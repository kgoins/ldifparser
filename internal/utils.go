package internal

import (
	"regexp"
	"strings"
)

var titleRegex *regexp.Regexp = regexp.MustCompile(`^# .*\.`)

func IsEntityTitle(line string) bool {
	return titleRegex.MatchString(line)
}

func IsComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

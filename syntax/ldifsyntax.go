package syntax

import (
	"regexp"
	"strings"
)

var titleRegex *regexp.Regexp = regexp.MustCompile(`^# .*\.`)

func IsEntityTitle(line string) bool {
	return titleRegex.MatchString(line)
}

func IsEntitySeparator(line string) bool {
	return strings.TrimSpace(line) == ""
}

func IsLdifComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

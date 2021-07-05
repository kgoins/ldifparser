package syntax

import (
	"regexp"
	"strings"
)

var titleRegex *regexp.Regexp = regexp.MustCompile(`^# .*\.`)
var attrRegex *regexp.Regexp = regexp.MustCompile(`^[[:alnum:].;-]+:{1,2} \S+$`)

func IsEntityTitle(line string) bool {
	return titleRegex.MatchString(line)
}

func IsEntitySeparator(line string) bool {
	return strings.TrimSpace(line) == ""
}

func IsLdifComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func IsLdifAttributeLine(line string) bool {
	if IsLdifComment(line) || IsEntitySeparator(line) {
		return false
	}

	return attrRegex.MatchString(line)
}

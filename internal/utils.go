package internal

import "regexp"

var titleRegex *regexp.Regexp = regexp.MustCompile(`^# .*\.`)

func IsEntityTitle(line string) bool {
	return titleRegex.MatchString(line)
}

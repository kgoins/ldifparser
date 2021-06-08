package ldifwriter

import (
	"regexp"

	"github.com/kgoins/entityfilter/entityfilter/filter"
	"github.com/kgoins/ldifwriter/entitybuilder"
	"github.com/kgoins/ldifwriter/internal"
)

const LDAPMaxLineSize int = 1024 * 1024 * 10

type ReaderConf struct {
	Logger internal.ILogger

	Filter          []filter.EntityFilter
	AttributeFilter entitybuilder.AttributeFilter

	TitleLineRegex    *regexp.Regexp
	ScannerBufferSize int
}

func NewReaderConf() ReaderConf {
	regex, _ := regexp.Compile(`^# .*\.`)

	return ReaderConf{
		Filter:            []filter.EntityFilter{},
		AttributeFilter:   entitybuilder.NewAttributeFilter(),
		TitleLineRegex:    regex,
		ScannerBufferSize: LDAPMaxLineSize,
	}
}

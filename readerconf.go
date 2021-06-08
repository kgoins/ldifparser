package ldifwriter

import (
	"regexp"

	hashset "github.com/kgoins/hashset/pkg"
	"github.com/kgoins/ldifwriter/internal"
)

const LDAPMaxLineSize int = 1024 * 1024 * 10

type ReaderConf struct {
	Logger internal.ILogger

	EntityFilter    []IEntityFilter
	AttributeFilter hashset.StrHashset

	TitleLineRegex    *regexp.Regexp
	ScannerBufferSize int
}

func NewReaderConf() ReaderConf {
	regex, _ := regexp.Compile(`^# .*\.`)

	return ReaderConf{
		EntityFilter:      []IEntityFilter{},
		AttributeFilter:   hashset.NewStrHashset(),
		TitleLineRegex:    regex,
		ScannerBufferSize: LDAPMaxLineSize,
	}
}

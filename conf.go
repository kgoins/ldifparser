package ldifparser

import (
	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/kgoins/ldifparser/internal"
)

const LDAPMaxLineSize int = 1000000

type ReaderConf struct {
	Logger            internal.ILogger
	AttributeFilter   entitybuilder.AttributeFilter
	ScannerBufferSize int
}

// NewReaderConf constructs a ReaderConf that has logging
// disabled and a scan buffer size of `LDAPMaxLineSize`
func NewReaderConf() ReaderConf {
	return ReaderConf{
		Logger:            internal.NewNopLogger(),
		AttributeFilter:   entitybuilder.NewAttributeFilter(),
		ScannerBufferSize: LDAPMaxLineSize,
	}
}

type WriterConf struct {
	Logger         internal.ILogger
	SortAttributes bool
}

func NewWriterConf() WriterConf {
	return WriterConf{
		Logger:         internal.NewNopLogger(),
		SortAttributes: false,
	}
}

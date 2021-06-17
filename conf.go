package ldifparser

import (
	"github.com/kgoins/entityfilter/entityfilter/filter"
	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/kgoins/ldifparser/internal"
)

const LDAPMaxLineSize int = 1024 * 1024 * 10

type ReaderConf struct {
	Logger internal.ILogger

	Filter          []filter.EntityFilter
	AttributeFilter entitybuilder.AttributeFilter

	ScannerBufferSize int
}

func NewReaderConf() ReaderConf {
	return ReaderConf{
		Logger:            internal.NewNopLogger(),
		Filter:            []filter.EntityFilter{},
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

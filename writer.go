package ldifparser

import (
	"fmt"
	"io"
	"sort"

	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser/syntax"
)

type LdifWriter struct {
	output io.Writer
	WriterConf
}

func NewLdifWriter(o io.Writer, conf ...WriterConf) LdifWriter {
	var actualConf WriterConf
	if len(conf) == 0 {
		actualConf = NewWriterConf()
	} else {
		actualConf = conf[0]
	}

	return LdifWriter{
		output:     o,
		WriterConf: actualConf,
	}
}

func StringifyAttribute(attr entity.Attribute) []string {
	vals := make([]string, 0, attr.Value.Size())

	for _, value := range attr.Value.Values() {
		vals = append(vals, fmt.Sprintf("%s: %s", attr.Name, value))
	}

	return vals
}

func (w LdifWriter) writeAttribute(attr entity.Attribute) {
	for _, line := range StringifyAttribute(attr) {
		fmt.Fprint(w.output, line+"\n")
	}
}

// WriteEntity will serialize an Entity to LDIF format and write
// it to the configured io.Writer. Attributes will be printed alphabetically.
func (w LdifWriter) WriteEntity(e entity.Entity) (err error) {

	titleLine, err := syntax.BuildTitleLine(e)
	if err != nil {
		return
	}

	fmt.Fprint(w.output, titleLine+"\n")

	attrNames := e.GetAllAttributeNames()
	if w.SortAttributes {
		sort.Strings(attrNames)
	}

	for _, name := range attrNames {
		attr, found := e.GetAttribute(name)
		if found {
			w.writeAttribute(attr)
		}
	}

	w.output.Write([]byte{'\n'})
	return
}

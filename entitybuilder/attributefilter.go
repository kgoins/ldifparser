package entitybuilder

import (
	"strings"

	hashset "github.com/kgoins/hashset/pkg"
	"github.com/kgoins/ldapentity/entity"
)

// AttributeFilters are used to specify which attributes
// should be included in the entity being built.
type AttributeFilter interface {
	Add(...string)
	Contains(string) bool
	IsEmpty() bool

	IsFiltered(attr entity.Attribute) bool
}

type HashsetAttrFilter struct {
	hashset.StrHashset
}

func (f *HashsetAttrFilter) Add(attr ...string) {
	for _, s := range attr {
		f.StrHashset.Add(strings.ToLower(s))
	}
}

// IsFiltered will return true if the filter specifies that
// the attribute should be excluded
func (f HashsetAttrFilter) IsFiltered(attr entity.Attribute) bool {
	if f.IsEmpty() {
		return false
	}

	return !f.Contains(strings.ToLower(attr.Name))
}

// NewAttributeFilter constructs an AttributeFilter with
// lowercase attribute names if any are present.
func NewAttributeFilter(filterParts ...string) AttributeFilter {
	s := hashset.NewStrHashset()
	filter := HashsetAttrFilter{s}
	filter.Add(filterParts...)
	return &filter
}

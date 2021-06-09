package entitybuilder

import (
	"strings"

	hashset "github.com/kgoins/hashset/pkg"
)

// AttributeFilters are used to specify which attributes
// should be included in the entity being built.
type AttributeFilter interface {
	Add(...string)
	Contains(string) bool
	IsEmpty() bool
}

func NewAttributeFilter(filterParts ...[]string) AttributeFilter {
	set := hashset.NewStrHashset()

	if len(filterParts) > 0 {
		for _, attr := range filterParts[0] {
			set.Add(strings.ToLower(attr))
		}
	}

	return &set
}

package entitybuilder

import (
	"errors"
	"strings"

	hashset "github.com/kgoins/hashset/pkg"
	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser/internal"
)

func BuildAttribute(name string, initValue string) entity.Attribute {
	return entity.Attribute{
		Name:  strings.TrimRight(name, ":"),
		Value: hashset.NewStrHashset(initValue),
	}
}

func BuildAttributeFromLine(attrLine string) (a entity.Attribute, err error) {
	if internal.IsComment(attrLine) {
		err = errors.New("unable to build attribute from LDIF comment")
		return
	}

	lineParts := strings.Split(attrLine, ": ")
	if len(lineParts) != 2 {
		err = errors.New("malformed attribute line")
		return
	}

	return BuildAttribute(lineParts[0], lineParts[1]), nil
}

// BuildEntity will construct an Entity from a list of attribute strings
// and filter out all attributes not in `includeAttrs`.
// A null or empty AttributeFilter will include all attributes.
// The `includeAttrs` argument must contain lowercase string values.
// `entityLines` are expected to be in LDIF attribute line format: ex) attrName: value
func BuildEntity(entityLines []string, includeAttrs ...AttributeFilter) (e entity.Entity, err error) {
	var attrFilter AttributeFilter
	if len(includeAttrs) == 0 || includeAttrs[0] == nil {
		attrFilter = NewAttributeFilter()
	} else {
		attrFilter = includeAttrs[0]
	}

	attrMap, err := newAttributeMap(entityLines)
	if err != nil {
		return
	}

	dn, found := attrMap.GetDN()
	if !found || len(dn.GetValues()) < 1 {
		err = errors.New("unable to find entity DN")
		return
	}

	e = entity.NewEntity(dn.GetValues()[0])
	for _, attr := range attrMap {
		if !attrFilter.IsFiltered(attr) {
			e.AddAttribute(attr)
		}
	}

	return
}

package entitybuilder

import (
	"errors"
	"strings"

	"github.com/kgoins/ldapentity/entity"
)

func splitAttrLine(attrLine string) ([]string, error) {
	lineParts := strings.Split(attrLine, ": ")

	if len(lineParts) != 2 {
		return nil, errors.New("malformed attribute line")
	}

	if lineParts[0] == "" || lineParts[1] == "" {
		return nil, errors.New("malformed attribute line")
	}

	lineParts[0] = strings.TrimRight(lineParts[0], ":")
	return lineParts, nil
}

// BuildAttributeFromLine constructs an LDAP attribute from
// an LDIF line, which are expected to be in `attrName: value` format.
func BuildAttributeFromLine(attrLine string) (a entity.Attribute, err error) {
	attrParts, err := splitAttrLine(attrLine)
	if err != nil {
		return
	}

	a = entity.NewEntityAttribute(attrParts[0], attrParts[1])
	return
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

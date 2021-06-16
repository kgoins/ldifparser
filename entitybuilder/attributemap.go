package entitybuilder

import (
	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser/internal"
)

type attributeMap map[string]entity.Attribute

func (m attributeMap) GetDN() (entity.Attribute, bool) {
	var dn entity.Attribute

	dn, found := m["dn"]
	if !found {
		dn, found = m["distinguishedname"]
	}

	return dn, found
}

func newAttributeMap(attrLines []string) (attributeMap, error) {
	attrs := make(map[string]entity.Attribute)

	for _, line := range attrLines {
		if internal.IsComment(line) {
			continue
		}

		attr, err := BuildAttributeFromLine(line)
		if err != nil {
			return nil, err
		}

		attrs[attr.Name] = attr
	}

	return attrs, nil
}

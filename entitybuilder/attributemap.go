package entitybuilder

import (
	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser/syntax"
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

func mergeAttributes(a1, a2 entity.Attribute) entity.Attribute {
	for _, val := range a2.GetValues() {
		a1.Value.Add(val)
	}

	return a1
}

func newAttributeMap(attrLines []string) (attributeMap, error) {
	attrs := make(map[string]entity.Attribute)

	for _, line := range attrLines {
		if syntax.IsLdifComment(line) {
			continue
		}

		attr, err := BuildAttributeFromLine(line)
		if err != nil {
			return nil, err
		}

		if _, attrExists := attrs[attr.Name]; attrExists {
			existing := attrs[attr.Name]
			attr = mergeAttributes(existing, attr)
		}

		attrs[attr.Name] = attr
	}

	return attrs, nil
}

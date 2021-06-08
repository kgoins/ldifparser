package entitybuilder

import (
	"errors"
	"strings"

	hashset "github.com/kgoins/hashset/pkg"
	"github.com/kgoins/ldapentity/entity"
)

func BuildAttribute(name string, initValue string) entity.Attribute {
	return entity.Attribute{
		Name:  strings.TrimRight(name, ":"),
		Value: hashset.NewStrHashset(initValue),
	}
}

func BuildAttributeFromLine(attrLine string) (entity.Attribute, error) {
	lineParts := strings.Split(attrLine, ": ")
	if len(lineParts) != 2 {
		return entity.Attribute{}, errors.New("malformed attribute line")
	}

	return BuildAttribute(lineParts[0], lineParts[1]), nil
}

// BuildEntity is a wrapper for BuildEntityFromAttrList
// that doesn't require an attribute filter and returns all entities.
func BuildEntity(entityLines []string) (entity.Entity, error) {
	return BuildFromAttrList(
		entityLines,
		nil,
	)
}

// BuildFromAttrList will construct an Entity from a list of attribute strings
// and filter out all attributes not in `includeAttrs`.
// Either a null or empty HashSetStr value in `includeAttrs` will include all attributes.
// The `includeAttrs` argument must contain lowercase string values.
func BuildFromAttrList(entityLines []string, includeAttrs AttributeFilter) (entity.Entity, error) {
	entity := entity.NewEntity()
	hasAttrFilter := (includeAttrs != nil) && !includeAttrs.IsEmpty()

	// Ensure that we always pull a DN if possible
	if hasAttrFilter {
		includeAttrs.Add("dn")
		includeAttrs.Add("distinguishedname")
	}

	for _, line := range entityLines {
		attr, err := BuildAttributeFromLine(line)
		if err != nil {
			continue
		}

		if !hasAttrFilter {
			entity.AddAttribute(attr)
		} else {
			if includeAttrs.Contains(strings.ToLower(attr.Name)) {
				entity.AddAttribute(attr)
			}
		}
	}

	_, found := entity.GetDN()
	if !found {
		return entity, errors.New("Unable to parse object DN")
	}

	return entity, nil
}

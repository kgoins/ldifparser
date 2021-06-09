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

func BuildAttributeFromLine(attrLine string) (entity.Attribute, error) {
	lineParts := strings.Split(attrLine, ": ")
	if len(lineParts) != 2 {
		return entity.Attribute{}, errors.New("malformed attribute line")
	}

	return BuildAttribute(lineParts[0], lineParts[1]), nil
}

// BuildEntity will construct an Entity from a list of attribute strings
// and filter out all attributes not in `includeAttrs`.
// Either a null or empty HashSetStr value in `includeAttrs` will include all attributes.
// The `includeAttrs` argument must contain lowercase string values.
func BuildEntity(entityLines []string, includeAttrs ...AttributeFilter) (entity.Entity, error) {
	entity := entity.NewEntity()

	hasAttrFilter := len(includeAttrs) > 0 && includeAttrs[0] != nil
	var attrFilter AttributeFilter
	if hasAttrFilter {
		attrFilter = includeAttrs[0]
	}

	// Ensure that we always pull a DN if possible
	if hasAttrFilter {
		attrFilter.Add("dn")
		attrFilter.Add("distinguishedname")
	}

	for _, line := range entityLines {
		attr, err := BuildAttributeFromLine(line)
		if err != nil {
			if internal.IsEntityTitle(line) {
				continue
			}

			return entity, err
		}

		if !hasAttrFilter {
			entity.AddAttribute(attr)
		} else {
			if attrFilter.Contains(strings.ToLower(attr.Name)) {
				entity.AddAttribute(attr)
			}
		}
	}

	_, found := entity.GetDN()
	if !found {
		return entity, errors.New("unable to parse object DN")
	}

	return entity, nil
}

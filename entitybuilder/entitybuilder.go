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
// A null or empty AttributeFilter will include all attributes.
// The `includeAttrs` argument must contain lowercase string values.
func BuildEntity(entityLines []string, includeAttrs ...AttributeFilter) (entity.Entity, error) {
	entity := entity.NewEntity()

	var attrFilter AttributeFilter
	if len(includeAttrs) == 0 || includeAttrs[0] == nil {
		attrFilter = NewAttributeFilter()
	} else {
		attrFilter = includeAttrs[0]
	}

	// If there is a user specified attr filter, make sure
	// the DN is pulled so a well-formed entity can be built
	if !attrFilter.IsEmpty() {
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

		if !attrFilter.IsFiltered(attr) {
			entity.AddAttribute(attr)
		}
	}

	_, found := entity.GetDN()
	if !found {
		return entity, errors.New("unable to parse object DN")
	}

	return entity, nil
}

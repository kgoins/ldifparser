package entitybuilder_test

import (
	"testing"

	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/stretchr/testify/require"
)

func TestEntityBuilder_BuildAttributeFromValidLine(t *testing.T) {
	attrLine := "userAccountControl: 66048"

	attr, err := entitybuilder.BuildAttributeFromLine(attrLine)
	if err != nil {
		t.Fatalf("Unable to build from valid attr line")
	}

	if attr.Name != "userAccountControl" {
		t.Fatalf("Failed to parse attr name")
	}

	if attr.Value.Size() != 1 {
		t.Fatalf("Failed to parse attr value")
	}

	if attr.Value.Values()[0] != "66048" {
		t.Fatalf("Failed to parse attr value")
	}
}

var defaultTestAttrLines = []string{
	"# MYPC, ContosoUsers, contoso.com",
	"dn: CN=MYPC,OU=ContosoUsers,DC=contoso,DC=com",
	"objectClass: top",
	"objectClass: computer",
	"cn: MYPC",
}

var defaultTestEntitySize = 3

func TestEntityBuilder_BuildFromAttrList_NoInclude(t *testing.T) {
	r := require.New(t)

	e, err := entitybuilder.BuildEntity(defaultTestAttrLines)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.Equal(defaultTestEntitySize, e.Size())

	cnAttr, found := e.GetSingleValuedAttribute("cn")
	r.True(found)
	r.Equal(cnAttr, "MYPC")
}

func TestEntityBuilder_BuildFromAttrList_IncludeList(t *testing.T) {
	r := require.New(t)

	attrFilter := entitybuilder.NewAttributeFilter([]string{"cn"})

	e, err := entitybuilder.BuildEntity(defaultTestAttrLines, attrFilter)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.Equal(2, e.Size()) // only CN and DN should remain

	cnAttr, found := e.GetSingleValuedAttribute("cn")
	r.True(found)
	r.Equal(cnAttr, "MYPC")

	_, found = e.GetAttribute("objectClass")
	r.False(found)
}

func TestEntityBuilder_BuildFromAttrList_NullOrEmptyAttrFilter(t *testing.T) {
	r := require.New(t)

	filters := []entitybuilder.AttributeFilter{
		nil,
		entitybuilder.NewAttributeFilter(),
	}

	for _, filter := range filters {
		e, err := entitybuilder.BuildEntity(defaultTestAttrLines, filter)
		r.NoError(err)
		r.Equal(defaultTestEntitySize, e.Size())
	}
}

func TestEntityBuilder_BuildFromAttrList_MalformedTitle(t *testing.T) {
	r := require.New(t)

	attrLines := []string{
		"MYPC, ContosoUsers, contoso.com",
	}

	_, err := entitybuilder.BuildEntity(attrLines)
	r.Error(err)
}

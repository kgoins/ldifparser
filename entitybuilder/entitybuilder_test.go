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

func TestEntityBuilder_BuildFromAttrList_NoInclude(t *testing.T) {
	r := require.New(t)

	attrLines := []string{
		"# MYPC, ContosoUsers, contoso.com",
		"dn: CN=MYPC,OU=ContosoUsers,DC=contoso,DC=com",
		"objectClass: top",
		"objectClass: computer",
		"cn: MYPC",
	}

	e, err := entitybuilder.BuildEntity(attrLines, nil)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.Equal(3, e.Size())

	cnAttr, found := e.GetSingleValuedAttribute("cn")
	r.True(found)
	r.Equal(cnAttr, "MYPC")
}

func TestEntityBuilder_BuildFromAttrList_IncludeList(t *testing.T) {
	r := require.New(t)

	attrLines := []string{
		"# MYPC, ContosoUsers, contoso.com",
		"dn: CN=MYPC,OU=ContosoUsers,DC=contoso,DC=com",
		"objectClass: top",
		"objectClass: computer",
		"cn: MYPC",
	}

	attrFilter := entitybuilder.NewAttributeFilter([]string{"cn"})

	e, err := entitybuilder.BuildEntity(attrLines, attrFilter)
	r.NoError(err)

	r.False(e.IsEmpty())
	r.Equal(2, e.Size()) // only CN and DN should remain

	cnAttr, found := e.GetSingleValuedAttribute("cn")
	r.True(found)
	r.Equal(cnAttr, "MYPC")

	_, found = e.GetAttribute("objectClass")
	r.False(found)
}

func TestEntityBuilder_BuildFromAttrList_MalformedTitle(t *testing.T) {
	r := require.New(t)

	attrLines := []string{
		"MYPC, ContosoUsers, contoso.com",
	}

	_, err := entitybuilder.BuildEntity(attrLines)
	r.Error(err)
}

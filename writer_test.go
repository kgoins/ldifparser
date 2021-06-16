package ldifparser_test

import (
	"strings"
	"testing"

	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser"
	"github.com/kgoins/ldifparser/entitybuilder"
	"github.com/stretchr/testify/require"
)

func buildTestEntity() entity.Entity {
	e := entity.NewEntity("CN=MYUSR,OU=ContosoUsers,DC=contoso,DC=com")
	e.AddAttribute(entity.NewEntityAttribute("objectClass", "top", "person", "user"))
	e.AddAttribute(entity.NewEntityAttribute("cn", "MYUSR"))
	e.AddAttribute(entity.NewEntityAttribute("whenCreated", "20120423175240.0Z"))
	e.AddAttribute(entity.NewEntityAttribute("sAMAccountName", "MYUSR"))

	return e
}

func TestWriter_StringifyAttribute_SingleValue(t *testing.T) {
	r := require.New(t)
	attr := entity.NewEntityAttribute("whenCreated", "20120423175240.0Z")

	attrStr := ldifparser.StringifyAttribute(attr)
	r.Len(attrStr, 1)

	a2, err := entitybuilder.BuildAttributeFromLine(attrStr[0])
	r.NoError(err)
	r.True(attr.Equals(a2))
}

func TestWriter_StringifyAttribute_MultiValue(t *testing.T) {
	r := require.New(t)
	attr := entity.NewEntityAttribute("objectClass", "top", "person", "user")

	attrStr := ldifparser.StringifyAttribute(attr)
	r.Len(attrStr, 3)
}

func TestWriter_WriteEntity(t *testing.T) {
	r := require.New(t)

	var outBuffer strings.Builder
	writer := ldifparser.NewLdifWriter(&outBuffer)

	e := buildTestEntity()

	err := writer.WriteEntity(e)
	r.NoError(err)

	outStr := outBuffer.String()
	inBuffer := strings.NewReader(outStr)
	reader := ldifparser.NewLdifReader(inBuffer)

	samaccountname, _ := e.GetSingleValuedAttribute("sAMAccountName")
	eOut, err := reader.ReadEntity("sAMAccountName", samaccountname)

	r.NoError(err)
	r.True(e.Equals(eOut))
}

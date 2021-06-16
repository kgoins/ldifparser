package ldifparser_test

import (
	"strings"
	"testing"

	"github.com/kgoins/ldapentity/entity"
	"github.com/kgoins/ldifparser"
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

func TestWriter_WriteEntity(t *testing.T) {
	r := require.New(t)

	var outBuffer strings.Builder
	writer := ldifparser.NewLdifWriter(&outBuffer)

	e := buildTestEntity()

	err := writer.WriteEntity(e)
	r.NoError(err)

	inBuffer := strings.NewReader(outBuffer.String())
	reader := ldifparser.NewLdifReader(inBuffer)

	samaccountname, _ := e.GetSingleValuedAttribute("sAMAccountName")
	eOut, err := reader.BuildEntity("sAMAccountName", samaccountname)

	r.NoError(err)
	r.True(e.Equals(eOut))
}

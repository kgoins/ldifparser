package syntax_test

import (
	"testing"

	"github.com/kgoins/ldifparser/syntax"
	"github.com/stretchr/testify/require"
)

func TestSyntax_IsEntityTitle(t *testing.T) {
	r := require.New(t)

	testMap := map[string]bool{
		"# DISABLEDUSER, ContosoUsers, contoso.com":         true,
		"# DISABLEDUSER\\, ME, ContosoUsers, contoso.com":   true,
		"# DISABLEDUSER (myusr), ContosoUsers, contoso.com": true,
		"DISABLEDUSER, ContosoUsers, contoso.com":           false,
	}

	for testTitle, expectedResp := range testMap {
		resp := syntax.IsEntityTitle(testTitle)
		r.Equal(expectedResp, resp)
	}
}

func TestSyntax_IsLdifAttrLine(t *testing.T) {
	r := require.New(t)

	testMap := map[string]bool{
		"# something something":         false,
		"\n":                            false,
		"1.1.123.5laskjdf: asl dfkj":    false,
		"dn: cn=me,dc=corp,dc=com":      true,
		"y-attr123: asldfkj":            true,
		"1.1.123.5laskjdf: asldfkj":     true,
		"sn;lang-en: Ogasawara":         true,
		"dn:: dWlkPXJvZ2FzYXdhcmEsb3==": true,
	}

	for testTitle, expectedResp := range testMap {
		resp := syntax.IsLdifAttributeLine(testTitle)
		r.Equal(expectedResp, resp)
	}
}

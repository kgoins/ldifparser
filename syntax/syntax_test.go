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

package internal_test

import (
	"testing"

	"github.com/kgoins/ldifparser/internal"
	"github.com/stretchr/testify/require"
)

func TestInternal_IsEntityTitle(t *testing.T) {
	r := require.New(t)

	testMap := map[string]bool{
		"# DISABLEDUSER, ContosoUsers, contoso.com":         true,
		"# DISABLEDUSER\\, ME, ContosoUsers, contoso.com":   true,
		"# DISABLEDUSER (myusr), ContosoUsers, contoso.com": true,
		"DISABLEDUSER, ContosoUsers, contoso.com":           false,
	}

	for testTitle, expectedResp := range testMap {
		resp := internal.IsEntityTitle(testTitle)
		r.Equal(expectedResp, resp)
	}
}

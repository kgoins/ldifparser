package entitybuilder_test

import (
	"testing"

	"github.com/kgoins/ldifparser/entitybuilder"
)

func TestAttribute_BuildFromValidLine(t *testing.T) {
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

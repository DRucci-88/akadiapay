package phase1_test

import (
	"akadia/internal/shared"
	"testing"
)

func TestUpdateMapSetIfNotNil(t *testing.T) {
	var nilBool *bool
	falseValue := false
	zeroValue := 0.0

	updateMap := shared.UpdateMap{}
	updateMap.SetIfNotNil("allow_partial", &falseValue)
	updateMap.SetIfNotNil("minimum_amount", &zeroValue)
	updateMap.SetIfNotNil("ignored", nilBool)

	allowPartial, exists := updateMap["allow_partial"]
	if !exists {
		t.Fatalf("expected false bool pointer to be included")
	}
	if value, ok := allowPartial.(bool); !ok || value {
		t.Fatalf("expected stored false bool, got %#v", allowPartial)
	}

	minimumAmount, exists := updateMap["minimum_amount"]
	if !exists {
		t.Fatalf("expected zero float pointer to be included")
	}
	if value, ok := minimumAmount.(float64); !ok || value != 0 {
		t.Fatalf("expected stored zero float, got %#v", minimumAmount)
	}

	if _, exists := updateMap["ignored"]; exists {
		t.Fatalf("expected nil pointer to be skipped")
	}
}

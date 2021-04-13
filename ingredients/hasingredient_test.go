package ingredients

import "testing"

func TestHasIngredient(t *testing.T) {
	hasItems := []string{"fish", "potato"}
	hasNotItems := []string{"salt", "gravy"}

	for _, item := range hasItems {
		if !HasIngredient(item) {
			t.Errorf("We should have the ingredient '%s'!", item)
		}
	}

	for _, item := range hasNotItems {
		if HasIngredient(item) {
			t.Errorf("We should NOT have the ingredient '%s'!", item)
		}
	}
}

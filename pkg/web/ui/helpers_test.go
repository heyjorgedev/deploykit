package ui

import "testing"

func TestTailwindMerge(t *testing.T) {
	result := tailwindMerge("w-12 h-10 sm:w-20", "w-10 sm:w-10")

	if result != "w-10 h-10 sm:w-10" {
		t.Errorf("Expected w-10 h-10 sm:w-10, got %s", result)
	}
}

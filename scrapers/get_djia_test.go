package scrapers

import "testing"

func TestGetDJIA(t *testing.T) {
	actual, err := GetDJIA()

	if err != nil {
		t.Fatal(err)
	}

	if actual == nil {
		t.Errorf("actual was nil")
	}

	if len(actual) < 100 {
		t.Errorf("expected min: %v, actual: %v", 100, len(actual))
	}
}

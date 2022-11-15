package indexes

import "testing"

func TestGetSP500(t *testing.T) {
	actual, err := ScrapeSP500()

	if err != nil {
		t.Fatal(err)
	}

	if actual == nil {
		t.Errorf("actual was nil")
	}

	if len(actual) < 500 {
		t.Errorf("expected: %v, actual: %v", 500, len(actual))
	}
}

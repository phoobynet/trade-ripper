package calendars

import "testing"

func Test_Previous(t *testing.T) {
	Initialize()

	previous := Previous()

	if previous.Date == "" {
		t.Error("Previous() returned an empty date")
	}
}

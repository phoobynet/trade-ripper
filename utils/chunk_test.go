package utils

import "testing"

func Test_Chunks(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	chunks := Chunk(items, 3)

	if len(chunks) != 4 {
		t.Errorf("Expected 4 chunks, got %d", len(chunks))
	}

	t.Log(chunks)
}

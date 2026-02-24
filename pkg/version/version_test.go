package version

import "testing"

func TestValue(t *testing.T) {
	if Value == "" {
		t.Fatal("version must not be empty")
	}
}

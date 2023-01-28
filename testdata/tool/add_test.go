package add

import (
	"testing"
)

func Test_add(t *testing.T) {
	got := Add(2, 3)
	want := 5

	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

package simpleset

import "testing"

func TestSet(t *testing.T) {
	set := NewSet()
	set.Add(6)
	set.Add(7)
	t.Log("size", set.Size())
	if set.Has(6) {
		t.Log("6 in")
	}
	if set.Has(7) {
		t.Log("7 in")
	}

	if set.Has(10) {
		t.Log("10 in")
	} else {
		t.Log("10 not in")
	}
	set.Rem(6)
	t.Log("size", set.Size())
}

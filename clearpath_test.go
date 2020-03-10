package hit

import "testing"

func TestClearPath_Contains(t *testing.T) {
	tests := []struct {
		haystack clearPath
		needle   clearPath
		want     bool
	}{
		{newClearPath("Foo", nil).Push("Bar", nil), newClearPath("Foo", nil).Push("Bar", nil), true},
		{newClearPath("Foo", nil).Push("Bar", nil), newClearPath("Foo", nil), true},
		{newClearPath("Foo", nil).Push("Bar", nil), newClearPath("Foo", []interface{}{1}), false},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", nil), newClearPath("Foo", nil), true},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", nil), newClearPath("Foo", nil).Push("Bar", nil), true},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", nil), newClearPath("Foo", []interface{}{1}), true},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", nil), newClearPath("Foo", []interface{}{1}).Push("Bar", nil), true},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), newClearPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), true},
		{newClearPath("Foo", []interface{}{1}).Push("Bar", nil), newClearPath("Foo", []interface{}{2}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

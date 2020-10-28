package hit

import "testing"

func TestClearPath_Contains(t *testing.T) {
	tests := []struct {
		haystack callPath
		needle   callPath
		want     bool
	}{
		{newCallPath("Foo", nil).Push("Bar", nil), newCallPath("Foo", nil).Push("Bar", nil), true},
		{newCallPath("Foo", nil).Push("Bar", nil), newCallPath("Foo", nil), true},
		{newCallPath("Foo", nil).Push("Bar", nil), newCallPath("Foo", []interface{}{1}), false},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", nil), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", nil).Push("Bar", nil), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", []interface{}{1}), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", []interface{}{1}).Push("Bar", nil), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), newCallPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", []interface{}{2}), false},
		{newCallPath("Foo", []interface{}{1, 2}).Push("Bar", nil), newCallPath("Foo", []interface{}{1}), true},
		{newCallPath("Foo", []interface{}{1}).Push("Bar", nil), newCallPath("Foo", []interface{}{1, 2}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

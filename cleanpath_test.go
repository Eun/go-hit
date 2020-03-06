package hit

import "testing"

func TestCleanPath_Contains(t *testing.T) {
	tests := []struct {
		haystack CleanPath
		needle   CleanPath
		want     bool
	}{
		{NewCleanPath("Foo", nil).Push("Bar", nil), NewCleanPath("Foo", nil).Push("Bar", nil), true},
		{NewCleanPath("Foo", nil).Push("Bar", nil), NewCleanPath("Foo", nil), true},
		{NewCleanPath("Foo", nil).Push("Bar", nil), NewCleanPath("Foo", []interface{}{1}), false},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), NewCleanPath("Foo", nil), true},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), NewCleanPath("Foo", nil).Push("Bar", nil), true},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), NewCleanPath("Foo", []interface{}{1}), true},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), true},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), NewCleanPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), true},
		{NewCleanPath("Foo", []interface{}{1}).Push("Bar", nil), NewCleanPath("Foo", []interface{}{2}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

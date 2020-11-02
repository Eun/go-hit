package hit

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
)

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

func TestClearPath_ContainsFunc(t *testing.T) {
	f1 := func(Hit) {}
	f2 := func(Hit) {}
	tests := []struct {
		haystack callPath
		needle   callPath
		want     bool
	}{
		{newCallPath("Foo", []interface{}{f1}), newCallPath("Foo", []interface{}{f1}), true},
		{newCallPath("Foo", []interface{}{f2}), newCallPath("Foo", []interface{}{f1}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearPath_ContainsReader(t *testing.T) {
	var r1 io.Reader = bytes.NewReader(nil)
	r2 := rand.Reader
	tests := []struct {
		haystack callPath
		needle   callPath
		want     bool
	}{
		{newCallPath("Foo", []interface{}{r1}), newCallPath("Foo", []interface{}{r1}), true},
		{newCallPath("Foo", []interface{}{r2}), newCallPath("Foo", []interface{}{r1}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

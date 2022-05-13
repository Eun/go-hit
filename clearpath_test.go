package hit_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"

	. "github.com/otto-eng/go-hit"
)

func TestClearPath_Contains(t *testing.T) {
	tests := []struct {
		haystack CallPath
		needle   CallPath
		want     bool
	}{
		{NewCallPath("Foo", nil).Push("Bar", nil), NewCallPath("Foo", nil).Push("Bar", nil), true},
		{NewCallPath("Foo", nil).Push("Bar", nil), NewCallPath("Foo", nil), true},
		{NewCallPath("Foo", nil).Push("Bar", nil), NewCallPath("Foo", []interface{}{1}), false},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", nil), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", nil).Push("Bar", nil), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", []interface{}{1}), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), NewCallPath("Foo", []interface{}{1}).Push("Bar", []interface{}{1}), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", []interface{}{2}), false},
		{NewCallPath("Foo", []interface{}{1, 2}).Push("Bar", nil), NewCallPath("Foo", []interface{}{1}), true},
		{NewCallPath("Foo", []interface{}{1}).Push("Bar", nil), NewCallPath("Foo", []interface{}{1, 2}), false},
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
		haystack CallPath
		needle   CallPath
		want     bool
	}{
		{NewCallPath("Foo", []interface{}{f1}), NewCallPath("Foo", []interface{}{f1}), true},
		{NewCallPath("Foo", []interface{}{f2}), NewCallPath("Foo", []interface{}{f1}), false},
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
		haystack CallPath
		needle   CallPath
		want     bool
	}{
		{NewCallPath("Foo", []interface{}{r1}), NewCallPath("Foo", []interface{}{r1}), true},
		{NewCallPath("Foo", []interface{}{r2}), NewCallPath("Foo", []interface{}{r1}), false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.haystack.Contains(tt.needle); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

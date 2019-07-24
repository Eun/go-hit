package hit_test

import (
	"testing"
)

func TestExamples(t *testing.T) {
	t.Run("", func(t *testing.T) {
		ExampleHead()
	})
	t.Run("", func(t *testing.T) {
		ExamplePost()
	})
	t.Run("", func(t *testing.T) {
		Example_statusCode()
	})
	t.Run("", func(t *testing.T) {
		Example_hash()
	})
	t.Run("", func(t *testing.T) {
		Example_extensibility()
	})
	t.Run("", func(t *testing.T) {
		Example_cookie()
	})
	t.Run("", func(t *testing.T) {
		Example_cookie_alternative()
	})
}

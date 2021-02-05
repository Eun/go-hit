package doctest

import (
	"testing"

	. "github.com/Eun/go-hit"
)

func TestRunTest(t *testing.T) {
	t.Run("http", func(t *testing.T) {
		RunTest(true, func() {
			MustDo(
				Post("http://example.com"),
				Send().Body().String("Hello World"),
				Expect().Body().String().Equal("Hello World"),
			)
		})
	})

	t.Run("https", func(t *testing.T) {
		RunTest(true, func() {
			MustDo(
				Post("https://example.com"),
				Send().Body().String("Hello World"),
				Expect().Body().String().Equal("Hello World"),
			)
		})
	})
}

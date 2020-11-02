package hit_test

import (
	"testing"

	"net/http"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestStoreHeadersTrailers(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runHeaderTrailerTest := func(f func(store func(...string) IStoreStep)) {
		for _, v := range []func(...string) IStoreStep{
			Store().Request().Headers, Store().Request().Trailers,
			Store().Response().Headers, Store().Response().Trailers,
		} {
			f(v)
		}
	}
	runHeaderTrailerTest(func(store func(...string) IStoreStep) {
		t.Run("all headers", func(t *testing.T) {
			t.Run("http.Headers", func(t *testing.T) {
				var headers http.Header
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store().In(&headers),
				)

				require.Equal(t, http.Header{
					"X-String":  {"Foo"},
					"X-Strings": {"Hello", "World"},
					"X-Int":     {"3"},
					"X-Ints":    {"3", "4"},
					"X-Mixed":   {"3.0", "4", "Foo"},
				}, headers)
			})

			t.Run("map", func(t *testing.T) {
				var v map[string]interface{}
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store().In(&v),
				)

				require.Equal(t, map[string]interface{}{
					"X-String":  []string{"Foo"},
					"X-Strings": []string{"Hello", "World"},
					"X-Int":     []string{"3"},
					"X-Ints":    []string{"3", "4"},
					"X-Mixed":   []string{"3.0", "4", "Foo"},
				}, v)
			})

			t.Run("nil", func(t *testing.T) {
				ExpectError(t,
					Do(
						Get(s.URL),
						mockHeadersAndTrailers(),
						store().In(nil),
					),
					PtrStr("destination type cannot be nil"),
				)
			})

			t.Run("not a pointer", func(t *testing.T) {
				var v map[string]interface{}
				ExpectError(t,
					Do(
						Get(s.URL),
						mockHeadersAndTrailers(),
						store().In(v),
					),
					PtrStr("destination must be a pointer"),
				)
			})
		})

		t.Run("specific header", func(t *testing.T) {
			t.Run("string", func(t *testing.T) {
				var v string
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store("X-String").In(&v),
				)

				require.Equal(t, "Foo", v)
			})

			t.Run("empty string", func(t *testing.T) {
				var v string
				Test(t,
					Get(s.URL),
					mockFormValues(),
					store("X-Unknown").In(&v),
				)

				require.Equal(t, "", v)
			})

			t.Run("[]string", func(t *testing.T) {
				var v []string
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store("X-Strings").In(&v),
				)

				require.Equal(t, []string{"Hello", "World"}, v)
			})

			t.Run("int", func(t *testing.T) {
				var v int
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store("X-Int").In(&v),
				)

				require.Equal(t, 3, v)
			})

			t.Run("[]int", func(t *testing.T) {
				var v []int
				Test(t,
					Get(s.URL),
					mockHeadersAndTrailers(),
					store("X-Ints").In(&v),
				)

				require.Equal(t, []int{3, 4}, v)
			})

			t.Run("multiple header values into one string", func(t *testing.T) {
				var v string
				ExpectError(t,
					Do(
						Get(s.URL),
						mockHeadersAndTrailers(),
						store("X-Strings").In(&v),
					),
					PtrStr(`could not put []string{"Hello", "World"} into *string`),
				)
			})

			t.Run("nil", func(t *testing.T) {
				ExpectError(t,
					Do(
						Get(s.URL),
						mockHeadersAndTrailers(),
						store("X-Strings").In(nil),
					),
					PtrStr("destination type cannot be nil"),
				)
			})

			t.Run("not a pointer", func(t *testing.T) {
				var v []string
				ExpectError(t,
					Do(
						Get(s.URL),
						mockHeadersAndTrailers(),
						store("X-Strings").In(v),
					),
					PtrStr("destination must be a pointer"),
				)
			})
		})
	})
}

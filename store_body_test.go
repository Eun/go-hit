package hit_test

import (
	"testing"

	"net/url"

	"io"
	"io/ioutil"

	"github.com/stretchr/testify/require"

	. "github.com/Eun/go-hit"
)

func TestStoreBody(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	runResponseRequestTest := func(f func(storeBody func() IStoreBody)) {
		for _, v := range []func() IStoreBody{Store().Request().Body, Store().Response().Body} {
			f(v)
		}
	}
	runResponseRequestTest(func(storeBody func() IStoreBody) {
		t.Run("bool", func(t *testing.T) {
			var v bool
			Test(t,
				Post(s.URL),
				Send().Body().String("true"),
				storeBody().Bool().In(&v),
			)
			require.Equal(t, true, v)
		})

		t.Run("bytes", func(t *testing.T) {
			var v []byte
			Test(t,
				Post(s.URL),
				Send().Body().String("Hello World"),
				storeBody().Bytes().In(&v),
			)
			require.Equal(t, []byte("Hello World"), v)
		})

		t.Run("float32", func(t *testing.T) {
			var v float32
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Float32().In(&v),
			)
			require.Equal(t, float32(1.0), v)
		})

		t.Run("float64", func(t *testing.T) {
			var v float64
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Float64().In(&v),
			)
			require.Equal(t, 1.0, v)
		})

		t.Run("float64", func(t *testing.T) {
			var v float64
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Float64().In(&v),
			)
			require.Equal(t, 1.0, v)
		})

		t.Run("FormValues", func(t *testing.T) {
			t.Run("all values", func(t *testing.T) {
				t.Run("url.Values", func(t *testing.T) {
					var values url.Values
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues().In(&values),
					)

					require.Equal(t, url.Values{
						"X-String":  {"Foo"},
						"X-Strings": {"Hello", "World"},
						"X-Int":     {"3"},
						"X-Ints":    {"3", "4"},
						"X-Mixed":   {"3.0", "4", "Foo"},
					}, values)
				})
			})

			t.Run("specific value", func(t *testing.T) {
				t.Run("string", func(t *testing.T) {
					var v string
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues("X-String").In(&v),
					)

					require.Equal(t, "Foo", v)
				})

				t.Run("empty string", func(t *testing.T) {
					var v string
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues("X-Unknown").In(&v),
					)

					require.Equal(t, "", v)
				})

				t.Run("[]string", func(t *testing.T) {
					var v []string
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues("X-Strings").In(&v),
					)

					require.Equal(t, []string{"Hello", "World"}, v)
				})

				t.Run("int", func(t *testing.T) {
					var v int
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues("X-Int").In(&v),
					)

					require.Equal(t, 3, v)
				})

				t.Run("[]int", func(t *testing.T) {
					var v []int
					Test(t,
						Get(s.URL),
						mockFormValues(),
						storeBody().FormValues("X-Ints").In(&v),
					)

					require.Equal(t, []int{3, 4}, v)
				})

				t.Run("multiple header values into one string", func(t *testing.T) {
					var v string
					ExpectError(t,
						Do(
							Get(s.URL),
							mockFormValues(),
							storeBody().FormValues("X-Strings").In(&v),
						),
						PtrStr(`could not put []string{"Hello", "World"} into *string`),
					)
				})
			})
		})

		t.Run("int", func(t *testing.T) {
			var n int
			Test(t,
				Post(s.URL),
				Send().Body().String("10"),
				storeBody().Int().In(&n),
			)
			require.Equal(t, 10, n)
		})

		t.Run("int8", func(t *testing.T) {
			var v int8
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Int8().In(&v),
			)
			require.Equal(t, int8(1), v)
		})

		t.Run("int16", func(t *testing.T) {
			var v int16
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Int16().In(&v),
			)
			require.Equal(t, int16(1), v)
		})

		t.Run("int32", func(t *testing.T) {
			var v int32
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Int32().In(&v),
			)
			require.Equal(t, int32(1), v)
		})

		t.Run("int64", func(t *testing.T) {
			var v int64
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Int64().In(&v),
			)
			require.Equal(t, int64(1), v)
		})

		t.Run("json", func(t *testing.T) {
			t.Run("map", func(t *testing.T) {
				var v map[string]interface{}
				Test(t,
					Post(s.URL),
					Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
					storeBody().JSON().In(&v),
				)
				require.Equal(t, map[string]interface{}{"Name": "Joe", "Id": 12.0}, v)
			})

			t.Run("struct", func(t *testing.T) {
				type User struct {
					Name string
					Id   int //nolint:golint,stylecheck //ignore struct field `Id` should be `ID` (golint)
				}
				var v User
				Test(t,
					Post(s.URL),
					Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
					storeBody().JSON().In(&v),
				)
				require.Equal(t, User{Name: "Joe", Id: 12}, v)
			})

			t.Run("JQ", func(t *testing.T) {
				t.Run("string", func(t *testing.T) {
					var str string

					Test(t,
						Post(s.URL),
						Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
						storeBody().JSON().JQ(".Name").In(&str),
					)

					require.Equal(t, `Joe`, str)
				})

				t.Run("int", func(t *testing.T) {
					var n int

					Test(t,
						Post(s.URL),
						Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Id": 12}),
						storeBody().JSON().JQ(".Id").In(&n),
					)

					require.Equal(t, 12, n)
				})
			})
		})

		t.Run("reader", func(t *testing.T) {
			var v io.Reader
			Test(t,
				Post(s.URL),
				Send().Body().String("Hello World"),
				storeBody().Reader().In(&v),
			)
			buf, err := ioutil.ReadAll(v)
			require.NoError(t, err)
			require.Equal(t, []byte(`Hello World`), buf)
		})

		t.Run("string", func(t *testing.T) {
			var v string
			Test(t,
				Post(s.URL),
				Send().Body().String("Hello World"),
				storeBody().String().In(&v),
			)
			require.Equal(t, `Hello World`, v)
		})

		t.Run("uint", func(t *testing.T) {
			var n uint
			Test(t,
				Post(s.URL),
				Send().Body().String("10"),
				storeBody().Uint().In(&n),
			)

			require.Equal(t, uint(10), n)
		})

		t.Run("uint8", func(t *testing.T) {
			var v uint8
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Uint8().In(&v),
			)
			require.Equal(t, uint8(1), v)
		})

		t.Run("uint16", func(t *testing.T) {
			var v uint16
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Uint16().In(&v),
			)
			require.Equal(t, uint16(1), v)
		})

		t.Run("uint32", func(t *testing.T) {
			var v uint32
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Uint32().In(&v),
			)
			require.Equal(t, uint32(1), v)
		})

		t.Run("uint64", func(t *testing.T) {
			var v uint64
			Test(t,
				Post(s.URL),
				Send().Body().String("1"),
				storeBody().Uint64().In(&v),
			)
			require.Equal(t, uint64(1), v)
		})
	})
}

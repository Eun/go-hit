package httpbody_test

import (
	"bytes"
	"io"
	"testing"

	"io/ioutil"

	. "github.com/Eun/go-hit"
	"github.com/stretchr/testify/require"
)

func TestHTTPBody_SetGet(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	prepare := func(setFunc func(hit Hit), old, expect interface{}) []IStep {
		return []IStep{
			Post(s.URL),
			Send().Body(old),
			Send().Body(setFunc),
			Expect().Body().Equal(expect),
		}
	}

	t.Run("Reader", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			buf, err := ioutil.ReadAll(hit.Request().Body().Reader())
			require.NoError(t, err)
			require.Equal(t, []byte("Some Old Data"), buf)
			hit.Request().Body().SetReader(bytes.NewReader([]byte("Hello World")))
		}, "Some Old Data", "Hello World")...)
	})

	t.Run("Bytes", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, []byte("Some Old Data"), hit.Request().Body().Bytes())
			hit.Request().Body().SetBytes([]byte("Hello World"))
		}, "Some Old Data", "Hello World")...)
	})

	t.Run("String", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, "Some Old Data", hit.Request().Body().String())
			hit.Request().Body().SetString("Hello World")
		}, "Some Old Data", "Hello World")...)
	})

	t.Run("Stringf", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, "Some Old Data", hit.Request().Body().String())
			hit.Request().Body().SetStringf("Hello %s", "World")
		}, "Some Old Data", "Hello World")...)
	})

	t.Run("Int", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, 5, hit.Request().Body().Int())
			hit.Request().Body().SetInt(10)
		}, 5, 10)...)
	})

	t.Run("Int8", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, int8(5), hit.Request().Body().Int8())
			hit.Request().Body().SetInt8(10)
		}, 5, 10)...)
	})

	t.Run("Int16", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, int16(5), hit.Request().Body().Int16())
			hit.Request().Body().SetInt16(10)
		}, 5, 10)...)
	})

	t.Run("Int32", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, int32(5), hit.Request().Body().Int32())
			hit.Request().Body().SetInt32(10)
		}, 5, 10)...)
	})

	t.Run("Int64", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, int64(5), hit.Request().Body().Int64())
			hit.Request().Body().SetInt64(10)
		}, 5, 10)...)
	})

	t.Run("Uint", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, uint(5), hit.Request().Body().Uint())
			hit.Request().Body().SetUint(10)
		}, 5, 10)...)
	})

	t.Run("Uint8", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, uint8(5), hit.Request().Body().Uint8())
			hit.Request().Body().SetUint8(10)
		}, 5, 10)...)
	})

	t.Run("Uint16", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, uint16(5), hit.Request().Body().Uint16())
			hit.Request().Body().SetUint16(10)
		}, 5, 10)...)
	})

	t.Run("Uint32", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, uint32(5), hit.Request().Body().Uint32())
			hit.Request().Body().SetUint32(10)
		}, 5, 10)...)
	})

	t.Run("Uint64", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, uint64(5), hit.Request().Body().Uint64())
			hit.Request().Body().SetUint64(10)
		}, 5, 10)...)
	})

	t.Run("Float32", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, float32(5.0), hit.Request().Body().Float32())
			hit.Request().Body().SetFloat32(10.0)
		}, 5.0, 10.0)...)
	})

	t.Run("Float64", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, 5.0, hit.Request().Body().Float64())
			hit.Request().Body().SetFloat64(10.0)
		}, 5.0, 10.0)...)
	})

	t.Run("Bool", func(t *testing.T) {
		Test(t, prepare(func(hit Hit) {
			require.Equal(t, false, hit.Request().Body().Bool())
			hit.Request().Body().SetBool(true)
		}, false, true)...)
	})

	t.Run("Generic", func(t *testing.T) {
		t.Run("Reader", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				buf, err := ioutil.ReadAll(hit.Request().Body().Reader())
				require.NoError(t, err)
				require.Equal(t, []byte("Some Old Data"), buf)
				hit.Request().Body().Set(bytes.NewReader([]byte("Hello World")))
			}, "Some Old Data", bytes.NewReader([]byte("Hello World")))...)

			ExpectError(t,
				Do(
					Post(s.URL),
					Send().Body("Hello World"),
					Expect().Body().Equal(instantErrorReader{io.ErrUnexpectedEOF}),
				),
				PtrStr(`unable to read data from reader: unexpected EOF`),
			)
		})

		t.Run("Bytes", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, []byte("Some Old Data"), hit.Request().Body().Bytes())
				hit.Request().Body().Set([]byte("Hello World"))
			}, "Some Old Data", "Hello World")...)
		})

		t.Run("String", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, "Some Old Data", hit.Request().Body().String())
				hit.Request().Body().Set("Hello World")
			}, "Some Old Data", "Hello World")...)
		})

		t.Run("Int", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, 5, hit.Request().Body().Int())
				hit.Request().Body().Set(10)
			}, 5, 10)...)
		})

		t.Run("Int8", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, int8(5), hit.Request().Body().Int8())
				hit.Request().Body().Set(int8(10))
			}, 5, int8(10))...)
		})

		t.Run("Int16", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, int16(5), hit.Request().Body().Int16())
				hit.Request().Body().Set(int16(10))
			}, 5, int16(10))...)
		})

		t.Run("Int32", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, int32(5), hit.Request().Body().Int32())
				hit.Request().Body().Set(int32(10))
			}, 5, int32(10))...)
		})

		t.Run("Int64", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, int64(5), hit.Request().Body().Int64())
				hit.Request().Body().Set(int64(10))
			}, 5, int64(10))...)
		})

		t.Run("Uint", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, uint(5), hit.Request().Body().Uint())
				hit.Request().Body().Set(uint(10))
			}, 5, uint(10))...)
		})

		t.Run("Uint8", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, uint8(5), hit.Request().Body().Uint8())
				hit.Request().Body().Set(uint8(10))
			}, 5, uint8(10))...)
		})

		t.Run("Uint16", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, uint16(5), hit.Request().Body().Uint16())
				hit.Request().Body().Set(uint16(10))
			}, 5, uint16(10))...)
		})

		t.Run("Uint32", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, uint32(5), hit.Request().Body().Uint32())
				hit.Request().Body().Set(uint32(10))
			}, 5, uint32(10))...)
		})

		t.Run("Uint64", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, uint64(5), hit.Request().Body().Uint64())
				hit.Request().Body().Set(uint64(10))
			}, 5, uint64(10))...)
		})

		t.Run("Float32", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, float32(5.0), hit.Request().Body().Float32())
				hit.Request().Body().Set(float32(10.0))
			}, 5.0, float32(10.0))...)
		})

		t.Run("Float64", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, 5.0, hit.Request().Body().Float64())
				hit.Request().Body().Set(10.0)
			}, 5.0, 10.0)...)
		})

		t.Run("Bool", func(t *testing.T) {
			Test(t, prepare(func(hit Hit) {
				require.Equal(t, false, hit.Request().Body().Bool())
				hit.Request().Body().Set(true)
			}, false, true)...)
		})
	})
}

func TestHTTPBody_Contains(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	t.Run("Reader", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains(bytes.NewReader([]byte("Hello"))),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains(bytes.NewReader([]byte("Bye"))),
			),
			PtrStr(`"[72 101 108 108 111 32 87 111 114 108 100]" does not contain "[66 121 101]"`),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains(instantErrorReader{io.ErrUnexpectedEOF}),
			),
			PtrStr(`unable to read data from reader: unexpected EOF`),
		)
	})

	t.Run("Bytes", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains([]byte("Hello")),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains([]byte("Bye")),
			),
			PtrStr(`"[72 101 108 108 111 32 87 111 114 108 100]" does not contain "[66 121 101]"`),
		)
	})

	t.Run("String", func(t *testing.T) {
		Test(t,
			Post(s.URL),
			Send().Body("Hello World"),
			Expect().Body().Contains("Hello"),
		)
		ExpectError(t,
			Do(
				Post(s.URL),
				Send().Body("Hello World"),
				Expect().Body().Contains("Bye"),
			),
			PtrStr(`"Hello World" does not contain "Bye"`),
		)
	})
}

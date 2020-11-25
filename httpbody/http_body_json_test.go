package httpbody

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpBodyJson_Set(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want interface{}
	}{
		{
			name: "string",
			data: "Hello World",
			want: "Hello World",
		},
		{
			name: "map",
			data: map[string]interface{}{"Name": "Joe", "Id": float64(10)},
			want: map[string]interface{}{"Name": "Joe", "Id": float64(10)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer("", func(body *HTTPBody) {
				body.JSON().Set(tt.data)
				var got interface{}
				body.JSON().MustDecode(&got)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

func TestHttpBodyJson_MustDecode(t *testing.T) {
	require.Panics(t, func() {
		testServer(`Hello World`, func(body *HTTPBody) {
			var slice []string
			body.JSON().MustDecode(&slice)
		})
	})
}

func TestHttpBodyJson_Decode(t *testing.T) {
	type User struct {
		Name string
	}
	type UserRenamedFields struct {
		SomeName string `json:"name"`
	}

	t.Run("string", func(t *testing.T) {
		testServer(`"Hello World"`, func(body *HTTPBody) {
			container := ""
			want := "Hello World"
			wantErr := false
			if err := body.JSON().Decode(&container); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			if !reflect.DeepEqual(container, want) {
				t.Errorf("container = %v, want %v", container, want)
			}
		})
	})

	t.Run("map", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := map[string]interface{}{}
			want := map[string]interface{}{"Name": "Joe", "Id": float64(10)}
			wantErr := false
			if err := body.JSON().Decode(&container); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			if !reflect.DeepEqual(container, want) {
				t.Errorf("container = %v, want %v", container, want)
			}
		})
	})

	t.Run("unable to unmarshal", func(t *testing.T) {
		testServer(`Hello World`, func(body *HTTPBody) {
			container := []string{}
			want := []string{}
			wantErr := true
			if err := body.JSON().Decode(&container); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			if !reflect.DeepEqual(container, want) {
				t.Errorf("container = %v, want %v", container, want)
			}
		})
	})

	// {
	// not supported
	// 	name:      "stream",
	// 	response:  `{"Name":"Joe"}{"Name":"Alice"}`,
	// 	container: []map[string]interface{}{},
	// 	want:      []map[string]interface{}{{"Name": "Joe"}, {"Name": "Alice"}},
	// 	wantErr:   false,
	// },

	t.Run("struct with missing data", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := User{}
			want := User{"Joe"}
			wantErr := false
			if err := body.JSON().Decode(&container); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			if !reflect.DeepEqual(container, want) {
				t.Errorf("container = %v, want %v", container, want)
			}
		})
	})

	t.Run("struct with other names and missing data", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := UserRenamedFields{}
			want := UserRenamedFields{"Joe"}
			wantErr := false
			if err := body.JSON().Decode(&container); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			if !reflect.DeepEqual(container, want) {
				t.Errorf("container = %v, want %v", container, want)
			}
		})
	})
}

func TestHttpBodyJson_JQ(t *testing.T) {
	type User struct {
		Name string
	}

	type UserRenamedFields struct {
		SomeName string `json:"name"`
	}

	t.Run("expression - string", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := ""
			expression := ".Name"
			wantErr := false
			want := "Joe"
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("expression - float64", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := float64(0)
			expression := ".Id"
			wantErr := false
			want := float64(10)
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("expression - int", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := 0
			expression := ".Id"
			wantErr := false
			want := 10
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("unable to get expr value", func(t *testing.T) {
		testServer(`"Hello World"`, func(body *HTTPBody) {
			container := []string{}
			expression := "A"
			wantErr := true
			want := []string{}
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("nil expr value", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := ""
			expression := ".Address"
			wantErr := false
			want := ""
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("array", func(t *testing.T) {
		testServer(`[{"Name":"Joe"}]`, func(body *HTTPBody) {
			container := ""
			expression := ".[0].Name"
			wantErr := false
			want := "Joe"
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("struct with missing data", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := User{}
			expression := "."
			wantErr := false
			want := User{"Joe"}
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("struct with other names and missing data", func(t *testing.T) {
		testServer(`{"Name":"Joe","Id":10}`, func(body *HTTPBody) {
			container := UserRenamedFields{}
			expression := "."
			wantErr := false
			want := UserRenamedFields{"Joe"}
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})

	t.Run("return array", func(t *testing.T) {
		testServer(`{"users":[{"Name":"Joe","Id":10},{"Name":"Alice","Id":11},{"Name":"Bob","Id":12}]}`, func(body *HTTPBody) {
			container := []User{}
			expression := ".users[]"
			wantErr := false
			want := []User{{"Joe"}, {"Alice"}, {"Bob"}}
			if err := body.JSON().JQ(&container, expression); (err != nil) != wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, wantErr)
			}
			require.Equal(t, want, container, want)
		})
	})
}

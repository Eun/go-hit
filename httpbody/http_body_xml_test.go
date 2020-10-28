package httpbody

import (
	"reflect"
	"testing"
)

func TestHttpBodyXml_Set(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want string
	}{
		{
			name: "string",
			data: "Hello World",
			want: "<string>Hello World</string>",
		},
		{
			name: "slice",
			data: []string{"Alice", "Bob"},
			want: "<string>Alice</string><string>Bob</string>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer("", func(body *HTTPBody) {
				body.XML().Set(tt.data)
				if got := body.MustString(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			})
		})
	}
}

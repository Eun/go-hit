package internal

import "testing"

func TestMakeURL(t *testing.T) {
	tests := []struct {
		Base     string
		URL      string
		Expected string
	}{
		{
			"http://example.com",
			"index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com/",
			"index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com",
			"/index.html",
			"http://example.com/index.html",
		},
		{
			"http://example.com",
			"",
			"http://example.com",
		},
		{
			"http://example.com/",
			"",
			"http://example.com/",
		},
		{
			"",
			"index.html",
			"index.html",
		},
		{
			"",
			"/index.html",
			"/index.html",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run("", func(t *testing.T) {
			if s := MakeURL(test.Base, test.URL); s != test.Expected {
				t.Errorf("expected `%s' got `%s'", test.Expected, s)
			}
		})
	}
}

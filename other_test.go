package hit

import "testing"

//
// func TestOther(t *testing.T) {
// 	client := &http.Client{}
// 	call := Connect(t, "http://127.0.0.1")
// 	call.SetHTTPClient(client)
// 	require.Equal(t, client, call.HTTPClient())
// }
//
// func TestPanicT_Errorf(t *testing.T) {
// 	require.Panics(t, func() {
// 		PanicT{}.Errorf("")
// 	})
// }
//
// func TestPanicT_FailNow(t *testing.T) {
// 	require.Panics(t, func() {
// 		PanicT{}.FailNow()
// 	})
// }

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
			if s := makeURL(test.Base, test.URL); s != test.Expected {
				t.Errorf("expected `%s' got `%s'", test.Expected, s)
			}
		})
	}
}

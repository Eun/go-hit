// Package implicit can be used for testing purposes. It spawns a http test server that provides
// functionality to test the hit framework. It also overwrites the http.DefaultTransport field. So use with care.
package implicit

import (
	"net/http"

	"github.com/Eun/go-hit/doctest/server"
)

//nolint:gochecknoinits // this is needed so we can overwrite the http.DefaultTransport implicitly.
func init() {
	srv := server.NewServer()
	http.DefaultTransport = srv.Transport()
}

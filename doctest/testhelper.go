// Package doctest is a package to help test the hit framework.
package doctest

import (
	"net/http"

	"github.com/otto-eng/go-hit/doctest/server"
)

// RunTest mocks an test http server.
func RunTest(expectRequest bool, test func()) {
	srv := server.NewServer()

	http.DefaultTransport = srv.Transport()

	test()

	_ = srv.Close()

	if expectRequest && srv.RequestCount() == 0 {
		panic("expected at least one request")
	}
}

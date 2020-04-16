package hit

import (
	"net/http"

	"io"

	"github.com/Eun/go-hit/internal/misc"
	"golang.org/x/xerrors"
)

type debugHeader struct {
	debug      *debug
	mode       debugHeaderMode
	expression []string
}

type debugHeaderMode int

const (
	debugHeaderRequest debugHeaderMode = iota
	debugHeaderResponse
	debugTrailerRequest
	debugTrailerResponse
)

func newDebugHeader(debug *debug, mode debugHeaderMode, expression []string) IStep {
	return &debugHeader{
		debug:      debug,
		mode:       mode,
		expression: expression,
	}
}

func (*debugHeader) when() StepTime {
	return BeforeExpectStep
}

func (d *debugHeader) exec(hit Hit) error {
	var headers http.Header
	switch d.mode {
	case debugHeaderRequest:
		headers = hit.Request().Header
	case debugHeaderResponse:
		headers = hit.Response().Header
	case debugTrailerRequest:
		headers = hit.Request().Trailer
	case debugTrailerResponse:
		// we have to read the body to get the trailers
		_, _ = io.Copy(misc.DevNullWriter(), hit.Response().Body().Reader())
		headers = hit.Response().Trailer
	default:
		return xerrors.New("unknown mode")
	}
	return d.debug.printWithExpression(hit, d.debug.getMap(headers), d.expression)
}

func (*debugHeader) clearPath() clearPath {
	return nil // not clearable
}

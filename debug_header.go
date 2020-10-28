package hit

import (
	"io/ioutil"
	"net/http"

	"io"

	"golang.org/x/xerrors"

	"github.com/Eun/go-hit/errortrace"
)

type debugHeader struct {
	cp     callPath
	debug  *debug
	mode   debugHeaderMode
	header string
}

type debugHeaderMode int

const (
	debugHeaderRequest debugHeaderMode = iota
	debugHeaderResponse
	debugTrailerRequest
	debugTrailerResponse
)

func newDebugHeaders(cp callPath, debug *debug, mode debugHeaderMode, header string) IStep {
	return &debugHeader{
		cp:     cp,
		debug:  debug,
		mode:   mode,
		header: header,
	}
}

func (*debugHeader) trace() *errortrace.ErrorTrace {
	return nil
}

func (*debugHeader) when() StepTime {
	return BeforeExpectStep
}

func (d *debugHeader) exec(hit *hitImpl) error {
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
		_, _ = io.Copy(ioutil.Discard, hit.Response().Body().Reader())
		headers = hit.Response().Trailer
	default:
		return xerrors.New("unknown mode")
	}
	if d.header == "" {
		return d.debug.print(d.debug.out(hit), d.debug.getMap(headers))
	}
	return d.debug.print(d.debug.out(hit), d.debug.getMap(headers)[d.header])
}

func (d *debugHeader) callPath() callPath {
	return d.cp
}

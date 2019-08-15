package hit

import (
	"encoding/json"
	"net/http"

	"io"

	"fmt"

	"time"

	"github.com/Eun/go-hit/errortrace"
	"github.com/tidwall/pretty"
)

type Callback func(hit Hit)

type State uint8

const (
	Ready State = iota
	Working
	Done
)

type Hit interface {
	Request() *HTTPRequest
	SetRequest(*http.Request)

	Response() *HTTPResponse

	HTTPClient() *http.Client
	SetHTTPClient(*http.Client)

	Stdout() io.Writer
	SetStdout(io.Writer)

	BaseURL() string
	SetBaseURL(string, ...interface{})

	AddSteps(...IStep)
	RemoveSteps(...IStep)

	Send(...interface{}) ISend
	Expect(...interface{}) IExpect
	Debug()
}

type defaultInstance struct {
	steps    []IStep
	request  *HTTPRequest
	response *HTTPResponse
	client   *http.Client
	state    State
	stdout   io.Writer
	baseURL  string
}

// Request returns the current request
func (hit *defaultInstance) Request() *HTTPRequest {
	return hit.request
}

// SetRequest sets the request for the current instance
func (hit *defaultInstance) SetRequest(request *http.Request) {
	hit.request = newHTTPRequest(hit, request)
}

// Response returns the current Response
func (hit *defaultInstance) Response() *HTTPResponse {
	return hit.response
}

// HTTPClient gets the current http.Client
func (hit *defaultInstance) HTTPClient() *http.Client {
	return hit.client
}

// SetHTTPClient sets the client for the request
func (hit *defaultInstance) SetHTTPClient(client *http.Client) {
	hit.client = client
}

// Stdout gets the current output
func (hit *defaultInstance) Stdout() io.Writer {
	return hit.stdout
}

// SetStdout sets the output to the specified writer
func (hit *defaultInstance) SetStdout(w io.Writer) {
	hit.stdout = w
}

// BaseURL returns the current base url
func (hit *defaultInstance) BaseURL() string {
	return hit.baseURL
}

// SetBaseURL sets the base url
func (hit *defaultInstance) SetBaseURL(url string, a ...interface{}) {
	hit.baseURL = fmt.Sprintf(url, a...)
}

// AddSteps adds the specified steps to the queue
func (hit *defaultInstance) AddSteps(steps ...IStep) {
	if len(steps) <= 0 {
		return
	}
	hit.steps = append(hit.steps, steps...)
}

// RemoveSteps removes the specified steps to the queue
func (hit *defaultInstance) RemoveSteps(steps ...IStep) {
	size := len(steps)
	if size <= 0 {
		return
	}
	for i := 0; i < size; i++ {
		for j := len(hit.steps) - 1; j >= 0; j++ {
			if hit.steps[j] == steps[i] {
				hit.steps = append(hit.steps[:j], hit.steps[j+1:]...)
				break
			}
		}
	}
}

func (hit *defaultInstance) collectSteps(state StepTime, offset int) []IStep {
	var steps []IStep
	for i := offset; i < len(hit.steps); i++ {
		w := hit.steps[i].when()
		if w&^CleanStep == state { // remove CleanStep flag
			if w&CleanStep == CleanStep { // if CleanStep is set
				steps = nil
				continue
			}
			steps = append(steps, hit.steps[i])
		}
	}
	return steps
}

func (hit *defaultInstance) runSteps(state StepTime) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			if _, ok := r.(errortrace.ErrorTraceError); !ok {
				err = errortrace.FormatError(fmt.Sprintf("%v", r))
				fmt.Println(err)
				return
			}
			err = fmt.Errorf("%v", r)
		}
	}()

	totalSteps := len(hit.steps)
	// find all steps we want to run
	stepsToRun := hit.collectSteps(state, 0)
	size := len(stepsToRun)
	i := 0
	for {
		if i >= size {
			return nil
		}
		stepsToRun[i].exec(hit)

		// maybe there is a new step added
		newTotalSteps := len(hit.steps)
		if totalSteps != newTotalSteps {
			// yes they have been modified
			// find all new steps after the last scan
			stepsToRun = append(stepsToRun[:i], hit.collectSteps(state, totalSteps-1)...)
			size = len(stepsToRun)
			i = 0
			totalSteps = newTotalSteps
			continue
		}
		i++
	}
}

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func (hit *defaultInstance) Send(data ...interface{}) ISend {
	step := Send(data...)
	hit.AddSteps(step)
	return step
}

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func (hit *defaultInstance) Expect(data ...interface{}) IExpect {
	step := Expect(data...)
	hit.AddSteps(step)
	return step
}

// Debug prints the current Request and Response to hit.Stdout()
func (hit *defaultInstance) Debug() {
	type M map[string]interface{}

	getBody := func(body *HTTPBody) interface{} {
		reader := body.JSON().body.Reader()
		// if there is a json reader
		if reader != nil {
			var container interface{}
			if err := json.NewDecoder(reader).Decode(&container); err == nil {
				return container
			}
		}
		s := body.String()
		if len(s) == 0 {
			return nil
		}
		return s
	}

	getHeader := func(header http.Header) map[string]interface{} {
		m := make(map[string]interface{})
		for key := range header {
			m[key] = header.Get(key)
		}
		return m
	}

	m := M{
		"Time": time.Now().String(),
		"Request": M{
			"Header":           getHeader(hit.Request().Header),
			"Trailer":          getHeader(hit.Request().Trailer),
			"Method":           hit.Request().Method,
			"URL":              hit.Request().URL,
			"Proto":            hit.Request().Proto,
			"ProtoMajor":       hit.Request().ProtoMajor,
			"ProtoMinor":       hit.Request().ProtoMinor,
			"ContentLength":    hit.Request().ContentLength,
			"TransferEncoding": hit.Request().TransferEncoding,
			"Host":             hit.Request().Host,
			"Form":             hit.Request().Form,
			"PostForm":         hit.Request().PostForm,
			"MultipartForm":    hit.Request().MultipartForm,
			"RemoteAddr":       hit.Request().RemoteAddr,
			"RequestURI":       hit.Request().RequestURI,
			"Body":             getBody(hit.Request().Body()),
		},
		"Response": M{
			"Header":           getHeader(hit.Response().Header),
			"Trailer":          getHeader(hit.Response().Trailer),
			"Proto":            hit.Response().Proto,
			"ProtoMajor":       hit.Response().ProtoMajor,
			"ProtoMinor":       hit.Response().ProtoMinor,
			"ContentLength":    hit.Response().ContentLength,
			"TransferEncoding": hit.Response().TransferEncoding,
			"Body":             getBody(hit.Response().body),
			"Status":           hit.Response().Status,
			"StatusCode":       hit.Response().StatusCode,
		},
	}

	bytes, err := json.Marshal(m)
	errortrace.NoError(err)
	_, _ = hit.Stdout().Write(pretty.Color(pretty.Pretty(bytes), nil))
}

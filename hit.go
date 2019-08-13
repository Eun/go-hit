package hit

import (
	"encoding/json"
	"net/http"

	"io"

	"os"

	"fmt"

	"strings"

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
	SetRequest(*http.Request)
	Request() *HTTPRequest
	Response() *HTTPResponse

	SetHTTPClient(*http.Client)
	HTTPClient() *http.Client

	SetStdout(io.Writer)
	Stdout() io.Writer

	SetBaseURL(string)
	// BaseURL returns the current base url
	BaseURL() string

	SetSteps([]Step)
	Steps() []Step
	AddSteps(...Step)
	RemoveSteps(...Step)

	Send(data ...interface{}) ISend
	Expect(data ...interface{}) IExpect
}

type defaultInstance struct {
	steps    []Step
	request  *HTTPRequest
	response *HTTPResponse
	client   *http.Client
	state    State
	stdout   io.Writer
	baseURL  string
}

// SetRequest sets the request for the current instance
func (hit *defaultInstance) SetRequest(request *http.Request) {
	hit.request = newHTTPRequest(hit, request)
}

func (hit *defaultInstance) Request() *HTTPRequest {
	return hit.request
}

func (hit *defaultInstance) Response() *HTTPResponse {
	return hit.response
}

func (hit *defaultInstance) SetHTTPClient(client *http.Client) {
	hit.client = client
}

func (hit *defaultInstance) HTTPClient() *http.Client {
	return hit.client
}

func (hit *defaultInstance) SetStdout(w io.Writer) {
	hit.stdout = w
}

func (hit *defaultInstance) Stdout() io.Writer {
	return hit.stdout
}

func (hit *defaultInstance) SetBaseURL(s string) {
	hit.baseURL = s
}

func (hit *defaultInstance) BaseURL() string {
	return hit.baseURL
}

func (hit *defaultInstance) SetSteps(calls []Step) {
	hit.steps = calls
}

func (hit *defaultInstance) Steps() []Step {
	return hit.steps
}

func (hit *defaultInstance) AddSteps(steps ...Step) {
	if len(steps) <= 0 {
		return
	}
	hit.steps = append(hit.steps, steps...)
}

func (hit *defaultInstance) RemoveSteps(steps ...Step) {
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

func (hit *defaultInstance) collectSteps(state StepTime, offset int) []Step {
	var steps []Step
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

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func (hit *defaultInstance) Expect(data ...interface{}) IExpect {
	step := Expect(data...)
	hit.AddSteps(step)
	return step
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

func makeURL(base, url string, a ...interface{}) string {
	if base != "" {
		base = strings.TrimRight(base, "/")
	}
	return base + strings.TrimLeft(fmt.Sprintf(url, a...), "/")
}

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func Expect(data ...interface{}) IExpect {
	if arg, ok := getLastArgument(data); ok {
		return &dummyExpect{newExpect().Interface(arg)}
	}
	return newExpect()
}

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func Send(data ...interface{}) ISend {
	if arg, ok := getLastArgument(data); ok {
		return &dummySend{newSend().Interface(arg)}
	}
	return newSend()
}

func Debug() Step {
	return MakeStep(BeforeExpectStep, func(hit Hit) {
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
			return body.Bytes()
		}

		getHeader := func(header http.Header) map[string]interface{} {
			m := make(map[string]interface{})
			for key := range header {
				m[key] = header.Get(key)
			}
			return m
		}

		m := M{
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
	})
}

// SetBaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, WithMethod call
func SetBaseURL(s string) Step {
	return MakeStep(BeforeSendStep, func(hit Hit) {
		hit.SetBaseURL(s)
	})
}

// WithRequest creates a new Hit instance with an existing http request
func WithRequest(request *http.Request) Step {
	return MakeStep(BeforeSendStep, func(hit Hit) {
		hit.SetRequest(request)
	})
}

// WithMethod creates a new Hit instance with the specified method and url
// WithMethod(t, "POST", "http://%s/%s", domain, path)
func WithMethod(method, url string, a ...interface{}) Step {
	et := errortrace.Prepare()
	return MakeStep(BeforeSendStep, func(hit Hit) {
		request, err := http.NewRequest(method, makeURL(hit.BaseURL(), url, a...), nil)
		et.NoError(err, "unable to create request")
		hit.SetRequest(request)
	})
}

// Connect creates a new Hit instance with CONNECT as the http method, use the optional arguments to format the url
// Connect(t, "http://%s/%s", domain, path)
func Connect(url string, a ...interface{}) Step {
	return WithMethod("CONNECT", url, a...)
}

// Delete creates a new Hit instance with DELETE as the http method, use the optional arguments to format the url
// Delete("http://%s/%s", domain, path)
func Delete(url string, a ...interface{}) Step {
	return WithMethod("DELETE", url, a...)
}

// Get creates a new Hit instance with GET as the http method, use the optional arguments to format the url
// Get("http://%s/%s", domain, path)
func Get(url string, a ...interface{}) Step {
	return WithMethod("GET", url, a...)
}

// Head creates a new Hit instance with HEAD as the http method, use the optional arguments to format the url
// Head("http://%s/%s", domain, path)
func Head(url string, a ...interface{}) Step {
	return WithMethod("HEAD", url, a...)
}

// Post creates a new Hit instance with POST as the http method, use the optional arguments to format the url
// Post("http://%s/%s", domain, path)
func Post(url string, a ...interface{}) Step {
	return WithMethod("POST", url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http method, use the optional arguments to format the url
// Options("http://%s/%s", domain, path)
func Options(url string, a ...interface{}) Step {
	return WithMethod("OPTIONS", url, a...)
}

// Put creates a new Hit instance with PUT as the http method, use the optional arguments to format the url
// Put("http://%s/%s", domain, path)
func Put(url string, a ...interface{}) Step {
	return WithMethod("PUT", url, a...)
}

// Trace creates a new Hit instance with TRACE as the http method, use the optional arguments to format the url
// Trace("http://%s/%s", domain, path)
func Trace(url string, a ...interface{}) Step {
	return WithMethod("TRACE", url, a...)
}

func Test(t TestingT, steps ...Step) {
	if err := Do(steps...); err != nil {
		t.Error(err.Error())
	}
}

func Do(steps ...Step) error {
	hit := &defaultInstance{
		client: http.DefaultClient,
		stdout: os.Stdout,
		steps:  steps,
	}
	if err := hit.runSteps(BeforeSendStep); err != nil {
		return err
	}
	if hit.request == nil {
		return fmt.Errorf("unable to perform request: no request set, did you called Post(), Get(), ...?")
	}
	if err := hit.runSteps(SendStep); err != nil {
		return err
	}
	if err := hit.runSteps(AfterSendStep); err != nil {
		return err
	}
	hit.request.Request.Body = hit.request.Body().Reader()
	res, err := hit.client.Do(hit.request.Request)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s", err.Error())
	}

	hit.response = newHTTPResponse(hit, res)
	if err := hit.runSteps(BeforeExpectStep); err != nil {
		return err
	}
	if err := hit.runSteps(ExpectStep); err != nil {
		return err
	}
	if err := hit.runSteps(AfterExpectStep); err != nil {
		return err
	}
	if hit.request.Request.Body != nil {
		_ = hit.request.Request.Body.Close()
	}
	return nil
}

func RemoveCalls(whenToRun, whichToRemove StepTime) Step {
	return MakeStep(whenToRun, func(hit Hit) {
		calls := hit.Steps()
		size := len(calls)
		newCalls := make([]Step, 0, size)
		for i := 0; i < size; i++ {
			if calls[i].when() == whichToRemove {
				newCalls = append(newCalls, calls[i])
			}
		}
		hit.SetSteps(calls)
	})
}

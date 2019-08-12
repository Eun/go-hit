package hit

import (
	"encoding/json"
	"net/http"

	"io"

	"os"

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
	T() TestingT
	SetRequest(r *http.Request) Hit
	Request() *HTTPRequest
	Response() *HTTPResponse
	Connect(url string) Hit
	Delete(url string) Hit
	Get(url string) Hit
	Head(url string) Hit
	Post(url string) Hit
	Options(url string) Hit
	Put(url string) Hit
	Trace(url string) Hit

	SetHTTPClient(client *http.Client) Hit
	HTTPClient() *http.Client
	SetStdout(w io.Writer) Hit
	Stdout() io.Writer
	// SetBaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, WithMethod call
	SetBaseURL(string) Hit

	// BaseURL returns the current base url
	BaseURL() string

	// Expect expects the body to be equal the specified value, omit the parameter to get more options
	// Examples:
	//           Expect("Hello World")
	//           Expect().Body().Contains("Hello World")
	Expect(data ...interface{}) Expect

	// Send sends the specified data as the body payload
	// Examples:
	//           Send("Hello World")
	//           Send().Body("Hello World")
	Send(data ...interface{}) Send
	Do() Hit
	State() State
	Debug() Hit
	Copy() Hit
}

type defaultInstance struct {
	t        TestingT
	request  *HTTPRequest
	response *HTTPResponse
	client   *http.Client
	state    State
	send     *defaultSend
	expect   *defaultExpect
	stdout   io.Writer
	baseURL  string
}

// SetRequest sets the request for the current instance
func (hit *defaultInstance) SetRequest(request *http.Request) Hit {
	hit.request = newHTTPRequest(hit, request)
	return hit
}

func (hit *defaultInstance) setMethodAndUrl(method, url string) Hit {
	request, err := http.NewRequest(method, hit.baseURL+url, nil)
	errortrace.Panic.NoError(hit.T(), err, "unable to create request")
	return hit.SetRequest(request)
}

func (hit *defaultInstance) Connect(url string) Hit {
	return hit.setMethodAndUrl("CONNECT", url)
}
func (hit *defaultInstance) Delete(url string) Hit {
	return hit.setMethodAndUrl("DELETE", url)
}

func (hit *defaultInstance) Get(url string) Hit {
	return hit.setMethodAndUrl("GET", url)
}

func (hit *defaultInstance) Head(url string) Hit {
	return hit.setMethodAndUrl("HEAD", url)
}

func (hit *defaultInstance) Post(url string) Hit {
	return hit.setMethodAndUrl("POST", url)
}

func (hit *defaultInstance) Options(url string) Hit {
	return hit.setMethodAndUrl("OPTIONS", url)
}

func (hit *defaultInstance) Put(url string) Hit {
	return hit.setMethodAndUrl("PUT", url)
}

func (hit *defaultInstance) Trace(url string) Hit {
	return hit.setMethodAndUrl("TRACE", url)
}

func (hit *defaultInstance) T() TestingT {
	return hit.t
}

func (hit *defaultInstance) Request() *HTTPRequest {
	return hit.request
}

func (hit *defaultInstance) Response() *HTTPResponse {
	return hit.response
}

func (hit *defaultInstance) SetHTTPClient(client *http.Client) Hit {
	hit.client = client
	return hit
}

func (hit *defaultInstance) HTTPClient() *http.Client {
	return hit.client
}

func (hit *defaultInstance) SetStdout(w io.Writer) Hit {
	hit.stdout = w
	return hit
}

func (hit *defaultInstance) Stdout() io.Writer {
	return hit.stdout
}

// Expects expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func (hit *defaultInstance) Expect(data ...interface{}) Expect {
	if arg := getLastArgument(data); arg != nil {
		hit.expect.Interface(arg)
	}
	return hit.expect
}

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func (hit *defaultInstance) Send(data ...interface{}) Send {
	if arg := getLastArgument(data); arg != nil {
		hit.send.Interface(arg)
	}
	return hit.send
}

func (hit *defaultInstance) Do() Hit {
	if hit.State() != Ready {
		errortrace.Panic.Errorf(hit.T(), "request already fired")
	}

	hit.state = Working
	hit.runSendCalls()

	hit.request.Request.Body = hit.request.Body().Reader()

	res, err := hit.client.Do(hit.request.Request)
	errortrace.Panic.NoError(hit.T(), err, "unable to perform request")
	hit.response = newHTTPResponse(hit, res)
	hit.runExpectCalls()
	if hit.request.Request.Body != nil {
		_ = hit.request.Request.Body.Close()
	}
	hit.state = Done
	return hit
}

func (hit *defaultInstance) State() State {
	return hit.state
}

func (hit *defaultInstance) runSendCalls() {
	send := hit.send.CollectedSends()
	for i := 0; i < len(send); i++ {
		send[i](hit)
	}
}

func (hit *defaultInstance) runExpectCalls() {
	expect := hit.expect.CollectedExpects()
	for i := 0; i < len(expect); i++ {
		expect[i](hit)
	}
}

func (hit *defaultInstance) Debug() Hit {
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
			"Header":           getHeader(hit.request.Header),
			"Trailer":          getHeader(hit.request.Trailer),
			"Method":           hit.request.Method,
			"URL":              hit.request.URL,
			"Proto":            hit.request.Proto,
			"ProtoMajor":       hit.request.ProtoMajor,
			"ProtoMinor":       hit.request.ProtoMinor,
			"ContentLength":    hit.request.ContentLength,
			"TransferEncoding": hit.request.TransferEncoding,
			"Host":             hit.request.Host,
			"Form":             hit.request.Form,
			"PostForm":         hit.request.PostForm,
			"MultipartForm":    hit.request.MultipartForm,
			"RemoteAddr":       hit.request.RemoteAddr,
			"RequestURI":       hit.request.RequestURI,
			"Body":             getBody(hit.request.body),
		},
	}

	if hit.response != nil {
		m["Response"] = M{
			"Header":           getHeader(hit.response.Header),
			"Trailer":          getHeader(hit.response.Trailer),
			"Proto":            hit.response.Proto,
			"ProtoMajor":       hit.response.ProtoMajor,
			"ProtoMinor":       hit.response.ProtoMinor,
			"ContentLength":    hit.response.ContentLength,
			"TransferEncoding": hit.response.TransferEncoding,
			"Body":             getBody(hit.response.body),
			"Status":           hit.response.Status,
			"StatusCode":       hit.response.StatusCode,
		}
	}

	bytes, err := json.Marshal(m)
	errortrace.Panic.NoError(hit.T(), err)
	_, _ = hit.stdout.Write(pretty.Color(pretty.Pretty(bytes), nil))

	return hit
}

// Copy creates a new instance with the values of the parent instance
func (hit *defaultInstance) Copy() Hit {
	e := &defaultInstance{
		t:       hit.t,
		client:  hit.client,
		stdout:  hit.stdout,
		baseURL: hit.baseURL,
	}
	// copy hit.request.Request

	if hit.request != nil {
		e.request = hit.request.copy(e)
	}

	// copy send
	e.send = hit.send.copy(e)
	// copy expect
	e.expect = hit.expect.copy(e)
	return e
}

// SetBaseURL sets the base url for each Connect, Delete, Get, Head, Post, Options, Put, Trace, WithMethod call
func (hit *defaultInstance) SetBaseURL(s string) Hit {
	hit.baseURL = s
	return hit
}

func (hit *defaultInstance) BaseURL() string {
	return hit.baseURL
}

func New(t TestingT) Hit {
	e := &defaultInstance{
		t:      t,
		client: http.DefaultClient,
		stdout: os.Stdout,
	}
	e.send = newSend(e)
	e.expect = newExpect(e)
	return e
}

// WithRequest creates a new Hit instance with an existing http request
func WithRequest(t TestingT, request *http.Request) Hit {
	return New(t).SetRequest(request)
}

func WithMethod(t TestingT, method, url string) Hit {
	request, err := http.NewRequest(method, url, nil)
	errortrace.Panic.NoError(t, err, "unable to create request")
	return WithRequest(t, request)
}

func Connect(t TestingT, url string) Hit {
	return WithMethod(t, "CONNECT", url)
}

func Delete(t TestingT, url string) Hit {
	return WithMethod(t, "DELETE", url)
}

func Get(t TestingT, url string) Hit {
	return WithMethod(t, "GET", url)
}

func Head(t TestingT, url string) Hit {
	return WithMethod(t, "HEAD", url)
}

func Post(t TestingT, url string) Hit {
	return WithMethod(t, "POST", url)
}

func Options(t TestingT, url string) Hit {
	return WithMethod(t, "OPTIONS", url)
}

func Put(t TestingT, url string) Hit {
	return WithMethod(t, "PUT", url)
}

func Trace(t TestingT, url string) Hit {
	return WithMethod(t, "TRACE", url)
}

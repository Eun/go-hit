package hit

import (
	"encoding/json"
	"net/http"

	"io"

	"os"

	"fmt"

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

	// Connect uses the CONNECT http method in the request, use the optional arguments to format the url
	// Connect("http://%s/%s", domain, path)
	Connect(url string, args ...interface{}) Hit

	// Delete uses the DELETE http method in the request, use the optional arguments to format the url
	// Delete("http://%s/%s", domain, path)
	Delete(url string, args ...interface{}) Hit

	// Get uses the GET http method in the request, use the optional arguments to format the url
	// Get("http://%s/%s", domain, path)
	Get(url string, args ...interface{}) Hit

	// Head uses the HEAD http method in the request, use the optional arguments to format the url
	// Head("http://%s/%s", domain, path)
	Head(url string, args ...interface{}) Hit

	// Post uses the POST http method in the request, use the optional arguments to format the url
	// Post("http://%s/%s", domain, path)
	Post(url string, args ...interface{}) Hit

	// Options uses the OPTIONS http method in the request, use the optional arguments to format the url
	// Options("http://%s/%s", domain, path)
	Options(url string, args ...interface{}) Hit

	// Put uses the PUT http method in the request, use the optional arguments to format the url
	// Put("http://%s/%s", domain, path)
	Put(url string, args ...interface{}) Hit

	// Trace uses the TRACE http method in the request, use the optional arguments to format the url
	// Trace("http://%s/%s", domain, path)
	Trace(url string, args ...interface{}) Hit

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

func (hit *defaultInstance) setMethodAndUrl(method, url string, a ...interface{}) Hit {
	request, err := http.NewRequest(method, hit.baseURL+fmt.Sprintf(url, a...), nil)
	errortrace.Panic.NoError(hit.T(), err, "unable to create request")
	return hit.SetRequest(request)
}

func (hit *defaultInstance) Connect(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("CONNECT", url, a...)
}
func (hit *defaultInstance) Delete(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("DELETE", url, a...)
}

func (hit *defaultInstance) Get(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("GET", url, a...)
}

func (hit *defaultInstance) Head(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("HEAD", url, a...)
}

func (hit *defaultInstance) Post(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("POST", url, a...)
}

func (hit *defaultInstance) Options(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("OPTIONS", url, a...)
}

func (hit *defaultInstance) Put(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("PUT", url, a...)
}

func (hit *defaultInstance) Trace(url string, a ...interface{}) Hit {
	return hit.setMethodAndUrl("TRACE", url, a...)
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

// WithMethod creates a new Hit instance with the specified method and url
// WithMethod(t, "POST", "http://%s/%s", domain, path)
func WithMethod(t TestingT, method, url string, a ...interface{}) Hit {
	request, err := http.NewRequest(method, fmt.Sprintf(url, a...), nil)
	errortrace.Panic.NoError(t, err, "unable to create request")
	return WithRequest(t, request)
}

// Connect creates a new Hit instance with CONNECT as the http method, use the optional arguments to format the url
// Connect(t, "http://%s/%s", domain, path)
func Connect(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "CONNECT", url, a...)
}

// Delete creates a new Hit instance with DELETE as the http method, use the optional arguments to format the url
// Delete(t, "http://%s/%s", domain, path)
func Delete(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "DELETE", url, a...)
}

// Get creates a new Hit instance with GET as the http method, use the optional arguments to format the url
// Get(t, "http://%s/%s", domain, path)
func Get(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "GET", url, a...)
}

// Head creates a new Hit instance with HEAD as the http method, use the optional arguments to format the url
// Head(t, "http://%s/%s", domain, path)
func Head(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "HEAD", url, a...)
}

// Post creates a new Hit instance with POST as the http method, use the optional arguments to format the url
// Post(t, "http://%s/%s", domain, path)
func Post(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "POST", url, a...)
}

// Options creates a new Hit instance with OPTIONS as the http method, use the optional arguments to format the url
// Options(t, "http://%s/%s", domain, path)
func Options(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "OPTIONS", url, a...)
}

// Put creates a new Hit instance with PUT as the http method, use the optional arguments to format the url
// Put(t, "http://%s/%s", domain, path)
func Put(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "PUT", url, a...)
}

// Trace creates a new Hit instance with TRACE as the http method, use the optional arguments to format the url
// Trace(t, "http://%s/%s", domain, path)
func Trace(t TestingT, url string, a ...interface{}) Hit {
	return WithMethod(t, "TRACE", url, a...)
}

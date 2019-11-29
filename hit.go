package hit

import (
	"io"
	"net/http"

	"fmt"
)

//go:generate go run generate_template_framework.go

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

	Send(...interface{}) ISend
	Expect(...interface{}) IExpect
	Debug(...string) IStep
	AddSteps(...IStep)

	Description() string
	SetDescription(string)
}

type defaultInstance struct {
	steps       []IStep
	request     *HTTPRequest
	response    *HTTPResponse
	client      *http.Client
	state       State
	stdout      io.Writer
	baseURL     string
	description string
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

func (hit *defaultInstance) runSteps(state StepTime) error {
	totalSteps := len(hit.steps)
	// find all steps we want to run
	stepsToRun := hit.collectSteps(state, 0)
	size := len(stepsToRun)
	i := 0
	for {
		if i >= size {
			return nil
		}
		if err := stepsToRun[i].exec(hit); err != nil {
			return err
		}

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
	return nil
}

// Send sends the specified data as the body payload
// Examples:
//           Send("Hello World")
//           Send().Body("Hello World")
func (hit *defaultInstance) Send(data ...interface{}) ISend {
	if arg, ok := getLastArgument(data); ok {
		return finalSend{newSend(hit).Interface(arg)}
	}
	return newSend(hit)
}

// Expect expects the body to be equal the specified value, omit the parameter to get more options
// Examples:
//           Expect("Hello World")
//           Expect().Body().Contains("Hello World")
func (hit *defaultInstance) Expect(data ...interface{}) IExpect {
	if arg, ok := getLastArgument(data); ok {
		return finalExpect{newExpect(hit).Interface(arg)}
	}
	return newExpect(hit)
}

// Debug prints the current Request and Response to hit.Stdout(), you can filter the output based on expressions
// Examples:
//           Debug()
//           Debug("Request")
func (hit *defaultInstance) Debug(expression ...string) IStep {
	return Debug(expression...)
}

// AddSteps add the specified steps to the queue
func (hit *defaultInstance) AddSteps(steps ...IStep) {
	for i := 0; i < len(steps); i++ {
		hit.steps = append(hit.steps, steps[i])
	}
}

// Description gets the current description that will be printed in an error case
func (hit *defaultInstance) Description() string {
	return hit.description
}

// SetDescription sets a custom description for this test.
// The description will be printed in an error case
func (hit *defaultInstance) SetDescription(description string) {
	hit.description = description
}

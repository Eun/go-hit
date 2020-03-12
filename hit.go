package hit

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/xerrors"
)

//go:generate go run generate_template_framework.go

// Callback will be used for Custom() functions
type Callback func(hit Hit)

type Hit interface {
	// Request returns the current request
	Request() *HTTPRequest

	// SetRequest sets the request for the current instance
	SetRequest(request *http.Request)

	// Response returns the current Response
	Response() *HTTPResponse

	// HTTPClient gets the current http.Client
	HTTPClient() *http.Client

	// SetHTTPClient sets the client for the request
	SetHTTPClient(client *http.Client)

	// Stdout gets the current output
	Stdout() io.Writer

	// SetStdout sets the output to the specified writer
	SetStdout(w io.Writer)

	// BaseURL returns the current base url
	BaseURL() string

	// SetBaseURL sets the base url
	SetBaseURL(url string, a ...interface{})

	// CurrentStep returns the current working step
	CurrentStep() IStep

	// Steps returns the current step list
	Steps() []IStep

	// AddSteps adds the specified steps to the step list
	AddSteps(steps ...IStep)

	// RemoveSteps removes the specified steps from the step list
	RemoveSteps(steps ...IStep)

	// Do runs the specified steps in in this context.
	//
	// Example:
	//     Expect().Custom(func(hit Hit) {
	//         err := Do(
	//             Expect().Status(200),
	//         )
	//         if err != nil {
	//             panic(err)
	//         }
	//     })
	Do(steps ...IStep) error

	// MustDo runs the specified steps in in this context and panics on failure
	// Example:
	//           Expect().Custom(func(hit Hit) {
	//               MustDo(
	//                   Expect().Status(200),
	//               )
	//           })
	MustDo(steps ...IStep)

	// Description gets the current description that will be printed in an error case
	Description() string

	// SetDescription sets a custom description for this test.
	// The description will be printed in an error case
	SetDescription(string)
}

type defaultInstance struct {
	steps       []IStep
	currentStep IStep
	request     *HTTPRequest
	response    *HTTPResponse
	client      *http.Client
	state       StepTime
	stdout      io.Writer
	baseURL     string
	description string
}

func (hit *defaultInstance) Request() *HTTPRequest {
	return hit.request
}

func (hit *defaultInstance) SetRequest(request *http.Request) {
	hit.request = newHTTPRequest(hit, request)
}

func (hit *defaultInstance) Response() *HTTPResponse {
	return hit.response
}

func (hit *defaultInstance) HTTPClient() *http.Client {
	return hit.client
}

func (hit *defaultInstance) SetHTTPClient(client *http.Client) {
	hit.client = client
}

func (hit *defaultInstance) Stdout() io.Writer {
	return hit.stdout
}

func (hit *defaultInstance) SetStdout(w io.Writer) {
	hit.stdout = w
}

func (hit *defaultInstance) BaseURL() string {
	return hit.baseURL
}

func (hit *defaultInstance) SetBaseURL(url string, a ...interface{}) {
	hit.baseURL = fmt.Sprintf(url, a...)
}

func (hit *defaultInstance) collectSteps(state StepTime, offset int) []IStep {
	var collectedSteps []IStep
	for i := 0; i < len(hit.steps); i++ {
		w := hit.steps[i].when()
		if w == state {
			// skip the offset
			if offset > 0 {
				offset--
				continue
			}
			collectedSteps = append(collectedSteps, hit.steps[i])
		}
	}
	return collectedSteps
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
		hit.currentStep = stepsToRun[i]
		if err := stepsToRun[i].exec(hit); err != nil {
			return err
		}

		// maybe there is a new step added
		newTotalSteps := len(hit.steps)
		if totalSteps != newTotalSteps {
			// yes they have been modified
			// find all new steps after the last scan
			stepsToRun = append(stepsToRun[:i], hit.collectSteps(state, i+1)...)
			size = len(stepsToRun)
			i = 0
			totalSteps = newTotalSteps
			continue
		}
		i++
	}
	return nil
}

func (hit *defaultInstance) CurrentStep() IStep {
	return hit.currentStep
}

func (hit *defaultInstance) Steps() []IStep {
	return hit.steps
}

func (hit *defaultInstance) AddSteps(steps ...IStep) {
	for i := 0; i < len(steps); i++ {
		hit.steps = append(hit.steps, steps[i])
	}
}

func (hit *defaultInstance) RemoveSteps(steps ...IStep) {
removeStep:
	for j := len(steps) - 1; j >= 0; j-- {
		for i := len(hit.steps) - 1; i >= 0; i-- {
			if hit.steps[i] == steps[j] {
				// remove the step from steps and hit.steps
				steps = append(steps[:j], steps[j+1:]...)
				hit.steps = append(hit.steps[:i], hit.steps[i+1:]...)
				continue removeStep
			}
		}
	}
}

func (hit *defaultInstance) Do(steps ...IStep) error {
	for _, step := range steps {
		if step.when() != hit.state && step.when() != CleanStep {
			return xerrors.Errorf("unable to execute `%s' during %s, can only be run during %s", step.clearPath().String(), hit.state.String(), step.when().String())
		}
		if err := step.exec(hit); err != nil {
			return err
		}
	}
	return nil
}

func (hit *defaultInstance) MustDo(steps ...IStep) {
	if err := hit.Do(steps...); err != nil {
		panic(err)
	}
}

func (hit *defaultInstance) Description() string {
	return hit.description
}

func (hit *defaultInstance) SetDescription(description string) {
	hit.description = description
}

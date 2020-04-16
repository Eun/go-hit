package hit

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// IClear provides a clear functionality to remove previous steps from running
type IClear interface {
	// Send removes all previous Send() steps and all steps chained to Send() e.g. Send().Body("Hello World").
	//
	// Usage:
	//     Clear().Send()                      // will remove all Send() steps and all chained steps to Send() e.g. Send().Body("Hello World")
	//     Clear().Send().Body()               // will remove all Send().Body() steps and all chained steps to Body() e.g. Send().Body().Equal("Hello World")
	//     Clear().Send().Body("Hello World")  // will remove all Send().Body("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body("Hello Earth"),
	//         Clear().Send(),
	//         Send().Body("Hello World"),
	//     )
	Send() IClearSend

	// Expect removes all previous Expect() steps and all steps chained to Expect() e.g. Expect().Body("Hello World").
	//
	// If you specify an argument it will only remove the Expect() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect()                      // will remove all Expect() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
	//     Clear().Expect().Body()               // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
	//     Clear().Expect().Body("Hello World")  // will remove all Expect().Body("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect().Body("Hello Earth"),
	//         Clear().Expect(),
	//         Expect().Body("Hello World"),
	//     )
	Expect() IClearExpect
}

type clear struct{}

func newClear() IClear {
	return &clear{}
}

func (clr *clear) Send() IClearSend {
	return newClearSend(newClearPath("Send", nil))
}

func (clr *clear) Expect() IClearExpect {
	return newClearExpect(newClearPath("Expect", nil))
}

func removeSteps(hit Hit, path clearPath) error {
	var stepsToRemove []IStep
	steps := hit.Steps()
	availableSteps := make([]IStep, 0, len(steps))

	for _, step := range steps {
		if step == hit.CurrentStep() {
			break
		}
		p := step.clearPath()
		if p == nil {
			continue
		}
		availableSteps = append(availableSteps, step)
		if p.Contains(path) {
			stepsToRemove = append(stepsToRemove, step)
		}
	}
	if len(stepsToRemove) == 0 {
		var sb strings.Builder
		fmt.Fprintf(&sb, "unable to find a step with %s\n", path.CallString())

		if len(availableSteps) > 0 {
			fmt.Fprintf(&sb, "got these steps:\n")
			for _, step := range availableSteps {
				fmt.Fprintf(&sb, "\t%s\n", step.clearPath().CallString())
			}
		}

		return xerrors.New(sb.String())
	}
	hit.RemoveSteps(stepsToRemove...)
	return nil
}

func removeStep(cleanPath clearPath) *hitStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			return removeSteps(hit, cleanPath)
		},
	}
}

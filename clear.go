// +build !generate

package hit

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

// IClear provides a clear functionality to remove previous steps from running.
type IClear interface {
	// Send removes all previous steps chained to Send()
	// e.g. Send().Body().String("Hello World") or Send().Headers("Content-Type", "application/json").
	//
	// Usage:
	//     // will remove all previous steps chained to Send() e.g. Send().Body().String("Hello World") or
	//     // Send().Headers("Content-Type", "application/json")
	//     Clear().Send()
	//
	//     // will remove all previous steps chained to Send().Body() e.g. Send().Body().String("Hello World") or
	//     // Send().Body().Int(127)
	//     Clear().Send().Body()
	//
	//     // will remove all previous Send().Body().String() steps e.g. Send().Body().String("Hello World") or
	//     // Send().Body().String("Hello Earth")
	//     Clear().Send().Body().String()
	//
	//     // will remove all previous Send().Body().String("Hello World") steps
	//     Clear().Send().Body().String("Hello World")
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send().Body().String("Hello Earth"),
	//         Clear().Send(),
	//         Send().Body().String("Hello World"),
	//     )
	Send() IClearSend

	// Expect removes all previous steps chained to Expect()
	// e.g. Expect().Body().String().Equal("Hello World") or Expect().Headers("Content-Type").Equal("application/json").
	//
	// Usage:
	//     // will remove all previous steps chained to Expect() e.g. Expect().Body().String().Equal("Hello World")
	//     // or Expect().Status().Equal(200)
	//     Clear().Expect()
	//
	//     // will remove all previous steps chained to Expect().Body().String() e.g.
	//     // Expect().Body().String().Equal("Hello World") or Expect().Body().String().Contains("Hello")
	//     Clear().Expect().Body().String()
	//
	//     // will remove all previous Expect().Body().String().Equal() steps e.g.
	//     // Expect().Body().String().Equal("Hello World") or Expect().Body().String().Contains("Hello")
	//     Clear().Expect().Body().String().Equal()
	//
	//     // will remove all previous Expect().Body().String().Equal("Hello World") steps
	//     Clear().Expect().Body().String().Equal("Hello World")
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Expect().Body().String().Equal("Hello Earth"),
	//         Clear().Expect(),
	//         Expect().Body().String().Equal("Hello World"),
	//     )
	Expect() IClearExpect
}

type clear struct {
	cp callPath
}

func newClear(cp callPath) IClear {
	return &clear{
		cp: cp,
	}
}

func (clr *clear) Send() IClearSend {
	return newClearSend(clr.cp.Push("Send", nil))
}

func (clr *clear) Expect() IClearExpect {
	return newClearExpect(clr.cp.Push("Expect", nil))
}

func removeSteps(hit Hit, path callPath) error {
	path = path[1:] // remove Clean()
	var stepsToRemove []IStep
	steps := hit.Steps()
	availableSteps := make([]IStep, 0, len(steps))

	for _, step := range steps {
		if step == hit.CurrentStep() {
			break
		}
		p := step.callPath()
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
		fmt.Fprintf(&sb, "unable to find a step with %s\n", path.CallString(true))

		if len(availableSteps) > 0 {
			fmt.Fprintf(&sb, "got these steps:\n")
			for _, step := range availableSteps {
				fmt.Fprintf(&sb, "\t%s\n", step.callPath().CallString(true))
			}
		}

		return xerrors.New(sb.String())
	}
	hit.RemoveSteps(stepsToRemove...)
	return nil
}

func removeStep(cleanPath callPath) *hitStep {
	return &hitStep{
		Trace:    ett.Prepare(),
		When:     cleanStep,
		CallPath: cleanPath,
		Exec: func(hit *hitImpl) error {
			return removeSteps(hit, cleanPath)
		},
	}
}

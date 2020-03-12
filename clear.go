package hit

// IClear provides a clear functionality to remove previous steps from running
type IClear interface {
	// Send removes all previous Send() steps and all steps chained to Send() e.g. Send().Body("Hello World").
	//
	// If you specify an argument it will only remove the Send() steps matching that argument.
	//
	// Usage:
	//     Clear().Send()                      // will remove all Send() steps and all chained steps to Send() e.g. Send().Body("Hello World")
	//     Clear().Send("Hello World")         // will remove all Send("Hello World") steps
	//     Clear().Send().Body()               // will remove all Send().Body() steps and all chained steps to Body() e.g. Send().Body().Equal("Hello World")
	//     Clear().Send().Body("Hello World")  // will remove all Send().Body("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Send("Hello Earth"),
	//         Clear().Send(),
	//         Send("Hello World"),
	//     )
	Send(value ...interface{}) IClearSend

	// Expect removes all previous Expect() steps and all steps chained to Expect() e.g. Expect().Body("Hello World").
	//
	// If you specify an argument it will only remove the Expect() steps matching that argument.
	//
	// Usage:
	//     Clear().Expect()                      // will remove all Expect() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
	//     Clear().Expect("Hello World")         // will remove all Expect("Hello World") steps
	//     Clear().Expect().Body()               // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
	//     Clear().Expect().Body("Hello World")  // will remove all Expect().Body("Hello World") steps
	//
	// Example:
	//     MustDo(
	//         Post("https://example.com"),
	//         Expect("Hello Earth"),
	//         Clear().Expect(),
	//         Expect("Hello World"),
	//     )
	Expect(value ...interface{}) IClearExpect
}

type clear struct{}

func newClear() IClear {
	return &clear{}
}

func (clr *clear) Send(value ...interface{}) IClearSend {
	return newClearSend(newClearPath("Send", value), value)
}

func (clr *clear) Expect(value ...interface{}) IClearExpect {
	return newClearExpect(newClearPath("Expect", value), value)
}

func removeSteps(hit Hit, path clearPath) {
	var stepsToRemove []IStep

	for _, step := range hit.Steps() {
		if step == hit.CurrentStep() {
			break
		}
		if step.clearPath().Contains(path) {
			stepsToRemove = append(stepsToRemove, step)
		}
	}
	hit.RemoveSteps(stepsToRemove...)
}

func removeStep(cleanPath clearPath) *hitStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil, // not clearable
		Exec: func(hit Hit) error {
			removeSteps(hit, cleanPath)
			return nil
		},
	}
}

package hit

type IClear interface {
	// Send removes all Send() steps and all steps chained to Send() e.g. Send().Body("Hello World")
	// if you specify an argument it will only remove the Send() steps matching that argument
	// Examples:
	//           Clear().Send()                      // will remove all Send() steps and all chained steps to Send() e.g. Send().Body("Hello World")
	//           Clear().Send("Hello World")         // will remove all Send("Hello World") steps
	//           Clear().Send().Body()               // will remove all Send().Body() steps and all chained steps to Body() e.g. Send().Body().Equal("Hello World")
	//           Clear().Send().Body("Hello World")  // will remove all Send().Body("Hello World") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Send("Hello World"),
	//               Clear().Send(),
	//               Send("Hello Universe"),
	//           )
	Send(...interface{}) IClearSend

	// Expect removes all Expect() steps and all steps chained to Expect() e.g. Expect().Body("Hello World")
	// if you specify an argument it will only remove the Expect() steps matching that argument
	// Examples:
	//           Clear().Expect()                      // will remove all Expect() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
	//           Clear().Expect("Hello World")         // will remove all Expect("Hello World") steps
	//           Clear().Expect().Body()               // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
	//           Clear().Expect().Body("Hello World")  // will remove all Expect().Body("Hello World") steps
	//
	//           Do(
	//               Post("https://example.com"),
	//               Expect("Hello World"),
	//               Clear().Expect(),
	//               Expect("Hello Universe"),
	//           )
	Expect(...interface{}) IClearExpect
}

type clear struct {
}

func newClear() IClear {
	return &clear{}
}

// Send removes all Send() steps and all steps chained to Send() e.g. Send().Body("Hello World")
// if you specify an argument it will only remove the Send() steps matching that argument
// Examples:
//           Clear().Send()                      // will remove all Send() steps and all chained steps to Send() e.g. Send().Body("Hello World")
//           Clear().Send("Hello World")         // will remove all Send("Hello World") steps
//           Clear().Send().Body()               // will remove all Send().Body() steps and all chained steps to Body() e.g. Send().Body().Equal("Hello World")
//           Clear().Send().Body("Hello World")  // will remove all Send().Body("Hello World") steps
//
//           Do(
//               Post("https://example.com"),
//               Send("Hello World"),
//               Clear().Send(),
//               Send("Hello Universe"),
//           )
func (clr *clear) Send(data ...interface{}) IClearSend {
	return newClearSend(newClearPath("Send", data), data)
}

// Expect removes all Expect() steps and all steps chained to Expect() e.g. Expect().Body("Hello World")
// if you specify an argument it will only remove the Expect() steps matching that argument
// Examples:
//           Clear().Expect()                      // will remove all Expect() steps and all chained steps to Expect() e.g. Expect().Body("Hello World")
//           Clear().Expect("Hello World")         // will remove all Expect("Hello World") steps
//           Clear().Expect().Body()               // will remove all Expect().Body() steps and all chained steps to Body() e.g. Expect().Body().Equal("Hello World")
//           Clear().Expect().Body("Hello World")  // will remove all Expect().Body("Hello World") steps
//
//           Do(
//               Post("https://example.com"),
//               Expect("Hello World"),
//               Clear().Expect(),
//               Expect("Hello Universe"),
//           )
func (clr *clear) Expect(data ...interface{}) IClearExpect {
	return newClearExpect(newClearPath("Expect", data), data)
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

func removeStep(cleanPath clearPath) IStep {
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

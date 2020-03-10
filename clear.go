package hit

type IClear interface {
	Send(...interface{}) IClearSend
	Expect(...interface{}) IClearExpect
}

type clear struct {
}

func newClear() IClear {
	return &clear{}
}

func (clr *clear) Send(data ...interface{}) IClearSend {
	return newClearSend(newClearPath("Send", data), data)
}

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

package hit

type IClear interface {
	// Send() IClearSend
	Expect(...interface{}) IClearExpect
}

type clear struct {
	hit Hit
}

func newClear(hit Hit) IClear {
	return &clear{
		hit: hit,
	}
}

// func (clr *clear) Send() IClearSend {
// 	return nil
// }

func (clr *clear) Expect(data ...interface{}) IClearExpect {
	return newClearExpect(clr, NewCleanPath("Expect", data), data)
}

func removeSteps(hit Hit, path CleanPath) {
	var stepsToRemove []IStep

	for _, step := range hit.Steps() {
		// if internal.StringSliceHasPrefixSlice(step.CleanPath(), path) {
		// 	stepsToRemove = append(stepsToRemove, step)
		// }
		if step.CleanPath().Contains(path) {
			stepsToRemove = append(stepsToRemove, step)
		}
	}
	hit.RemoveSteps(stepsToRemove...)
}

func removeStep(cleanPath CleanPath) IStep {
	return custom(Step{
		When:      CleanStep,
		CleanPath: nil, // not clearable
		Exec: func(hit Hit) {
			removeSteps(hit, cleanPath)
		},
	})
}

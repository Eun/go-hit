package hit

import (
	"github.com/Eun/go-hit/internal"
)

type IClear interface {
	// Send() IClearSend
	Expect() IClearExpect
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

func (clr *clear) Expect() IClearExpect {
	return newClearExpect(clr, clr.hit, NewCleanPath("Expect"))
}

func removeSteps(hit Hit, path []string) {
	var stepsToRemove []IStep

	for _, step := range hit.Steps() {
		if internal.StringSliceHasPrefixSlice(step.CleanPath(), path) {
			stepsToRemove = append(stepsToRemove, step)
		}
	}
	hit.RemoveSteps(stepsToRemove...)
}

func removeStep(hit Hit, cleanPath []string) IStep {
	return custom(Step{
		When:      CleanStep,
		CleanPath: nil, // not clearable
		Instance:  hit,
		Exec: func(hit Hit) {
			removeSteps(hit, cleanPath)
		},
	})
}

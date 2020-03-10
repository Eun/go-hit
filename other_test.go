package hit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollectSteps(t *testing.T) {
	a := &hitStep{
		When:      SendStep,
		ClearPath: newClearPath("a", nil),
	}
	b := &hitStep{
		When:      ExpectStep,
		ClearPath: newClearPath("b", nil),
	}
	c := &hitStep{
		When:      ExpectStep,
		ClearPath: newClearPath("c", nil),
	}
	d := &hitStep{
		When:      SendStep,
		ClearPath: newClearPath("d", nil),
	}
	hit := defaultInstance{steps: []IStep{a, b, c, d}}
	steps := hit.collectSteps(SendStep, 0)
	require.Equal(t, []IStep{a, d}, steps)

	steps = hit.collectSteps(SendStep, 1)
	require.Equal(t, []IStep{d}, steps)

	steps = hit.collectSteps(ExpectStep, 1)
	require.Equal(t, []IStep{c}, steps)
}

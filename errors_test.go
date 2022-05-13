package hit_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/otto-eng/go-hit"
)

func TestError_FailingStepIs(t *testing.T) {
	s := EchoServer()
	defer s.Close()

	err := Do(
		Post(s.URL),
		Send().Body().Int8(16),
		Expect().Body().Int8().Equal(16),
		Expect().Body().Int8().Equal(17),
	)

	var hitError *Error
	require.True(t, errors.As(err, &hitError))
	require.False(t, hitError.FailingStepIs(Expect().Body().Int8().Equal(16)))
	require.True(t, hitError.FailingStepIs(Expect().Body().Int8().Equal(17)))
}

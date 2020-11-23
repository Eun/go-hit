package hit

import (
	"github.com/Eun/go-hit/errortrace"
)

// Error represents the error that will be returned during an execution.
type Error struct {
	callPath callPath
	et       *errortrace.ErrorTrace
}

// Error returns the string representation for the error.
func (e *Error) Error() string {
	return e.et.Error()
}

// Implement xerrors

// Is implements the xerrors interface so we can use the xerrors.Is() function.
func (e *Error) Is(err error) bool {
	return e.et == err
}

// Unwrap implements the xerrors.Wrapper interface.
func (e *Error) Unwrap() error {
	return e.et
}

// FailingStepIs returns true if the specified step is the same as the failing.
//
// Example:
//     err := Do(
//         Get("https://example.com"),
//         Expect().Status().Equal(http.StatusNoContent),
//     )
//     var hitError *Error
//     if errors.As(err, &hitError) {
//         if hitError.FailingStepIs(Expect().Status().Equal(http.StatusNoContent)) {
//             fmt.Printf("Expected StatusNoContent")
//             return
//         }
//     }
func (e *Error) FailingStepIs(s IStep) bool {
	if e.callPath == nil || s == nil {
		return false
	}
	cp := s.callPath()
	if cp == nil {
		return false
	}
	return e.callPath.Equal(cp)
}

func wrapError(hit Hit, err error) *Error {
	t := ett.Prepare()
	t.SetError(err)
	t.SetDescription(hit.Description())
	return &Error{
		callPath: nil,
		et:       t,
	}
}

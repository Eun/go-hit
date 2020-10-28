package hit

// IExpectFormValues provides assertions on the http response body FormValues.
type IExpectFormValues interface {
	// Contains expects the specific header to contain all of the specified values.
	//
	// Usage:
	//     Expect().Body().FormValues("username").Contains("joe")
	Contains(values ...interface{}) IStep

	// NotContains expects the specified header to not contain all of the specified values.
	//
	// Usage:
	//     Expect().Body().FormValues("username").NotContains("joe")
	NotContains(values ...interface{}) IStep

	// OneOf expects the specified header to contain one of the specified values.
	//
	// Usage:
	//     Expect().Body().FormValues("username").OneOf("joe", "alice")
	OneOf(values ...interface{}) IStep

	// NotOneOf expects the specified header to not contain one of the specified values.
	//
	// Usage:
	//     Expect().Body().FormValues("username").NotOneOf("joe", "alice")
	NotOneOf(values ...interface{}) IStep

	// Empty expects the specified header to be empty.
	//
	// Usage:
	//     Expect().Body().FormValues("username").Empty()
	Empty() IStep

	// NotEmpty expects the specified header to be empty.
	//
	// Usage:
	//     Expect().Body().FormValues("username").NotEmpty()
	NotEmpty() IStep

	// Len expects the specified header to be the same length then specified.
	//
	// Usage:
	//     Expect().Body().FormValues("username").Len().GreaterThan(0)
	Len() IExpectInt

	// Equal expects the specified header to be equal the specified value.
	//
	// Usage:
	//     Expect().Body().FormValues("username").Equal("joe")
	//     Expect().Body().FormValues("usernames").Equal("joe", "alice")
	//     Expect().Body().FormValues("length").Equal(10)
	Equal(value ...interface{}) IStep

	// NotEqual expects the specified header to be not equal the specified value.
	//
	// Usage:
	//     Expect().Body().FormValues("username").NotEqual("joe")
	//     Expect().Body().FormValues("usernames").NotEqual("joe", "alice")
	//     Expect().Body().FormValues("length").NotEqual(10)
	NotEqual(value ...interface{}) IStep
}

func newExpectFormValues(cleanPath callPath, valueCallback expectHeaderValueCallback) IExpectFormValues {
	return &expectHeader{
		cleanPath:     cleanPath,
		valueCallback: valueCallback,
	}
}

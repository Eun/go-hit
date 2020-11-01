package hit

// TestingT is the minimum interface that is required for Test().
type TestingT interface {
	FailNow()
}

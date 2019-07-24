package expr

// Option is the type that will be used for Options
type Option uint8

const (
	// IgnoreCase can be used to turn of case sensitivity
	IgnoreCase Option = iota
	// IgnoreNotFound ignores not found values and return nil (only MustGetValue)
	IgnoreNotFound
	// IgnoreError ignores errors and return nil (only MustGetValue)
	IgnoreError
)

type options []Option

func (opts options) HasOption(option Option) bool {
	for _, v := range opts {
		if v == option {
			return true
		}
	}
	return false
}

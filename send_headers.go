package hit

type ISendHeaders interface {
	Set(name, value string) IStep
	Delete(name string) IStep
	Clear() IStep
}

type sendHeaders struct {
	cleanPath CleanPath
}

func newSendHeaders(path CleanPath) ISendHeaders {
	return &sendHeaders{
		cleanPath: path,
	}
}

// Set sets a header to the specified value
func (hdr *sendHeaders) Set(name, value string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Set", []interface{}{name, value}),
		Exec: func(hit Hit) {
			hit.Request().Header.Set(name, value)
		},
	})
}

// Delete deletes a header
func (hdr *sendHeaders) Delete(name string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Delete", []interface{}{name}),
		Exec: func(hit Hit) {
			hit.Request().Header.Del(name)
		},
	})
}

// Clear removes all headers
func (hdr *sendHeaders) Clear() IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Clear", nil),
		Exec: func(hit Hit) {
			var names []string
			for name := range hit.Request().Header {
				names = append(names, name)
			}
			for _, name := range names {
				hit.Request().Header.Del(name)
			}
		},
	})
}

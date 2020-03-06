package hit

type ISendHeaders interface {
	Set(name, value string) IStep
	Delete(name string) IStep
	Clear() IStep
}

type sendHeaders struct {
	send      ISend
	hit       Hit
	cleanPath CleanPath
}

func newSendHeaders(send ISend, hit Hit, path CleanPath) ISendHeaders {
	return &sendHeaders{
		send:      send,
		hit:       hit,
		cleanPath: path,
	}
}

// Set sets a header to the specified value
func (hdr *sendHeaders) Set(name, value string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Set"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			hit.Request().Header.Set(name, value)
		},
	})
}

// Delete deletes a header
func (hdr *sendHeaders) Delete(name string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Delete"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			hit.Request().Header.Del(name)
		},
	})
}

// Clear removes all headers
func (hdr *sendHeaders) Clear() IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Clear"),
		Instance:  hdr.hit,
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

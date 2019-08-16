package hit

type ISendHeaders interface {
	IStep
	Set(name, value string) IStep
	Delete(name string) IStep
	Clear() IStep
}

type sendHeaders struct {
	send ISend
}

func newSendHeaders(send ISend) ISendHeaders {
	return &sendHeaders{
		send: send,
	}
}

func (hdr *sendHeaders) when() StepTime {
	return hdr.send.when()
}

func (hdr *sendHeaders) exec(hit Hit) error {
	return hdr.send.exec(hit)
}

// Set sets a header to the specified value
func (hdr *sendHeaders) Set(name, value string) IStep {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Set(name, value)
	})
}

// Delete deletes a header
func (hdr *sendHeaders) Delete(name string) IStep {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Del(name)
	})
}

// Clear removes all headers
func (hdr *sendHeaders) Clear() IStep {
	return hdr.send.Custom(func(hit Hit) {
		var names []string
		for name := range hit.Request().Header {
			names = append(names, name)
		}
		for _, name := range names {
			hit.Request().Header.Del(name)
		}
	})
}

package hit

type sendHeaders struct {
	send *defaultSend
}

func newSendHeaders(send *defaultSend) *sendHeaders {
	return &sendHeaders{
		send: send,
	}
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

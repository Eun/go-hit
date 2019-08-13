package hit

type sendHeaders struct {
	*defaultSend
}

func newSendHeaders(send *defaultSend) *sendHeaders {
	return &sendHeaders{
		defaultSend: send,
	}
}

func (hdr *sendHeaders) Set(name, value string) Step {
	return hdr.Custom(func(hit Hit) {
		hit.Request().Header.Set(name, value)
	})
}

func (hdr *sendHeaders) Delete(name string) Step {
	return hdr.Custom(func(hit Hit) {
		hit.Request().Header.Del(name)
	})
}

func (hdr *sendHeaders) Clear() Step {
	return hdr.Custom(func(hit Hit) {
		var names []string
		for name := range hit.Request().Header {
			names = append(names, name)
		}
		for _, name := range names {
			hit.Request().Header.Del(name)
		}
	})
}

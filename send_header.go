package hit

type sendHeaders struct {
	Hit
	send *defaultSend
}

func newSendHeaders(send *defaultSend) *sendHeaders {
	return &sendHeaders{
		Hit:  send.Hit,
		send: send,
	}
}

func (hdr *sendHeaders) Set(name, value string) Hit {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Set(name, value)
	})
}

func (hdr *sendHeaders) Delete(name string) Hit {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Del(name)
	})
}

func (hdr *sendHeaders) Clear() Hit {
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

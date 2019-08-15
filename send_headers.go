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

func (snd *sendHeaders) when() StepTime {
	return snd.send.when()
}

func (snd *sendHeaders) exec(hit Hit) {
	snd.send.exec(hit)
}

// Set sets a header to the specified value
func (snd *sendHeaders) Set(name, value string) IStep {
	return snd.send.Custom(func(hit Hit) {
		hit.Request().Header.Set(name, value)
	})
}

// Delete deletes a header
func (snd *sendHeaders) Delete(name string) IStep {
	return snd.send.Custom(func(hit Hit) {
		hit.Request().Header.Del(name)
	})
}

// Clear removes all headers
func (snd *sendHeaders) Clear() IStep {
	return snd.send.Custom(func(hit Hit) {
		var names []string
		for name := range hit.Request().Header {
			names = append(names, name)
		}
		for _, name := range names {
			hit.Request().Header.Del(name)
		}
	})
}

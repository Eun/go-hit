package hit

type ISendSpecificHeader interface {
	IStep
	Set(value string) IStep
	Delete() IStep
}

type sendSpecificHeader struct {
	send   ISend
	header string
}

func newSendSpecificHeader(send ISend, header string) ISendSpecificHeader {
	return &sendSpecificHeader{
		send:   send,
		header: header,
	}
}

func (snd *sendSpecificHeader) when() StepTime {
	return snd.send.when()
}

func (snd *sendSpecificHeader) exec(hit Hit) {
	snd.send.exec(hit)
}

// Set sets the header to the specified value
func (snd *sendSpecificHeader) Set(value string) IStep {
	return snd.send.Custom(func(hit Hit) {
		hit.Request().Header.Set(snd.header, value)
	})
}

// Delete deletes the header
func (snd *sendSpecificHeader) Delete() IStep {
	return snd.send.Custom(func(hit Hit) {
		hit.Request().Header.Del(snd.header)
	})
}

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

func (hdr *sendSpecificHeader) when() StepTime {
	return hdr.send.when()
}

func (hdr *sendSpecificHeader) exec(hit Hit) error {
	return hdr.send.exec(hit)
}

// Set sets the header to the specified value
func (hdr *sendSpecificHeader) Set(value string) IStep {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Set(hdr.header, value)
	})
}

// Delete deletes the header
func (hdr *sendSpecificHeader) Delete() IStep {
	return hdr.send.Custom(func(hit Hit) {
		hit.Request().Header.Del(hdr.header)
	})
}

package hit

type sendSpecificHeader struct {
	send   *defaultSend
	header string
}

func newSendSpecificHeader(send *defaultSend, header string) *sendSpecificHeader {
	return &sendSpecificHeader{
		send:   send,
		header: header,
	}
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

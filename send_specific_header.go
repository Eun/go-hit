package hit

type ISendSpecificHeader interface {
	Set(value string) IStep
	Delete() IStep
}

type sendSpecificHeader struct {
	send      ISend
	hit       Hit
	cleanPath CleanPath
	header    string
}

func newSendSpecificHeader(send ISend, hit Hit, cleanPath CleanPath, header string) ISendSpecificHeader {
	return &sendSpecificHeader{
		send:      send,
		hit:       hit,
		cleanPath: cleanPath,
		header:    header,
	}
}

// Set sets the header to the specified value
func (hdr *sendSpecificHeader) Set(value string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Set"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			hit.Request().Header.Set(hdr.header, value)
		},
	})
}

// Delete deletes the header
func (hdr *sendSpecificHeader) Delete() IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Delete"),
		Instance:  hdr.hit,
		Exec: func(hit Hit) {
			hit.Request().Header.Del(hdr.header)
		},
	})
}

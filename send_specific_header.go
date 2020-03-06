package hit

type ISendSpecificHeader interface {
	Set(value string) IStep
	Delete() IStep
}

type sendSpecificHeader struct {
	cleanPath CleanPath
	header    string
}

func newSendSpecificHeader(cleanPath CleanPath, header string) ISendSpecificHeader {
	return &sendSpecificHeader{
		cleanPath: cleanPath,
		header:    header,
	}
}

// Set sets the header to the specified value
func (hdr *sendSpecificHeader) Set(value string) IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Set", []interface{}{value}),
		Exec: func(hit Hit) {
			hit.Request().Header.Set(hdr.header, value)
		},
	})
}

// Delete deletes the header
func (hdr *sendSpecificHeader) Delete() IStep {
	return custom(Step{
		When:      SendStep,
		CleanPath: hdr.cleanPath.Push("Delete", nil),
		Exec: func(hit Hit) {
			hit.Request().Header.Del(hdr.header)
		},
	})
}

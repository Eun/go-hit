package httpbody

import (
	"bufio"
	"io"
	"io/ioutil"
	"strings"
	"unicode"
)

type bodyDecoder struct {
	name       string
	decodeFunc func(body *HTTPBody) (interface{}, error)
}

type bodyDecoders []bodyDecoder

func jsonDecoder() bodyDecoder {
	return bodyDecoder{
		name: "JSON",
		decodeFunc: func(body *HTTPBody) (interface{}, error) {
			reader := body.Reader()
			// if there is a reader
			if reader != nil {
				var container interface{}
				err := json.NewDecoder(reader).Decode(&container)
				return container, err
			}
			return nil, nil
		},
	}
}

// func xmlDecoder() bodyDecoder {
// 	return bodyDecoder{
// 		name: "XML",
// 		decodeFunc: func(body *HTTPBody) interface{} {
// 			reader := body.Reader()
// 			// if there is a reader
// 			if reader != nil {
// 				var container interface{}
// 				if err := xml.NewDecoder(reader).Decode(&container); err == nil {
// 					return container
// 				}
// 			}
// 			return nil
// 		},
// 	}
// }

func stringDecoder() bodyDecoder {
	return bodyDecoder{
		name: "String",
		decodeFunc: func(body *HTTPBody) (interface{}, error) {
			// we read the body and put it into a strings.Builder
			// if we found a rune that is not printable we abort
			bodyReader := body.Reader()
			if bodyReader == nil {
				return nil, nil
			}

			var sb strings.Builder
			r := bufio.NewReader(bodyReader)

			for {
				rn, _, err := r.ReadRune()
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}

				if !unicode.IsPrint(rn) {
					return nil, nil
				}
				if _, err = sb.WriteRune(rn); err != nil {
					return nil, err
				}
			}
			return sb.String(), nil
		},
	}
}

func byteDecoder() bodyDecoder {
	return bodyDecoder{
		name: "Bytes",
		decodeFunc: func(body *HTTPBody) (interface{}, error) {
			r := body.Reader()
			if r == nil {
				return nil, nil
			}
			buf, err := ioutil.ReadAll(r)
			if err != nil {
				return nil, err
			}
			return buf, r.Close()
		},
	}
}

type decoder struct {
	decoders bodyDecoders
}

func (decs bodyDecoders) Has(v bodyDecoder) bool {
	for _, dec := range decs {
		if dec.name == v.name {
			return true
		}
	}
	return false
}

func newDecoder(preferredDecoders ...bodyDecoder) *decoder {
	var fallbackDecoders = []bodyDecoder{jsonDecoder(), stringDecoder(), byteDecoder()}

	decoders := make(bodyDecoders, len(preferredDecoders))
	copy(decoders, preferredDecoders)

	for _, dec := range fallbackDecoders {
		if !decoders.Has(dec) {
			decoders = append(decoders, dec)
		}
	}

	return &decoder{decoders: decoders}
}

func (d *decoder) Run(body *HTTPBody) interface{} {
	for _, dec := range d.decoders {
		if v, err := dec.decodeFunc(body); v != nil && err == nil {
			return v
		}
	}
	return nil
}

package httpbody

import (
	"encoding/json"
	"unicode"
)

type bodyDecoder struct {
	name       string
	decodeFunc func(body *HttpBody) interface{}
}

type bodyDecoders []bodyDecoder

func jsonDecoder() bodyDecoder {
	return bodyDecoder{
		name: "JSON",
		decodeFunc: func(body *HttpBody) interface{} {
			reader := body.Reader()
			// if there is a reader
			if reader != nil {
				var container interface{}
				if err := json.NewDecoder(reader).Decode(&container); err == nil {
					return container
				}
			}
			return nil
		},
	}
}

// func xmlDecoder() bodyDecoder {
// 	return bodyDecoder{
// 		name: "XML",
// 		decodeFunc: func(body *HttpBody) interface{} {
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
		decodeFunc: func(body *HttpBody) interface{} {
			s := body.String()
			for _, r := range s {
				if !unicode.IsPrint(r) {
					return nil
				}
			}
			return s
		},
	}
}

func byteDecoder() bodyDecoder {
	return bodyDecoder{
		name: "Bytes",
		decodeFunc: func(body *HttpBody) interface{} {
			return body.Bytes()
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

func (d *decoder) Run(body *HttpBody) interface{} {
	for _, dec := range d.decoders {
		if v := dec.decodeFunc(body); v != nil {
			return v
		}
	}
	return nil
}

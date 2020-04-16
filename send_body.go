package hit

import (
	"strings"

	"io"

	"fmt"

	"github.com/Eun/go-hit/errortrace"
	"github.com/Eun/go-hit/internal/misc"
	"golang.org/x/xerrors"
)

type ISendBody interface {
	IStep
	// JSON sets the request body to the specified json value.
	//
	// Usage:
	//     Send().Body().JSON(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().JSON(map[string]interface{}{"Name": "Joe"}),
	//     )
	JSON(value interface{}) IStep

	// XML sets the request body to the specified xml value.
	//
	// Usage:
	//     Send().Body().XML([]string{"A", "B"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().XML([]string{"A", "B"}),
	//     )
	XML(value interface{}) IStep

	// Interface sets the request body to the specified json value.
	//
	// Usage:
	//     Send().Body().Interface("Hello World")
	//     Send().Body().Interface(map[string]interface{}{"Name": "Joe"})
	//
	// Example:
	//     MustDo(
	//         Get("https://example.com"),
	//         Send().Body().Interface("Hello World"),
	//     )
	Interface(value interface{}) IStep
}

type sendBody struct {
	cleanPath clearPath
	trace     *errortrace.ErrorTrace
}

func newSendBody(clearPath clearPath, params []interface{}) ISendBody {
	snd := &sendBody{
		cleanPath: clearPath,
		trace:     ett.Prepare(),
	}

	if param, ok := misc.GetLastArgument(params); ok {
		return &finalSendBody{
			&hitStep{
				Trace:     snd.trace,
				When:      SendStep,
				ClearPath: clearPath,
				Exec:      snd.Interface(param).exec,
			},
			"only usable with Send().Body() not with Send().Body(value)",
		}
	}

	return snd
}

func (*sendBody) when() StepTime {
	return SendStep
}

func (body *sendBody) exec(hit Hit) error {
	return body.trace.Format(hit.Description(), "unable to run Send().Body() without an argument or without a chain. Please use Send().Body(something) or Send().Body().Something")
}

func (body *sendBody) clearPath() clearPath {
	return body.cleanPath
}

func (body *sendBody) JSON(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: body.clearPath().Push("JSON", []interface{}{value}),
		Exec: func(hit Hit) error {
			hit.Request().Body().JSON().Set(value)
			return nil
		},
	}
}

func (body *sendBody) XML(value interface{}) IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      SendStep,
		ClearPath: body.clearPath().Push("XML", []interface{}{value}),
		Exec: func(hit Hit) error {
			hit.Request().Body().XML().Set(value)
			return nil
		},
	}
}

func (body *sendBody) Interface(value interface{}) IStep {
	switch x := value.(type) {
	case func(e Hit):
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.clearPath().Push("Interface", []interface{}{value}),
			Exec: func(hit Hit) error {
				x(hit)
				return nil
			},
		}
	case func(e Hit) error:
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.clearPath().Push("Interface", []interface{}{value}),
			Exec:      x,
		}
	default:
		if f := misc.GetGenericFunc(value); f.IsValid() {
			return &hitStep{
				Trace:     ett.Prepare(),
				When:      SendStep,
				ClearPath: body.clearPath().Push("Interface", []interface{}{value}),
				Exec: func(hit Hit) error {
					misc.CallGenericFunc(f)
					return nil
				},
			}
		}
		return &hitStep{
			Trace:     ett.Prepare(),
			When:      SendStep,
			ClearPath: body.clearPath().Push("Interface", []interface{}{value}),
			Exec: func(hit Hit) error {
				switch strings.ToLower(hit.Request().Header.Get("Content-Type")) {
				case "application/json", "text/json":
					hit.Request().Body().JSON().Set(value)
				case "application/xml", "text/xml":
					hit.Request().Body().XML().Set(value)
				default:
					switch v := value.(type) {
					case string:
						hit.Request().Body().SetString(v)
					case []byte:
						hit.Request().Body().SetBytes(v)
					case io.Reader:
						hit.Request().Body().SetReader(v)
					case int:
						hit.Request().Body().SetInt(v)
					case int8:
						hit.Request().Body().SetInt8(v)
					case int16:
						hit.Request().Body().SetInt16(v)
					case int32:
						hit.Request().Body().SetInt32(v)
					case int64:
						hit.Request().Body().SetInt64(v)
					case uint:
						hit.Request().Body().SetUint(v)
					case uint8:
						hit.Request().Body().SetUint8(v)
					case uint16:
						hit.Request().Body().SetUint16(v)
					case uint32:
						hit.Request().Body().SetUint32(v)
					case uint64:
						hit.Request().Body().SetUint64(v)
					case float32:
						hit.Request().Body().SetFloat32(v)
					case float64:
						hit.Request().Body().SetFloat64(v)
					case bool:
						hit.Request().Body().SetBool(v)
					default:
						return fmt.Errorf("unable to set http body to %#v, either specify a Content-Type or use a stringable value", value)
					}
				}
				return nil
			},
		}
	}
}

type finalSendBody struct {
	IStep
	message string
}

func (body *finalSendBody) fail() IStep {
	return &hitStep{
		Trace:     ett.Prepare(),
		When:      CleanStep,
		ClearPath: nil,
		Exec: func(hit Hit) error {
			return xerrors.New(body.message)
		},
	}
}

func (body *finalSendBody) JSON(interface{}) IStep {
	return body.fail()
}

func (body *finalSendBody) XML(interface{}) IStep {
	return body.fail()
}

func (body *finalSendBody) Interface(interface{}) IStep {
	return body.fail()
}

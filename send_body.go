package hit

type sendBody struct {
	*defaultSend
}

func newSendBody(send *defaultSend) *sendBody {
	return &sendBody{
		defaultSend: send,
	}
}

func (body *sendBody) JSON(data interface{}) Step {
	return body.Custom(func(hit Hit) {
		hit.Request().Body().JSON().Set(data)
	})
}

func (body *sendBody) Interface(data interface{}) Step {
	return body.Custom(func(hit Hit) {
		hit.Request().Body().Set(data)
	})
}

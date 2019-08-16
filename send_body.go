package hit

type ISendBody interface {
	IStep
	JSON(data interface{}) IStep
	Interface(data interface{}) IStep
}

type sendBody struct {
	send ISend
}

func newSendBody(send ISend) ISendBody {
	return &sendBody{send}
}

func (body *sendBody) when() StepTime {
	return body.send.when()
}

func (body *sendBody) exec(hit Hit) error {
	return body.send.exec(hit)
}

func (body *sendBody) JSON(data interface{}) IStep {
	return body.send.Custom(func(hit Hit) {
		hit.Request().Body().JSON().Set(data)
	})
}

func (body *sendBody) Interface(data interface{}) IStep {
	return body.send.Custom(func(hit Hit) {
		hit.Request().Body().Set(data)
	})
}

type finalSendBody struct {
	IStep
}

func (d finalSendBody) JSON(data interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}

func (d finalSendBody) Interface(data interface{}) IStep {
	panic("only usable with Send().Body() not with Send().Body(value)")
}

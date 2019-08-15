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

func (snd *sendBody) when() StepTime {
	return snd.send.when()
}

func (snd *sendBody) exec(hit Hit) {
	snd.send.exec(hit)
}

func (snd *sendBody) JSON(data interface{}) IStep {
	return snd.send.Custom(func(hit Hit) {
		hit.Request().Body().JSON().Set(data)
	})
}

func (snd *sendBody) Interface(data interface{}) IStep {
	return snd.send.Custom(func(hit Hit) {
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

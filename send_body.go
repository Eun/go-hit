package hit

type sendBody struct {
	Hit
	send *defaultSend
}

func newSendBody(send *defaultSend) *sendBody {
	return &sendBody{
		Hit:  send.Hit,
		send: send,
	}
}

func (body *sendBody) JSON(data interface{}) Hit {
	return body.send.Custom(func(hit Hit) {
		hit.Request().Body().JSON().Set(data)
	})
}

func (body *sendBody) Interface(data interface{}) Hit {
	return body.send.Custom(func(hit Hit) {
		hit.Request().Body().Set(data)
	})
}

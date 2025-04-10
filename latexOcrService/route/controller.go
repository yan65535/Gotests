package route

import (
	"latexOcrService/adaptor"
	"latexOcrService/api/basic"
	"latexOcrService/config"
)

type Controller struct {
	basic *basic.Ctrl
}

func NewController(adaptors adaptor.Adaptors, conf config.Conf) *Controller {
	return &Controller{}
}

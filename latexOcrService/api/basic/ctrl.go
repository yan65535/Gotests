package basic

import (
	"latexOcrService/adaptor"
	"latexOcrService/config"
	"latexOcrService/service/basic"
)

type Ctrl struct {
	subject *basic.Service
	grade   *basic.Service
	chapter *basic.Service
	stage   *basic.Service
	label   *basic.Service
}

func NewBasicCtrl(adaptors adaptor.Adaptors, conf config.Conf) *Ctrl {
	return &Ctrl{}
}

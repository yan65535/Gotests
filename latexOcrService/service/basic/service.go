package basic

import (
	"latexOcrService/adaptor"
	"latexOcrService/config"
)

type Service struct {
	adaptor adaptor.Adaptors
	conf    config.Conf
}

func NewBasicService(adaptor adaptor.Adaptors, conf config.Conf) *Service {
	return &Service{
		adaptor: adaptor,
		conf:    conf,
	}
}

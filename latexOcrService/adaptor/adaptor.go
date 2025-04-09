package adaptor

import (
	"gorm.io/gorm"
	"latexOcrService/adaptor/repo"
)

type Adaptors struct {
	ILatexOcr *repo.ILatex
	LatexOcr  *repo.Adaptor
}

func NewAdaptors(client *gorm.DB) (*Adaptors, error) {
	return &Adaptors{
		LatexOcr: repo.NewAdaptor(client),
	}, nil
}

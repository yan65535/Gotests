package repo

import "gorm.io/gorm"

type Adaptor struct {
	ILatex
	db *gorm.DB
}

func NewAdaptor(db *gorm.DB) *Adaptor {
	return &Adaptor{
		db: db,
	}
}

package usecase

import "github.com/lovoo/goka"

type Usecase struct {
	gv *goka.View
}

func NewUsecase(gv *goka.View) *Usecase {
	return &Usecase{
		gv,
	}
}

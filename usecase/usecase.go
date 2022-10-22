package usecase

import "github.com/lovoo/goka"

type Usecase struct {
	gv  *goka.View
	gv2 *goka.View
}

func NewUsecase(gv *goka.View, gv2 *goka.View) *Usecase {
	return &Usecase{
		gv,
		gv2,
	}
}

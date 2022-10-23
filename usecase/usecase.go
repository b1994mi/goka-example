package usecase

import "github.com/lovoo/goka"

type Usecase struct {
	aboveThresholdView *goka.View
	balanceView        *goka.View
}

func NewUsecase(aboveThresholdView *goka.View, balanceView *goka.View) *Usecase {
	return &Usecase{
		aboveThresholdView,
		balanceView,
	}
}

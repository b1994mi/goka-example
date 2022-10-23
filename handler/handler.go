package handler

import (
	"goka-example/usecase"

	"github.com/lovoo/goka"
)

type handler struct {
	uc *usecase.Usecase
}

func NewHandler(aboveThresholdView *goka.View, balanceView *goka.View) *handler {
	uc := usecase.NewUsecase(aboveThresholdView, balanceView)

	return &handler{
		uc,
	}
}

package handler

import (
	"goka-example/usecase"

	"github.com/lovoo/goka"
)

type handler struct {
	uc *usecase.Usecase
}

func NewHandler(gv *goka.View) *handler {
	uc := usecase.NewUsecase(gv)

	return &handler{
		uc,
	}
}

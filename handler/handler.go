package handler

import (
	"goka-example/usecase"
)

type handler struct {
	uc *usecase.Usecase
}

func NewHandler() *handler {
	uc := usecase.NewUsecase()

	return &handler{
		uc,
	}
}

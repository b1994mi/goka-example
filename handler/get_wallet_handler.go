package handler

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (h *handler) GetWalletHandler(w http.ResponseWriter, req bunrouter.Request) error {
	resp, err := h.uc.GetWalletUsecase()
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
	}

	bunrouter.JSON(w, resp)
	return nil
}

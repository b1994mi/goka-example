package handler

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"goka-example/request"
)

func (h *handler) GetWalletHandler(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseForm()
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	resp, err := h.uc.GetWalletUsecase(request.GetWallet{
		WalletID: req.Form["wallet_id"][0],
	})
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	bunrouter.JSON(w, resp)
	return nil
}

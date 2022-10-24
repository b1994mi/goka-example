package handler

import (
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/b1994mi/goka-example/request"
)

func (h *handler) GetWalletHandler(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseForm()
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	wIDwID, ok := req.Form["wallet_id"]
	if !ok {
		bunrouter.JSON(w, bunrouter.H{"message": "must include wallet_id in query params"})
		return nil
	}

	withTrxWithTrx, ok := req.Form["with_trx"]
	if !ok {
		withTrxWithTrx = []string{""}
	}

	resp, err := h.uc.GetWalletUsecase(request.GetWallet{
		WalletID: wIDwID[0],
		WithTrx:  withTrxWithTrx[0],
	})
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	bunrouter.JSON(w, resp)
	return nil
}

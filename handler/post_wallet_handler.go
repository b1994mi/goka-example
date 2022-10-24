package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/b1994mi/goka-example/request"
)

func (h *handler) PostWalletHandler(w http.ResponseWriter, req bunrouter.Request) error {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	var reqBody request.PostWallet
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	resp, err := h.uc.PostWalletUsecase(reqBody)
	if err != nil {
		bunrouter.JSON(w, bunrouter.H{"message": err})
		return nil
	}

	bunrouter.JSON(w, resp)
	return nil
}

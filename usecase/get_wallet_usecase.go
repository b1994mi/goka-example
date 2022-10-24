package usecase

import (
	"fmt"
	"time"

	"goka-example/model"
	"goka-example/request"
	"goka-example/response"
)

func (uc *Usecase) GetWalletUsecase(req request.GetWallet) (interface{}, error) {
	aboveThresholdRes, err := uc.aboveThresholdView.Get(req.WalletID)
	if err != nil {
		return nil, err
	}

	wt, ok := aboveThresholdRes.(*model.WalletThreshold)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshall aboveThresholdRes to WalletThreshold")
	}

	balanceRes, err := uc.balanceView.Get(req.WalletID)
	if err != nil {
		return nil, err
	}

	bal, ok := balanceRes.(*model.Wallet)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshall aboveThresholdRes to WalletThreshold")
	}

	res := response.GetWallet{
		WalletID:       req.WalletID,
		Balance:        bal.Balance,
		AboveThreshold: wt.IsAboveThreshold,
	}

	if req.WithTrx == "true" {
		for _, v := range wt.Transactions {
			res.Transactions = append(res.Transactions, &struct {
				Amount float64
				Time   time.Time
			}{
				Amount: v.Amount,
				Time:   v.Time.AsTime(),
			})
		}
	}

	return res, nil
}

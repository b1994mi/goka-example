package usecase

import (
	"goka-example/request"
)

func (uc *Usecase) GetWalletUsecase(req request.GetWallet) (interface{}, error) {
	aboveThresholdRes, err := uc.aboveThresholdView.Get(req.WalletID)
	if err != nil {
		return nil, err
	}

	balanceRes, err := uc.balanceView.Get(req.WalletID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"aboveThreshold": aboveThresholdRes,
		"balance":        balanceRes,
	}, nil
}

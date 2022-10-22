package usecase

func (uc *Usecase) GetWalletUsecase() (interface{}, error) {
	gv1, err := uc.gv.Get("user")
	if err != nil {
		return nil, err
	}

	gv2, err := uc.gv2.Get("user")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"gv":  gv1,
		"gv2": gv2,
	}, nil
}

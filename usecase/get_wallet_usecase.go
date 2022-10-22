package usecase

func (uc *Usecase) GetWalletUsecase() (interface{}, error) {
	return uc.gv.Get("user")
}

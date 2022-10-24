package request

type GetWallet struct {
	WalletID string `uri:"wallet_id"`
	WithTrx  string `uri:"with_trx"`
}

type PostWallet struct {
	WalletID int     `json:"wallet_id"`
	Amount   float64 `json:"amount"`
}

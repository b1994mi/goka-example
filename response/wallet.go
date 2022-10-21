package response

type GetWallet struct {
	WalletID       int     `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`
}

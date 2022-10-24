package response

import "time"

type GetWallet struct {
	WalletID       string  `json:"wallet_id"`
	Balance        float64 `json:"balance"`
	AboveThreshold bool    `json:"above_threshold"`

	Transactions []*struct {
		Amount float64
		Time   time.Time
	} `json:"transactions,omitempty"`
}

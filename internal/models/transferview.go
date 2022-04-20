package models

import (
	"time"
)

type TransferView struct {
	// block fields
	BlockNumber uint64 `json:"block_number"`

	// tx fields
	TxHash      string    `json:"tx_hash"`
	Status      string    `json:"status"`
	Fee         uint64    `json:"fee"`
	CreatedDate time.Time `json:"created_date"`

	// transfer fields
	ID          string `json:"id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      uint64 `json:"amount"`
	Currency    string `json:"currency"`
}

func (TransferView) TableName() string {
	return "transfer_view"
}

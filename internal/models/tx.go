package models

import "time"

type Tx struct {
	Hash            string `gorm:"primaryKey"`
	BlockNumber     uint64
	Status          string
	Fee             uint64
	FeeCurrency     string
	FeePayerAddress string
	CreatedDate     time.Time
	Transfers       []*Transfer `gorm:"foreignKey:TxHash"`
}

func (Tx) TableName() string {
	return "txs"
}

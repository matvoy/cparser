package models

import "time"

type Block struct {
	BlockNumber     uint64 `gorm:"primaryKey"`
	Hash            string
	ProposerAddress string
	CreatedDate     time.Time
	Txs             []*Tx `gorm:"foreignKey:BlockNumber"`
}

func (Block) TableName() string {
	return "blocks"
}

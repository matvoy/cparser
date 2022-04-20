package models

type Transfer struct {
	ID          string `gorm:"primaryKey"`
	TxHash      string
	FromAddress string
	ToAddress   string
	Amount      uint64
	Currency    string
}

func (Transfer) TableName() string {
	return "transfers"
}

package ch

import (
	"github.com/matvoy/cparser/internal/models"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type repository struct {
	db1 *gorm.DB
	db2 *gorm.DB
	db3 *gorm.DB
	db4 *gorm.DB
}

func (r *repository) Cleanup() error {
	if err := r.db1.Exec("ALTER TABLE transfers ON CLUSTER cluster_1 DELETE WHERE 1=1;").Error; err != nil {
		return err
	}
	if err := r.db1.Exec("ALTER TABLE txs ON CLUSTER cluster_1 DELETE WHERE 1=1;").Error; err != nil {
		return err
	}
	if err := r.db1.Exec("ALTER TABLE blocks ON CLUSTER cluster_1 DELETE WHERE 1=1;").Error; err != nil {
		return err
	}
	return nil
}

func NewRepository(dsn1, dsn2, dsn3, dsn4 string) (*repository, error) {
	db1, err := gorm.Open(clickhouse.Open(dsn1), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db2, err := gorm.Open(clickhouse.Open(dsn2), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db3, err := gorm.Open(clickhouse.Open(dsn3), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db4, err := gorm.Open(clickhouse.Open(dsn4), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &repository{db1, db2, db3, db4}, nil
}

func (r *repository) InsertBlockData(block *models.Block) error {
	if err := r.db1.Omit("Txs").Create(block).Error; err != nil {
		return err
	}
	for _, tx := range block.Txs {
		if err := r.db2.Omit("Transfers").Create(tx).Error; err != nil {
			return err
		}
		for _, transfer := range tx.Transfers {
			if err := r.db3.Create(transfer).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *repository) SelectTransfers() ([]*models.TransferView, error) {
	var transfers []*models.TransferView
	return transfers, r.db4.Find(&transfers).Error
}

func (r *repository) Close() {
	db1, _ := r.db1.DB()
	db1.Close()
	db2, _ := r.db2.DB()
	db2.Close()
	db3, _ := r.db3.DB()
	db3.Close()
	db4, _ := r.db4.DB()
	db4.Close()
}

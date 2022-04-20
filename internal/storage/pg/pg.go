package pg

import (
	"github.com/matvoy/cparser/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func (r *repository) Cleanup() error {
	return r.db.Exec("delete from blocks;").Error
}

func NewRepository(dsn string) (*repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &repository{db}, nil
}

func (r *repository) InsertBlockData(block *models.Block) error {
	if block == nil {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(block).Error
}

func (r *repository) SelectTransfers() ([]*models.TransferView, error) {
	var transfers []*models.TransferView
	return transfers, r.db.Find(&transfers).Error
}

func (r *repository) Close() {
	db, _ := r.db.DB()
	db.Close()
}

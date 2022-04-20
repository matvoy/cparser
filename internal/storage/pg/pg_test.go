package pg

import (
	"reflect"
	"testing"

	"github.com/matvoy/cparser/internal/models"
	"gorm.io/gorm"
)

func Test_repository_InsertBlockData(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		block *models.Block
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}
			if err := r.InsertBlockData(tt.args.block); (err != nil) != tt.wantErr {
				t.Errorf("repository.InsertBlockData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_SelectTransfers(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.TransferView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}
			got, err := r.SelectTransfers()
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.SelectTransfers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.SelectTransfers() = %v, want %v", got, tt.want)
			}
		})
	}
}

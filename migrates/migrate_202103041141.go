package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/models"
	"gorm.io/gorm"
)

type Migrate202103041141 struct {}

func (m *Migrate202103041141) GetId() string {
	return "Migrate202103041141"
}

func (m *Migrate202103041141) Migrate() *gormigrate.Gormigrate {
	db := config.GetDbDefault()
	merchant := &models.Merchant{}
	return gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       m.GetId(),
			Migrate: func(tx *gorm.DB) error {
				err := tx.Migrator().AutoMigrate(merchant)
				if err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				err := tx.Migrator().DropColumn(merchant, "ws_room_id")
				if err != nil {
					return err
				}
				return nil
			},
		},
	})
}
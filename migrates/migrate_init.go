package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/models"
	"gorm.io/gorm"
)

type MigrateInit struct {}

func (m *MigrateInit) GetId() string {
	return "init"
}

func (m *MigrateInit) Migrate() *gormigrate.Gormigrate {
	db := config.GetDbDefault()
	mArr := []ModelTableNameInterface{
		&models.Merchant{},
		&models.MchAccount{},
	}
	return gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       "init",
			Migrate: func(tx *gorm.DB) error {
				for _, ma := range mArr {
					err := tx.Migrator().AutoMigrate(ma)
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				for _, ma := range mArr {
					err := tx.Migrator().DropTable(ma.TableName())
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	})
}

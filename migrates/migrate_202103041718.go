package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/zqjzqj/instantCustomer/config"
	"github.com/zqjzqj/instantCustomer/models"
	"gorm.io/gorm"
)

type Migrate202103041718 struct {}

func (m *Migrate202103041718) GetId() string {
	return "Migrate202103041718"
}

func (m *Migrate202103041718) Migrate() *gormigrate.Gormigrate {
	db := config.GetDbDefault()
	v := &models.Visitors{}
	ma := &models.MchAccount{}
	return gormigrate.New(db.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID:       m.GetId(),
			Migrate: func(tx *gorm.DB) error {
				err := tx.Migrator().AutoMigrate(v)
				if err != nil {
					return err
				}
				err = tx.Migrator().AutoMigrate(ma)
				if err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				err := tx.Migrator().DropTable(v.TableName())
				if err != nil {
					return err
				}
				err = tx.Migrator().DropColumn(ma, "session_num")
				if err != nil {
					return err
				}
				return nil
			},
		},
	})
}
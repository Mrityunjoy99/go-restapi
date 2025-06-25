package migration

import (
	"github.com/Mrityunjoy99/sample-go/src/infrastructure/database"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type V20250529214848Employee struct {
	database.BaseModel
	FirstName string `gorm:"column:first_name;not null"`
	LastName  string `gorm:"column:last_name;not null"`
}

func (V20250529214848Employee) TableName() string {
	return "employees"
}

var V20250529214848 *gormigrate.Migration = &gormigrate.Migration{
	ID: "20250529214848_Employee",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&V20250529214848Employee{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("employees")
	},
}

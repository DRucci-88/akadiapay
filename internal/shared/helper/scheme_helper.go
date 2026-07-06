package helper

import (
	"akadia/model"
	"fmt"

	"gorm.io/gorm"
)

func CreateSchemas(db *gorm.DB) error {
	schemas := []string{
		model.SchemaMaster,
		model.SchemaPayment,
	}

	for _, schema := range schemas {
		if err := db.Exec(
			fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema),
		).Error; err != nil {
			return err
		}
	}

	return nil
}

package seeder

import (
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {

	log.Println("========================================")
	log.Println("🌱 Database Seeder")
	log.Println("========================================")

	// if err := seedDepartments(db); err != nil {
	// 	return err
	// }

	log.Println("========================================")
	log.Println("✅ Seeder Completed")
	log.Println("========================================")

	return nil
}

package seed

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {

	log.Println("========================================")
	log.Println("🌱 Database Seeder")
	log.Println("========================================")

	if err := SeedRoles(db); err != nil {
		return err
	}

	if err := SeedTenants(db); err != nil {
		return err
	}

	if err := SeedUsers(db); err != nil {
		return err
	}

	if err := SeedUserTenantRoles(db); err != nil {
		return err
	}

	if err := SeedStudents(db); err != nil {
		return err
	}

	if err := SeedParentStudents(db); err != nil {
		return err
	}

	log.Println("========================================")
	log.Println("✅ Seeder Completed")
	log.Println("========================================")

	return nil
}

func ptr[T any](v T) *T {
	return &v
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ptrDate(year int, month time.Month, day int) *time.Time {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &t
}

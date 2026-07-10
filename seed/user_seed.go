package seed

import (
	"akadia/internal/platform/security"
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUsers(db *gorm.DB) error {
	password, err := security.HashPassword("password")
	if err != nil {
		return err
	}

	users := []model.User{
		{
			Email:       "admin@akadia.id",
			Password:    password,
			DisplayName: "Platform Administrator",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "admin@sman1.id",
			Password:    password,
			DisplayName: "SMAN 1 Administrator",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "treasurer@sman1.id",
			Password:    password,
			DisplayName: "SMAN 1 Treasurer",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "admin@harapan.id",
			Password:    password,
			DisplayName: "Harapan Bangsa Administrator",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "finance@harapan.id",
			Password:    password,
			DisplayName: "Harapan Bangsa Finance Officer",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "budi.parent@gmail.com",
			Password:    password,
			DisplayName: "Budi Santoso",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "asep.parent@gmail.com",
			Password:    password,
			DisplayName: "Asep Budiman",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "kevin@student.id",
			Password:    password,
			DisplayName: "Kevin Wijaya",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "rucco@student.id",
			Password:    password,
			DisplayName: "Rucco Le Amor",
			Status:      model.UserStatusActive,
		},
		{
			Email:       "gilis@student.id",
			Password:    password,
			DisplayName: "Gilis Kilis",
			Status:      model.UserStatusActive,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&users).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Users (%d records)", len(users))

	return nil
}

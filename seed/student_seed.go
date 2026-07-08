package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedStudents(db *gorm.DB) error {
	seed := NewSeeder(db)

	students := []model.Student{
		{
			TenantID: seed.MustTenantByCode("SMAN1").ID,
			UserID:   ptr(seed.MustUserByEmail("kevin@student.id").ID),

			NISN:      "1000000001",
			FullName:  "Kevin Wijaya",
			Gender:    model.StudentGenderMale,
			Phone:     "081234567801",
			Address:   "Jakarta",
			BirthDate: date(2008, 5, 10),
			Status:    model.StudentStatusActive,
		},
		{
			TenantID: seed.MustTenantByCode("SMAN1").ID,
			UserID:   ptr(seed.MustUserByEmail("rucco@student.id").ID),

			NISN:      "1000000002",
			FullName:  "Rucco Le Amor",
			Gender:    model.StudentGenderMale,
			Phone:     "081234567802",
			Address:   "Jakarta",
			BirthDate: date(2007, 8, 21),
			Status:    model.StudentStatusActive,
		},
		{
			TenantID: seed.MustTenantByCode("SMAHB").ID,
			UserID:   ptr(seed.MustUserByEmail("gilis@student.id").ID),

			NISN:      "1000000003",
			FullName:  "Gilis Kilis",
			Gender:    model.StudentGenderFemale,
			Phone:     "081234567803",
			Address:   "Bandung",
			BirthDate: date(2008, 2, 15),
			Status:    model.StudentStatusActive,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&students).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Students (%d records)", len(students))

	return nil
}

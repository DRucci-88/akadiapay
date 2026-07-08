package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedParentStudents(db *gorm.DB) error {
	seed := NewSeeder(db)

	parentStudents := []model.ParentStudent{
		// ======================================
		// Budi Santoso
		// ======================================
		{
			ParentUserID: seed.MustUserByEmail("budi.parent@gmail.com").ID,
			StudentID:    seed.MustStudentByNISN("1000000001").ID,
			Relationship: model.ParentStudentRelationshipFather,
			IsPrimary:    true,
		},
		{
			ParentUserID: seed.MustUserByEmail("budi.parent@gmail.com").ID,
			StudentID:    seed.MustStudentByNISN("1000000002").ID,
			Relationship: model.ParentStudentRelationshipFather,
			IsPrimary:    true,
		},

		// ======================================
		// Asep Budiman
		// ======================================
		{
			ParentUserID: seed.MustUserByEmail("asep.parent@gmail.com").ID,
			StudentID:    seed.MustStudentByNISN("1000000003").ID,
			Relationship: model.ParentStudentRelationshipFather,
			IsPrimary:    true,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&parentStudents).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Parent Students (%d records)", len(parentStudents))

	return nil
}

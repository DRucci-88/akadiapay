package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRoles(db *gorm.DB) error {
	roles := []model.Role{
		{
			Code:        model.RoleCodeSuperAdmin,
			Name:        "Super Administrator",
			Description: ptr("Akadia Platform Administrator"),
		},
		{
			Code:        model.RoleCodeSchoolAdmin,
			Name:        "School Administrator",
			Description: ptr("School Administrator"),
		},
		{
			Code:        model.RoleCodeTreasurer,
			Name:        "Treasurer / Finance Officer",
			Description: ptr("School finance officer"),
		},
		{
			Code:        model.RoleCodeParent,
			Name:        "Parent",
			Description: ptr("Student Parent"),
		},
		{
			Code:        model.RoleCodeStudent,
			Name:        "Student",
			Description: ptr("Student"),
		},
		{
			Code:        model.RoleCodeTeacher,
			Name:        "Teacher",
			Description: ptr("Teacher"),
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&roles).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Roles (%d records)", len(roles))

	return nil
}

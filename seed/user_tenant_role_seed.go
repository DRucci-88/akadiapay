package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedUserTenantRoles(db *gorm.DB) error {
	seed := NewSeeder(db)

	userTenantRoles := []model.UserTenantRole{
		// ============================
		// Platform
		// ============================
		{
			UserID:    seed.MustUserByEmail("admin@akadia.id").ID,
			TenantID:  seed.MustTenantByCode("AKADIA").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeSuperAdmin).ID,
			IsDefault: true,
		},

		// ============================
		// SMAN 1
		// ============================
		{
			UserID:    seed.MustUserByEmail("admin@sman1.sch.id").ID,
			TenantID:  seed.MustTenantByCode("SMAN1").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeSchoolAdmin).ID,
			IsDefault: true,
		},
		{
			UserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			TenantID:  seed.MustTenantByCode("SMAN1").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeParent).ID,
			IsDefault: true,
		},
		{
			UserID:    seed.MustUserByEmail("kevin@student.sch.id").ID,
			TenantID:  seed.MustTenantByCode("SMAN1").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeStudent).ID,
			IsDefault: true,
		},
		{
			UserID:    seed.MustUserByEmail("rucco@student.sch.id").ID,
			TenantID:  seed.MustTenantByCode("SMAN1").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeStudent).ID,
			IsDefault: true,
		},

		// ============================
		// SMA Harapan Bangsa
		// ============================
		{
			UserID:    seed.MustUserByEmail("admin@harapan.sch.id").ID,
			TenantID:  seed.MustTenantByCode("SMAHB").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeSchoolAdmin).ID,
			IsDefault: true,
		},
		{
			UserID:    seed.MustUserByEmail("asep.parent@gmail.com").ID,
			TenantID:  seed.MustTenantByCode("SMAHB").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeParent).ID,
			IsDefault: true,
		},
		{
			UserID:    seed.MustUserByEmail("gilis@student.sch.id").ID,
			TenantID:  seed.MustTenantByCode("SMAHB").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeStudent).ID,
			IsDefault: true,
		},

		// ============================
		// Multi Role Demo
		// ============================
		{
			UserID:    seed.MustUserByEmail("admin@sman1.sch.id").ID,
			TenantID:  seed.MustTenantByCode("AKADIA").ID,
			RoleID:    seed.MustRoleByCode(model.RoleCodeSuperAdmin).ID,
			IsDefault: false,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&userTenantRoles).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed User Tenant Roles (%d records)", len(userTenantRoles))

	return nil
}

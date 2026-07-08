package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedTenants(db *gorm.DB) error {
	tenants := []model.Tenant{
		{
			Code:    "AKADIA",
			Name:    "Akadia Platform",
			Email:   ptr("admin@akadia.id"),
			Phone:   "0210000000",
			Address: "Jakarta",
			Status:  model.TenantStatusActive,
		},
		{
			Code:    "SMAN1",
			Name:    "SMAN 1 Jakarta",
			Email:   ptr("admin@sman1.sch.id"),
			Phone:   "0211111111",
			Address: "Jakarta",
			Status:  model.TenantStatusActive,
		},
		{
			Code:    "SMAHB",
			Name:    "SMA Harapan Bangsa",
			Email:   ptr("admin@harapan.sch.id"),
			Phone:   "0222222222",
			Address: "Bandung",
			Status:  model.TenantStatusTrial,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Create(&tenants).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Tenants (%d records)", len(tenants))

	return nil
}

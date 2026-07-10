package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPaymentProducts(db *gorm.DB) error {
	seed := NewSeeder(db)

	sman1 := seed.MustTenantByCode("SMAN1")
	smahb := seed.MustTenantByCode("SMAHB")

	paymentProducts := []model.PaymentProduct{
		// =====================================================
		// SMAN 1 Jakarta
		// =====================================================
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "FULL_PAYMENT").ID,
			Code:               "SMAN1_SPP_JUL_2026",
			Name:               "SPP July 2026",
			Description:        "Monthly tuition bill for July 2026. Full payment only.",
			RevenueAccountCode: "4101",
			RevenueAccountName: "Tuition Revenue",
			Price:              500000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "FULL_PAYMENT").ID,
			Code:               "SMAN1_SPP_AUG_2026",
			Name:               "SPP August 2026",
			Description:        "Monthly tuition bill for August 2026. Full payment only.",
			RevenueAccountCode: "4101",
			RevenueAccountName: "Tuition Revenue",
			Price:              500000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "PARTIAL_20_MIN_100K").ID,
			Code:               "SMAN1_STUDY_TOUR_2026",
			Name:               "Study Tour 2026",
			Description:        "Study tour bill. Partial payment is allowed with minimum rules.",
			RevenueAccountCode: "4199",
			RevenueAccountName: "Other Education Revenue",
			Price:              800000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "INSTALLMENT_10_MIN_50K").ID,
			Code:               "SMAN1_BUILDING_FEE_2026",
			Name:               "Building Fee 2026",
			Description:        "Building contribution. Installment policy is used for partial-payment testing.",
			RevenueAccountCode: "4103",
			RevenueAccountName: "Building Revenue",
			Price:              1500000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "FULL_PAYMENT").ID,
			Code:               "SMAN1_BOOK_PACKAGE_2026",
			Name:               "Book Package 2026",
			Description:        "Book package. Full payment only.",
			RevenueAccountCode: "4105",
			RevenueAccountName: "Book Revenue",
			Price:              250000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "FULL_PAYMENT").ID,
			Code:               "SMAN1_UNIFORM_2026",
			Name:               "Uniform Package 2026",
			Description:        "Uniform package. Used for cancelled obligation scenario.",
			RevenueAccountCode: "4104",
			RevenueAccountName: "Uniform Revenue",
			Price:              350000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           sman1.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAN1", "FULL_PAYMENT").ID,
			Code:               "SMAN1_OLD_LAB_FEE_2025",
			Name:               "Old Laboratory Fee 2025",
			Description:        "Inactive product for status-filter testing.",
			RevenueAccountCode: "4199",
			RevenueAccountName: "Other Education Revenue",
			Price:              150000,
			Status:             model.PaymentProductStatusInactive,
		},

		// =====================================================
		// SMA Harapan Bangsa
		// =====================================================
		{
			TenantID:           smahb.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAHB", "FULL_PAYMENT").ID,
			Code:               "SMAHB_SPP_JUL_2026",
			Name:               "SPP July 2026",
			Description:        "Monthly tuition bill for July 2026. Full payment only.",
			RevenueAccountCode: "4101",
			RevenueAccountName: "Tuition Revenue",
			Price:              600000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           smahb.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAHB", "PARTIAL_20_MIN_100K").ID,
			Code:               "SMAHB_STUDY_TOUR_2026",
			Name:               "Study Tour 2026",
			Description:        "Study tour bill. Partial payment is allowed with minimum rules.",
			RevenueAccountCode: "4199",
			RevenueAccountName: "Other Education Revenue",
			Price:              900000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           smahb.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAHB", "INSTALLMENT_10_MIN_50K").ID,
			Code:               "SMAHB_BUILDING_FEE_2026",
			Name:               "Building Fee 2026",
			Description:        "Building contribution. Installment policy is used for partial-payment testing.",
			RevenueAccountCode: "4103",
			RevenueAccountName: "Building Revenue",
			Price:              2000000,
			Status:             model.PaymentProductStatusActive,
		},
		{
			TenantID:           smahb.ID,
			PaymentPolicyID:    seed.MustPaymentPolicyByCode("SMAHB", "FULL_PAYMENT").ID,
			Code:               "SMAHB_BOOK_PACKAGE_2026",
			Name:               "Book Package 2026",
			Description:        "Book package. Used for cancelled obligation scenario.",
			RevenueAccountCode: "4105",
			RevenueAccountName: "Book Revenue",
			Price:              300000,
			Status:             model.PaymentProductStatusActive,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"tenant_id",
				"payment_policy_id",
				"name",
				"description",
				"revenue_account_code",
				"revenue_account_name",
				"price",
				"status",
			}),
		}).
		Create(&paymentProducts).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Payment Products (%d records)", len(paymentProducts))
	return nil
}

package seed

import (
	"akadia/model"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPaymentPolicies(db *gorm.DB) error {
	seed := NewSeeder(db)

	schoolTenantCodes := []string{"SMAN1", "SMAHB"}
	paymentPolicies := make([]model.PaymentPolicy, 0, len(schoolTenantCodes)*4)

	for _, tenantCode := range schoolTenantCodes {
		tenant := seed.MustTenantByCode(tenantCode)
		paymentPolicies = append(paymentPolicies,
			model.PaymentPolicy{
				TenantID:            tenant.ID,
				Code:                "FULL_PAYMENT",
				Name:                "Full Payment",
				Description:         "Must be paid fully in one allocation. Used to test full-payment-only policy.",
				AllowPartial:        false,
				MinimumAmount:       0,
				MinimumPercentage:   0,
				AllowOverPayment:    false,
				AutoCloseObligation: true,
			},
			model.PaymentPolicy{
				TenantID:            tenant.ID,
				Code:                "PARTIAL_20_MIN_100K",
				Name:                "Partial Payment - 20 Percent / Minimum 100K",
				Description:         "Partial payment is allowed. Allocation must pass both minimum amount and percentage rules.",
				AllowPartial:        true,
				MinimumAmount:       100000,
				MinimumPercentage:   20,
				AllowOverPayment:    false,
				AutoCloseObligation: true,
			},
			model.PaymentPolicy{
				TenantID:            tenant.ID,
				Code:                "INSTALLMENT_10_MIN_50K",
				Name:                "Installment - 10 Percent / Minimum 50K",
				Description:         "Installment-style payment. Fully paid obligations are not auto-closed to test non-closing policy behavior.",
				AllowPartial:        true,
				MinimumAmount:       50000,
				MinimumPercentage:   10,
				AllowOverPayment:    false,
				AutoCloseObligation: false,
			},
			model.PaymentPolicy{
				TenantID:            tenant.ID,
				Code:                "OVERPAY_ALLOWED_DEMO",
				Name:                "Overpayment Allowed Demo",
				Description:         "Reserved for future student deposit / customer credit testing. Not used by MVP payment order seed.",
				AllowPartial:        true,
				MinimumAmount:       0,
				MinimumPercentage:   0,
				AllowOverPayment:    true,
				AutoCloseObligation: true,
			},
		)
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "code"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"description",
				"allow_partial",
				"minimum_amount",
				"minimum_percentage",
				"allow_over_payment",
				"auto_close_obligation",
			}),
		}).
		Create(&paymentPolicies).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Payment Policies (%d records)", len(paymentPolicies))
	return nil
}

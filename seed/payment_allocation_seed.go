package seed

import (
	"akadia/model"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPaymentAllocations(db *gorm.DB) error {
	seed := NewSeeder(db)

	julyPeriod := date(2026, time.July, 1)
	annualPeriod := date(2026, time.January, 1)

	paymentAllocations := []model.PaymentAllocation{
		// PAY-SEED-SMAN1-0001: Kevin, multi-allocation, multi-revenue ledger.
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAN1-0001").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000001", "SMAN1_SPP_JUL_2026", julyPeriod).ID,
			AllocatedAmount:     500000,
		},
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAN1-0001").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000001", "SMAN1_BOOK_PACKAGE_2026", annualPeriod).ID,
			AllocatedAmount:     250000,
		},

		// PAY-SEED-SMAN1-0002: Kevin, partial payment with minimum policy satisfied.
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAN1-0002").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000001", "SMAN1_STUDY_TOUR_2026", annualPeriod).ID,
			AllocatedAmount:     300000,
		},

		// PAY-SEED-SMAN1-0005: Rucco, installment partial payment.
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAN1-0005").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000002", "SMAN1_BUILDING_FEE_2026", annualPeriod).ID,
			AllocatedAmount:     500000,
		},

		// PAY-SEED-SMAHB-0001: Gilis, QRIS full payment.
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAHB-0001").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000003", "SMAHB_SPP_JUL_2026", julyPeriod).ID,
			AllocatedAmount:     600000,
		},

		// PAY-SEED-SMAHB-0002: Gilis, Virtual Account partial payment.
		{
			PaymentOrderID:      seed.MustPaymentOrderByNumber("PAY-SEED-SMAHB-0002").ID,
			StudentObligationID: seed.MustStudentObligationByStudentProductPeriod("1000000003", "SMAHB_BUILDING_FEE_2026", annualPeriod).ID,
			AllocatedAmount:     500000,
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "payment_order_id"},
				{Name: "student_obligation_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"allocated_amount"}),
		}).
		Create(&paymentAllocations).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Payment Allocations (%d records)", len(paymentAllocations))
	return nil
}

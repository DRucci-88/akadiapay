package seed

import (
	"akadia/model"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedStudentObligations(db *gorm.DB) error {
	seed := NewSeeder(db)

	issuedAt := date(2026, time.June, 20)
	julyPeriod := date(2026, time.July, 1)
	augustPeriod := date(2026, time.August, 1)
	annualPeriod := date(2026, time.January, 1)

	studentObligations := []model.StudentObligation{
		// =====================================================
		// Kevin Wijaya - SMAN 1
		// Covers: closed full payment, open full payment, closed book package,
		// partial payment, and cancelled bill.
		// =====================================================
		{
			StudentID:         seed.MustStudentByNISN("1000000001").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_SPP_JUL_2026").ID,
			Period:            julyPeriod,
			OriginalAmount:    500000,
			OutstandingAmount: 0,
			DueDate:           date(2026, time.July, 10),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusClosed,
			Notes:             "Seed: fully paid through PAY-SEED-SMAN1-0001.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000001").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_SPP_AUG_2026").ID,
			Period:            augustPeriod,
			OriginalAmount:    500000,
			OutstandingAmount: 500000,
			DueDate:           date(2026, time.August, 10),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPending,
			Notes:             "Seed: unpaid full-payment-only obligation.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000001").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_BOOK_PACKAGE_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    250000,
			OutstandingAmount: 0,
			DueDate:           date(2026, time.July, 15),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusClosed,
			Notes:             "Seed: fully paid together with SPP July in one payment order.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000001").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_STUDY_TOUR_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    800000,
			OutstandingAmount: 500000,
			DueDate:           date(2026, time.September, 30),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPartial,
			Notes:             "Seed: partial payment of 300000 through PAY-SEED-SMAN1-0002.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000001").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_UNIFORM_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    350000,
			OutstandingAmount: 350000,
			DueDate:           date(2026, time.July, 25),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusCancelled,
			Notes:             "Seed: cancelled obligation for cancellation/filter edge case.",
		},

		// =====================================================
		// Rucco Le Amor - SMAN 1
		// Covers: installment partial, unpaid full payment, unpaid partial, unpaid book.
		// =====================================================
		{
			StudentID:         seed.MustStudentByNISN("1000000002").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_BUILDING_FEE_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    1500000,
			OutstandingAmount: 1000000,
			DueDate:           date(2026, time.December, 31),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPartial,
			Notes:             "Seed: installment payment of 500000 through PAY-SEED-SMAN1-0005.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000002").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_SPP_JUL_2026").ID,
			Period:            julyPeriod,
			OriginalAmount:    500000,
			OutstandingAmount: 500000,
			DueDate:           date(2026, time.July, 10),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPending,
			Notes:             "Seed: unpaid SPP for full-payment validation.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000002").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_STUDY_TOUR_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    800000,
			OutstandingAmount: 800000,
			DueDate:           date(2026, time.September, 30),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPending,
			Notes:             "Seed: unpaid partial-payment obligation. Minimum valid allocation is 160000.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000002").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAN1_BOOK_PACKAGE_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    250000,
			OutstandingAmount: 250000,
			DueDate:           date(2026, time.July, 15),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPending,
			Notes:             "Seed: unpaid book package for multi-allocation testing.",
		},

		// =====================================================
		// Gilis Kilis - SMA Harapan Bangsa
		// Covers: cross-tenant data isolation, QRIS payment, VA payment,
		// and cancelled bill.
		// =====================================================
		{
			StudentID:         seed.MustStudentByNISN("1000000003").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAHB_SPP_JUL_2026").ID,
			Period:            julyPeriod,
			OriginalAmount:    600000,
			OutstandingAmount: 0,
			DueDate:           date(2026, time.July, 10),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusClosed,
			Notes:             "Seed: fully paid by QRIS through PAY-SEED-SMAHB-0001.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000003").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAHB_BUILDING_FEE_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    2000000,
			OutstandingAmount: 1500000,
			DueDate:           date(2026, time.December, 31),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPartial,
			Notes:             "Seed: partial VA payment of 500000 through PAY-SEED-SMAHB-0002.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000003").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAHB_STUDY_TOUR_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    900000,
			OutstandingAmount: 900000,
			DueDate:           date(2026, time.September, 30),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusPending,
			Notes:             "Seed: unpaid SMAHB partial-payment obligation for tenant isolation testing.",
		},
		{
			StudentID:         seed.MustStudentByNISN("1000000003").ID,
			PaymentProductID:  seed.MustPaymentProductByCode("SMAHB_BOOK_PACKAGE_2026").ID,
			Period:            annualPeriod,
			OriginalAmount:    300000,
			OutstandingAmount: 300000,
			DueDate:           date(2026, time.July, 15),
			IssuedAt:          issuedAt,
			Status:            model.StudentObligationStatusCancelled,
			Notes:             "Seed: cancelled SMAHB obligation.",
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "student_id"},
				{Name: "payment_product_id"},
				{Name: "period"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"original_amount",
				"outstanding_amount",
				"due_date",
				"issued_at",
				"status",
				"notes",
			}),
		}).
		Create(&studentObligations).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Student Obligations (%d records)", len(studentObligations))
	return nil
}

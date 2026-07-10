package seed

import (
	"akadia/model"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPaymentOrders(db *gorm.DB) error {
	seed := NewSeeder(db)

	ledgerPostedAt := date(2026, time.July, 5)

	paymentOrders := []model.PaymentOrder{
		// =====================================================
		// SMAN 1 Jakarta - Budi pays for Kevin and Rucco
		// =====================================================
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000001").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0001",
			OrderDate:       date(2026, time.July, 5),
			TotalAmount:     750000,
			Status:          model.PaymentOrderStatusCompleted,
			PaymentMethod:   model.PaymentMethodBankTransfer,
			ReferenceNumber: ptr("TRF-SMAN1-SEED-0001"),
			LedgerPostedAt:  &ledgerPostedAt,
			Notes:           "Seed: completed multi-allocation payment for Kevin SPP July and Book Package.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000001").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0002",
			OrderDate:       date(2026, time.July, 8),
			TotalAmount:     300000,
			Status:          model.PaymentOrderStatusCompleted,
			PaymentMethod:   model.PaymentMethodCash,
			ReferenceNumber: ptr("CASH-SMAN1-SEED-0002"),
			LedgerPostedAt:  &ledgerPostedAt,
			Notes:           "Seed: completed partial payment for Kevin Study Tour.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000001").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0003",
			OrderDate:       date(2026, time.July, 12),
			TotalAmount:     500000,
			Status:          model.PaymentOrderStatusPending,
			PaymentMethod:   model.PaymentMethodVirtualAccount,
			ReferenceNumber: ptr("VA-SMAN1-SEED-0003"),
			Notes:           "Seed: pending order with no allocation yet.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000001").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0004",
			OrderDate:       date(2026, time.July, 13),
			TotalAmount:     350000,
			Status:          model.PaymentOrderStatusCancelled,
			PaymentMethod:   model.PaymentMethodQRIS,
			ReferenceNumber: ptr("QRIS-SMAN1-SEED-0004"),
			Notes:           "Seed: cancelled order for status and cancellation testing.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000002").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0005",
			OrderDate:       date(2026, time.July, 9),
			TotalAmount:     500000,
			Status:          model.PaymentOrderStatusCompleted,
			PaymentMethod:   model.PaymentMethodBankTransfer,
			ReferenceNumber: ptr("TRF-SMAN1-SEED-0005"),
			LedgerPostedAt:  &ledgerPostedAt,
			Notes:           "Seed: completed installment payment for Rucco Building Fee.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAN1").ID,
			StudentID:       seed.MustStudentByNISN("1000000002").ID,
			PaidByUserID:    seed.MustUserByEmail("budi.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAN1-0006",
			OrderDate:       date(2026, time.July, 14),
			TotalAmount:     250000,
			Status:          model.PaymentOrderStatusPending,
			PaymentMethod:   model.PaymentMethodCash,
			ReferenceNumber: ptr("CASH-SMAN1-SEED-0006"),
			Notes:           "Seed: pending cash order for Rucco.",
		},

		// =====================================================
		// SMA Harapan Bangsa - Asep pays for Gilis
		// =====================================================
		{
			TenantID:        seed.MustTenantByCode("SMAHB").ID,
			StudentID:       seed.MustStudentByNISN("1000000003").ID,
			PaidByUserID:    seed.MustUserByEmail("asep.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAHB-0001",
			OrderDate:       date(2026, time.July, 6),
			TotalAmount:     600000,
			Status:          model.PaymentOrderStatusCompleted,
			PaymentMethod:   model.PaymentMethodQRIS,
			ReferenceNumber: ptr("QRIS-SMAHB-SEED-0001"),
			LedgerPostedAt:  &ledgerPostedAt,
			Notes:           "Seed: completed QRIS payment for Gilis SPP July.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAHB").ID,
			StudentID:       seed.MustStudentByNISN("1000000003").ID,
			PaidByUserID:    seed.MustUserByEmail("asep.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAHB-0002",
			OrderDate:       date(2026, time.July, 10),
			TotalAmount:     500000,
			Status:          model.PaymentOrderStatusCompleted,
			PaymentMethod:   model.PaymentMethodVirtualAccount,
			ReferenceNumber: ptr("VA-SMAHB-SEED-0002"),
			LedgerPostedAt:  &ledgerPostedAt,
			Notes:           "Seed: completed VA payment for Gilis Building Fee.",
		},
		{
			TenantID:        seed.MustTenantByCode("SMAHB").ID,
			StudentID:       seed.MustStudentByNISN("1000000003").ID,
			PaidByUserID:    seed.MustUserByEmail("asep.parent@gmail.com").ID,
			OrderNumber:     "PAY-SEED-SMAHB-0003",
			OrderDate:       date(2026, time.July, 15),
			TotalAmount:     900000,
			Status:          model.PaymentOrderStatusExpired,
			PaymentMethod:   model.PaymentMethodVirtualAccount,
			ReferenceNumber: ptr("VA-SMAHB-SEED-0003"),
			Notes:           "Seed: expired order with no allocation.",
		},
	}

	if err := db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "order_number"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"tenant_id",
				"student_id",
				"paid_by_user_id",
				"order_date",
				"total_amount",
				"status",
				"payment_method",
				"reference_number",
				"ledger_posted_at",
				"notes",
			}),
		}).
		Create(&paymentOrders).Error; err != nil {
		return err
	}

	log.Printf("✓ Seed Payment Orders (%d records)", len(paymentOrders))
	return nil
}

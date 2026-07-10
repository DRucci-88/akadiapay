package seed

import (
	"akadia/model"
	"log"
	"time"

	"gorm.io/gorm"
)

type seedLedgerLine struct {
	AccountCode string
	AccountName string
	Debit       float64
	Credit      float64
}

func SeedLedgerEntries(db *gorm.DB) error {
	seed := NewSeeder(db)

	seededOrderCount := 0
	seededEntryCount := 0

	ledgerSeeds := map[string][]seedLedgerLine{
		"PAY-SEED-SMAN1-0001": {
			{AccountCode: "1102", AccountName: "Bank", Debit: 750000},
			{AccountCode: "4101", AccountName: "Tuition Revenue", Credit: 500000},
			{AccountCode: "4105", AccountName: "Book Revenue", Credit: 250000},
		},
		"PAY-SEED-SMAN1-0002": {
			{AccountCode: "1101", AccountName: "Cash", Debit: 300000},
			{AccountCode: "4199", AccountName: "Other Education Revenue", Credit: 300000},
		},
		"PAY-SEED-SMAN1-0005": {
			{AccountCode: "1102", AccountName: "Bank", Debit: 500000},
			{AccountCode: "4103", AccountName: "Building Revenue", Credit: 500000},
		},
		"PAY-SEED-SMAHB-0001": {
			{AccountCode: "1103", AccountName: "QRIS Clearing", Debit: 600000},
			{AccountCode: "4101", AccountName: "Tuition Revenue", Credit: 600000},
		},
		"PAY-SEED-SMAHB-0002": {
			{AccountCode: "1104", AccountName: "Virtual Account Clearing", Debit: 500000},
			{AccountCode: "4103", AccountName: "Building Revenue", Credit: 500000},
		},
	}

	for orderNumber, lines := range ledgerSeeds {
		paymentOrder := seed.MustPaymentOrderByNumber(orderNumber)

		var existingCount int64
		if err := db.
			Model(&model.LedgerEntry{}).
			Where("payment_order_id = ?", paymentOrder.ID).
			Count(&existingCount).Error; err != nil {
			return err
		}
		if existingCount > 0 {
			continue
		}

		entries := make([]model.LedgerEntry, 0, len(lines))
		for _, line := range lines {
			entries = append(entries, model.LedgerEntry{
				PaymentOrderID: paymentOrder.ID,
				EntryDate:      normalizeSeedLedgerDate(paymentOrder.OrderDate),
				AccountCode:    line.AccountCode,
				AccountName:    line.AccountName,
				Debit:          line.Debit,
				Credit:         line.Credit,
				Description:    "Seed ledger posting for " + paymentOrder.OrderNumber,
			})
		}

		if err := db.Create(&entries).Error; err != nil {
			return err
		}

		seededOrderCount++
		seededEntryCount += len(entries)
	}

	log.Printf("✓ Seed Ledger Entries (%d entries for %d completed orders)", seededEntryCount, seededOrderCount)
	return nil
}

func normalizeSeedLedgerDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

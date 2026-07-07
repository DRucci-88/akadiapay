package model

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	PaymentOrderID uuid.UUID     `gorm:"type:uuid;not null;index:idx_ledger_entry_order"`
	PaymentOrder   *PaymentOrder `gorm:"foreignKey:PaymentOrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	EntryDate   time.Time `gorm:"not null"`
	AccountCode string    `gorm:"type:varchar(30);not null"`
	AccountName string    `gorm:"type:varchar(150);not null"`
	Debit       float64   `gorm:"type:numeric(18,2);default:0;not null"`
	Credit      float64   `gorm:"type:numeric(18,2);default:0;not null"`
	Description string    `gorm:"type:text"`

	BaseModel
}

func (LedgerEntry) TableName() string {
	return SchemaPayment + ".ledger_entries"
}

/*
DEBIT (Menambah tagihan siswa), CREDIT (Mengurangi tagihan siswa)
*/

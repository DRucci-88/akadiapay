package domain

import (
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LedgerEntryHandler interface {
	FindAll(c *gin.Context)
	FindByPaymentOrderID(c *gin.Context)
	PostPayment(c *gin.Context)
}

type LedgerEntryService interface {
	PostPayment(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
	) error
	FindByPaymentOrderID(
		ctx context.Context,
		authContext *security.AuthContext,
		paymentOrderID uuid.UUID,
	) ([]LedgerEntryResponse, error)
	FindPaginate(
		ctx context.Context,
		authContext *security.AuthContext,
		pageable *shared.Pageable,
		filter *LedgerEntryFilter,
	) (*shared.Page[LedgerEntryResponse], error)
}

type LedgerEntryRepository interface {
	FindByPaymentOrderID(
		ctx context.Context,
		paymentOrderID uuid.UUID,
		tenantID uuid.UUID,
	) ([]model.LedgerEntry, error)
	ExistsByPaymentOrderID(
		ctx context.Context,
		paymentOrderID uuid.UUID,
		tenantID uuid.UUID,
	) (bool, error)
	CreateInBatches(
		ctx context.Context,
		entries []model.LedgerEntry,
		batchSize int,
	) error
	FindPaginate(
		ctx context.Context,
		tenantID uuid.UUID,
		pageable *shared.Pageable,
		filter *LedgerEntryFilter,
	) (*shared.Page[model.LedgerEntry], error)
}

type LedgerEntryResponse struct {
	ID             uuid.UUID `json:"id"`
	PaymentOrderID uuid.UUID `json:"payment_order_id"`
	EntryDate      time.Time `json:"entry_date"`
	AccountCode    string    `json:"account_code"`
	AccountName    string    `json:"account_name"`
	Debit          float64   `json:"debit"`
	Credit         float64   `json:"credit"`
	Description    string    `json:"description"`
}

type LedgerEntryFilter struct {
	PaymentOrderID *uuid.UUID `form:"payment_order_id"`
	AccountCode    *string    `form:"account_code"`
	EntryDateFrom  *time.Time `form:"entry_date_from" time_format:"2006-01-02"`
	EntryDateTo    *time.Time `form:"entry_date_to" time_format:"2006-01-02"`
	Keyword        *string    `form:"keyword"`
}

func NewLedgerEntryResponse(model *model.LedgerEntry) *LedgerEntryResponse {
	return &LedgerEntryResponse{
		ID:             model.ID,
		PaymentOrderID: model.PaymentOrderID,
		EntryDate:      model.EntryDate,
		AccountCode:    model.AccountCode,
		AccountName:    model.AccountName,
		Debit:          model.Debit,
		Credit:         model.Credit,
		Description:    model.Description,
	}
}

func NewLedgerEntryResponses(models []model.LedgerEntry) []LedgerEntryResponse {
	resList := make([]LedgerEntryResponse, 0, len(models))
	for i := range models {
		resList = append(resList, *NewLedgerEntryResponse(&models[i]))
	}
	return resList
}

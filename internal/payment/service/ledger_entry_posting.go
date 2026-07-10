package service

import (
	"akadia/domain"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
)

type ledgerAccount struct {
	Code string
	Name string
}

func postLedgerEntriesForPaymentOrder(
	ctx context.Context,
	repo domain.RepositoryManagerPayment,
	tenantID uuid.UUID,
	paymentOrderID uuid.UUID,
) error {
	paymentOrder, err := repo.PaymentOrder().LockByID(ctx, paymentOrderID, tenantID)
	if err != nil {
		return err
	}
	if paymentOrder.Status != model.PaymentOrderStatusCompleted {
		return shared.ErrPaymentOrderStatusInvalid
	}
	if paymentOrder.LedgerPostedAt != nil {
		return nil
	}

	exists, err := repo.LedgerEntry().ExistsByPaymentOrderID(ctx, paymentOrderID, tenantID)
	if err != nil {
		return err
	}
	if exists {
		now := time.Now().UTC()
		_, err := repo.PaymentOrder().MarkLedgerPosted(ctx, paymentOrderID, &now)
		return err
	}

	allocations, err := repo.PaymentAllocation().FindByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return err
	}
	if len(allocations) == 0 {
		return shared.ErrPaymentAllocationNotFound
	}

	totalAllocated := 0.0
	obligationIDs := make([]uuid.UUID, 0, len(allocations))
	for _, allocation := range allocations {
		if !shared.FloatGreater(allocation.AllocatedAmount, 0) {
			return shared.ErrPaymentAllocationAmountInvalid
		}
		totalAllocated += allocation.AllocatedAmount
		obligationIDs = append(obligationIDs, allocation.StudentObligationID)
	}
	if !shared.FloatEqual(totalAllocated, paymentOrder.TotalAmount) {
		return shared.ErrLedgerUnbalancedSource
	}

	obligationIDs = uniqueUUIDs(obligationIDs)
	obligations, err := repo.StudentObligation().FindByIDs(ctx, obligationIDs)
	if err != nil {
		return err
	}
	if len(obligations) != len(obligationIDs) {
		return shared.ErrStudentObligationNotFound
	}

	productIDs := make([]uuid.UUID, 0, len(obligations))
	for _, obligation := range obligations {
		productIDs = append(productIDs, obligation.PaymentProductID)
	}

	productIDs = uniqueUUIDs(productIDs)
	paymentProducts, err := repo.PaymentProduct().FindByIDsIncludingDeleted(ctx, productIDs)
	if err != nil {
		return err
	}
	if len(paymentProducts) != len(productIDs) {
		return shared.ErrPaymentProductNotFound
	}

	entries, err := buildLedgerEntries(paymentOrder, allocations, obligations, paymentProducts)
	if err != nil {
		return err
	}

	if err := repo.LedgerEntry().CreateInBatches(ctx, entries, 100); err != nil {
		return err
	}

	now := time.Now().UTC()
	_, err = repo.PaymentOrder().MarkLedgerPosted(ctx, paymentOrderID, &now)
	return err
}

func buildLedgerEntries(
	paymentOrder *model.PaymentOrder,
	allocations []model.PaymentAllocation,
	obligations []model.StudentObligation,
	paymentProducts []model.PaymentProduct,
) ([]model.LedgerEntry, error) {
	debitAccount, err := resolveLedgerDebitAccount(paymentOrder.PaymentMethod)
	if err != nil {
		return nil, err
	}

	obligationMap := make(map[uuid.UUID]model.StudentObligation, len(obligations))
	for _, obligation := range obligations {
		obligationMap[obligation.ID] = obligation
	}

	paymentProductMap := make(map[uuid.UUID]model.PaymentProduct, len(paymentProducts))
	for _, paymentProduct := range paymentProducts {
		paymentProductMap[paymentProduct.ID] = paymentProduct
	}

	entryDate := normalizeLedgerEntryDate(paymentOrder.OrderDate)
	debitTotal := 0.0
	creditTotals := make(map[string]float64)
	creditNames := make(map[string]string)

	for _, allocation := range allocations {
		obligation, exists := obligationMap[allocation.StudentObligationID]
		if !exists {
			return nil, shared.ErrStudentObligationNotFound
		}

		paymentProduct, exists := paymentProductMap[obligation.PaymentProductID]
		if !exists {
			return nil, shared.ErrPaymentProductNotFound
		}

		creditAccount, err := resolvePaymentProductRevenueAccount(paymentProduct)
		if err != nil {
			return nil, err
		}
		creditTotals[creditAccount.Code] += allocation.AllocatedAmount
		creditNames[creditAccount.Code] = creditAccount.Name
		debitTotal += allocation.AllocatedAmount
	}

	entries := make([]model.LedgerEntry, 0, len(creditTotals)+1)
	entries = append(entries, model.LedgerEntry{
		PaymentOrderID: paymentOrder.ID,
		EntryDate:      entryDate,
		AccountCode:    debitAccount.Code,
		AccountName:    debitAccount.Name,
		Debit:          debitTotal,
		Description:    fmt.Sprintf("Payment receipt %s", paymentOrder.OrderNumber),
	})

	accountCodes := make([]string, 0, len(creditTotals))
	for accountCode := range creditTotals {
		accountCodes = append(accountCodes, accountCode)
	}
	sort.Strings(accountCodes)

	for _, accountCode := range accountCodes {
		entries = append(entries, model.LedgerEntry{
			PaymentOrderID: paymentOrder.ID,
			EntryDate:      entryDate,
			AccountCode:    accountCode,
			AccountName:    creditNames[accountCode],
			Credit:         creditTotals[accountCode],
			Description:    fmt.Sprintf("Revenue posting %s", paymentOrder.OrderNumber),
		})
	}

	if err := validateLedgerEntries(entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func validateLedgerEntries(entries []model.LedgerEntry) error {
	if len(entries) == 0 {
		return shared.ErrLedgerUnbalanced
	}

	totalDebit := 0.0
	totalCredit := 0.0

	for _, entry := range entries {
		hasDebit := shared.FloatGreater(entry.Debit, 0)
		hasCredit := shared.FloatGreater(entry.Credit, 0)

		if hasDebit == hasCredit {
			return shared.ErrLedgerUnbalanced
		}

		totalDebit += entry.Debit
		totalCredit += entry.Credit
	}

	if !shared.FloatEqual(totalDebit, totalCredit) {
		return shared.ErrLedgerUnbalanced
	}

	return nil
}

func resolveLedgerDebitAccount(
	paymentMethod model.PaymentOrderPaymentMethod,
) (*ledgerAccount, error) {
	switch paymentMethod {
	case model.PaymentMethodCash:
		return &ledgerAccount{Code: "1101", Name: "Cash"}, nil
	case model.PaymentMethodBankTransfer:
		return &ledgerAccount{Code: "1102", Name: "Bank"}, nil
	case model.PaymentMethodQRIS:
		return &ledgerAccount{Code: "1103", Name: "QRIS Clearing"}, nil
	case model.PaymentMethodVirtualAccount:
		return &ledgerAccount{Code: "1104", Name: "Virtual Account Clearing"}, nil
	case model.PaymentMethodCreditCard:
		return &ledgerAccount{Code: "1102", Name: "Bank"}, nil
	default:
		return nil, shared.ErrLedgerDebitAccountNotConfigured
	}
}

func normalizeLedgerEntryDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

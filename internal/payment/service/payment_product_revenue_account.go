package service

import (
	"akadia/internal/shared"
	"akadia/model"
	"log"
	"strings"
)

var paymentProductRevenueAccountNames = map[string]string{
	"4101": "Tuition Revenue",
	"4102": "Registration Revenue",
	"4103": "Building Revenue",
	"4104": "Uniform Revenue",
	"4105": "Book Revenue",
	"4199": "Other Education Revenue",
}

var paymentProductRevenueAccountCodesByName = map[string]string{
	"TUITION REVENUE":         "4101",
	"REGISTRATION REVENUE":    "4102",
	"BUILDING REVENUE":        "4103",
	"UNIFORM REVENUE":         "4104",
	"BOOK REVENUE":            "4105",
	"OTHER EDUCATION REVENUE": "4199",
}

func normalizePaymentProductRevenueAccount(paymentProduct *model.PaymentProduct) {
	revenueAccountCode, revenueAccountName := resolvePaymentProductRevenueAccountFields(*paymentProduct)
	paymentProduct.RevenueAccountCode = revenueAccountCode
	paymentProduct.RevenueAccountName = revenueAccountName
}

func resolvePaymentProductRevenueAccount(
	paymentProduct model.PaymentProduct,
) (*ledgerAccount, error) {
	revenueAccountCode := strings.TrimSpace(paymentProduct.RevenueAccountCode)
	revenueAccountName := strings.TrimSpace(paymentProduct.RevenueAccountName)

	if revenueAccountCode == "" || revenueAccountName == "" {
		log.Printf(
			"ledger credit mapping fallback payment_product_id=%s code=%s name=%s",
			paymentProduct.ID,
			paymentProduct.Code,
			paymentProduct.Name,
		)
	}

	revenueAccountCode, revenueAccountName = resolvePaymentProductRevenueAccountFields(paymentProduct)
	if revenueAccountCode == "" || revenueAccountName == "" {
		return nil, shared.ErrLedgerCreditAccountNotConfigured
	}

	return &ledgerAccount{
		Code: revenueAccountCode,
		Name: revenueAccountName,
	}, nil
}

func resolvePaymentProductRevenueAccountFields(
	paymentProduct model.PaymentProduct,
) (string, string) {
	revenueAccountCode := strings.TrimSpace(paymentProduct.RevenueAccountCode)
	revenueAccountName := strings.TrimSpace(paymentProduct.RevenueAccountName)

	if revenueAccountCode != "" && revenueAccountName == "" {
		if resolvedName, ok := paymentProductRevenueAccountNames[revenueAccountCode]; ok {
			revenueAccountName = resolvedName
		}
	}
	if revenueAccountName != "" && revenueAccountCode == "" {
		if resolvedCode, ok := paymentProductRevenueAccountCodesByName[strings.ToUpper(revenueAccountName)]; ok {
			revenueAccountCode = resolvedCode
		}
	}
	if revenueAccountCode != "" && revenueAccountName != "" {
		return revenueAccountCode, revenueAccountName
	}

	return resolveDefaultPaymentProductRevenueAccount(paymentProduct.Code, paymentProduct.Name)
}

func resolveDefaultPaymentProductRevenueAccount(
	productCode string,
	productName string,
) (string, string) {
	normalized := strings.ToUpper(strings.TrimSpace(productCode + " " + productName))

	switch {
	case strings.Contains(normalized, "SPP"), strings.Contains(normalized, "TUITION"):
		return "4101", "Tuition Revenue"
	case strings.Contains(normalized, "REGISTRATION"), strings.Contains(normalized, "DAFTAR"):
		return "4102", "Registration Revenue"
	case strings.Contains(normalized, "BUILDING"), strings.Contains(normalized, "GEDUNG"):
		return "4103", "Building Revenue"
	case strings.Contains(normalized, "UNIFORM"), strings.Contains(normalized, "SERAGAM"):
		return "4104", "Uniform Revenue"
	case strings.Contains(normalized, "BOOK"), strings.Contains(normalized, "BUKU"):
		return "4105", "Book Revenue"
	default:
		return "4199", "Other Education Revenue"
	}
}

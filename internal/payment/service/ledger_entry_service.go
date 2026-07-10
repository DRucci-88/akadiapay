package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type ledgerEntryServiceImpl struct {
	repo                domain.RepositoryManagerPayment
	ledgerEntryRepo     domain.LedgerEntryRepository
	paymentOrderService domain.PaymentOrderService
}

func NewLedgerEntryService(
	repo domain.RepositoryManagerPayment,
	paymentOrderService domain.PaymentOrderService,
) domain.LedgerEntryService {
	return &ledgerEntryServiceImpl{
		repo:                repo,
		ledgerEntryRepo:     repo.LedgerEntry(),
		paymentOrderService: paymentOrderService,
	}
}

func (s *ledgerEntryServiceImpl) PostPayment(
	ctx context.Context,
	authContext *security.AuthContext,
	paymentOrderID uuid.UUID,
) error {
	if err := validateLedgerAccessRole(authContext); err != nil {
		return err
	}

	if _, err := s.paymentOrderService.FindByID(ctx, authContext, paymentOrderID); err != nil {
		return err
	}

	return s.repo.Transaction(ctx, func(repo domain.RepositoryManagerPayment) error {
		return postLedgerEntriesForPaymentOrder(ctx, repo, authContext.TenantID, paymentOrderID)
	})
}

func (s *ledgerEntryServiceImpl) FindByPaymentOrderID(
	ctx context.Context,
	authContext *security.AuthContext,
	paymentOrderID uuid.UUID,
) ([]domain.LedgerEntryResponse, error) {
	if err := validateLedgerAccessRole(authContext); err != nil {
		return nil, err
	}

	if _, err := s.paymentOrderService.FindByID(ctx, authContext, paymentOrderID); err != nil {
		return nil, err
	}

	entries, err := s.ledgerEntryRepo.FindByPaymentOrderID(ctx, paymentOrderID, authContext.TenantID)
	if err != nil {
		return nil, err
	}

	return domain.NewLedgerEntryResponses(entries), nil
}

func (s *ledgerEntryServiceImpl) FindPaginate(
	ctx context.Context,
	authContext *security.AuthContext,
	pageable *shared.Pageable,
	filter *domain.LedgerEntryFilter,
) (*shared.Page[domain.LedgerEntryResponse], error) {
	if err := validateLedgerAccessRole(authContext); err != nil {
		return shared.NewPageEmpty[domain.LedgerEntryResponse](pageable), err
	}

	page, err := s.ledgerEntryRepo.FindPaginate(ctx, authContext.TenantID, pageable, filter)
	if err != nil {
		return shared.NewPageEmpty[domain.LedgerEntryResponse](pageable), err
	}

	return &shared.Page[domain.LedgerEntryResponse]{
		Data:       domain.NewLedgerEntryResponses(page.Data),
		Pagination: page.Pagination,
	}, nil
}

func validateLedgerAccessRole(authContext *security.AuthContext) error {
	switch authContext.RoleCode {
	case model.RoleCodeSuperAdmin, model.RoleCodeSchoolAdmin, model.RoleCodeTreasurer:
		return nil
	default:
		return shared.ErrAuthUnauthorized
	}
}

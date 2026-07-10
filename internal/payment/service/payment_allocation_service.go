package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"

	"github.com/google/uuid"
)

type paymentAllocationServiceImpl struct {
	repo                  domain.RepositoryManagerPayment
	paymentAllocationRepo domain.PaymentAllocationRepository
	paymentOrderRepo      domain.PaymentOrderRepository
	studentObligationRepo domain.StudentObligationRepository
	paymentOrderService   domain.PaymentOrderService
}

func NewPaymentAllocationService(
	repo domain.RepositoryManagerPayment,
	paymentOrderService domain.PaymentOrderService,
) domain.PaymentAllocationService {
	return &paymentAllocationServiceImpl{
		repo:                  repo,
		paymentAllocationRepo: repo.PaymentAllocation(),
		paymentOrderRepo:      repo.PaymentOrder(),
		studentObligationRepo: repo.StudentObligation(),
		paymentOrderService:   paymentOrderService,
	}
}

func (s *paymentAllocationServiceImpl) Allocate(
	ctx context.Context,
	authContext *security.AuthContext,
	paymentOrderID uuid.UUID,
	req *domain.PaymentAllocationAllocate,
) (*domain.PaymentAllocationResult, error) {
	paymentOrder, err := s.paymentOrderService.FindByID(ctx, authContext, paymentOrderID)
	if err != nil {
		return nil, err
	}
	if paymentOrder.Status != model.PaymentOrderStatusPending {
		return nil, shared.ErrPaymentOrderStatusInvalid
	}
	if len(req.Allocations) == 0 {
		return nil, shared.ErrPaymentAllocationRequired
	}

	requestedObligationIDs := make([]uuid.UUID, 0, len(req.Allocations))
	requestedAmounts := make(map[uuid.UUID]float64, len(req.Allocations))
	totalRequested := 0.0
	for _, allocation := range req.Allocations {
		if allocation.AllocatedAmount <= 0 {
			return nil, shared.ErrPaymentAllocationAmountInvalid
		}
		if _, exists := requestedAmounts[allocation.StudentObligationID]; exists {
			return nil, shared.ErrPaymentAllocationDuplicateObligation
		}
		requestedObligationIDs = append(requestedObligationIDs, allocation.StudentObligationID)
		requestedAmounts[allocation.StudentObligationID] = allocation.AllocatedAmount
		totalRequested += allocation.AllocatedAmount
	}

	existingAllocations, err := s.paymentAllocationRepo.FindByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return nil, err
	}

	existingTotalAllocated := 0.0
	existingAllocationMap := make(map[uuid.UUID]struct{}, len(existingAllocations))
	for _, allocation := range existingAllocations {
		existingTotalAllocated += allocation.AllocatedAmount
		existingAllocationMap[allocation.StudentObligationID] = struct{}{}
	}
	for obligationID := range requestedAmounts {
		if _, exists := existingAllocationMap[obligationID]; exists {
			return nil, shared.ErrPaymentAllocationDuplicateObligation
		}
	}

	if existingTotalAllocated+totalRequested > paymentOrder.TotalAmount {
		return nil, shared.ErrPaymentAllocationTotalExceedsOrder
	}

	obligations, err := s.studentObligationRepo.FindByIDs(ctx, requestedObligationIDs)
	if err != nil {
		return nil, err
	}
	if len(obligations) != len(requestedObligationIDs) {
		return nil, shared.ErrStudentObligationNotFound
	}

	obligationMap := make(map[uuid.UUID]model.StudentObligation, len(obligations))
	for _, obligation := range obligations {
		obligationMap[obligation.ID] = obligation
	}

	newAllocations := make([]model.PaymentAllocation, 0, len(requestedObligationIDs))
	type settlement struct {
		id                uuid.UUID
		outstandingAmount float64
		status            model.StudentObligationStatus
	}
	settlements := make([]settlement, 0, len(requestedObligationIDs))

	for _, obligationID := range requestedObligationIDs {
		obligation, exists := obligationMap[obligationID]
		if !exists {
			return nil, shared.ErrStudentObligationNotFound
		}
		if obligation.StudentID != paymentOrder.StudentID {
			return nil, shared.ErrStudentObligationNotFound
		}

		allocatedAmount := requestedAmounts[obligationID]
		if allocatedAmount > obligation.OutstandingAmount {
			return nil, shared.ErrPaymentAllocationAmountExceedsOutstanding
		}

		outstandingAmount := obligation.OutstandingAmount - allocatedAmount
		status := model.StudentObligationStatusPartial
		if outstandingAmount == 0 {
			status = model.StudentObligationStatusPaid
		}

		newAllocations = append(newAllocations, model.PaymentAllocation{
			PaymentOrderID:      paymentOrderID,
			StudentObligationID: obligationID,
			AllocatedAmount:     allocatedAmount,
		})
		settlements = append(settlements, settlement{
			id:                obligationID,
			outstandingAmount: outstandingAmount,
			status:            status,
		})
	}

	orderStatus := model.PaymentOrderStatusPending
	if existingTotalAllocated+totalRequested == paymentOrder.TotalAmount {
		orderStatus = model.PaymentOrderStatusCompleted
	}

	if err := s.repo.Transaction(ctx, func(repo domain.RepositoryManagerPayment) error {
		if err := repo.PaymentAllocation().CreateBatch(ctx, newAllocations); err != nil {
			return err
		}
		for _, settlement := range settlements {
			if _, err := repo.StudentObligation().UpdateSettlement(
				ctx,
				settlement.id,
				settlement.outstandingAmount,
				settlement.status,
			); err != nil {
				return err
			}
		}
		if _, err := repo.PaymentOrder().UpdateStatus(ctx, paymentOrderID, orderStatus); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return s.FindByPaymentOrderID(ctx, authContext, paymentOrderID)
}

func (s *paymentAllocationServiceImpl) FindByPaymentOrderID(
	ctx context.Context,
	authContext *security.AuthContext,
	paymentOrderID uuid.UUID,
) (*domain.PaymentAllocationResult, error) {
	paymentOrder, err := s.paymentOrderService.FindByID(ctx, authContext, paymentOrderID)
	if err != nil {
		return nil, err
	}

	allocations, err := s.paymentAllocationRepo.FindByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return nil, err
	}

	totalAllocated := 0.0
	for _, allocation := range allocations {
		totalAllocated += allocation.AllocatedAmount
	}

	return &domain.PaymentAllocationResult{
		PaymentOrderID:  paymentOrderID,
		TotalAllocated:  totalAllocated,
		RemainingAmount: paymentOrder.TotalAmount - totalAllocated,
		OrderStatus:     paymentOrder.Status,
		Allocations:     domain.NewPaymentAllocationResponses(allocations),
	}, nil
}

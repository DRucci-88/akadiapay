package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"sort"

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
		if !shared.FloatGreater(allocation.AllocatedAmount, 0) {
			return nil, shared.ErrPaymentAllocationAmountInvalid
		}
		if _, exists := requestedAmounts[allocation.StudentObligationID]; exists {
			return nil, shared.ErrPaymentAllocationDuplicateObligation
		}
		requestedObligationIDs = append(requestedObligationIDs, allocation.StudentObligationID)
		requestedAmounts[allocation.StudentObligationID] = allocation.AllocatedAmount
		totalRequested += allocation.AllocatedAmount
	}
	sort.Slice(requestedObligationIDs, func(i, j int) bool {
		return requestedObligationIDs[i].String() < requestedObligationIDs[j].String()
	})

	newAllocations := make([]model.PaymentAllocation, 0, len(requestedObligationIDs))
	type settlement struct {
		id                uuid.UUID
		outstandingAmount float64
		status            model.StudentObligationStatus
	}
	settlements := make([]settlement, 0, len(requestedObligationIDs))

	if err := s.repo.Transaction(ctx, func(repo domain.RepositoryManagerPayment) error {
		lockedPaymentOrder, err := repo.PaymentOrder().LockByID(ctx, paymentOrderID, authContext.TenantID)
		if err != nil {
			return err
		}
		if lockedPaymentOrder.Status != model.PaymentOrderStatusPending {
			return shared.ErrPaymentOrderStatusInvalid
		}

		existingAllocations, err := repo.PaymentAllocation().FindByPaymentOrderID(ctx, paymentOrderID)
		if err != nil {
			return err
		}

		existingTotalAllocated := 0.0
		existingAllocationMap := make(map[uuid.UUID]struct{}, len(existingAllocations))
		for _, allocation := range existingAllocations {
			existingTotalAllocated += allocation.AllocatedAmount
			existingAllocationMap[allocation.StudentObligationID] = struct{}{}
		}
		for obligationID := range requestedAmounts {
			if _, exists := existingAllocationMap[obligationID]; exists {
				return shared.ErrPaymentAllocationDuplicateObligation
			}
		}

		if shared.FloatGreater(existingTotalAllocated+totalRequested, lockedPaymentOrder.TotalAmount) {
			return shared.ErrPaymentAllocationTotalExceedsOrder
		}

		obligations, err := repo.StudentObligation().LockByIDs(ctx, requestedObligationIDs)
		if err != nil {
			return err
		}
		if len(obligations) != len(requestedObligationIDs) {
			return shared.ErrStudentObligationNotFound
		}

		obligationMap := make(map[uuid.UUID]model.StudentObligation, len(obligations))
		for _, obligation := range obligations {
			obligationMap[obligation.ID] = obligation
		}

		paymentProductIDs := make([]uuid.UUID, 0, len(obligations))
		for _, obligation := range obligations {
			paymentProductIDs = append(paymentProductIDs, obligation.PaymentProductID)
		}

		paymentProducts, err := repo.PaymentProduct().FindByIDsIncludingDeleted(ctx, uniqueUUIDs(paymentProductIDs))
		if err != nil {
			return err
		}

		paymentProductMap := make(map[uuid.UUID]model.PaymentProduct, len(paymentProducts))
		for _, paymentProduct := range paymentProducts {
			paymentProductMap[paymentProduct.ID] = paymentProduct
		}

		newAllocations = newAllocations[:0]
		settlements = settlements[:0]
		for _, obligationID := range requestedObligationIDs {
			obligation, exists := obligationMap[obligationID]
			if !exists {
				return shared.ErrStudentObligationNotFound
			}
			if obligation.StudentID != lockedPaymentOrder.StudentID {
				return shared.ErrStudentObligationNotFound
			}

			paymentProduct, exists := paymentProductMap[obligation.PaymentProductID]
			if !exists {
				return shared.ErrPaymentProductNotFound
			}
			if paymentProduct.PaymentPolicy == nil {
				return shared.ErrPaymentPolicyNotFound
			}

			allocatedAmount := requestedAmounts[obligationID]
			if err := validateAllocationAgainstPaymentPolicy(
				allocatedAmount,
				obligation.OutstandingAmount,
				paymentProduct.PaymentPolicy,
			); err != nil {
				return err
			}

			outstandingAmount := obligation.OutstandingAmount - allocatedAmount
			status := model.StudentObligationStatusPartial
			if shared.FloatIsZero(outstandingAmount) {
				status = model.StudentObligationStatusPaid
				if paymentProduct.PaymentPolicy.AutoCloseObligation {
					status = model.StudentObligationStatusClosed
				}
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
		if shared.FloatEqual(existingTotalAllocated+totalRequested, lockedPaymentOrder.TotalAmount) {
			orderStatus = model.PaymentOrderStatusCompleted
		}

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
		if orderStatus == model.PaymentOrderStatusCompleted {
			if err := postLedgerEntriesForPaymentOrder(ctx, repo, authContext.TenantID, paymentOrderID); err != nil {
				return err
			}
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

func validateAllocationAgainstPaymentPolicy(
	allocatedAmount float64,
	outstandingAmount float64,
	paymentPolicy *model.PaymentPolicy,
) error {
	if shared.FloatGreater(allocatedAmount, outstandingAmount) {
		return shared.ErrPaymentAllocationAmountExceedsOutstanding
	}

	if !paymentPolicy.AllowPartial && !shared.FloatEqual(allocatedAmount, outstandingAmount) {
		return shared.ErrPaymentAllocationFullPaymentRequired
	}

	if paymentPolicy.AllowPartial {
		if shared.FloatGreater(paymentPolicy.MinimumAmount, 0) &&
			shared.FloatLess(allocatedAmount, paymentPolicy.MinimumAmount) {
			return shared.ErrPaymentAllocationBelowMinimumAmount
		}

		minimumByPercentage := outstandingAmount * paymentPolicy.MinimumPercentage / 100
		if shared.FloatGreater(paymentPolicy.MinimumPercentage, 0) &&
			shared.FloatLess(allocatedAmount, minimumByPercentage) {
			return shared.ErrPaymentAllocationBelowMinimumPercentage
		}
	}

	if !paymentPolicy.AllowOverPayment && shared.FloatGreater(allocatedAmount, outstandingAmount) {
		return shared.ErrPaymentAllocationAmountExceedsOutstanding
	}

	return nil
}

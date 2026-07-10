package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type paymentOrderServiceImpl struct {
	repo                  domain.RepositoryManagerPayment
	paymentAllocationRepo domain.PaymentAllocationRepository
	paymentOrderRepo      domain.PaymentOrderRepository
	studentObligationRepo domain.StudentObligationRepository
	studentService        domain.StudentService
	parentStudentService  domain.ParentStudentService
}

func NewPaymentOrderService(
	repo domain.RepositoryManagerPayment,
	studentService domain.StudentService,
	parentStudentService domain.ParentStudentService,
) domain.PaymentOrderService {
	return &paymentOrderServiceImpl{
		repo:                  repo,
		paymentAllocationRepo: repo.PaymentAllocation(),
		paymentOrderRepo:      repo.PaymentOrder(),
		studentObligationRepo: repo.StudentObligation(),
		studentService:        studentService,
		parentStudentService:  parentStudentService,
	}
}

func (s *paymentOrderServiceImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentOrderFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentOrderResponse], error) {
	page, err := s.paymentOrderRepo.FindPaginate(ctx, pageable, filter, authContext)
	if err != nil {
		return shared.NewPageEmpty[domain.PaymentOrderResponse](pageable), err
	}

	return &shared.Page[domain.PaymentOrderResponse]{
		Data:       domain.NewPaymentOrderResponses(page.Data),
		Pagination: page.Pagination,
	}, nil
}

func (s *paymentOrderServiceImpl) FindByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentOrder, error) {
	paymentOrder, err := s.paymentOrderRepo.FirstByID(ctx, id, authContext.TenantID)
	if err != nil {
		return nil, err
	}

	if _, err := s.validateStudentAccess(ctx, authContext, paymentOrder.StudentID, shared.ErrPaymentOrderNotFound); err != nil {
		return nil, err
	}

	return paymentOrder, nil
}

func (s *paymentOrderServiceImpl) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentOrderCreate,
) (*model.PaymentOrder, error) {
	if req.PaymentDate.IsZero() {
		return nil, shared.ErrPaymentOrderDateRequired
	}
	if !shared.FloatGreater(req.Amount, 0) {
		return nil, shared.ErrPaymentOrderAmountInvalid
	}
	if !isValidPaymentOrderMethod(req.PaymentMethod) {
		return nil, shared.ErrPaymentOrderMethodInvalid
	}

	if _, err := s.validateStudentAccess(ctx, authContext, req.StudentID, shared.ErrStudentNotFound); err != nil {
		return nil, err
	}

	totalOutstanding, err := s.studentObligationRepo.SumOutstandingByStudentID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if !shared.FloatGreater(totalOutstanding, 0) {
		return nil, shared.ErrPaymentOrderOutstandingRequired
	}
	if shared.FloatGreater(req.Amount, totalOutstanding) {
		return nil, shared.ErrPaymentOrderAmountExceedsOutstanding
	}

	now := time.Now().UTC()
	paymentOrder := &model.PaymentOrder{
		TenantID:      authContext.TenantID,
		StudentID:     req.StudentID,
		PaidByUserID:  authContext.UserID,
		OrderNumber:   generatePaymentOrderNumber(now),
		OrderDate:     req.PaymentDate,
		TotalAmount:   req.Amount,
		Status:        model.PaymentOrderStatusPending,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
	}

	if err := s.paymentOrderRepo.Create(ctx, paymentOrder); err != nil {
		return nil, err
	}

	return paymentOrder, nil
}

func (s *paymentOrderServiceImpl) Cancel(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentOrder, error) {
	paymentOrder, err := s.FindByID(ctx, authContext, id)
	if err != nil {
		return nil, err
	}
	if paymentOrder.LedgerPostedAt != nil {
		return nil, shared.ErrPostedPaymentCannotBeCancelled
	}
	if paymentOrder.Status != model.PaymentOrderStatusPending {
		return nil, shared.ErrPaymentOrderStatusInvalid
	}

	paymentAllocations, err := s.paymentAllocationRepo.FindByPaymentOrderID(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(paymentAllocations) > 0 {
		return nil, shared.ErrPaymentOrderAllocated
	}

	if _, err := s.paymentOrderRepo.UpdateStatus(ctx, id, model.PaymentOrderStatusCancelled); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, authContext, id)
}

func generatePaymentOrderNumber(now time.Time) string {
	return fmt.Sprintf(
		"PO-%s-%s",
		now.Format("20060102150405"),
		uuid.NewString()[:8],
	)
}

func isValidPaymentOrderMethod(method model.PaymentOrderPaymentMethod) bool {
	switch method {
	case model.PaymentMethodCash,
		model.PaymentMethodBankTransfer,
		model.PaymentMethodVirtualAccount,
		model.PaymentMethodQRIS,
		model.PaymentMethodCreditCard:
		return true
	default:
		return false
	}
}

func (s *paymentOrderServiceImpl) validateStudentAccess(
	ctx context.Context,
	authContext *security.AuthContext,
	studentID uuid.UUID,
	notFoundErr error,
) (*model.Student, error) {
	student, err := s.studentService.FirstByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if student.TenantID != authContext.TenantID {
		return nil, notFoundErr
	}

	switch authContext.RoleCode {
	case model.RoleCodeParent:
		exists, err := s.parentStudentService.ExistsByParentUserIDAndStudentID(ctx, authContext.UserID, studentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, notFoundErr
		}
	case model.RoleCodeStudent:
		if authContext.StudentID == nil || *authContext.StudentID != studentID {
			return nil, notFoundErr
		}
	}

	return student, nil
}

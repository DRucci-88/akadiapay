package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"time"
)

type studentObligationServiceImpl struct {
	repo                  domain.RepositoryManagerPayment
	studentObligationRepo domain.StudentObligationRepository
	studentService        domain.StudentService
	paymentProductService domain.PaymentProductService
}

func NewStudentObligationService(
	repo domain.RepositoryManagerPayment,
	studentService domain.StudentService,
	paymentProductService domain.PaymentProductService,
) domain.StudentObligationService {
	return &studentObligationServiceImpl{
		repo:                  repo,
		studentObligationRepo: repo.StudentObligation(),
		studentService:        studentService,
		paymentProductService: paymentProductService,
	}
}

func (s *studentObligationServiceImpl) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.StudentObligationCreate,
) (*model.StudentObligation, error) {
	if req.DueDate.IsZero() {
		return nil, shared.ErrStudentObligationDueDateRequired
	}

	student, err := s.studentService.FirstByID(ctx, req.StudentID)
	if err != nil {
		return nil, err
	}
	if student.TenantID != authContext.TenantID {
		return nil, shared.ErrStudentNotFound
	}

	paymentProduct, err := s.paymentProductService.FindByID(
		ctx,
		authContext,
		req.PaymentProductID,
		model.PaymentProductPreloadPaymentPolicy,
	)
	if err != nil {
		return nil, err
	}
	if paymentProduct.PaymentPolicy == nil {
		return nil, shared.ErrPaymentPolicyNotFound
	}

	amount := paymentProduct.Price
	if req.Amount != nil {
		amount = *req.Amount
	}
	if amount <= 0 {
		return nil, shared.ErrStudentObligationAmountInvalid
	}

	period := normalizeStudentObligationPeriod(req.DueDate)
	exists, err := s.studentObligationRepo.ExistsActiveByStudentIDAndPaymentProductIDAndPeriod(
		ctx,
		req.StudentID,
		req.PaymentProductID,
		period,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, shared.ErrStudentObligationAlreadyExists
	}

	now := time.Now().UTC()
	studentObligation := &model.StudentObligation{
		StudentID:         req.StudentID,
		PaymentProductID:  req.PaymentProductID,
		Period:            period,
		OriginalAmount:    amount,
		OutstandingAmount: amount,
		DueDate:           req.DueDate,
		IssuedAt:          now,
		Status:            model.StudentObligationStatusPending,
		Notes:             req.Notes,
	}

	if err := s.studentObligationRepo.Create(ctx, studentObligation); err != nil {
		return nil, err
	}

	return studentObligation, nil
}

func (s *studentObligationServiceImpl) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.StudentObligationFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.StudentObligationResponse], error) {
	page, err := s.studentObligationRepo.FindPaginate(ctx, pageable, filter, authContext)
	if err != nil {
		return shared.NewPageEmpty[domain.StudentObligationResponse](pageable), err
	}

	return &shared.Page[domain.StudentObligationResponse]{
		Data:       domain.NewStudentObligationResponses(page.Data),
		Pagination: page.Pagination,
	}, nil
}

func normalizeStudentObligationPeriod(dueDate time.Time) time.Time {
	return time.Date(dueDate.Year(), dueDate.Month(), 1, 0, 0, 0, 0, time.UTC)
}

package service

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"time"

	"github.com/google/uuid"
)

type studentObligationServiceImpl struct {
	repo                  domain.RepositoryManagerPayment
	studentObligationRepo domain.StudentObligationRepository
	studentService        domain.StudentService
	paymentProductService domain.PaymentProductService
	parentStudentService  domain.ParentStudentService
}

func NewStudentObligationService(
	repo domain.RepositoryManagerPayment,
	studentService domain.StudentService,
	paymentProductService domain.PaymentProductService,
	parentStudentService domain.ParentStudentService,
) domain.StudentObligationService {
	return &studentObligationServiceImpl{
		repo:                  repo,
		studentObligationRepo: repo.StudentObligation(),
		studentService:        studentService,
		paymentProductService: paymentProductService,
		parentStudentService:  parentStudentService,
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

	paymentProduct, err := s.getValidatedPaymentProduct(ctx, authContext, req.PaymentProductID)
	if err != nil {
		return nil, err
	}

	studentObligation, err := s.prepareStudentObligation(
		ctx,
		authContext,
		s.studentObligationRepo,
		req.StudentID,
		paymentProduct,
		req.DueDate,
		req.Amount,
		req.Notes,
	)
	if err != nil {
		return nil, err
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

func (s *studentObligationServiceImpl) CreateBulk(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.StudentObligationBulkCreate,
) ([]model.StudentObligation, error) {
	if len(req.StudentIDs) == 0 {
		return nil, shared.ErrStudentObligationStudentIDsEmpty
	}
	if req.DueDate.IsZero() {
		return nil, shared.ErrStudentObligationDueDateRequired
	}

	paymentProduct, err := s.getValidatedPaymentProduct(ctx, authContext, req.PaymentProductID)
	if err != nil {
		return nil, err
	}

	studentIDs := uniqueUUIDs(req.StudentIDs)
	studentObligations := make([]model.StudentObligation, 0, len(studentIDs))
	for _, studentID := range studentIDs {
		studentObligation, err := s.prepareStudentObligation(
			ctx,
			authContext,
			s.studentObligationRepo,
			studentID,
			paymentProduct,
			req.DueDate,
			nil,
			"",
		)
		if err != nil {
			return nil, err
		}

		studentObligations = append(studentObligations, *studentObligation)
	}

	if err := s.repo.Transaction(ctx, func(repo domain.RepositoryManagerPayment) error {
		return repo.StudentObligation().CreateBatch(ctx, studentObligations)
	}); err != nil {
		return nil, err
	}

	return studentObligations, nil
}

func (s *studentObligationServiceImpl) FindOutstandingByStudentID(
	ctx context.Context,
	authContext *security.AuthContext,
	studentID uuid.UUID,
) (*domain.StudentOutstandingResponse, error) {
	_, err := s.validateStudentAccess(ctx, authContext, studentID)
	if err != nil {
		return nil, err
	}

	obligations, err := s.studentObligationRepo.FindOutstandingByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	totalOutstanding, err := s.studentObligationRepo.SumOutstandingByStudentID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return &domain.StudentOutstandingResponse{
		StudentID:        studentID,
		TotalOutstanding: totalOutstanding,
		Obligations:      domain.NewStudentObligationResponses(obligations),
	}, nil
}

func (s *studentObligationServiceImpl) FindByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.StudentObligation, error) {
	studentObligation, err := s.studentObligationRepo.FirstByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.validateStudentAccess(ctx, authContext, studentObligation.StudentID)
	if err != nil {
		return nil, err
	}

	return studentObligation, nil
}

func (s *studentObligationServiceImpl) Update(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	req *domain.StudentObligationUpdate,
) (*model.StudentObligation, error) {
	studentObligation, err := s.FindByID(ctx, authContext, id)
	if err != nil {
		return nil, err
	}
	if req.DueDate == nil && req.Notes == nil {
		return studentObligation, nil
	}

	if req.DueDate != nil && req.DueDate.IsZero() {
		return nil, shared.ErrStudentObligationDueDateRequired
	}

	if _, err := s.studentObligationRepo.Update(ctx, studentObligation.ID, req); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, authContext, id)
}

func (s *studentObligationServiceImpl) Delete(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) error {
	studentObligation, err := s.FindByID(ctx, authContext, id)
	if err != nil {
		return err
	}

	hasPaymentAllocations, err := s.studentObligationRepo.HasPaymentAllocations(ctx, studentObligation.ID)
	if err != nil {
		return err
	}
	if hasPaymentAllocations {
		return shared.ErrStudentObligationAllocated
	}

	_, err = s.studentObligationRepo.Delete(ctx, studentObligation.ID)
	return err
}

func (s *studentObligationServiceImpl) getValidatedPaymentProduct(
	ctx context.Context,
	authContext *security.AuthContext,
	paymentProductID uuid.UUID,
) (*model.PaymentProduct, error) {
	paymentProduct, err := s.paymentProductService.FindByID(
		ctx,
		authContext,
		paymentProductID,
		model.PaymentProductPreloadPaymentPolicy,
	)
	if err != nil {
		return nil, err
	}
	if paymentProduct.PaymentPolicy == nil {
		return nil, shared.ErrPaymentPolicyNotFound
	}

	return paymentProduct, nil
}

func (s *studentObligationServiceImpl) validateStudentAccess(
	ctx context.Context,
	authContext *security.AuthContext,
	studentID uuid.UUID,
) (*model.Student, error) {
	student, err := s.studentService.FirstByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if student.TenantID != authContext.TenantID {
		return nil, shared.ErrStudentNotFound
	}

	switch authContext.RoleCode {
	case model.RoleCodeParent:
		exists, err := s.parentStudentService.ExistsByParentUserIDAndStudentID(ctx, authContext.UserID, studentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, shared.ErrStudentNotFound
		}
	case model.RoleCodeStudent:
		if authContext.StudentID == nil || *authContext.StudentID != studentID {
			return nil, shared.ErrStudentNotFound
		}
	}

	return student, nil
}

func (s *studentObligationServiceImpl) prepareStudentObligation(
	ctx context.Context,
	authContext *security.AuthContext,
	studentObligationRepo domain.StudentObligationRepository,
	studentID uuid.UUID,
	paymentProduct *model.PaymentProduct,
	dueDate time.Time,
	amountRequest *float64,
	notes string,
) (*model.StudentObligation, error) {
	student, err := s.studentService.FirstByID(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if student.TenantID != authContext.TenantID {
		return nil, shared.ErrStudentNotFound
	}

	amount := paymentProduct.Price
	if amountRequest != nil {
		amount = *amountRequest
	}
	if amount <= 0 {
		return nil, shared.ErrStudentObligationAmountInvalid
	}

	period := normalizeStudentObligationPeriod(dueDate)
	exists, err := studentObligationRepo.ExistsActiveByStudentIDAndPaymentProductIDAndPeriod(
		ctx,
		studentID,
		paymentProduct.ID,
		period,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, shared.ErrStudentObligationAlreadyExists
	}

	now := time.Now().UTC()
	return &model.StudentObligation{
		StudentID:         studentID,
		PaymentProductID:  paymentProduct.ID,
		Period:            period,
		OriginalAmount:    amount,
		OutstandingAmount: amount,
		DueDate:           dueDate,
		IssuedAt:          now,
		Status:            model.StudentObligationStatusPending,
		Notes:             notes,
	}, nil
}

func uniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	res := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		res = append(res, id)
	}
	return res
}

func normalizeStudentObligationPeriod(dueDate time.Time) time.Time {
	return time.Date(dueDate.Year(), dueDate.Month(), 1, 0, 0, 0, 0, time.UTC)
}

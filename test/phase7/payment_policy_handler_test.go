package phase7_test

import (
	"akadia/domain"
	paymenthandler "akadia/internal/payment/handler"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	tenantID = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	userID   = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	policyID = uuid.MustParse("44444444-4444-4444-4444-444444444444")
)

func adminAuthContext() *security.AuthContext {
	return &security.AuthContext{
		UserID:         userID,
		TenantID:       tenantID,
		RoleCode:       model.RoleCodeSchoolAdmin,
		TokenExpiredAt: time.Now().UTC().Add(time.Hour),
	}
}

func makePaymentPolicy() *model.PaymentPolicy {
	return &model.PaymentPolicy{
		TenantID: tenantID,
		Code:     "FULL_PAYMENT",
		Name:     "Full Payment",
		BaseModel: model.BaseModel{
			ID: policyID,
		},
	}
}

type fakePaymentPolicyService struct {
	firstByIDFn func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentPolicy, error)
	createFn    func(ctx context.Context, authContext *security.AuthContext, req *domain.PaymentPolicyCreate) (*model.PaymentPolicy, error)
}

func (f *fakePaymentPolicyService) FirstByID(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
) (*model.PaymentPolicy, error) {
	if f.firstByIDFn == nil {
		panic("unexpected call: FirstByID")
	}
	return f.firstByIDFn(ctx, authContext, id)
}

func (f *fakePaymentPolicyService) FindPaginate(
	ctx context.Context,
	pageable *shared.Pageable,
	filter *domain.PaymentPolicyFilter,
	authContext *security.AuthContext,
) (*shared.Page[domain.PaymentPolicyResponse], error) {
	panic("unexpected call: FindPaginate")
}

func (f *fakePaymentPolicyService) Create(
	ctx context.Context,
	authContext *security.AuthContext,
	req *domain.PaymentPolicyCreate,
) (*model.PaymentPolicy, error) {
	if f.createFn == nil {
		panic("unexpected call: Create")
	}
	return f.createFn(ctx, authContext, req)
}

func (f *fakePaymentPolicyService) Update(
	ctx context.Context,
	authContext *security.AuthContext,
	id uuid.UUID,
	req *domain.PaymentPolicyUpdate,
) (*model.PaymentPolicy, error) {
	panic("unexpected call: Update")
}

func newTestContext(method string, target string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(method, target, bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(domain.ContextKeyAuth, adminAuthContext())
	return c, recorder
}

func TestPaymentPolicyHandlerCreateInvalidJSONReturns422(t *testing.T) {
	service := &fakePaymentPolicyService{
		createFn: func(ctx context.Context, authContext *security.AuthContext, req *domain.PaymentPolicyCreate) (*model.PaymentPolicy, error) {
			t.Fatalf("service create should not be called for invalid json")
			return nil, nil
		},
	}
	handler := paymenthandler.NewPaymentPolicyHandler(service)
	c, recorder := newTestContext(http.MethodPost, "/payment-policies", []byte(`{`))

	handler.Create(c)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", recorder.Code)
	}
}

func TestPaymentPolicyHandlerFindByIDInvalidUUIDReturns422(t *testing.T) {
	service := &fakePaymentPolicyService{
		firstByIDFn: func(ctx context.Context, authContext *security.AuthContext, id uuid.UUID) (*model.PaymentPolicy, error) {
			t.Fatalf("service lookup should not be called for invalid uuid")
			return nil, nil
		},
	}
	handler := paymenthandler.NewPaymentPolicyHandler(service)
	c, recorder := newTestContext(http.MethodGet, "/payment-policies/not-a-uuid", nil)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	handler.FindByID(c)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", recorder.Code)
	}
}

func TestPaymentPolicyHandlerCreateSuccessReturnsResponseData(t *testing.T) {
	service := &fakePaymentPolicyService{
		createFn: func(ctx context.Context, authContext *security.AuthContext, req *domain.PaymentPolicyCreate) (*model.PaymentPolicy, error) {
			if authContext == nil || authContext.TenantID != tenantID {
				t.Fatalf("expected auth context to be forwarded")
			}
			if req.Code != "FULL_PAYMENT" || req.Name != "Full Payment" {
				t.Fatalf("expected request payload to be bound and forwarded")
			}
			return makePaymentPolicy(), nil
		},
	}
	handler := paymenthandler.NewPaymentPolicyHandler(service)
	body := []byte(`{"code":"FULL_PAYMENT","name":"Full Payment","allow_partial":false}`)
	c, recorder := newTestContext(http.MethodPost, "/payment-policies", body)

	handler.Create(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var payload map[string]map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	data := payload["data"]
	if data["id"] != policyID.String() {
		t.Fatalf("expected response id %s, got %v", policyID, data["id"])
	}
	if data["code"] != "FULL_PAYMENT" {
		t.Fatalf("expected response code FULL_PAYMENT, got %v", data["code"])
	}
}

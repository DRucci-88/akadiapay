package seed

import (
	"akadia/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB

	tenants            map[string]model.Tenant
	users              map[string]model.User
	roles              map[model.RoleCode]model.Role
	students           map[string]model.Student
	paymentPolicies    map[string]model.PaymentPolicy
	paymentProducts    map[string]model.PaymentProduct
	paymentOrders      map[string]model.PaymentOrder
	studentObligations map[string]model.StudentObligation
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,

		tenants:            make(map[string]model.Tenant),
		users:              make(map[string]model.User),
		roles:              make(map[model.RoleCode]model.Role),
		students:           make(map[string]model.Student),
		paymentPolicies:    make(map[string]model.PaymentPolicy),
		paymentProducts:    make(map[string]model.PaymentProduct),
		paymentOrders:      make(map[string]model.PaymentOrder),
		studentObligations: make(map[string]model.StudentObligation),
	}
}

func (s *Seeder) MustTenantByCode(code string) model.Tenant {
	if tenant, ok := s.tenants[code]; ok {
		return tenant
	}

	var tenant model.Tenant

	if err := s.db.
		Where("code = ?", code).
		First(&tenant).Error; err != nil {
		panic(fmt.Errorf("tenant '%s' not found: %w", code, err))
	}

	s.tenants[code] = tenant

	return tenant
}

func (s *Seeder) MustUserByEmail(email string) model.User {
	if user, ok := s.users[email]; ok {
		return user
	}

	var user model.User

	if err := s.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {
		panic(fmt.Errorf("user '%s' not found: %w", email, err))
	}

	s.users[email] = user

	return user
}

func (s *Seeder) MustStudentByNISN(nisn string) model.Student {
	if student, ok := s.students[nisn]; ok {
		return student
	}

	var student model.Student

	if err := s.db.
		Where("nisn = ?", nisn).
		First(&student).Error; err != nil {
		panic(fmt.Errorf("student '%s' not found: %w", nisn, err))
	}

	s.students[nisn] = student

	return student
}

func (s *Seeder) MustRoleByCode(code model.RoleCode) model.Role {
	if role, ok := s.roles[code]; ok {
		return role
	}

	var role model.Role

	if err := s.db.
		Where("code = ?", code).
		First(&role).Error; err != nil {
		panic(fmt.Errorf("role '%s' not found: %w", code, err))
	}

	s.roles[code] = role

	return role
}

func (s *Seeder) MustPaymentPolicyByCode(tenantCode string, code string) model.PaymentPolicy {
	key := tenantCode + ":" + code
	if paymentPolicy, ok := s.paymentPolicies[key]; ok {
		return paymentPolicy
	}

	tenant := s.MustTenantByCode(tenantCode)
	var paymentPolicy model.PaymentPolicy

	if err := s.db.
		Where("tenant_id = ?", tenant.ID).
		Where("code = ?", code).
		First(&paymentPolicy).Error; err != nil {
		panic(fmt.Errorf("payment policy '%s' for tenant '%s' not found: %w", code, tenantCode, err))
	}

	s.paymentPolicies[key] = paymentPolicy
	return paymentPolicy
}

func (s *Seeder) MustPaymentProductByCode(code string) model.PaymentProduct {
	if paymentProduct, ok := s.paymentProducts[code]; ok {
		return paymentProduct
	}

	var paymentProduct model.PaymentProduct
	if err := s.db.
		Where("code = ?", code).
		First(&paymentProduct).Error; err != nil {
		panic(fmt.Errorf("payment product '%s' not found: %w", code, err))
	}

	s.paymentProducts[code] = paymentProduct
	return paymentProduct
}

func (s *Seeder) MustPaymentOrderByNumber(orderNumber string) model.PaymentOrder {
	if paymentOrder, ok := s.paymentOrders[orderNumber]; ok {
		return paymentOrder
	}

	var paymentOrder model.PaymentOrder
	if err := s.db.
		Where("order_number = ?", orderNumber).
		First(&paymentOrder).Error; err != nil {
		panic(fmt.Errorf("payment order '%s' not found: %w", orderNumber, err))
	}

	s.paymentOrders[orderNumber] = paymentOrder
	return paymentOrder
}

func (s *Seeder) MustStudentObligationByStudentProductPeriod(nisn string, paymentProductCode string, period time.Time) model.StudentObligation {
	key := nisn + ":" + paymentProductCode + ":" + period.Format("2006-01-02")
	if studentObligation, ok := s.studentObligations[key]; ok {
		return studentObligation
	}

	student := s.MustStudentByNISN(nisn)
	paymentProduct := s.MustPaymentProductByCode(paymentProductCode)

	var studentObligation model.StudentObligation
	if err := s.db.
		Where("student_id = ?", student.ID).
		Where("payment_product_id = ?", paymentProduct.ID).
		Where("period = ?", period).
		First(&studentObligation).Error; err != nil {
		panic(fmt.Errorf("student obligation student='%s' product='%s' period='%s' not found: %w", nisn, paymentProductCode, period.Format("2006-01-02"), err))
	}

	s.studentObligations[key] = studentObligation
	return studentObligation
}

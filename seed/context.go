package seed

import (
	"akadia/model"
	"fmt"

	"gorm.io/gorm"
)

type Seeder struct {
	db *gorm.DB

	tenants  map[string]model.Tenant
	users    map[string]model.User
	roles    map[model.RoleCode]model.Role
	students map[string]model.Student
}

func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{
		db: db,

		tenants:  make(map[string]model.Tenant),
		users:    make(map[string]model.User),
		roles:    make(map[model.RoleCode]model.Role),
		students: make(map[string]model.Student),
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

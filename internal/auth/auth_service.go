package auth

import (
	"akadia/domain"
	"akadia/internal/platform/security"
	"akadia/internal/shared"
	"akadia/model"
	"context"
	"log"

	"github.com/google/uuid"
)

type authServiceImpl struct {
	appConfig             domain.AppConfigProvider
	userService           domain.UserService
	userTenantRoleService domain.UserTenantRoleService
	studentService        domain.StudentService
	tenantService         domain.TenantService
}

func NewAuthService(
	appConfig domain.AppConfigProvider,
	userService domain.UserService,
	userTenantRoleService domain.UserTenantRoleService,
	studentService domain.StudentService,
	tenantService domain.TenantService,
) domain.AuthService {
	return &authServiceImpl{
		appConfig:             appConfig,
		userService:           userService,
		userTenantRoleService: userTenantRoleService,
		studentService:        studentService,
		tenantService:         tenantService,
	}
}

func (s *authServiceImpl) Login(
	ctx context.Context,
	req *domain.AuthLoginRequest,
) ([]domain.AuthLoginResponse, error) {
	user, err := s.userService.FirstByEmail(ctx, req.Email)
	if err != nil {
		return nil, shared.ErrInvalidCredential
	}

	if !security.CheckPasswordHash(req.Password, user.Password) {
		return nil, shared.ErrInvalidCredential
	}

	log.Printf("User %+v", user)

	userTenantRoleList, err := s.userTenantRoleService.FindByUserID(ctx, user.ID)

	jwtSecretKeyByte := s.appConfig.JWT_SECRET_BYTE()
	resList := make([]domain.AuthLoginResponse, 0, len(userTenantRoleList))
	for i := range userTenantRoleList {

		temp := &userTenantRoleList[i]

		var studentID *uuid.UUID
		if temp.Role.Code == model.RoleCodeStudent {
			student, err := s.studentService.FindByUserID(ctx, temp.UserID)
			if err != nil {
				return make([]domain.AuthLoginResponse, 0), err
			}
			studentID = &student.ID
		}

		// log.Println(string(jwtSecretKeyByte))
		// log.Println(temp.User.Email)
		// log.Println(temp.UserID)
		// log.Println(temp.TenantID)
		// log.Println(studentID)
		// log.Println(temp.Role.Code)
		token, err := security.GenerateJWT(jwtSecretKeyByte, temp.User.Email, temp.UserID, temp.TenantID, studentID, temp.Role.Code)

		if err != nil {
			return make([]domain.AuthLoginResponse, 0), err
		}

		resList = append(resList, domain.AuthLoginResponse{
			Token:      token,
			UserID:     temp.UserID,
			StudentID:  studentID,
			TenantID:   temp.TenantID,
			TenantCode: temp.Tenant.Code,
			TenantName: temp.Tenant.Name,
			RoleCode:   string(temp.Role.Code),
			IsDefault:  temp.IsDefault,
		})
	}

	return resList, nil
}

func (s *authServiceImpl) Profile(
	ctx context.Context,
	authContext *security.AuthContext,
) (*domain.AuthProfileResponse, error) {
	switch authContext.RoleCode {
	case model.RoleCodeStudent:
		return s.profileStudent(ctx, authContext)
	case model.RoleCodeParent:
		return s.profileParent(ctx, authContext)
	default:
		return s.profileTenant(ctx, authContext)
	}

}

func (s *authServiceImpl) profileStudent(
	ctx context.Context,
	authContext *security.AuthContext,
) (*domain.AuthProfileResponse, error) {
	user, err := s.userService.FirstByID(ctx, authContext.UserID)

	if err != nil {
		return nil, err
	}
	student, err := s.studentService.FindByUserID(ctx, authContext.UserID)

	if err != nil {
		return nil, err
	}
	return &domain.AuthProfileResponse{
		UserID:      *student.UserID,
		RoleCode:    model.RoleCodeStudent,
		Email:       student.User.Email,
		TenantID:    student.TenantID,
		TenantCode:  student.Tenant.Code,
		TenantName:  student.Tenant.Name,
		FullName:    &student.FullName,
		DisplayName: user.DisplayName,
	}, nil
}

func (s *authServiceImpl) profileParent(
	ctx context.Context,
	authContext *security.AuthContext,
) (*domain.AuthProfileResponse, error) {
	user, err := s.userService.FirstByID(ctx, authContext.UserID)

	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantService.FirstByID(ctx, authContext.TenantID)

	if err != nil {
		return nil, err
	}

	return &domain.AuthProfileResponse{
		UserID:      user.ID,
		RoleCode:    model.RoleCodeParent,
		Email:       user.Email,
		TenantID:    tenant.ID,
		TenantCode:  tenant.Code,
		TenantName:  tenant.Name,
		FullName:    nil,
		DisplayName: user.DisplayName,
	}, nil
}

func (s *authServiceImpl) profileTenant(
	ctx context.Context,
	authContext *security.AuthContext,
) (*domain.AuthProfileResponse, error) {
	user, err := s.userService.FirstByID(ctx, authContext.UserID)

	if err != nil {
		return nil, err
	}

	tenant, err := s.tenantService.FirstByID(ctx, authContext.TenantID)

	if err != nil {
		return nil, err
	}

	userTenantRole, err := s.userTenantRoleService.FirstByUserIDAndTenantID(ctx, authContext.UserID, authContext.TenantID)

	return &domain.AuthProfileResponse{
		UserID:      user.ID,
		RoleCode:    userTenantRole.Role.Code,
		Email:       user.Email,
		TenantID:    tenant.ID,
		TenantCode:  tenant.Code,
		TenantName:  tenant.Name,
		FullName:    nil,
		DisplayName: user.DisplayName,
	}, nil
}

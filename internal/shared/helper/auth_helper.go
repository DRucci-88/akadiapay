package helper

import (
	"akadia/model"
	"akadia/plarform/security"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Cost 12 adalah standar industri yang aman dan efisien
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(
	jwtSecretKey string,
	email string,
	userID uuid.UUID,
	tenantID uuid.UUID,
	studentID *uuid.UUID,
	roleCode model.RoleCode,
) (string, error) {

	claims := security.JWTClaims{
		UserID:    userID,
		TenantID:  tenantID,
		StudentID: studentID,
		RoleCode:  roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

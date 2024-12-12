package service

import (
	"fmt"
	"learnDB/internal/domain"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	Get(int) (*domain.User, error)
	GetAll() ([]domain.User, error)
	GetUserByUsername(string) (*domain.User, error)
}

type AuthService struct {
	storage         UserStorage
	salt            string
	secretKey       []byte
	expirationTime  time.Duration
	adminCredential string
}

func NewAuthService(st UserStorage, salt string, sk []byte, exp time.Duration, ac string) *AuthService {
	return &AuthService{storage: st, salt: salt, secretKey: sk, expirationTime: exp, adminCredential: ac}
}

func (a *AuthService) CheckUserCreds(u *domain.User) (bool, OperationResult) {
	dbu, err := a.storage.GetUserByUsername(u.Username)
	if err != nil {
		log.Printf("auth service check user error: %s", err)
		return false, InternalError
	}

	if dbu == nil {
		return false, BadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbu.Password), []byte(u.Password+a.salt)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, Ok
		} else {
			return false, InternalError
		}
	}
	return true, Ok
}

func (a *AuthService) CreateAccessToken(u *domain.User) (string, bool) {
	claims := jwt.MapClaims{
		"sub":   u.Username,
		"admin": a.adminCredential == u.Username,
		"exp":   time.Now().Add(a.expirationTime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(a.secretKey)

	if err != nil {
		log.Printf("auth service create access token error: %s", err)
		return "", false
	}
	return t, true
}

func (a *AuthService) KeyFunc(token *jwt.Token) (interface{}, error) {

	switch token.Method.Alg() {
	case jwt.SigningMethodHS256.Name:
		return []byte(a.secretKey), nil
	default:
		return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
	}

}

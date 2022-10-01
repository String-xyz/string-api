package service

import (
	"net/mail"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var TOKEN_SECRET = os.Getenv("JWT_SECRET_KEY")

type JWT struct {
	ExpAt        time.Time `json:"expAt"`
	IssuedAt     time.Time `json:"issuedAt"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refreshToken"`
}

type JWTClaims struct {
	ID string
	jwt.StandardClaims
}

type AuthValidator interface {
	Validate(string) (bool, error)
}

type Auth interface {
	GenerateJWT(model.User) (JWT, error)
	GenerateAPIKey(model.Platform) error
	LoginEmail(email string, password string) error
	RefreshToken()
	LoginPK() error
	LoginOTP() error
}

type auth struct {
	repo repository.AuthStrategy
}

func NewAuth(repo repository.AuthStrategy) Auth {
	return &auth{repo}
}

func (a auth) GenerateJWT(m model.User) (JWT, error) {
	claims := JWTClaims{}
	t := &JWT{
		IssuedAt:     time.Now(),
		ExpAt:        time.Now().Add(time.Hour * 24),
		RefreshToken: uuid.NewString(),
	}

	claims.ID = m.ID
	claims.ExpiresAt = t.ExpAt.Unix()
	claims.IssuedAt = t.IssuedAt.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(TOKEN_SECRET))
	if err != nil {
		return *t, err
	}
	t.Token = signed

	return *t, nil
}

func (a auth) GenerateAPIKey(model.Platform) error {
	return nil
}

func (a auth) LoginEmail(email string, password string) error {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return errors.Wrap(err, "Invalid email")
	}
	m, err := a.repo.Get(addr.Address)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(m.Data), []byte(password))
	if err != nil {
		return errors.Wrap(err, "invalid or wrong password")
	}
	return nil
}

func (a auth) Validate(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_SECRET), nil
	})
	return t.Valid, err
}

func (a auth) LoginPK() error {
	return nil
}

func (a auth) LoginOTP() error {
	return nil
}

func (a auth) RefreshToken() {

}

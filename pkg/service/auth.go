package service

import (
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type EntityType string
type AuthType string

var TOKEN_SECRET = os.Getenv("JWT_SECRET_KEY")

const (
	Platform   = EntityType("PLATFORM")
	User       = EntityType("USER")
	JWTAuth    = AuthType("JWT")
	APIKeyAuth = AuthType("API_KEY")
)

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
	LoginPK() error
	LoginOTP() error
	RefreshToken()
}

type auth struct {
	repo repository.Auth
}

func NewAuth() Auth {
	return &auth{}
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

func (a auth) LoginPK() error {
	return nil
}

func (a auth) LoginOTP() error {
	return nil
}

func (a auth) RefreshToken() {

}

func (a auth) Validate(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(TOKEN_SECRET), nil
	})
	return t.Valid, err
}

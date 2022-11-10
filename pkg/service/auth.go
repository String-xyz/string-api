package service

import (
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserRegister = model.UserRegister
type UserLoginEmail = model.UserEmailLogin
type UserPKLogin = model.UserPKLogin

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

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
	Challenge(publicAddres string) (string, error)
	ValidateAPIKey(key string) bool
	RefreshToken(string)
}

type auth struct {
	authRepo repository.AuthStrategy
}

func NewAuth(a repository.AuthStrategy) Auth {
	return &auth{a}
}

// GenerateJWT generates a jwt token and a refresh token which is saved on redis
func (a auth) GenerateJWT(m model.User) (JWT, error) {
	claims := JWTClaims{}
	refreshToken := uuidWithoutHyphens()
	t := &JWT{
		IssuedAt:     time.Now(),
		ExpAt:        time.Now().Add(time.Hour * 24),
		RefreshToken: refreshToken,
	}

	claims.ID = m.ID
	claims.ExpiresAt = t.ExpAt.Unix()
	claims.IssuedAt = t.IssuedAt.Unix()
	// replace this signing method with RSA or something similar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return *t, err
	}
	t.Token = signed
	return *t, a.authRepo.CreateJWTRefresh(common.ToSha256(refreshToken), m.ID)
}

func (a auth) Validate(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return t.Valid, err
}

func (a auth) Challenge(publicAddress string) (string, error) {
	if !hexRegex.MatchString(publicAddress) {
		return "", errors.New("invalid address")
	}
	nonce := uuid.NewString()
	return nonce, a.authRepo.CreateAny(publicAddress, nonce, time.Minute*10)
}

func (a auth) ValidateAPIKey(key string) bool {
	hashed := common.ToSha256(key)
	authKey, err := a.authRepo.Get(hashed)
	if err != nil {
		return false
	}
	return authKey.Data == hashed
}

func (a auth) RefreshToken(token string) {
	//stra, err := a.authRepo.Get(common.ToSha256(token))
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

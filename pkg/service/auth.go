package service

import (
	netmail "net/mail"
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

type SignablePayload struct {
	Nonce string `json:"payload"`
}

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

type Auth interface {
	// PayloadToSign returns a payload to be sign by a wallet
	// to authenticate an user, the payload expires in 15 minutes
	PayloadToSign(walletAdress string) (SignablePayload, error)

	// VerifySignedPayload receives a signed payload from the user and verifies the signature
	// if signaure is valid it returns a JWT to authenticate the user
	VerifySignedPayload(model.WalletSignaturePayloadSigned) (UserCreateResponse, error)

	GenerateJWT(model.User) (JWT, error)
	ValidateAPIKey(key string) bool
}

type auth struct {
	repos repository.Repositories
}

// reusing UserRepos here
func NewAuth(r repository.Repositories) Auth {
	return &auth{r}
}

func (a auth) PayloadToSign(walletAddress string) (SignablePayload, error) {
	payload := model.WalletSignaturePayload{}
	signable := SignablePayload{}

	if !hexRegex.MatchString(walletAddress) {
		return signable, common.StringError(errors.New("missing or invalid address"))
	}
	payload.Address = walletAddress
	payload.Timestamp = time.Now().Unix()
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	encrypted, err := common.Encrypt(payload, key)
	if err != nil {
		return signable, common.StringError(err)
	}
	return SignablePayload{encrypted}, nil
}

func (a auth) VerifySignedPayload(request model.WalletSignaturePayloadSigned) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	payload, err := common.Decrypt[model.WalletSignaturePayload](request.Nonce, key)
	if err != nil {
		return resp, common.StringError(err)
	}
	err = verifyWalletAuthentication(request)
	if err != nil {
		return resp, common.StringError(err)
	}

	// Verify user is registered to this wallet address
	instrument, err := a.repos.Instrument.GetWallet(payload.Address)
	if err != nil {
		return resp, common.StringError(err)
	}
	user, err := a.repos.User.GetById(instrument.UserID)
	if err != nil {
		return resp, common.StringError(err)
	}

	// Create the JWT
	jwt, err := a.GenerateJWT(user)
	if err != nil {
		return resp, common.StringError(err)
	}
	return UserCreateResponse{JWT: jwt, User: user}, nil
}

// GenerateJWT generates a jwt token and a refresh token which is saved on redis
func (a auth) GenerateJWT(m model.User) (JWT, error) {
	claims := JWTClaims{}
	refreshToken := uuidWithoutHyphens()
	t := &JWT{
		IssuedAt:     time.Now(),
		ExpAt:        time.Now().Add(time.Minute * 15),
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
	return *t, a.repos.Auth.CreateJWTRefresh(common.ToSha256(refreshToken), m.ID)
}

func (a auth) ValidateJWT(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return t.Valid, err
}

func (a auth) ValidateAPIKey(key string) bool {
	hashed := common.ToSha256(key)
	authKey, err := a.repos.Auth.Get(hashed)
	if err != nil {
		return false
	}
	return authKey.Data == hashed
}

func verifyWalletAuthentication(request model.WalletSignaturePayloadSigned) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	preSignedPayload, err := common.Decrypt[model.WalletSignaturePayload](request.Nonce, key)
	if err != nil {
		return common.StringError(err)
	}
	// Verify users signature
	valid, err := common.ValidateExternalEVMSignature(request.Signature, preSignedPayload.Address, request.Nonce)
	if err != nil {
		return common.StringError(err)
	}
	if !valid {
		return common.StringError(errors.New("user signature invalid"))
	}

	// Verify timestamp is not expired past 15 minutes
	if time.Now().Unix() > preSignedPayload.Timestamp+(15*60) {
		return common.StringError(errors.New("login payload expired"))
	}

	return nil
}

// Use native mail package to check if email a valid email
func validEmail(email string) bool {
	_, err := netmail.ParseAddress(email)
	return err == nil
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

package service

import (
	"net/mail"
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
	"golang.org/x/crypto/bcrypt"
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
	Register(UserRegister) (JWT, error)
	LoginEmail(UserLoginEmail) (JWT, error)
	GenerateJWT(model.User) (JWT, error)
	LoginPK(UserPKLogin) (JWT, error)
	Challenge(publicAddres string) (string, error)
	GenerateAPIKey(model.Platform) error
	ValidateAPIKey(key string) bool
	RefreshToken(string)
	LoginOTP() error
}

type auth struct {
	authRepo    repository.AuthStrategy
	userRepo    repository.User
	contactRepo repository.Contact
}

func NewAuth(a repository.AuthStrategy, u repository.User, c repository.Contact) Auth {
	return &auth{a, u, c}
}

// Register registers an user with authentication (email/password)
// this is a rudimentary implementation of onboarding, will later have proper
// onboarding process.
func (a auth) Register(m UserRegister) (JWT, error) {
	tx := a.userRepo.MustBegin()
	user, err := a.userRepo.Create(model.User{FirstName: m.FirstName, LastName: m.LastName, Status: "registered", Type: "client"})
	if err != nil {
		a.userRepo.Rollback()
		return JWT{}, err
	}
	a.contactRepo.SetTx(tx)
	defer a.contactRepo.Reset()
	contact, err := a.contactRepo.Create(model.Contact{UserID: user.ID, Data: m.Email})
	if err != nil {
		a.userRepo.Rollback()
		return JWT{}, err
	}

	err = a.authRepo.Create(repository.AuthTypeEmail, model.AuthStrategy{
		EntityID:    user.ID,
		ContactID:   contact.ID,
		CreatedAt:   time.Now(),
		Type:        string(repository.AuthTypeEmail),
		EntityType:  string(repository.EntityTypeUser),
		Data:        m.Password,
		ContactData: contact.Data,
	})

	if err != nil {
		a.userRepo.Rollback()
		return JWT{}, err
	}

	err = a.contactRepo.Commit()
	if err != nil {
		return JWT{}, err
	}

	// TODO: Now share this with Unit21!!!!!!!!!! EntityCreate(...)

	return a.GenerateJWT(user)

}

func (a auth) LoginEmail(login UserLoginEmail) (JWT, error) {
	_, err := mail.ParseAddress(login.Email)
	if err != nil {
		return JWT{}, errors.Wrap(err, "invalid email")
	}
	m, err := a.authRepo.Get(login.Email)
	if err != nil {
		return JWT{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(m.Data), []byte(login.Password))
	if err != nil {
		return JWT{}, errors.Wrap(err, "invalid or wrong password")
	}

	return a.GenerateJWT(model.User{ID: m.EntityID})
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

func (a auth) LoginPK(login UserPKLogin) (JWT, error) {
	if !hexRegex.MatchString(login.PublicAddress) {
		return JWT{}, errors.New("invalid address")
	}
	if login.Signature == "" {
		return JWT{}, errors.New("invalid signature")
	}
	nonce, err := a.authRepo.GetKeyString(login.PublicAddress)
	if err != nil {
		return JWT{}, err
	}

	recoveredAddr, err := common.RecoverAddress(nonce, login.Signature)
	if login.PublicAddress != recoveredAddr.Hex() {
		return JWT{}, err
	}

	newNonce := uuid.NewString()
	err = a.authRepo.CreateAny(login.PublicAddress, newNonce, time.Minute*10)
	if err != nil {
		return JWT{}, err
	}

	return a.GenerateJWT(model.User{ID: login.PublicAddress})
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

func (a auth) GenerateAPIKey(m model.Platform) error {
	return nil
}

func (a auth) LoginOTP() error {
	return nil
}

func (a auth) RefreshToken(token string) {
	//stra, err := a.authRepo.Get(common.ToSha256(token))
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

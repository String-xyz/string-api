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

type UserRegister = model.UserRegister
type UserLogin = model.UserLogin

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
	LoginEmail(UserLogin) (JWT, error)
	GenerateJWT(model.User) (JWT, error)
	GenerateAPIKey(model.Platform) error
	RefreshToken()
	LoginPK() error
	LoginOTP() error
}

type auth struct {
	authRepo    repository.AuthStrategy
	userRepo    repository.User
	contactRepo repository.UserContact
}

func NewAuth(a repository.AuthStrategy, u repository.User, c repository.UserContact) Auth {
	return &auth{a, u, c}
}

// Register registers an user with authentication (email/password)
// this is a rudimentary implementation of onboarding, will later have proper
// onboarding process.
func (a auth) Register(m UserRegister) (JWT, error) {
	tx := a.userRepo.MustBegin()
	user, err := a.userRepo.Create(model.User{FirstName: m.FirstNname, LastName: m.LastName, Status: "registered", Type: "client"})
	if err != nil {
		a.userRepo.Rollback()
		return JWT{}, err
	}
	a.contactRepo.SetTx(tx)
	contact, err := a.contactRepo.Create(model.Contact{UserID: user.ID, Data: m.Email})
	if err != nil {
		a.contactRepo.Rollback()
		return JWT{}, err
	}
	a.contactRepo.Commit()
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
		return JWT{}, err
	}
	return a.GenerateJWT(user)

}

func (a auth) LoginEmail(login UserLogin) (JWT, error) {
	_, err := mail.ParseAddress(login.Email)
	if err != nil {
		return JWT{}, errors.Wrap(err, "Invalid email")
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

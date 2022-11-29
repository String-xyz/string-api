package service

import (
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/internal/unit21"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type UserRequest = model.UserRequest

type EmailVerification struct {
	Timestamp int64
	Email     string
	UserID    string
}
type EmailLogin struct {
	Timestamp int64
	UserID    string
}

type UserRepos struct {
	Auth         repository.AuthStrategy
	User         repository.User
	Contact      repository.Contact
	Instrument   repository.Instrument
	Device       repository.Device
	UserPlatform repository.UserPlatform
}

type User interface {
	GetStatus(request UserRequest) (model.UserOnboardingStatus, error) // If wallet addr is associated with user, return current state of their onboarding
	Create(request UserRequest) (JWT, error)                           // Receive new wallet addr, email, signature and send verification email
	RequestEmailAuthentication(request UserRequest, userId string) error
	ReceiveEmailAuthentication(encrypted string) error // Decrypts query and validates e-mail, wallet of user
	Name(request UserRequest) error                    // Takes name and wallet addr of user, associates name with wallet addr
	RequestEmailLogin(request UserRequest) error       // Takes wallet addr of user and sends login email
	ReceiveEmailLogin(encrypted string) (JWT, error)   // Decrypts query and logs user in, returning JWT
}

type user struct {
	repos UserRepos
}

func NewUser(repos UserRepos) User {
	return &user{repos: repos}
}

func (u user) GetStatus(request UserRequest) (model.UserOnboardingStatus, error) {
	res := model.UserOnboardingStatus{Status: "not found"}
	addr := request.WalletAddress
	if addr == "" {
		return res, common.StringError(errors.New("no wallet address provided"))
	}
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil {
		return res, common.StringError(err)
	}
	associatedUser, err := u.repos.User.GetById(instrument.UserID)
	if err != nil {
		return res, common.StringError(err)
	}
	if associatedUser.Status != "" {
		res.Status = associatedUser.Status
		return res, nil
	}
	return res, common.StringError(errors.New("not found"))
}

func (u user) Create(request UserRequest) (JWT, error) {
	addr := request.WalletAddress
	if addr == "" {
		return JWT{}, common.StringError(errors.New("no wallet address provided"))
	}
	signature := request.Signature
	if signature == "" {
		return JWT{}, common.StringError(errors.New("no signature provided"))
	}

	// Make sure wallet does not already exist
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return JWT{}, common.StringError(err)
	} else if err == nil && instrument.UserID != "" {
		return JWT{}, common.StringError(errors.New("wallet already associated with user"))
	} else if err == nil && instrument.PublicKey == addr {
		return JWT{}, common.StringError(errors.New("wallet already exists"))
	}

	// Make sure address is a wallet and not a smart contract
	if !common.IsWallet(addr) {
		return JWT{}, common.StringError(errors.New("address provided is not a valid wallet"))
	}

	// Verify signature
	valid, err := common.ValidateExternalEVMSignature(request.Signature, addr, addr) // they signed their own address.
	// it's like writing your name on your hand and then xeroxing it
	if err != nil {
		return JWT{}, common.StringError(err)
	}
	if !valid {
		return JWT{}, common.StringError(errors.New("signature invalid"))
	}

	// Initialize a new user
	user := model.User{Type: "string-user", Status: "unverified"} // Validated status pertains to specific instrument
	user, err = u.repos.User.Create(user)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	// Create a new wallet instrument and associate it with the new user
	instrument = model.Instrument{Type: "crypto-wallet", Status: "verified", Network: "EVM", PublicKey: addr, UserID: user.ID}
	instrument, err = u.repos.Instrument.Create(instrument)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	jwt, err := u.generateJWT(user)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	go u.createUnit21Entity(user)

	return jwt, nil
}

func (u user) RequestEmailAuthentication(request UserRequest, userId string) error {
	email := request.EmailAddress
	if email == "" {
		return common.StringError(errors.New("no email provided"))
	}

	user, err := u.repos.User.GetById(userId)
	if err != nil || user.ID != userId {
		return common.StringError(errors.New("invalid or expired JWT"))
	}

	contact, _ := u.repos.Contact.GetByData(email)
	if contact.Status == "validated" {
		return common.StringError(errors.New("email is already authenticated"))
	}

	// Encrypt required data to Base64 string and insert it in an email hyperlink
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	code, err := common.Encrypt(EmailVerification{Timestamp: time.Now().Unix(), Email: email, UserID: userId}, key)
	if err != nil {
		return common.StringError(err)
	}
	code = url.QueryEscape(code) // make sure special characters are browser friendly

	baseURL := common.GetBaseURL()
	from := mail.NewEmail("String Authentication", "auth@string.xyz")
	subject := "String Email Authentication"
	to := mail.NewEmail("New String User", email)
	textContent := "Click the link below to complete your e-mail authentication!"
	htmlContent := `<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='` + baseURL + `login/email?token=` + code + `' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Verify Email Now</a></div>`

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return common.StringError(err)
	}

	// Wait for up to 15 minutes, final timeout TBD
	now, lastPolled := time.Now().Unix(), time.Now().Unix()
	until := now + (60 * 15)
	for now < until {
		now = time.Now().Unix()
		if now-lastPolled < 3 {
			continue // throttle following logic in 3 second interval
		}
		lastPolled = now
		contact, err := u.repos.Contact.GetByData(email)
		if err != nil && errors.Cause(err).Error() != "not found" {
			return common.StringError(err)
		} else if err == nil && contact.Data == email {
			return nil // success
		}
	}

	// timed out
	return common.StringError(errors.New("link expired"))
}

func (u user) ReceiveEmailAuthentication(encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[EmailVerification](encrypted, key)
	if err != nil {
		return common.StringError(err)
	}

	// Wait for up to 15 minutes, final timeout TBD
	now := time.Now().Unix()
	if now-received.Timestamp > (60 * 15) {
		return common.StringError(errors.New("link expired"))
	}

	contact := model.Contact{UserID: received.UserID, Type: "email", Status: "validated", Data: received.Email}
	contact, _ = u.repos.Contact.Create(contact)

	return nil
}

func (u user) Name(request UserRequest) error {
	addr := request.WalletAddress
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	user, err := u.repos.User.GetById(instrument.UserID)
	if err != nil {
		return common.StringError(err)
	}
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	err = u.repos.User.Update(user.ID, updates)
	if err != nil {
		return common.StringError(err)
	}

	go u.updateUnit21Entity(user)

	return nil
}

func (u user) RequestEmailLogin(request UserRequest) error {
	email := request.EmailAddress
	if email == "" {
		return common.StringError(errors.New("no email address provided"))
	}

	// Ensure email exists in database
	model, err := u.repos.Contact.GetByData(email)
	if err != nil || model.Data != email {
		return nil // As this endpoint requires no authentication, do not provide information on the user such as "not found"
	}

	// Encrypt required data to Base64 string and insert it in an email hyperlink
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	code, err := common.Encrypt(EmailLogin{Timestamp: time.Now().Unix(), UserID: model.UserID}, key)
	if err != nil {
		return common.StringError(err)
	}
	code = url.QueryEscape(code) // make sure special characters are browser friendly

	baseURL := common.GetBaseURL()
	from := mail.NewEmail("String Authentication", "auth@string.xyz")
	subject := "String Email Login"
	to := mail.NewEmail(email, email)
	textContent := "Click the link below to log in!"
	htmlContent := "<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='" + baseURL + "login?token=" + code + "' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Login</a></div>"
	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return common.StringError(err)
	}
	// success
	return nil
}

func (u user) ReceiveEmailLogin(encrypted string) (JWT, error) {
	// encrypted, err := url.QueryUnescape(encrypted)
	// if err != nil {
	// 	return JWT{}, common.StringError(err)
	// }
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[EmailLogin](encrypted, key)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	// Wait for up to 15 minutes, final timeout TBD
	now := time.Now().Unix()
	if now-received.Timestamp > (60 * 15) {
		return JWT{}, common.StringError(errors.New("link expired"))
	}

	user, err := u.repos.User.GetById(received.UserID)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	jwt, err := u.generateJWT(user)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	go u.updateUnit21Entity(user)

	return jwt, nil
}

func (u user) generateJWT(m model.User) (JWT, error) {
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
		return JWT{}, common.StringError(err)
	}
	t.Token = signed
	err = u.repos.Auth.CreateJWTRefresh(common.ToSha256(refreshToken), m.ID)
	if err != nil {
		return JWT{}, common.StringError(err)
	}
	return *t, nil
}

func (u user) createUnit21Entity(user model.User) {
	// Createing a User Entity in Unit21
	u21Repo := unit21.EntityRepos{
		Device:       u.repos.Device,
		Contact:      u.repos.Contact,
		UserPlatform: u.repos.UserPlatform,
	}

	u21Entity := unit21.NewEntity(u21Repo)
	_, err := u21Entity.Create(user)
	if err != nil {
		log.Printf("Error creating Entity in Unit21: %s", err)
	}
}

func (u user) updateUnit21Entity(user model.User) {
	// Createing a User Entity in Unit21
	u21Repo := unit21.EntityRepos{
		Device:       u.repos.Device,
		Contact:      u.repos.Contact,
		UserPlatform: u.repos.UserPlatform,
	}

	u21Entity := unit21.NewEntity(u21Repo)
	_, err := u21Entity.Update(user)
	if err != nil {
		log.Printf("Error updating Entity in Unit21: %s", err)
	}
}

package service

import (
	"fmt"
	netMail "net/mail"
	"os"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
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

type UserRepos struct {
	User       repository.User
	Contact    repository.Contact
	Instrument repository.Instrument
}

type User interface {
	GetStatus(request UserRequest) (model.UserOnboardingStatus, error) // If wallet addr is associated with user, return current state of their onboarding
	Create(request UserRequest) error                                  // Create new user using wallet addr, optionally mark as validated if signature is provided
	Sign(request UserRequest) error                                    // Takes in a signed timestamp from user, validating their wallet
	Authenticate(request UserRequest) error                            // Takes e-mail and wallet addr of user, validates email with twilio
	ReceiveEmailAuthentication(encrypted string) error                 // Decrypts query and validates e-mail address of user
	Name(request UserRequest) error                                    // Takes name and wallet addr of user, associates name with wallet addr
}

type user struct {
	repos UserRepos
}

func NewUser(repos UserRepos) User {
	return &user{repos: repos}
}

func (u user) GetStatus(request UserRequest) (model.UserOnboardingStatus, error) {
	res := model.UserOnboardingStatus{Status: "Not Found"}
	addr := request.WalletAddress
	if addr == "" {
		return res, common.StringError(errors.New("no wallet address provided"))
	}
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil {
		return res, common.StringError(err)
	}
	associatedUser, err := u.repos.User.GetID(instrument.UserID)
	if err != nil {
		return res, common.StringError(err)
	}
	res.Status = associatedUser.Status
	return res, nil
}

func (u user) Create(request UserRequest) error {
	addr := request.WalletAddress
	if addr == "" {
		return common.StringError(errors.New("no wallet address provided"))
	}

	// Make sure wallet does not already exist
	// TODO: Revisit this logic.  Parsing error string for "not found" is bad.
	instrument, err := u.repos.Instrument.GetWallet(addr)
	fmt.Printf("\nINSTRUMENT=%+v", instrument)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return common.StringError(err)
	} else if err == nil && instrument.UserID != "" {
		return common.StringError(errors.New("wallet already associated with user"))
	} else if err == nil && instrument.PublicKey == addr {
		return common.StringError(errors.New("wallet already exists"))
	}

	// Make sure address is a wallet and not a smart contract
	if !common.IsWallet(addr) {
		return common.StringError(errors.New("address provided is not a valid wallet"))
	}

	// Optionally verify signature
	status := "Created"
	if request.Signature != "" {
		valid, err := common.ValidateExternalEVMSignature(request.Signature, addr, addr) // they signed their own address.
		// it's like writing your name on your hand and then xeroxing it
		if err != nil {
			return common.StringError(err)
		}
		if !valid {
			return common.StringError(errors.New("signature invalid"))
		}
		status = "Validated"
	}

	// Initialize a new user
	user := model.User{Type: "String User", Status: "Created"} // Validated status pertains to specific instrument
	user, err = u.repos.User.Create(user)
	if err != nil {
		return common.StringError(err)
	}

	// Create a new wallet instrument and associate it with the new user
	instrument = model.Instrument{Type: "Crypto Wallet", Status: status, Network: "EVM", PublicKey: addr, UserID: user.ID}
	instrument, err = u.repos.Instrument.Create(instrument)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (u user) Sign(request UserRequest) error {
	addr := request.WalletAddress
	// Make sure wallet exists
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	// Make sure wallet is associated with a user
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	_, err = u.repos.User.GetID(instrument.UserID) // Don't need to update user, just verify they exist
	if err != nil {
		return common.StringError(err)
	}
	if instrument.Status == "Validated" {
		return common.StringError(errors.New("wallet already validated"))
	}

	// TESTING
	// signed, _ := common.EVMSign(addr)
	// fmt.Printf("\n\nSECRET SIGNATURE=%+v", signed)

	// Verify signature
	valid, err := common.ValidateExternalEVMSignature(request.Signature, addr, addr)
	if err != nil {
		return common.StringError(err)
	}
	if !valid {
		return common.StringError(errors.New("signature invalid"))
	}
	validated := "Validated"
	status := model.UpdateStatus{Status: &validated}
	err = u.repos.Instrument.Update(instrument.ID, status)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (u user) Authenticate(request UserRequest) error {
	addr := request.WalletAddress
	email := request.EmailAddress
	if addr == "" || email == "" {
		return common.StringError(errors.New("missing wallet/email"))
	}
	_, err := netMail.ParseAddress(email)
	if err != nil {
		return common.StringError(err)
	}

	// Make sure wallet exists
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	// Make sure wallet is associated with a user
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	// user, err := u.repos.User.GetID(instrument.UserID)
	// if err != nil {
	// 	return common.StringError(err)
	// }

	// twilio / sendgrid stuff
	// key := os.Getenv("STRING_ENCRYPTION_KEY")
	// expectedResponse, err := common.Encrypt(EmailVerification{Timestamp: time.Now().Unix(), Email: email, UserID: instrument.UserID}, key)
	// if err != nil {
	// 	return common.StringError(err)
	// }
	// var AUTH_SID = os.Getenv("TWILIO_AUTH_SID")
	// client := twilio.NewRestClient()
	// params := &verify.CreateVerificationParams{}
	// params.SetTo(email)
	// params.SetChannel("email")
	// params.SetCustomCode(expectedResponse) // Code is needed to check if email was verified
	// _, err = client.VerifyV2.CreateVerification(AUTH_SID, params)
	// if err != nil {
	// 	return common.StringError(err)
	// }
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	code, err := common.Encrypt(EmailVerification{Timestamp: time.Now().Unix(), Email: email, UserID: instrument.UserID}, key)
	if err != nil {
		return common.StringError(err)
	}
	from := mail.NewEmail("String Authentication", "auth@string.xyz")
	subject := "String Email Authentication"
	to := mail.NewEmail("New String User", email)
	textContent := "Click the link below to complete your e-mail authentication!"
	htmlContent := "<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='http://localhost:5555/user/email?token=" + code + "' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Verify Email Now</a></div>"
	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return common.StringError(err)
	}
	fmt.Printf("\nResponse Status Code = %+v", response.StatusCode)
	fmt.Printf("\nResponse Body = %+v", response.Body)
	fmt.Printf("\nResponse Headers = %+v", response.Headers)
	// success
	return nil
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
	user, err := u.repos.User.GetID(instrument.UserID)
	if err != nil {
		return common.StringError(err)
	}
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	err = u.repos.User.Update(user.ID, updates)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

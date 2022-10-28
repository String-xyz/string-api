package service

import (
	"math"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

type UserRequest = model.UserRequest

type User interface {
	GetStatus(request UserRequest) (model.UserOnboardingStatus, error) // If wallet addr is associated with user, return current state of their onboarding
	Create(request UserRequest) error                                  // Create new user using wallet addr, optionally mark as validated if signature is provided
	Sign(request UserRequest) error                                    // Takes in a signed timestamp from user, validating their wallet
	Authenticate(request UserRequest) error                            // Takes e-mail and wallet addr of user, validates email with twilio
	Name(request UserRequest) error                                    // Takes name and wallet addr of user, associates name with wallet addr
}

type user struct {
	userRepo       repository.User
	contactRepo    repository.Contact
	instrumentRepo repository.Instrument
}

func NewUser(u repository.User, c repository.Contact, i repository.Instrument) User {
	return &user{userRepo: u, contactRepo: c, instrumentRepo: i}
}

func (u user) GetStatus(request UserRequest) (model.UserOnboardingStatus, error) {
	res := model.UserOnboardingStatus{Status: "Not Found"}
	addr := request.WalletAddress
	if addr == "" {
		return res, common.StringError(errors.New("no wallet address provided"))
	}
	instrument, err := u.instrumentRepo.GetWallet(addr)
	if err != nil {
		return res, common.StringError(err)
	}
	associatedUser, err := u.userRepo.GetID(instrument.UserID)
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
	instrument, err := u.instrumentRepo.GetWallet(addr)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return common.StringError(err)
	} else if err == nil && instrument.UserID != "" {
		return common.StringError(errors.New("wallet already associated with user"))
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
	user, err = u.userRepo.Create(user)
	if err != nil {
		return common.StringError(err)
	}

	// Create a new wallet instrument and associate it with the new user
	instrument = model.Instrument{Type: "Crypto Wallet", Status: status, Network: "EVM", PublicKey: addr, UserID: user.ID}
	instrument, err = u.instrumentRepo.Create(instrument)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func (u user) Sign(request UserRequest) error {
	addr := request.WalletAddress
	// Make sure wallet exists
	instrument, err := u.instrumentRepo.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	// Make sure wallet is associated with a user
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	_, err = u.userRepo.GetID(instrument.UserID) // Don't need to update user, just verify they exist
	if err != nil {
		return common.StringError(err)
	}
	if instrument.Status == "Validated" {
		return common.StringError(errors.New("wallet already validated"))
	}
	// Verify signature
	valid, err := common.ValidateExternalEVMSignature(request.Signature, addr, addr)
	if err != nil {
		return common.StringError(err)
	}
	if !valid {
		return common.StringError(errors.New("signature invalid"))
	}
	status := model.UpdateStatus{Status: "Validated"}
	err = u.instrumentRepo.Update(instrument.ID, status)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func randomNumericString(digits int) string {
	randomNumber := time.Now().Nanosecond()
	randomString := ""
	for i := digits - 1; i > -1; i++ {
		// iterating through digits prevents truncating leading 0s from string
		randomString += strconv.Itoa(randomNumber & int(math.Pow10(i)))
	}
	return randomString
}

func (u user) Authenticate(request UserRequest) error {
	addr := request.WalletAddress
	email := request.EmailAddress
	if addr == "" || email == "" {
		return common.StringError(errors.New("missing wallet/email"))
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return common.StringError(err)
	}

	// Make sure wallet exists
	instrument, err := u.instrumentRepo.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	// Make sure wallet is associated with a user
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	user, err := u.userRepo.GetID(instrument.UserID)
	if err != nil {
		return common.StringError(err)
	}

	// twilio / sendgrid stuff
	var AUTH_SID = os.Getenv("TWILIO_AUTH_SID")
	client := twilio.NewRestClient()
	params := &verify.CreateVerificationParams{}
	params.SetTo(email)
	params.SetChannel("email")
	code := randomNumericString(6)
	params.SetCustomCode(code) // Code is needed to check if email was verified
	_, err = client.VerifyV2.CreateVerification(AUTH_SID, params)
	if err != nil {
		return common.StringError(err)
	}

	go u.waitForEmailAuthentication(addr, email, user.ID, code)
	return nil
}

func (u user) waitForEmailAuthentication(addr string, email string, userId string, code string) {
	var AUTH_SID = os.Getenv("TWILIO_AUTH_SID")
	client := twilio.NewRestClient()
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(email)
	params.SetCode(code)

	// Wait for up to 15 minutes, final timeout TBD
	start := time.Now().Unix()
	end := start + (60 * 15)
	currentTime := start
	lastChecked := start
	for currentTime < end {
		currentTime = time.Now().Unix()
		if currentTime > lastChecked {
			resp, _ := client.VerifyV2.CreateVerificationCheck(AUTH_SID, params)
			lastChecked = currentTime
			if resp.Status != nil && *resp.Status == "approved" && resp.Valid != nil && *resp.Valid {
				// Only create email entry in table once it's validated
				contact := model.Contact{UserID: userId, Type: "email", Status: "validated", Data: email}
				contact, _ = u.contactRepo.Create(contact)
				return
			}
		}
	}

}

func (u user) Name(request UserRequest) error {
	addr := request.WalletAddress
	instrument, err := u.instrumentRepo.GetWallet(addr)
	if err != nil {
		return common.StringError(err)
	}
	if instrument.UserID == "" {
		return common.StringError(errors.New("wallet not associated with user"))
	}
	user, err := u.userRepo.GetID(instrument.UserID)
	if err != nil {
		return common.StringError(err)
	}
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	err = u.userRepo.Update(user.ID, updates)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

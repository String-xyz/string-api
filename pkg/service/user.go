package service

import (
	"os"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
)

type UserRequest = model.UserRequest
type UserUpdates = model.UpdateUserName

type UserCreateResponse struct {
	JWT  JWT        `json:"authToken"`
	User model.User `json:"user"`
}

type User interface {
	//GetStatus returns the onboarding status of an user
	GetStatus(userID string) (model.UserOnboardingStatus, error)

	// Create creates an user from a wallet signed payload
	// It associates the wallet to the user and also sets its status as verified
	// This payload usually comes from a previous requested one using (Auth.PayloadToSign) service
	Create(request model.WalletSignaturePayloadSigned) (UserCreateResponse, error)

	//Update updates the user firstname lastname middlename.
	// It fetches the user using the walletAddress provided
	Update(userID string, request UserUpdates) (model.User, error)
}

type user struct {
	repos       repository.Repositories
	auth        Auth
	fingerprint Fingerprint
	device      Device
	unit21      Unit21
}

func NewUser(repos repository.Repositories, auth Auth, fprint Fingerprint, device Device, unit21 Unit21) User {
	return &user{repos, auth, fprint, device, unit21}
}

func (u user) GetStatus(userID string) (model.UserOnboardingStatus, error) {
	res := model.UserOnboardingStatus{Status: "not found"}

	user, err := u.repos.User.GetById(userID)
	if err != nil {
		return res, common.StringError(err)
	}

	if user.Status != "" {
		res.Status = user.Status
		return res, nil
	}
	return res, common.StringError(errors.New("not found"))
}

func (u user) Create(request model.WalletSignaturePayloadSigned) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	payload, err := common.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return resp, common.StringError(err)
	}

	addr := payload.Address
	if addr == "" {
		return resp, common.StringError(errors.New("no wallet address provided"))
	}

	// Make sure wallet does not already exist
	exists, err := u.repos.Instrument.WalletAlreadyExists(addr)
	if err != nil {
		return resp, common.StringError(err)
	}

	if exists {
		return resp, common.StringError(errors.New("wallet already exists"))
	}

	// Make sure address is a wallet and not a smart contract
	if !common.IsWallet(addr) {
		return resp, common.StringError(errors.New("address provided is not a valid wallet"))
	}

	// Verify payload integrity
	if err := verifyWalletAuthentication(request); err != nil {
		return resp, common.StringError(err)
	}

	user, err := u.createUserData(addr)
	if err != nil {
		return resp, err
	}

	// create device only if there is a visitor
	device, err := u.device.CreateDeviceIfNeeded(user.ID, request.Fingerprint.VisitorID, request.Fingerprint.RequestID)

	if err != nil && errors.Cause(err).Error() != "not found" {
		return resp, common.StringError(err)
	}

	jwt, err := u.auth.GenerateJWT(user.ID, device)
	if err != nil {
		return resp, common.StringError(err)
	}

	// deviceService.RegisterNewUserDevice()
	go u.unit21.Entity.Create(user)

	return UserCreateResponse{JWT: jwt, User: user}, nil
}

func (u user) createUserData(addr string) (model.User, error) {
	tx := u.repos.User.MustBegin()
	u.repos.Instrument.SetTx(tx)
	u.repos.Device.SetTx(tx)

	defer u.repos.User.Reset(u.repos.Instrument, u.repos.Device)
	// Initialize a new user
	// Validated status pertains to specific instrument
	user := model.User{Type: "string-user", Status: "unverified"}
	user, err := u.repos.User.Create(user)
	if err != nil {
		u.repos.User.Rollback()
		return user, common.StringError(err)
	}
	// Create a new wallet instrument and associate it with the new user
	instrument := model.Instrument{Type: "Crypto Wallet", Status: "verified", Network: "EVM", PublicKey: addr, UserID: user.ID}
	instrument, err = u.repos.Instrument.Create(instrument)
	if err != nil {
		u.repos.Instrument.Rollback()
		return user, common.StringError(err)
	}
	if err := u.repos.User.Commit(); err != nil {
		return user, common.StringError(errors.New("error commiting transaction"))
	}

	go u.unit21.Instrument.Create(instrument)

	return user, nil
}

func (u user) Update(userID string, request UserUpdates) (model.User, error) {
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	user, err := u.repos.User.Update(userID, updates)
	if err != nil {
		return user, common.StringError(err)
	}

	go u.unit21.Entity.Update(user)

	return user, nil
}

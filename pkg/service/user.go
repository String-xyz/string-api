package service

import (
	"strings"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/internal/unit21"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type UserRequest = model.UserRequest
type UserUpdates = model.UpdateUserName

type EmailVerification struct {
	Timestamp int64
	Email     string
	UserID    string
}

type UserCreateResponse struct {
	JWT  JWT        `json:"authToken"`
	User model.User `json:"user"`
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
	//GetStatus returns the onboarding status of an user
	GetStatus(userID string) (model.UserOnboardingStatus, error)

	// Create creates an user from a wallet signed payload
	// It associates the wallet to the user and also sets its status as verified
	// This payload usually comes from a previous requested one using (Auth.PayloadToSign) service
	Create(request model.WalletSignaturePayload) (UserCreateResponse, error)

	//Update updates the user firstname lastname middlename.
	// It fetches the user using the walletAddress provided
	Update(userID string, request UserUpdates) (model.User, error)
}

type user struct {
	repos UserRepos
}

func NewUser(repos UserRepos) User {
	return &user{repos: repos}
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

func (u user) Create(request model.WalletSignaturePayload) (UserCreateResponse, error) {
	addr := request.Address
	resp := UserCreateResponse{}
	if addr == "" {
		return resp, common.StringError(errors.New("no wallet address provided"))
	}

	// Make sure wallet does not already exist
	instrument, err := u.repos.Instrument.GetWallet(addr)
	if err != nil && !strings.Contains(err.Error(), "not found") { // because we are wrapping error and care about its value
		return resp, common.StringError(err)
	} else if err == nil && instrument.UserID != "" {
		return resp, common.StringError(errors.New("wallet already associated with user"))
	} else if err == nil && instrument.PublicKey == addr {
		return resp, common.StringError(errors.New("wallet already exists"))
	}

	// Make sure address is a wallet and not a smart contract
	if !common.IsWallet(addr) {
		return resp, common.StringError(errors.New("address provided is not a valid wallet"))
	}

	// Verify payload integrity
	err = verifyWalletAuthentication(request)
	if err != nil {
		return resp, common.StringError(err)
	}

	// Initialize a new user
	user := model.User{Type: "string-user", Status: "unverified"} // Validated status pertains to specific instrument
	user, err = u.repos.User.Create(user)
	if err != nil {
		return resp, common.StringError(err)
	}

	// Create a new wallet instrument and associate it with the new user
	instrument = model.Instrument{Type: "crypto-wallet", Status: "verified", Network: "EVM", PublicKey: addr, UserID: user.ID}
	instrument, err = u.repos.Instrument.Create(instrument)
	if err != nil {
		return resp, common.StringError(err)
	}

	jwt, err := NewAuth(u.repos).GenerateJWT(user)
	if err != nil {
		return resp, common.StringError(err)
	}

	go u.createUnit21Entity(user)

	return UserCreateResponse{JWT: jwt, User: user}, nil
}

func (u user) Update(userID string, request UserUpdates) (model.User, error) {
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	user, err := u.repos.User.Update(userID, updates)
	if err != nil {
		return user, common.StringError(err)
	}

	go u.updateUnit21Entity(user)

	return user, nil
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
		log.Err(err).Msg("Error creating Entity in Unit21")
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
		log.Err(err).Msg("Error updating Entity in Unit21")
	}
}

package service

import (
	"context"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/env"
	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"

	"github.com/rs/zerolog/log"
)

type UserRequest = model.UserRequest
type UserUpdates = model.UpdateUserName

type UserCreateResponse struct {
	JWT  JWT        `json:"authToken"`
	User model.User `json:"user"`
}

type User interface {
	//GetStatus returns the onboarding status of an user
	GetStatus(ctx context.Context, userId string) (model.UserOnboardingStatus, error)

	// Create creates an user from a wallet signed payload
	// It associates the wallet to the user and also sets its status as verified
	// This payload usually comes from a previous requested one using (Auth.PayloadToSign) service
	Create(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string) (UserCreateResponse, error)

	//Update updates the user firstname lastname middlename.
	// It fetches the user using the walletAddress provided
	Update(ctx context.Context, userId string, request UserUpdates) (model.User, error)
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

func (u user) GetStatus(ctx context.Context, userId string) (model.UserOnboardingStatus, error) {
	res := model.UserOnboardingStatus{Status: "not found"}

	user, err := u.repos.User.GetById(ctx, userId)
	if err != nil {
		return res, libcommon.StringError(err)
	}

	if user.Status != "" {
		res.Status = user.Status
		return res, nil
	}
	return res, libcommon.StringError(serror.NOT_FOUND)
}

func (u user) Create(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	key, err := env.Get("STRING_ENCRYPTION_KEY")
	if err != nil {
		return resp, libcommon.StringError(err)
	}
	payload, err := libcommon.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// Make sure address is a wallet and not a smart contract
	addr := payload.Address
	if addr == "" || !common.IsWallet(addr) {
		return resp, libcommon.StringError(serror.INVALID_DATA)
	}

	// Make sure wallet does not already exist
	exists, err := u.repos.Instrument.WalletAlreadyExists(ctx, addr)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	if exists {
		return resp, libcommon.StringError(serror.ALREADY_IN_USE)
	}

	// Verify payload integrity
	if err := verifyWalletAuthentication(request); err != nil {
		return resp, libcommon.StringError(err)
	}

	user, err := u.createUserData(ctx, addr)
	if err != nil {
		return resp, err
	}

	// Associate user to platform
	err = u.repos.Platform.AssociateUser(ctx, user.Id, platformId)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// create device only if there is a visitor
	device, err := u.device.CreateDeviceIfNeeded(ctx, user.Id, request.Fingerprint.VisitorId, request.Fingerprint.RequestId)
	if err != nil && serror.Is(err, serror.NOT_FOUND) {
		return resp, libcommon.StringError(err)
	}

	if device.Fingerprint != "" {
		// validate that device on user creation
		now := time.Now()
		err = u.repos.Device.Update(ctx, device.Id, model.DeviceUpdates{ValidatedAt: &now})
		if err == nil {
			log.Err(err).Msg("Failed to verify user device")
		}
	}

	jwt, err := u.auth.GenerateJWT(user.Id, platformId, device)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// deviceService.RegisterNewUserDevice()
	// Create a new context since this will run in background
	ctx2 := context.Background()
	go u.unit21.Entity.Create(ctx2, user)

	return UserCreateResponse{JWT: jwt, User: user}, nil
}

func (u user) createUserData(ctx context.Context, addr string) (model.User, error) {
	tx := u.repos.User.MustBegin()
	u.repos.Instrument.SetTx(tx)
	u.repos.Device.SetTx(tx)

	defer u.repos.User.Reset(u.repos.Instrument, u.repos.Device)
	// Initialize a new user
	// Validated status pertains to specific instrument
	user := model.User{Type: "string-user", Status: "unverified"}
	user, err := u.repos.User.Create(ctx, user)
	if err != nil {
		u.repos.User.Rollback()
		return user, libcommon.StringError(err)
	}

	// Create a new wallet instrument and associate it with the new user
	instrument := model.Instrument{Type: "crypto wallet", Status: "verified", Network: "EVM", PublicKey: addr, UserId: user.Id}
	instrument, err = u.repos.Instrument.Create(ctx, instrument)
	if err != nil {
		u.repos.Instrument.Rollback()
		return user, libcommon.StringError(err)
	}
	if err := u.repos.User.Commit(); err != nil {
		return user, libcommon.StringError(err)
	}

	// Create a new context since this will run in background
	ctx2 := context.Background()
	go u.unit21.Instrument.Create(ctx2, instrument)

	return user, nil
}

func (u user) Update(ctx context.Context, userId string, request UserUpdates) (model.User, error) {
	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	user, err := u.repos.User.Update(ctx, userId, updates)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	// Create a new context since this will run in background
	ctx2 := context.Background()
	go u.unit21.Entity.Update(ctx2, user)

	return user, nil
}

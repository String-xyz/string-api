package service

import (
	"context"
	"strings"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"

	"github.com/rs/zerolog/log"
)

type UserRequest = model.UserRequest
type UserUpdates = model.UpdateUserName

type User interface {
	// GetStatus returns the onboarding status of a user
	GetStatus(ctx context.Context, userId string) (model.UserOnboardingStatus, error)

	// Create creates an user from a wallet signed payload
	// It associates the wallet to the user and also sets its status as verified
	// This payload usually comes from a previous requested one using (Auth.PayloadToSign) service
	Create(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string) (resp model.UserLoginResponse, err error)

	// Update updates the user's name fields (firstname, lastname, middlename).
	// It fetches the user using the walletAddress provided
	Update(ctx context.Context, userId string, request UserUpdates) (model.User, error)

	// GetUserByLoginPayload get the user without actually logging them in
	GetUserByLoginPayload(ctx context.Context, request model.WalletSignaturePayloadSigned) (user model.User, err error)

	// PreviewEmail returns a partially obfuscated version of the user's email address
	PreviewEmail(ctx context.Context, request model.WalletSignaturePayloadSigned) (email model.EmailPreview, err error)

	// RequestDeviceVerification sends verify device email without needing user id
	RequestDeviceVerification(ctx context.Context, request model.WalletSignaturePayloadSigned) (err error)

	// GetDeviceStatus checks the status of the device verification
	GetDeviceStatus(ctx context.Context, request model.WalletSignaturePayloadSigned) (model.UserOnboardingStatus, error)
}

type user struct {
	repos        repository.Repositories
	auth         Auth
	fingerprint  Fingerprint
	device       Device
	unit21       Unit21
	verification Verification
}

func NewUser(repos repository.Repositories, auth Auth, fprint Fingerprint, device Device, unit21 Unit21, verificationSrv Verification) User {
	return &user{repos, auth, fprint, device, unit21, verificationSrv}
}

func (u user) GetStatus(ctx context.Context, userId string) (model.UserOnboardingStatus, error) {
	_, finish := Span(ctx, "service.user.GetStatus")
	defer finish()

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

func (u user) Create(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string) (resp model.UserLoginResponse, err error) {
	_, finish := Span(ctx, "service.user.Create", SpanTag{"platformId": platformId})
	defer finish()

	// Verify payload integrity
	payload, err := verifyWalletAuthentication(request)
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

	return model.UserLoginResponse{JWT: jwt, User: user}, nil
}

func (u user) createUserData(ctx context.Context, addr string) (model.User, error) {
	_, finish := Span(ctx, "service.user.createUserData")
	defer finish()

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
	_, finish := Span(ctx, "service.user.Update")
	defer finish()

	updates := model.UpdateUserName{FirstName: request.FirstName, MiddleName: request.MiddleName, LastName: request.LastName}
	user, err := u.repos.User.Update(ctx, userId, updates)
	if err != nil {
		return user, libcommon.StringError(err)
	}
	// Create customer on checkout so we can use it when processing payments
	go func(user model.User, email string, platformId string) {
		createCustomer(user, email, platformId)
	}(user, "", "")

	// Create a new context since this will run in background
	ctx2 := context.Background()
	go u.unit21.Entity.Update(ctx2, user)

	return user, nil
}

func (u user) GetUserByLoginPayload(ctx context.Context, request model.WalletSignaturePayloadSigned) (user model.User, err error) {
	_, finish := Span(ctx, "service.device.GetUserByLoginPayload")
	defer finish()

	// Get wallet address from payload
	payload, err := verifyWalletAuthentication(request)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	// Verify there is a user registered to this wallet address
	instrument, err := u.repos.Instrument.GetWalletByAddr(ctx, payload.Address)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	user, err = u.repos.User.GetById(ctx, instrument.UserId)
	if err != nil {
		return user, libcommon.StringError(err)
	}

	return user, nil
}

func (u user) RequestDeviceVerification(ctx context.Context, request model.WalletSignaturePayloadSigned) error {
	_, finish := Span(ctx, "service.user.RequestDeviceVerification")
	defer finish()

	user, err := u.GetUserByLoginPayload(ctx, request)
	if err != nil {
		return libcommon.StringError(err)
	}

	user.Email = getValidatedEmailOrEmpty(ctx, u.repos.Contact, user.Id)

	if user.Email == "" {
		return libcommon.StringError(serror.NOT_FOUND)
	}

	device, err := u.device.CreateDeviceIfNeeded(ctx, user.Id, request.Fingerprint.VisitorId, request.Fingerprint.RequestId)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return libcommon.StringError(err)
	}

	if !isDeviceValidated(device) {
		u.verification.SendDeviceVerification(ctx, user.Id, user.Email, device.Id, device.Description)
	}

	return nil
}

func (u user) GetDeviceStatus(ctx context.Context, request model.WalletSignaturePayloadSigned) (model.UserOnboardingStatus, error) {
	_, finish := Span(ctx, "service.user.GetDeviceStatus")
	defer finish()

	resp := model.UserOnboardingStatus{Status: "unverified"}

	user, err := u.GetUserByLoginPayload(ctx, request)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	device, err := u.device.CreateDeviceIfNeeded(ctx, user.Id, request.Fingerprint.VisitorId, request.Fingerprint.RequestId)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return resp, libcommon.StringError(err)
	}

	if !isDeviceValidated(device) {
		return resp, nil
	}

	resp.Status = "verified"

	return resp, nil
}

func (u user) PreviewEmail(ctx context.Context, request model.WalletSignaturePayloadSigned) (email model.EmailPreview, err error) {
	_, finish := Span(ctx, "service.user.PreviewEmail")
	defer finish()

	user, err := u.GetUserByLoginPayload(ctx, request)
	if err != nil {
		return email, libcommon.StringError(err)
	}

	user.Email = getValidatedEmailOrEmpty(ctx, u.repos.Contact, user.Id)

	if user.Email == "" {
		return email, libcommon.StringError(serror.NOT_FOUND)
	}

	// Partially obfuscate email address
	// Ex. an****@g***l.com
	// This could possibly be done as a regex, but also readability is important
	address := strings.Split(user.Email, "@")

	// Allow for .co.uk, .co.jp, etc to be counted in the extension
	domain := strings.SplitN(address[1], ".", 2)

	ext := domain[1]

	// Allow for single character addresses, if those exist
	first_name_char := ""
	if len(address[0]) >= 2 {
		first_name_char = address[0][0:2]
	} else if len(address[0]) == 1 {
		first_name_char = address[0][0:1]
	}

	// If the name is 2 characters or less, add two stars minimum
	num_name_stars := len(address[0]) - 2
	if num_name_stars <= 0 {
		num_name_stars = 2
	}

	name_stars := strings.Repeat("*", num_name_stars)

	num_domain_stars := len(domain[0]) - 2
	if num_domain_stars <= 0 {
		num_domain_stars = 2
	}

	domain_stars := strings.Repeat("*", num_domain_stars)

	first_domain_char := string(domain[0][0])
	last_domain_char := string(domain[0][len(domain[0])-1])

	obfs_domain := first_domain_char + domain_stars + last_domain_char + "." + ext

	email.Email = first_name_char + name_stars + "@" + obfs_domain

	return email, nil
}

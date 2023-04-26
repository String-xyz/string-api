package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	libcommon "github.com/String-xyz/go-lib/common"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/string-api/env"
	"github.com/String-xyz/string-api/pkg/internal/common"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type SignablePayload struct {
	Nonce string `json:"nonce"`
}

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

var walletAuthenticationPrefix string = "Thank you for using String! By signing this message you are:\n\n1) Authorizing String to initiate off-chain transactions on your behalf, including your bank account, credit card, or debit card.\n\n2) Confirming that this wallet is owned by you.\n\nThis request will not trigger any blockchain transaction or cost any gas.\n\nNonce: "

type RefreshTokenResponse struct {
	Token string    `json:"token"`
	ExpAt time.Time `json:"expAt"`
}

type JWT struct {
	ExpAt        time.Time            `json:"expAt"`
	IssuedAt     time.Time            `json:"issuedAt"`
	Token        string               `json:"token"`
	RefreshToken RefreshTokenResponse `json:"refreshToken"`
}

type JWTClaims struct {
	UserId     string `json:"userId"`
	PlatformId string `json:"platformId"`
	DeviceId   string `json:"deviceId"`
	jwt.StandardClaims
}

type Auth interface {
	// PayloadToSign returns a payload to be sign by a wallet
	// to authenticate an user, the payload expires in 15 minutes
	PayloadToSign(walletAddress string) (SignablePayload, error)

	// VerifySignedPayload receives a signed payload from the user and verifies the signature
	// if signature is valid it returns a JWT to authenticate the user
	VerifySignedPayload(ctx context.Context, signature model.WalletSignaturePayloadSigned, platformId string, bypassDevice bool) (UserCreateResponse, error)

	GenerateJWT(string, string, ...model.Device) (JWT, error)
	ValidateAPIKeyPublic(key string) (string, error)
	ValidateAPIKeySecret(key string) (string, error)
	RefreshToken(ctx context.Context, token string, walletAddress string, platformId string) (UserCreateResponse, error)
	InvalidateRefreshToken(token string) error
}

type auth struct {
	repos        repository.Repositories
	verification Verification
	device       Device
}

// reusing UserRepos here
func NewAuth(r repository.Repositories, v Verification, d Device) Auth {
	return &auth{r, v, d}
}

func (a auth) PayloadToSign(walletAddress string) (SignablePayload, error) {
	payload := model.WalletSignaturePayload{}
	signable := SignablePayload{}

	if !hexRegex.MatchString(walletAddress) {
		return signable, libcommon.StringError(errors.New("missing or invalid address"))
	}
	payload.Address = walletAddress
	payload.Timestamp = time.Now().Unix()
	key := env.Var.STRING_ENCRYPTION_KEY
	encrypted, err := libcommon.Encrypt(payload, key)
	if err != nil {
		return signable, libcommon.StringError(err)
	}
	return SignablePayload{walletAuthenticationPrefix + encrypted}, nil
}

func (a auth) VerifySignedPayload(ctx context.Context, request model.WalletSignaturePayloadSigned, platformId string, bypassDevice bool) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	key := env.Var.STRING_ENCRYPTION_KEY
	payload, err := libcommon.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	if err := verifyWalletAuthentication(request); err != nil {
		return resp, libcommon.StringError(err)
	}

	// Verify user is registered to this wallet address
	instrument, err := a.repos.Instrument.GetWalletByAddr(ctx, payload.Address)
	if err != nil {
		return resp, libcommon.StringError(err)
	}
	user, err := a.repos.User.GetById(ctx, instrument.UserId)
	if err != nil {
		return resp, libcommon.StringError(err)
	}
	// TODO: remove user.Email and replace with association with contact via user and platform
	user.Email = getValidatedEmailOrEmpty(ctx, a.repos.Contact, user.Id)

	device, err := a.device.CreateDeviceIfNeeded(ctx, user.Id, request.Fingerprint.VisitorId, request.Fingerprint.RequestId)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return resp, libcommon.StringError(err)
	}

	// Send verification email if device is unknown and user has a validated email
	// and if verification is not bypassed
	if !bypassDevice && user.Email != "" && !isDeviceValidated(device) {
		go a.verification.SendDeviceVerification(user.Id, user.Email, device.Id, device.Description)
		return resp, libcommon.StringError(serror.UNKNOWN_DEVICE)
	}

	// Create the JWT
	jwt, err := a.GenerateJWT(user.Id, platformId, device)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// Invalidate device if it is unknown and was validated so it cannot be used again
	err = a.device.InvalidateUnknownDevice(ctx, device)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	return UserCreateResponse{JWT: jwt, User: user}, nil
}

// GenerateJWT generates a jwt token and a refresh token which is saved on redis
func (a auth) GenerateJWT(userId string, platformId string, m ...model.Device) (JWT, error) {
	claims := JWTClaims{}
	refreshToken := uuidWithoutHyphens()
	t := &JWT{
		IssuedAt: time.Now(),
		ExpAt:    time.Now().Add(time.Minute * 15),
	}

	// set device id if available
	if len(m) > 0 {
		claims.DeviceId = m[0].Id
	}

	claims.UserId = userId
	claims.PlatformId = platformId
	claims.ExpiresAt = t.ExpAt.Unix()
	claims.IssuedAt = t.IssuedAt.Unix()
	// replace this signing method with RSA or something similar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(env.Var.JWT_SECRET_KEY))
	if err != nil {
		return *t, err
	}
	t.Token = signed

	// create and save
	refreshObj, err := a.repos.Auth.CreateJWTRefresh(libcommon.ToSha256(refreshToken), userId)
	if err != nil {
		return *t, err
	}
	t.RefreshToken = RefreshTokenResponse{
		Token: refreshToken,
		ExpAt: refreshObj.ExpiresAt,
	}

	return *t, nil
}

func (a auth) ValidateJWT(token string) (bool, error) {
	var claims = &JWTClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.Var.JWT_SECRET_KEY), nil
	})
	return t.Valid, err
}

func (a auth) ValidateAPIKeyPublic(key string) (string, error) {
	ctx := context.Background()
	authKey, err := a.repos.Apikey.GetByData(ctx, key, "public")
	if err != nil {
		return "", libcommon.StringError(err)
	}

	if authKey.Id == "" {
		return "", libcommon.StringError(errors.New("invalid api key"))
	}

	if authKey.Data != key {
		return "", libcommon.StringError(errors.New("invalid api key"))
	}

	return *authKey.PlatformId, nil
}

func (a auth) ValidateAPIKeySecret(key string) (string, error) {
	ctx := context.Background()

	data := libcommon.ToSha256(key)
	authKey, err := a.repos.Apikey.GetByData(ctx, data, "secret")
	if err != nil {
		return "", libcommon.StringError(err)
	}

	if authKey.Id == "" {
		return "", libcommon.StringError(errors.New("invalid secret key"))
	}

	if authKey.Data != data {
		return "", libcommon.StringError(errors.New("invalid secret key"))
	}

	return *authKey.PlatformId, nil
}

func (a auth) InvalidateRefreshToken(refreshToken string) error {
	return a.repos.Auth.Delete(libcommon.ToSha256(refreshToken))
}

func (a auth) RefreshToken(ctx context.Context, refreshToken string, walletAddress string, platformId string) (UserCreateResponse, error) {
	resp := UserCreateResponse{}

	// get user id from refresh token
	userId, err := a.repos.Auth.GetUserIdFromRefreshToken(libcommon.ToSha256(refreshToken))
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// verify wallet address
	// Verify user is registered to this wallet address
	instrument, err := a.repos.Instrument.GetWalletByAddr(ctx, walletAddress)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	if instrument.UserId != userId {
		return resp, libcommon.StringError(serror.NOT_FOUND)
	}

	// get device
	device, err := a.repos.Device.GetByUserId(ctx, userId)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// create new jwt
	jwt, err := a.GenerateJWT(userId, platformId, device)
	if err != nil {
		return resp, libcommon.StringError(err)
	}
	resp.JWT = jwt

	// delete old refresh token
	err = a.InvalidateRefreshToken(refreshToken)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	user, err := a.repos.User.GetById(ctx, instrument.UserId)
	if err != nil {
		return resp, libcommon.StringError(err)
	}

	// get email
	user.Email = getValidatedEmailOrEmpty(ctx, a.repos.Contact, user.Id)
	resp.User = user

	return resp, nil
}

func verifyWalletAuthentication(request model.WalletSignaturePayloadSigned) error {
	key := env.Var.STRING_ENCRYPTION_KEY
	preSignedPayload, err := libcommon.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return libcommon.StringError(err)
	}
	// Verify users signature
	bytes := []byte(request.Nonce)
	valid, err := common.ValidateExternalEVMSignature(request.Signature, preSignedPayload.Address, bytes, true) // true: expect eip131
	if err != nil {
		return libcommon.StringError(err)
	}
	if !valid {
		return libcommon.StringError(errors.New("user signature invalid"))
	}

	// Verify timestamp is not expired past 15 minutes
	if time.Now().Unix() > preSignedPayload.Timestamp+(15*60) {
		return libcommon.StringError(serror.EXPIRED)
	}

	return nil
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

func getValidatedEmailOrEmpty(ctx context.Context, contactRepo repository.Contact, userId string) string {
	contact, err := contactRepo.GetByUserIdAndStatus(ctx, userId, "validated")
	if err != nil {
		return ""
	}

	return contact.Data
}

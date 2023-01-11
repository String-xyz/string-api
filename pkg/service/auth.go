package service

import (
	netmail "net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type SignablePayload struct {
	Nonce string `json:"nonce"`
}

var hexRegex *regexp.Regexp = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// var walletAuthenticationPrefix string = "" // For testing locally

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
	UserId   string
	DeviceId string
	jwt.StandardClaims
}

type Auth interface {
	// PayloadToSign returns a payload to be sign by a wallet
	// to authenticate an user, the payload expires in 15 minutes
	PayloadToSign(walletAdress string) (SignablePayload, error)

	// VerifySignedPayload receives a signed payload from the user and verifies the signature
	// if signaure is valid it returns a JWT to authenticate the user
	VerifySignedPayload(model.WalletSignaturePayloadSigned) (UserCreateResponse, error)

	GenerateJWT(model.Device) (JWT, error)
	ValidateAPIKey(key string) bool
	RefreshToken(token string, walletAddress string) (JWT, error)
	InvalidateRefreshToken(token string) error
}

type auth struct {
	repos        repository.Repositories
	fingerprint  Fingerprint
	verification Verification
}

// reusing UserRepos here
func NewAuth(r repository.Repositories, f Fingerprint, v Verification) Auth {
	return &auth{r, f, v}
}

func (a auth) PayloadToSign(walletAddress string) (SignablePayload, error) {
	payload := model.WalletSignaturePayload{}
	signable := SignablePayload{}

	if !hexRegex.MatchString(walletAddress) {
		return signable, common.StringError(errors.New("missing or invalid address"))
	}
	payload.Address = walletAddress
	payload.Timestamp = time.Now().Unix()
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	encrypted, err := common.Encrypt(payload, key)
	if err != nil {
		return signable, common.StringError(err)
	}
	return SignablePayload{walletAuthenticationPrefix + encrypted}, nil
}

func (a auth) VerifySignedPayload(request model.WalletSignaturePayloadSigned) (UserCreateResponse, error) {
	resp := UserCreateResponse{}
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	payload, err := common.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return resp, common.StringError(err)
	}

	if err := verifyWalletAuthentication(request); err != nil {
		return resp, common.StringError(err)
	}

	// Verify user is registered to this wallet address
	instrument, err := a.repos.Instrument.GetWallet(payload.Address)
	if err != nil {
		return resp, common.StringError(err)
	}
	user, err := a.repos.User.GetById(instrument.UserID)
	if err != nil {
		return resp, common.StringError(err)
	}

	created, device, err := a.createDeviceIfNeeded(user.ID, request.Fingerprint.VisitorID, request.Fingerprint.RequestID)
	if err != nil {
		return resp, common.StringError(err)
	}

	if created || device.ValidatedAt == nil {
		go a.verification.SendDeviceVerification(user.ID, device.ID, device.Description)
		return resp, common.StringError(errors.New("unknown device"))
	}

	// Create the JWT
	jwt, err := a.GenerateJWT(device)
	if err != nil {
		return resp, common.StringError(err)
	}
	return UserCreateResponse{JWT: jwt, User: user}, nil
}

func (a auth) createDeviceIfNeeded(userID, visitorID, requestID string) (bool, model.Device, error) {
	device, err := a.repos.Device.GetByUserIdAndFingerprint(userID, visitorID)
	if err == nil {
		return false, device, nil
	}
	// create device only if the error is not found
	if err != nil && err == repository.ErrNotFound {
		visitor, fpErr := a.fingerprint.GetVisitor(visitorID, requestID)
		if fpErr != nil {
			return false, model.Device{}, common.StringError(fpErr)
		}
		device, dErr := a.createDevice(userID, visitor)
		return dErr == nil, device, dErr
	}

	return false, device, common.StringError(err)
}

func (a auth) createDevice(userID string, visitor model.FPVisitor) (model.Device, error) {
	return a.repos.Device.Create(model.Device{
		UserID:      userID,
		Fingerprint: visitor.VisitorID,
		Type:        visitor.Type,
		IpAddresses: pq.StringArray{visitor.IPAddress},
		Description: visitor.UserAgent,
		LastUsedAt:  time.Now(),
		ValidatedAt: nil,
	})
}

// GenerateJWT generates a jwt token and a refresh token which is saved on redis
func (a auth) GenerateJWT(m model.Device) (JWT, error) {
	claims := JWTClaims{}
	refreshToken := uuidWithoutHyphens()
	t := &JWT{
		IssuedAt: time.Now(),
		ExpAt:    time.Now().Add(time.Minute * 15),
	}

	claims.DeviceId = m.ID
	claims.UserId = m.UserID
	claims.ExpiresAt = t.ExpAt.Unix()
	claims.IssuedAt = t.IssuedAt.Unix()
	// replace this signing method with RSA or something similar
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return *t, err
	}
	t.Token = signed

	// create and save
	refreshObj, err := a.repos.Auth.CreateJWTRefresh(common.ToSha256(refreshToken), m.UserID)
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
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	return t.Valid, err
}

func (a auth) ValidateAPIKey(key string) bool {
	hashed := common.ToSha256(key)
	authKey, err := a.repos.Auth.Get(hashed)
	if err != nil {
		return false
	}
	return authKey.Data == hashed
}

func (a auth) InvalidateRefreshToken(refreshToken string) error {
	return a.repos.Auth.Delete(common.ToSha256(refreshToken))
}

func (a auth) RefreshToken(refreshToken string, walletAddress string) (JWT, error) {
	// get user id from refresh token
	userId, err := a.repos.Auth.GetUserIdFromRefreshToken(common.ToSha256(refreshToken))
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	// verify wallet address
	// Verify user is registered to this wallet address
	instrument, err := a.repos.Instrument.GetWallet(walletAddress)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return JWT{}, common.StringError(errors.New("wallet address not associated with this user: " + walletAddress))
		}
		return JWT{}, common.StringError(err)
	}

	if instrument.UserID != userId {
		return JWT{}, common.StringError(errors.New("wallet address not associated with this user: " + walletAddress))
	}

	// get device
	device, err := a.repos.Device.GetByUserId(userId)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	// create new jwt
	jwt, err := a.GenerateJWT(device)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	// delete old refresh token
	err = a.InvalidateRefreshToken(refreshToken)
	if err != nil {
		return JWT{}, common.StringError(err)
	}

	return jwt, nil
}

func verifyWalletAuthentication(request model.WalletSignaturePayloadSigned) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	preSignedPayload, err := common.Decrypt[model.WalletSignaturePayload](request.Nonce[len(walletAuthenticationPrefix):], key)
	if err != nil {
		return common.StringError(err)
	}
	// Verify users signature
	bytes := []byte(request.Nonce)
	valid, err := common.ValidateExternalEVMSignature(request.Signature, preSignedPayload.Address, bytes, true) // true: expect eip131
	if err != nil {
		return common.StringError(err)
	}
	if !valid {
		return common.StringError(errors.New("user signature invalid"))
	}

	// Verify timestamp is not expired past 15 minutes
	if time.Now().Unix() > preSignedPayload.Timestamp+(15*60) {
		return common.StringError(errors.New("login payload expired"))
	}

	return nil
}

// Use native mail package to check if email a valid email
func validEmail(email string) bool {
	_, err := netmail.ParseAddress(email)
	return err == nil
}

func uuidWithoutHyphens() string {
	s := uuid.New().String()
	return strings.Replace(s, "-", "", -1)
}

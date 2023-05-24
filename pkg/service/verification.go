package service

import (
	"context"
	"net/url"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/go-lib/v2/validator"
	"github.com/String-xyz/string-api/config"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
)

type EmailVerification struct {
	Timestamp  int64
	Email      string
	UserId     string
	PlatformId string
}

type DeviceVerification struct {
	Timestamp int64
	DeviceId  string
	UserId    string
}

type Verification interface {
	// SendEmailVerification sends a link to the provided email for verification purpose, link expires in 15 minutes
	SendEmailVerification(ctx context.Context, platformId string, userId string, email string) error
	SendDeviceVerification(ctx context.Context, userId string, email string, deviceId string, deviceDescription string) error
	// VerifyEmail verifies the provided email and creates a contact
	VerifyEmail(ctx context.Context, platformId string, userId string, email string) error
	VerifyEmailWithEncryptedToken(ctx context.Context, encrypted string) error
	PreValidateEmail(ctx context.Context, platformId, userId, email string) error
}

type verification struct {
	repos  repository.Repositories
	unit21 Unit21
}

func NewVerification(repos repository.Repositories, unit21 Unit21) Verification {
	return &verification{repos, unit21}
}

func (v verification) SendEmailVerification(ctx context.Context, platformId string, userId string, email string) error {
	_, finish := Span(ctx, "service.verification.SendEmailVerification", SpanTag{"platformId": platformId})
	defer finish()

	if !validator.ValidEmail(email) {
		return libcommon.StringError(serror.INVALID_DATA)
	}

	user, err := v.repos.User.GetById(ctx, userId)
	if err != nil || user.Id != userId {
		return libcommon.StringError(serror.INVALID_DATA) // JWT expiration will not be hit here
	}

	contact, _ := v.repos.Contact.GetByData(ctx, email)
	if contact.Status == "validated" {
		return libcommon.StringError(serror.ALREADY_IN_USE)
	}

	// Encrypt required data to Base64 string and insert it in an email hyperlink
	key := config.Var.STRING_ENCRYPTION_KEY
	code, err := libcommon.Encrypt(EmailVerification{Timestamp: time.Now().Unix(), Email: email, UserId: userId, PlatformId: platformId}, key)
	if err != nil {
		return libcommon.StringError(err)
	}
	code = url.QueryEscape(code) // make sure special characters are browser friendly

	e := NewEmail()

	return e.SendEmailVerification(ctx, email, code)
}

func (v verification) SendDeviceVerification(ctx context.Context, userId string, email string, deviceId string, deviceDescription string) error {
	log.Info().Str("email", email)

	key := config.Var.STRING_ENCRYPTION_KEY

	code, err := libcommon.Encrypt(DeviceVerification{Timestamp: time.Now().Unix(), DeviceId: deviceId, UserId: userId}, key)
	if err != nil {
		return libcommon.StringError(err)
	}
	code = url.QueryEscape(code)

	link := config.Var.BASE_URL + "verification?type=device&token=" + code

	textContent := "We noticed that you attempted to log in from " + deviceDescription + " at " + time.Now().Local().Format(time.RFC1123) + ". Is this you?"

	e := NewEmail()
	return e.SendDeviceVerification(ctx, email, link, textContent)
}

func (v verification) VerifyEmail(ctx context.Context, userId string, email string, platformId string) error {
	_, finish := Span(ctx, "services.verification.VerifyEmail", SpanTag{"platformId": platformId})
	defer finish()

	now := time.Now()
	// 1. Create contact with email
	contact := model.Contact{UserId: userId, Type: "email", Status: "validated", Data: email, ValidatedAt: &now}
	contact, err := v.repos.Contact.Create(ctx, contact)
	if err != nil {
		return libcommon.StringError(err)
	}

	// 2. Update user status
	user, err := v.repos.User.UpdateStatus(ctx, userId, "email_verified")
	if err != nil {
		// TODO: Log error errors.New("User email verify error - userId: " + user.Id)
		return libcommon.StringError(err)
	}

	// 3. Associate contact with platform
	err = v.repos.Platform.AssociateContact(ctx, contact.Id, platformId)
	if err != nil {
		return libcommon.StringError(err)
	}

	// 4. update user in unit21
	ctx2 := context.Background() // Create a new context since this will run in background
	go v.unit21.Entity.Update(ctx2, user)

	return nil
}

func (v verification) VerifyEmailWithEncryptedToken(ctx context.Context, encrypted string) error {
	_, finish := Span(ctx, "services.verification.VerifyEmailWithEncryptedToken")
	defer finish()

	key := config.Var.STRING_ENCRYPTION_KEY

	received, err := libcommon.Decrypt[EmailVerification](encrypted, key)
	if err != nil {
		return libcommon.StringError(err)
	}
	// Wait for up to 15 minutes, final timeout TBD
	now := time.Now()
	if now.Unix()-received.Timestamp > (60 * 15) {
		return libcommon.StringError(serror.EXPIRED)
	}

	err = v.VerifyEmail(ctx, received.UserId, received.Email, received.PlatformId)
	if err != nil {
		return libcommon.StringError(err)
	}

	return nil
}

func (v verification) PreValidateEmail(ctx context.Context, platformId, userId, email string) error {
	_, finish := Span(ctx, "services.verification.PreValidateEmail", SpanTag{"platformId": platformId})
	defer finish()

	return v.VerifyEmail(ctx, userId, email, platformId)
}

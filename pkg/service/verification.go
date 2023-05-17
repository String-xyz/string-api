package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/go-lib/v2/validator"
	"github.com/String-xyz/string-api/config"
	"github.com/String-xyz/string-api/pkg/internal/common"

	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	// VerifyEmail verifies the provided email and creates a contact
	VerifyEmail(ctx context.Context, platformId string, userId string, email string) error
	SendDeviceVerification(userId, email string, deviceId string, deviceDescription string) error
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

	fromAddress := config.Var.AUTH_EMAIL_ADDRESS

	baseURL := common.GetBaseURL()
	from := mail.NewEmail("String Authentication", fromAddress)
	subject := "String Email Verification"
	to := mail.NewEmail("New String User", email)
	textContent := "Click the link below to complete your e-email verification!"
	htmlContent := `<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='` + baseURL + `verification?type=email&token=` + code + `' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Verify Email Now</a></div>`

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)

	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	_, err = client.Send(message)
	if err != nil {
		return libcommon.StringError(err)
	}

	return nil

}

func (v verification) SendDeviceVerification(userId, email, deviceId, deviceDescription string) error {
	log.Info().Str("email", email)

	key := config.Var.STRING_ENCRYPTION_KEY

	code, err := libcommon.Encrypt(DeviceVerification{Timestamp: time.Now().Unix(), DeviceId: deviceId, UserId: userId}, key)
	if err != nil {
		return libcommon.StringError(err)
	}
	code = url.QueryEscape(code)

	baseURL := common.GetBaseURL()
	fromAddress := config.Var.AUTH_EMAIL_ADDRESS
	from := mail.NewEmail("String XYZ", fromAddress)
	subject := "New Device Login Verification"
	to := mail.NewEmail("New Device Login", email)
	link := baseURL + "verification?type=device&token=" + code

	textContent := "We noticed that you attempted to log in from " + deviceDescription + " at " + time.Now().Local().Format(time.RFC1123) + ". Is this you?"
	htmlContent := fmt.Sprintf(`<div style="font-family: inherit; text-align: inherit"><span style="color: #172b4d; font-family: -apple-system, &quot;system-ui&quot;, &quot;Segoe UI&quot;, Roboto, Oxygen, Ubuntu, &quot;Fira Sans&quot;, &quot;Droid Sans&quot;, &quot;Helvetica Neue&quot;, sans-serif; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: -0.07px; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: pre-wrap; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; float: none; display: inline; font-size: 14px">%s</span></div>
	<td align="center" bgcolor="#b7c23e" class="inner-td" style="border-radius:6px; font-size:16px; text-align:center; background-color:inherit;"><a href="%s" style="background-color:#b7c23e; border:1px solid 0; border-color:0; border-radius:6px; border-width:1px; color:#ffffff; display:inline-block; font-size:14px; font-weight:normal; letter-spacing:0px; line-height:normal; padding:12px 18px 12px 18px; text-align:center; text-decoration:none; border-style:solid;" target="_blank">Yes</a></td>`,
		textContent, link)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	_, err = client.Send(message)
	if err != nil {
		log.Err(err).Msg("error sending device validation")
		return libcommon.StringError(err)
	}

	return nil
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

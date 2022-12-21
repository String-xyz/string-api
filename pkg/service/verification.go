package service

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailVerification struct {
	Timestamp int64
	Email     string
	UserID    string
}

type DeviceVerification struct {
	Timestamp int64
	DeviceID  string
	UserID    string
}

type Verification interface {
	// SendEmailVerification sends a link to the provided email for verification purpose, link expires in 15 minutes
	SendEmailVerification(userID string, email string) error

	// VerifyEmail verifies the provided email and creates a contact
	VerifyEmail(encrypted string) error

	SendDeviceVerification(userID string, deviceID string, deviceDescription string) error
	VerifyDevice(encrypted string) error
}

type verification struct {
	repos repository.Repositories
}

func NewVerification(repos repository.Repositories) Verification {
	return &verification{repos}
}

func (v verification) SendEmailVerification(userID, email string) error {
	if !validEmail(email) {
		return common.StringError(errors.New("missing or invalid email"))
	}

	user, err := v.repos.User.GetById(userID)
	if err != nil || user.ID != userID {
		return common.StringError(errors.New("invalid user")) // JWT expiration will not be hit here
	}

	contact, _ := v.repos.Contact.GetByData(email)
	if contact.Status == "validated" {
		return common.StringError(errors.New("email is already authenticated"))
	}

	// Encrypt required data to Base64 string and insert it in an email hyperlink
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	code, err := common.Encrypt(EmailVerification{Timestamp: time.Now().Unix(), Email: email, UserID: userID}, key)
	if err != nil {
		return common.StringError(err)
	}
	code = url.QueryEscape(code) // make sure special characters are browser friendly

	baseURL := common.GetBaseURL()
	from := mail.NewEmail("String Authentication", "auth@string.xyz")
	subject := "String Email Verification"
	to := mail.NewEmail("New String User", email)
	textContent := "Click the link below to complete your e-email verification!"
	htmlContent := `<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='` + baseURL + `verification?type=email&token=` + code + `' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Verify Email Now</a></div>`

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return common.StringError(err)
	}
	// Wait for up to 15 minutes, final timeout TBD
	now, lastPolled := time.Now().Unix(), time.Now().Unix()
	until := now + (60 * 15)
	for now < until {
		now = time.Now().Unix()
		if now-lastPolled < 3 {
			continue // throttle following logic in 3 second interval
		}
		lastPolled = now
		contact, err := v.repos.Contact.GetByData(email)
		if err != nil && errors.Cause(err).Error() != "not found" {
			return common.StringError(err)
		} else if err == nil && contact.Data == email {
			return nil // success
		}
	}
	// timed out
	return common.StringError(errors.New("link expired"))
}

func (v verification) SendDeviceVerification(userID, deviceID, deviceDescription string) error {
	email, err := v.repos.Contact.GetByUserIdAndStatus(userID, "validated")
	if err != nil {
		log.Err(err).Msg("Error getting a valid email")
		return err
	}
	log.Info().Str("email", email.Data)
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	code, err := common.Encrypt(DeviceVerification{Timestamp: time.Now().Unix(), DeviceID: deviceID, UserID: userID}, key)
	if err != nil {
		return common.StringError(err)
	}
	code = url.QueryEscape(code)

	baseURL := common.GetBaseURL()
	from := mail.NewEmail("String XYZ", "auth@string.xyz")
	subject := "New Device Login Verification"
	to := mail.NewEmail("New Device Login", email.Data)
	link := baseURL + "verification?type=device&token=" + code
	htmlContent := fmt.Sprintf(`<div style="font-family: inherit; text-align: inherit"><span style="color: #172b4d; font-family: -apple-system, &quot;system-ui&quot;, &quot;Segoe UI&quot;, Roboto, Oxygen, Ubuntu, &quot;Fira Sans&quot;, &quot;Droid Sans&quot;, &quot;Helvetica Neue&quot;, sans-serif; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; font-weight: 400; letter-spacing: -0.07px; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: pre-wrap; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial; float: none; display: inline; font-size: 14px">We noticed that you attempted to log in from a new device %s. Is this you?</span></div>
	<td align="center" bgcolor="#b7c23e" class="inner-td" style="border-radius:6px; font-size:16px; text-align:center; background-color:inherit;"><a href="%s" style="background-color:#b7c23e; border:1px solid 0; border-color:0; border-radius:6px; border-width:1px; color:#ffffff; display:inline-block; font-size:14px; font-weight:normal; letter-spacing:0px; line-height:normal; padding:12px 18px 12px 18px; text-align:center; text-decoration:none; border-style:solid;" target="_blank">Yes</a></td>`,
		deviceDescription, link)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		log.Err(err).Msg("error sending device validation")
		return common.StringError(err)
	}

	return nil
}

func (v verification) VerifyEmail(encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[EmailVerification](encrypted, key)
	if err != nil {
		return common.StringError(err)
	}
	// Wait for up to 15 minutes, final timeout TBD
	now := time.Now()
	if now.Unix()-received.Timestamp > (60 * 15) {
		return common.StringError(errors.New("link expired"))
	}
	contact := model.Contact{UserID: received.UserID, Type: "email", Status: "validated", Data: received.Email, ValidatedAt: &now}
	contact, err = v.repos.Contact.Create(contact)
	return common.StringError(err)
}

func (v verification) VerifyDevice(encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[DeviceVerification](encrypted, key)
	if err != nil {
		return common.StringError(err)
	}
	now := time.Now()
	if now.Unix()-received.Timestamp > (60 * 15) {
		return common.StringError(errors.New("link expired"))
	}
	err = v.repos.Device.Update(received.DeviceID, model.DeviceUpdates{ValidatedAt: &now})
	return err
}

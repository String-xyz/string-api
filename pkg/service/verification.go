package service

import (
	"net/url"
	"os"
	"time"

	"github.com/String-xyz/string-api/pkg/internal/common"
	"github.com/String-xyz/string-api/pkg/model"
	"github.com/String-xyz/string-api/pkg/repository"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Verification interface {
	// SendEmailVerification sends a link to the provided email for verification purpose, link expires in 15 minutes
	SendEmailVerification(userID string, email string) error

	// VerifyEmail verifies the provided email and creates a contact
	VerifyEmail(encrypted string) error
}

type verification struct {
	contact repository.Contact
	user    repository.User
}

func NewVerification(contact repository.Contact, user repository.User) Verification {
	return &verification{contact: contact, user: user}
}

func (v verification) SendEmailVerification(userID, email string) error {
	if !validEmail(email) {
		return common.StringError(errors.New("missing or invalid email"))
	}

	user, err := v.user.GetById(userID)
	if err != nil || user.ID != userID {
		return common.StringError(errors.New("invalid or expired JWT"))
	}

	contact, _ := v.contact.GetByData(email)
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
	htmlContent := `<div style='font-family: inherit; text-align: inherit; margin-left: 0px'><br><a href='` + baseURL + `verification/?type=email,token=` + code + `' style='background-color:#ffbe00; color:#000000; display:inline-block; padding:12px 40px 12px 40px; text-align:center; text-decoration:none;' target='_blank'>Verify Email Now</a></div>`

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
		contact, err := v.contact.GetByData(email)
		if err != nil && errors.Cause(err).Error() != "not found" {
			return common.StringError(err)
		} else if err == nil && contact.Data == email {
			return nil // success
		}
	}
	// timed out
	return common.StringError(errors.New("link expired"))
}

func (v verification) VerifyEmail(encrypted string) error {
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	received, err := common.Decrypt[EmailVerification](encrypted, key)
	if err != nil {
		return common.StringError(err)
	}
	// Wait for up to 15 minutes, final timeout TBD
	now := time.Now().Unix()
	if now-received.Timestamp > (60 * 15) {
		return common.StringError(errors.New("link expired"))
	}
	contact := model.Contact{UserID: received.UserID, Type: "email", Status: "validated", Data: received.Email}
	contact, _ = v.contact.Create(contact)
	return nil
}

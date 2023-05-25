package emailer

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Emailer interface {
	SendReceipt(ctx context.Context, email string, params ReceiptGenerationParams) error
	SendEmailVerification(ctx context.Context, email string, code string) error
	SendDeviceVerification(ctx context.Context, email string, link string, textContent string) error
}

//go:embed templates/*
var templatesFS embed.FS

type emailer struct {
}

func New() Emailer {
	return &emailer{}
}

type ReceiptGenerationParams struct {
	ReceiptType         string
	CustomerName        string
	PaymentDescriptor   string
	TransactionDate     string
	StringPaymentId     string
	TransactionId       string
	TransactionExplorer string
	DestinationAddress  string
	DestinationExplorer string
	PaymentMethod       string
	Platform            string
	ItemOrdered         string
	TokenId             string
	Subtotal            string
	NetworkFee          string
	ProcessingFee       string
	Total               string
}

func (e emailer) SendEmailVerification(ctx context.Context, email string, code string) error {
	// tmpl, err := template.ParseFiles("pkg/templates/email_verification.html")
	tmpl, err := template.ParseFS(templatesFS, "templates/email_verification.tpl")
	if err != nil {
		return err
	}

	baseURL := config.Var.BASE_URL

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "email_verification.tpl", map[string]interface{}{
		"BaseURL": baseURL,
		"Code":    code,
	})
	if err != nil {
		return err
	}

	from := mail.NewEmail("String Authentication", config.Var.AUTH_EMAIL_ADDRESS)
	subject := "String Email Verification"
	to := mail.NewEmail("New String User", email)
	textContent := "Click the link below to complete your e-email verification!"

	return sendEmail(ctx, from, subject, to, textContent, buf.String())
}

func (e emailer) SendDeviceVerification(ctx context.Context, email string, link string, textContent string) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/device_verification.tpl")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "device_verification.html", map[string]interface{}{
		"TextContent": textContent,
		"Link":        link,
	})
	if err != nil {
		return err
	}

	from := mail.NewEmail("String XYZ", config.Var.AUTH_EMAIL_ADDRESS)
	subject := "New Device Login Verification"
	to := mail.NewEmail("New Device Login", email)

	return sendEmail(ctx, from, subject, to, textContent, buf.String())
}

func (e emailer) SendReceipt(ctx context.Context, email string, params ReceiptGenerationParams) error {
	tmpl, err := template.ParseFS(templatesFS, "templates/receipt.tpl")
	// tmpl, err := template.ParseFiles("pkg/templates/receipt.html")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "receipt.html", map[string]interface{}{
		"Params": params,
	})
	if err != nil {
		return err
	}

	from := mail.NewEmail("String Receipt", config.Var.RECEIPTS_EMAIL_ADDRESS)
	subject := "Your " + params.ReceiptType + " Receipt from String"
	to := mail.NewEmail(params.CustomerName, email)
	return sendEmail(ctx, from, subject, to, "", buf.String())
}

func sendEmail(ctx context.Context, from *mail.Email, subject string, to *mail.Email, text string, html string) error {
	message := mail.NewSingleEmail(from, subject, to, text, html)
	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	_, err := client.Send(message)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

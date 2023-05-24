package service

import (
	"bytes"
	"context"
	"text/template"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email interface {
	SendReceipt(ctx context.Context, email string, params ReceiptGenerationParams) error
}

type email struct {
}

func NewEmail() Email {
	return &email{}
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

func (e email) SendReceipt(ctx context.Context, email string, params ReceiptGenerationParams) error {
	tmpl, err := template.ParseFiles("pkg/templates/receipt.html")
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

	fromAddress := config.Var.RECEIPTS_EMAIL_ADDRESS

	from := mail.NewEmail("String Receipt", fromAddress)
	subject := "Your " + params.ReceiptType + " Receipt from String"
	to := mail.NewEmail(params.CustomerName, email)
	textContent := ""
	message := mail.NewSingleEmail(from, subject, to, textContent, buf.String())
	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	_, err = client.Send(message)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

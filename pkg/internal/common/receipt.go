package common

import (
	"os"

	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type emailReceipt struct {
	Header string
	Body   [][2]string
	Footer string
}

type ReceiptGenerationParams struct {
	ReceiptType       string
	CustomerName      string
	PaymentDescriptor string
	TransactionDate   string
	StringPaymentId   string
	RecipientAddress  string
}

func StringifyEmailReceipt(e emailReceipt) string {
	res := e.Header + "<p><br>"

	for index := range e.Body {
		res += e.Body[index][0] + ": " + e.Body[index][1] + "<br>"
	}
	res += "</p>" + e.Footer
	return res
}

func GenerateReceipt(params ReceiptGenerationParams, body [][2]string) (StringEmail, error) {
	return GenerateEmail(EmailGenerationRequest{
		ProductName:    "PLATFORM NAME HERE",
		ProductLink:    "PLATFORM URL HERE",
		ProductLogoURL: "PLATFORM LOGO URL HERE",
		CustomerName:   params.CustomerName,
		Intros: []string{
			"Your" + params.ReceiptType + "Details:",
			"Transaction Date: " + params.TransactionDate,
			"String Payment ID: " + params.StringPaymentId,
		},
		Entries: GenerateEntries(body),
		Outros: []string{
			"The transaction will appear on your card statement as " + params.PaymentDescriptor,
			"All sales are final.  Please see our Terms of Service (https://www.string.xyz/terms-of-service)",
			"Please reference your String Payment ID " + params.StringPaymentId,
			"Service powered by String",
			"String XYZ LLC | 490 43rd St, #86, Oakland CA 94609. | NMLS ID: 2400614",
			"Please visit us at string.xyz.  Should you need to reach us, please contact us at support@string.xyz",
			"Consumer Fraud Warning",
			"If you feel you have been the victim of a scam you can contact the FTC at 1-877-FTC-HELP (382-4357)",
			"or online at www.ftc.gov (link is external); or the Consumer Financial Protection Bureau (CFPB) at 1-855-411-CFPB (2372)",
			"or online at www.consumerfinance.gov",
		},
		SendParams: SendEmailParameters{
			FromAddress: "auth@string.xyz", // TODO: create a new sender for receipts
			ToAddress:   params.RecipientAddress,
			FromName:    "String Receipt",
			ToName:      params.CustomerName,
			Subject:     "Your " + params.ReceiptType + " Receipt from String",
		},
	})
}

func EmailReceipt(email string, params ReceiptGenerationParams, body [][2]string) error {
	from := mail.NewEmail("String Receipt", "auth@string.xyz") // TODO: create a new sender for receipts
	subject := "Your " + params.ReceiptType + " Receipt from String"
	to := mail.NewEmail(params.CustomerName, email)
	htmlContent, err := GenerateReceipt(params, body)
	if err != nil {
		return libcommon.StringError(err)
	}
	message := mail.NewSingleEmail(from, subject, to, htmlContent.BodyText, htmlContent.BodyHTML)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err = client.Send(message)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

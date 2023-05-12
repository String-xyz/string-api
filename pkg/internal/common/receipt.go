package common

import (
	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
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
}

func StringifyEmailReceipt(e emailReceipt) string {
	res := e.Header + "<p><br>"

	for index := range e.Body {
		res += e.Body[index][0] + ": " + e.Body[index][1] + "<br>"
	}
	res += "</p>" + e.Footer
	return res
}

func GenerateReceipt(params ReceiptGenerationParams, body [][2]string) string {
	header := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>Your " + params.ReceiptType + " Details</header>" +
		"<br>Dear " + params.CustomerName + "," +
		"<br>Thank you for using String.  Here is your transaction receipt:" +
		"<br>Transaction Date: " + params.TransactionDate +
		"<br>String Payment ID: " + params.StringPaymentId

	footer := "" +
		"<br>The transaction will appear on your card statement as " + params.PaymentDescriptor +
		"<br>All sales are final.  Please see our <a href='https://www.string.xyz/terms-of-service'>Terms of Service</a>" +
		"<br>Please reference your String Payment ID " + params.StringPaymentId +
		"<br><br>Service powered by String" +
		"<br>String XYZ LLC | 490 43rd St, #86, Oakland CA 94609. | NMLS ID: 2400614" +
		"<br>Please visit us at string.xyz.  Should you need to reach us, please contact us at <a href='mailto:support@string.xyz'>support@string.xyz</a>." +
		"<br><br>Consumer Fraud Warning" +
		"<br>If you feel you have been the victim of a scam you can contact the FTC at 1-877-FTC-HELP (382-4357)" +
		"<br>or online at <a href='www.ftc.gov'>www.ftc.gov</a> (link is external); or the Consumer Financial Protection Bureau (CFPB) at 1-855-411-CFPB (2372)" +
		"<br>or online at <a href='www.consumerfinance.gov'>www.consumerfinance.gov</a>"
	return StringifyEmailReceipt(emailReceipt{Header: header, Body: body, Footer: footer})
}

func EmailReceipt(email string, params ReceiptGenerationParams, body [][2]string) error {
	fromAddress := config.Var.RECEIPTS_EMAIL_ADDRESS

	from := mail.NewEmail("String Receipt", fromAddress)
	subject := "Your " + params.ReceiptType + " Receipt from String"
	to := mail.NewEmail(params.CustomerName, email)
	textContent := ""
	htmlContent := GenerateReceipt(params, body)
	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(config.Var.SENDGRID_API_KEY)
	_, err := client.Send(message)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

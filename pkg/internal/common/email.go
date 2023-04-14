package common

import (
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/matcornic/hermes"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type StringEmail struct {
	BodyHTML   string
	BodyText   string
	SendParams SendEmailParameters
}

type EmailGenerationRequest struct {
	ProductName    string
	ProductLink    string
	ProductLogoURL string
	CustomerName   string
	Intros         []string
	Outros         []string
	Actions        []hermes.Action
	Entries        []hermes.Entry
	SendParams     SendEmailParameters
}

type SendEmailParameters struct {
	FromAddress string
	ToAddress   string
	FromName    string
	ToName      string
	Subject     string
}

func SendEmail(email StringEmail) error {
	from := mail.NewEmail(email.SendParams.FromName, email.SendParams.FromAddress)
	subject := email.SendParams.Subject
	to := mail.NewEmail(email.SendParams.ToName, email.SendParams.ToAddress)
	plainTextContent := email.BodyText
	htmlContent := email.BodyHTML

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}

func GenerateEmail(request EmailGenerationRequest) (StringEmail, error) {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: request.ProductName,
			Link: request.ProductLink,
			Logo: request.ProductLogoURL,
		},
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name:       request.CustomerName,
			Intros:     request.Intros,
			Actions:    request.Actions,
			Dictionary: request.Entries,
			Outros:     request.Outros,
		},
	}

	body, err := h.GenerateHTML(email)
	if err != nil {
		return StringEmail{}, common.StringError(err)
	}
	text, err := h.GeneratePlainText(email)
	if err != nil {
		return StringEmail{}, common.StringError(err)
	}
	return StringEmail{
		BodyHTML:   body,
		BodyText:   text,
		SendParams: request.SendParams,
	}, nil
}

func GenerateEntries(entries [][2]string) []hermes.Entry {
	var entriesList []hermes.Entry
	for _, entry := range entries {
		e := hermes.Entry{
			Key:   entry[0],
			Value: entry[1],
		}
		entriesList = append(entriesList, e)
	}
	return entriesList
}

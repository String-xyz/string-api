package service

import (
	"strings"

	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/string-api/config"
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSMS(message string, recipients []string) error {
	var SMS_SID = config.Var.TWILIO_SMS_SID
	client := twilio.NewRestClient() // TWILIO_ACCOUNT_SID and TWILIO_AUTH_TOKEN are loaded from env in constructor
	params := &twilioApi.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom(SMS_SID) // Select a number automatically with our SID

	var errs error
	for _, recipient := range recipients {
		params.SetTo(recipient)
		_, err := client.Api.CreateMessage(params)
		if err != nil {
			if errs != nil {
				errs = errors.New(errs.Error() + err.Error()) // Concatenate errors if there are multiple
			} else {
				errs = err
			}
		}
	}
	if errs != nil {
		return libcommon.StringError(errs)
	}
	return nil
}

func MessageTeam(message string) error {
	var teamNumbers = config.Var.TEAM_PHONE_NUMBERS
	recipients := strings.Split(teamNumbers, ",")
	err := SendSMS(message, recipients)
	if err != nil {
		return libcommon.StringError(err)
	}
	return nil
}

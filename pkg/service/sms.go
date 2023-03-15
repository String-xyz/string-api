package service

import (
	"os"
	"strings"

	libCommon "github.com/String-xyz/go-lib/common"
	"github.com/pkg/errors"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSMS(message string, recipients []string) error {
	var SMS_SID = os.Getenv("TWILIO_SMS_SID")
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
		return libCommon.StringError(errs)
	}
	return nil
}

func MessageStaff(message string) error {
	var devNumbers = os.Getenv("DEV_PHONE_NUMBERS")
	recipients := strings.Split(devNumbers, ",")
	err := SendSMS(message, recipients)
	if err != nil {
		return libCommon.StringError(err)
	}
	return nil
}

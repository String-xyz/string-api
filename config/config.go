package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

type vars struct {
	BASE_URL                 string
	ENV                      string
	PORT                     string
	STRING_HOTWALLET_ADDRESS string
	COINGECKO_API_URL        string
	COINCAP_API_URL          string
	OWLRACLE_API_URL         string
	OWLRACLE_API_KEY         string
	OWLRACLE_API_SECRET      string
	AWS_REGION               string
	AWS_ACCT                 string
	AWS_ACCESS_KEY_ID        string
	AWS_SECRET_ACCESS_KEY    string
	AWS_KMS_KEY_ID           string
	CHECKOUT_PUBLIC_KEY      string
	CHECKOUT_SECRET_KEY      string
	CHECKOUT_ENV             string
	EVM_PRIVATE_KEY          string
	DB_NAME                  string
	DB_USERNAME              string
	DB_PASSWORD              string
	DB_HOST                  string
	DB_PORT                  string
	REDIS_PASSWORD           string
	REDIS_HOST               string
	REDIS_PORT               string
	JWT_SECRET_KEY           string
	UNIT21_API_KEY           string
	UNIT21_ENV               string
	UNIT21_ORG_NAME          string
	UNIT21_RTR_URL           string
	TWILIO_ACCOUNT_SID       string
	TWILIO_AUTH_TOKEN        string
	TWILIO_SMS_SID           string
	TEAM_PHONE_NUMBERS       string
	STRING_ENCRYPTION_KEY    string
	SENDGRID_API_KEY         string
	FINGERPRINT_API_KEY      string
	FINGERPRINT_API_URL      string
	STRING_INTERNAL_ID       string
	STRING_WALLET_ID         string
	STRING_BANK_ID           string
	SERVICE_NAME             string
	DEBUG_MODE               string
	AUTH_EMAIL_ADDRESS       string
	RECEIPTS_EMAIL_ADDRESS   string
}

var Var vars

func LoadEnv() error {
	godotenv.Load(".env")
	missing := []string{}
	stype := reflect.ValueOf(&Var).Elem()
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		key := stype.Type().Field(i).Name
		value := os.Getenv(key)
		if value == "" {
			missing = append(missing, key)
		}
		field.SetString(value)
	}
	if len(missing) > 0 {
		panic("Missing environment variable: " + strings.Join(missing, ", "))
	}
	return nil
}

package env

import (
	"errors"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/joho/godotenv"
)

var keys = []string{
	"BASE_URL",
	"ENV",
	"PORT",
	"STRING_HOTWALLET_ADDRESS",
	"COINGECKO_API_URL",
	"COINCAP_API_URL",
	"OWLRACLE_API_URL",
	"OWLRACLE_API_KEY",
	"OWLRACLE_API_SECRET",
	"AWS_REGION",
	"AWS_ACCT",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_KMS_KEY_ID",
	"CHECKOUT_PUBLIC_KEY",
	"CHECKOUT_SECRET_KEY",
	"CHECKOUT_ENV",
	"EVM_PRIVATE_KEY",
	"DB_NAME",
	"DB_USERNAME",
	"DB_PASSWORD",
	"DB_HOST",
	"DB_PORT",
	"REDIS_PASSWORD",
	"REDIS_HOST",
	"REDIS_PORT",
	"JWT_SECRET_KEY",
	"UNIT21_API_KEY",
	"UNIT21_ENV",
	"UNIT21_ORG_NAME",
	"UNIT21_RTR_URL",
	"TWILIO_ACCOUNT_SID",
	"TWILIO_AUTH_TOKEN",
	"TWILIO_SMS_SID",
	"TEAM_PHONE_NUMBERS",
	"STRING_ENCRYPTION_KEY",
	"SENDGRID_API_KEY",
	"IPSTACK_API_KEY",
	"FINGERPRINT_API_KEY",
	"FINGERPRINT_API_URL",
	"STRING_INTERNAL_ID",
	"STRING_WALLET_ID",
	"STRING_BANK_ID",
	"SERVICE_NAME",
	"DEBUG_MODE",
	"DB_RESET",
	"AUTH_EMAIL_ADDRESS",
	"RECEIPTS_EMAIL_ADDRESS",
}

var envMap map[string]string

func LoadEnv() error {
	godotenv.Load(".env")
	envMap = make(map[string]string)
	for _, key := range keys {
		val := os.Getenv(key)
		if val == "" {
			return common.StringError(errors.New("Missing env var: " + key))
		}
		envMap[key] = val
	}
	return nil
}

func Get(key string) (string, error) {
	res, ok := envMap[key]
	if !ok {
		return "", common.StringError(errors.New("No such env var: " + key))
	}
	return res, nil
}

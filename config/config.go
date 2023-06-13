package config

type vars struct {
	AWS_ACCT                 string `required:"false"`
	AWS_ACCESS_KEY_ID        string `required:"false"`
	AWS_SECRET_ACCESS_KEY    string `required:"false"`
	DEBUG_MODE               string `required:"false"`
	SERVICE_NAME             string `required:"false"`
	STRING_HOTWALLET_ADDRESS string `required:"false"`
	CARD_FAIL_PROBABILITY    string `required:"false"`
	BASE_URL                 string `required:"true"`
	ENV                      string `required:"true"`
	PORT                     string `required:"true"`
	COINCAP_API_URL          string `required:"true"`
	COINGECKO_API_URL        string `required:"true"`
	OWLRACLE_API_URL         string `required:"true"`
	OWLRACLE_API_KEY         string `required:"true"`
	OWLRACLE_API_SECRET      string `required:"true"`
	AWS_REGION               string `required:"true"`
	AWS_KMS_KEY_ID           string `required:"true"`
	CHECKOUT_PUBLIC_KEY      string `required:"true"`
	CHECKOUT_SECRET_KEY      string `required:"true"`
	CHECKOUT_ENV             string `required:"true"`
	EVM_PRIVATE_KEY          string `required:"true"`
	DB_NAME                  string `required:"true"`
	DB_USERNAME              string `required:"true"`
	DB_PASSWORD              string `required:"true"`
	DB_HOST                  string `required:"true"`
	DB_PORT                  string `required:"true"`
	REDIS_PASSWORD           string `required:"true"`
	REDIS_HOST               string `required:"true"`
	REDIS_PORT               string `required:"true"`
	JWT_SECRET_KEY           string `required:"true"`
	UNIT21_API_KEY           string `required:"true"`
	UNIT21_ENV               string `required:"true"`
	UNIT21_ORG_NAME          string `required:"true"`
	UNIT21_RTR_URL           string `required:"true"`
	TWILIO_ACCOUNT_SID       string `required:"true"`
	TWILIO_AUTH_TOKEN        string `required:"true"`
	TWILIO_SMS_SID           string `required:"true"`
	TEAM_PHONE_NUMBERS       string `required:"true"`
	STRING_ENCRYPTION_KEY    string `required:"true"`
	SENDGRID_API_KEY         string `required:"true"`
	FINGERPRINT_API_KEY      string `required:"true"`
	FINGERPRINT_API_URL      string `required:"true"`
	STRING_INTERNAL_ID       string `required:"true"`
	STRING_WALLET_ID         string `required:"true"`
	STRING_BANK_ID           string `required:"true"`
	AUTH_EMAIL_ADDRESS       string `required:"true"`
	RECEIPTS_EMAIL_ADDRESS   string `required:"true"`
}

var Var vars

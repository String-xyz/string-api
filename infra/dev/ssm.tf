data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
}

data "aws_ssm_parameter" "evm_private_key" {
  name = "string-encrypted-sk"
}

data "aws_ssm_parameter" "string_encryption_secret" {
  name = "string-encryption-secret"  
}

data "aws_ssm_parameter" "string_internal_id" {
  name = "string-internal-id"  
}

data "aws_ssm_parameter" "string_wallet_id" {
  name = "string-wallet-id"  
}

data "aws_ssm_parameter" "string_bank_id" {
  name = "string-bank-id"     
}

data "aws_ssm_parameter" "string_platform_id" {
  name = "string-placeholder-platform-id"  
}

data "aws_ssm_parameter" "user_jwt_secret" {
  name = "user-jwt-secret"
}

data "aws_ssm_parameter" "unit21_api_key" {
  name = "unit21-api-key"
}

data "aws_ssm_parameter" "ipstack_api_key" {
  name = "ipstack-api-key"
}

data "aws_ssm_parameter" "customer_jwt_secret" {
  name = "customer-jwt-secret"
}

data "aws_ssm_parameter" "checkout_public_key" {
  name = "dev-checkout-public-key"
}

data "aws_ssm_parameter" "checkout_private_key" {
  name = "dev-checkout-private-key"
}

data "aws_ssm_parameter" "owlracle_api_key" {
  name = "dev-owlracle-api-key"
}

data "aws_ssm_parameter" "fingerprint_api_key" { 
  name = "fingerprint-api-key"
}

data "aws_ssm_parameter" "sendgrid_api_key" {
  name = "sendgrid-api-key"
}

data "aws_ssm_parameter" "twilio_sms_sid" {
  name = "twilio-sms-sid"  
}

data "aws_ssm_parameter" "twilio_account_sid" {
 name = "twilio-account-sid"  
}

data "aws_ssm_parameter" "twilio_auth_token" {
  name = "twilio-auth-token"
}

data "aws_ssm_parameter" "owlracle_api_secret" {
  name = "dev-owlracle-api-secret"
}

data "aws_ssm_parameter" "db_password" {
  name = "string-rds-pg-db-password"
}

data "aws_ssm_parameter" "db_username" {
  name = "string-rds-pg-db-username"
}

data "aws_ssm_parameter" "db_name" {
  name = "string-rds-pg-db-name"
}

data "aws_ssm_parameter" "db_host" {
  name = "${local.env}-string-write-db-host-url"
}

data "aws_ssm_parameter" "redis_auth_token" {
  name = "redis-auth-token"
}

data "aws_ssm_parameter" "redis_host_url" {
  name  = "redis-host-url"
}

data "aws_ssm_parameter" "team_phone_numbers" {
  name = "team-phone-numbers"
}

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}


data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
}

data "aws_ssm_parameter" "hot_wallet" {
  name = "private-key"
}

data "aws_ssm_parameter" "user_jwt_secret" {
  name = "user-jwt-secret"
}

data "aws_ssm_parameter" "customer_jwt_secret" {
  name = "customer-jwt-secret"
}

data "aws_ssm_parameter" "checkout_public_key" {
  name = "checkout-public-key"
}

data "aws_ssm_parameter" "checkout_private_key" {
  name = "checkout-private-key"
}

data "aws_ssm_parameter" "owlracle_api_key" {
  name = "owlracle-api-key"
}

data "aws_ssm_parameter" "owlracle_api_secret" {
  name = "owlracle-api-secret"
}

data "aws_ssm_parameter" "db_password" {
  name = "string-pg-db-password"
}

data "aws_ssm_parameter" "db_username" {
  name = "string-pg-db-username"
}

data "aws_ssm_parameter" "db_name" {
  name = "string-pg-db-name"
}

data "aws_ssm_parameter" "db_host" {
  name = "string-write-db-host-url"
}

data "aws_ssm_parameter" "redis_auth_token" {
  name = "redis-auth-token"
}

data "aws_ssm_parameter" "redis_host_url" {
  name  = "redis-host-url"
}

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}

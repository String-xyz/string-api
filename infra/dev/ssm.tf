data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
}

data "aws_ssm_parameter" "hot_wallet" {
  name = "dev-private-key" 
}

data "aws_ssm_parameter" "user_jwt_secret" {
  name = "user-jwt-secret" 
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

data "aws_ssm_parameter" "owlracle_api_secret" {
  name = "dev-owlracle-api-secret" 
}

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}

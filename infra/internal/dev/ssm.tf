data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
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

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}

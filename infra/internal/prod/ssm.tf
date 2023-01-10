data "aws_ssm_parameter" "datadog" {
  name = "datadog-key"
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
  name = "pg-cluster-write-host-url"
}

data "aws_ssm_parameter" "redis_auth_token" {
  name = "redis-auth-token"
}

data "aws_ssm_parameter" "redis_host_url" {
  name = "redis-host-url"
}

data "aws_kms_key" "kms_key" {
  key_id = "alias/main-kms-key"
}

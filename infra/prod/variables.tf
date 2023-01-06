locals {
  cluster_name       = "string-core"
  env                = "prod"
  service_name       = "string-api"
  domain             = "api.string-api.xyz"
  container_port     = "3000"
  origin_id          = "string-api"
  desired_task_count = "1"
  db_port            = "5432"
  redis_port         = "6379"
  memory             = 1024
  cpu                = 512
  region             = "us-west-2"
}

variable "versioning" {
  type    = string
  default = "v1.0.0-alpha"
}

locals {
  task_definition = jsonencode([
    {
      name      = local.service_name
      image     = "${aws_ecr_repository.repo.repository_url}:${var.versioning}"
      essential = true,
      dockerLabels = {
        "com.datadoghq.ad.instances" : "[{\"host\":\"%%host%%\"}]",
        "com.datadoghq.ad.check_names" : "[\"${local.service_name}\"]",
      },
      portMappings = [
        { containerPort = 3000 }
      ],
      secrets = [
        {
          name      = "EVM_PRIVATE_KEY"
          valueFrom = data.aws_ssm_parameter.evm_private_key.arn
        },
        {
          name      = "STRING_ENCRYPTION_KEY"
          valueFrom = data.aws_ssm_parameter.string_encryption_secret.arn
        },
        {
          name      = "STRING_INTERNAL_ID"
          valueFrom = data.aws_ssm_parameter.string_internal_id.arn
        },
        {
          name      = "STRING_WALLET_ID"
          valueFrom = data.aws_ssm_parameter.string_wallet_id.arn
        },
        {
          name      = "STRING_BANK_ID"
          valueFrom = data.aws_ssm_parameter.string_bank_id.arn
        },
        {
          name      = "STRING_PLACEHOLDER_PLATFORM_ID"
          valueFrom = data.aws_ssm_parameter.string_platform_id.arn
        },
        {
          name      = "UNIT21_API_KEY"
          valueFrom = data.aws_ssm_parameter.unit21_api_key.arn
        },
        {
          name      = "IPSTACK_API_KEY"
          valueFrom = data.aws_ssm_parameter.ipstack_api_key.arn
        },
        {
          name      = "CHECKOUT_PUBLIC_KEY"
          valueFrom = data.aws_ssm_parameter.checkout_public_key.arn
        },
        {
          name      = "CHECKOUT_SECRET_KEY"
          valueFrom = data.aws_ssm_parameter.checkout_private_key.arn
        },
        {
          name      = "OWLRACLE_API_KEY"
          valueFrom = data.aws_ssm_parameter.owlracle_api_key.arn
        },
        {
          name      = "OWLRACLE_API_SECRET"
          valueFrom = data.aws_ssm_parameter.owlracle_api_secret.arn
        },
        {
          name      = "FINGERPRINT_API_KEY"
          valueFrom = data.aws_ssm_parameter.fingerprint_api_key.arn
        },
        {
          name      = "SENDGRID_API_KEY"
          valuefrom = data.aws_ssm_parameter.sendgrid_api_key.arn
        },
        {
          name      = "TWILIO_ACCOUNT_SID"
          valueFrom = data.aws_ssm_parameter.twilio_account_sid.arn
        },
        {
          name      = "TWILIO_SMS_SID"
          valuefrom = data.aws_ssm_parameter.twilio_sms_sid.arn
        },
        {
          name      = "TWILIO_AUTH_TOKEN"
          valuefrom = data.aws_ssm_parameter.twilio_auth_token.arn
        },
        {
          name      = "DB_USERNAME"
          valueFrom = data.aws_ssm_parameter.db_username.arn
        },
        {
          name      = "DB_PASSWORD"
          valueFrom = data.aws_ssm_parameter.db_password.arn
        },
        {
          name      = "DB_HOST"
          valueFrom = data.aws_ssm_parameter.db_host.arn
        },
        {
          name      = "DB_NAME"
          valueFrom = data.aws_ssm_parameter.db_name.arn
        },
        {
          name      = "REDIS_HOST",
          valuefrom = data.aws_ssm_parameter.redis_host_url.arn
        },
        {
          name      = "REDIS_PASSWORD",
          valuefrom = data.aws_ssm_parameter.redis_auth_token.arn
        }
      ]
      environment = [
        {
          name  = "PORT"
          value = local.container_port
        },
        {
          name  = "REDIS_PORT"
          value = local.redis_port
        },
        {
          name  = "DB_PORT",
          value = local.db_port
        },
        {
          name  = "ENV"
          value = local.env
        },
        {
          name  = "AWS_REGION"
          value = local.region
        },
        {
          name  = "AWS_KMS_KEY_ID"
          value = data.aws_kms_key.kms_key.key_id
        },
        {
          name  = "OWLRACLE_API_URL"
          value = "https://api.owlracle.info/v3/"
        },
        {
          name  = "COINGECKO_API_URL"
          value = "https://api.coingecko.com/api/v3/"
        },
        {
          name  = "FINGERPRINT_API_URL"
          value = "https://api.fpjs.io/"
        },
        {
          name  = "BASE_URL"
          value = "https://api.string-api.xyz/"
        },
        {
          name  = "UNIT21_ENV"
          value = "api.prod2"
        },
        {
          name  = "UNIT21_ORG_NAME"
          value = "string"
        }
      ],
      logConfiguration = {
        logDriver = "awsfirelens"
        secretOptions = [{
          name      = "apiKey",
          valueFrom = data.aws_ssm_parameter.datadog.arn
        }]
        options = {
          Name             = "datadog"
          "dd_service"     = "${local.service_name}"
          "Host"           = "http-intake.logs.datadoghq.com"
          "dd_source"      = "${local.service_name}"
          "dd_message_key" = "log"
          "dd_tags"        = "project:${local.service_name}"
          "TLS"            = "on"
          "provider"       = "ecs"
        }
      }
    },
    {
      name      = "datadog-agent"
      image     = "public.ecr.aws/datadog/agent:latest"
      essential = true
      secrets = [{
        name      = "DD_API_KEY"
        valueFrom = data.aws_ssm_parameter.datadog.arn
      }],
      environment = [
        {
          name  = "DD_APM_ENABLED"
          value = "true"
        },
        {
          name  = "DD_SITE"
          value = "datadoghq.com"
        },
        {
          name  = "DD_SERVICE"
          value = local.service_name
        },
        {
          name  = "DD_VERSION"
          value = var.versioning
        },
        {
          name  = "DD_ENV"
          value = local.env
        },
        {
          name  = "ECS_FARGATE"
          value = "true"
        },
      ]
      portMappings = [{
        hostPort      = 8126,
        protocol      = "tcp",
        containerPort = 8126
        }
      ]
    },
    {
      name      = "log_router"
      image     = "public.ecr.aws/aws-observability/aws-for-fluent-bit:stable"
      essential = true
      firelensConfiguration = {
        type = "fluentbit"
        options = {
          "enable-ecs-log-metadata" = "true"
        }
      }
    }
  ])
}

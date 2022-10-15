locals {
  cluster_name       = "string-core"
  env                = "dev"
  service_name       = "string-api"
  root_domain        = "dev.string-api.xyz"
  container_port     = "3000"
  origin_id          = "string-api"
  desired_task_count = "1"
  db_port            = "5432"
  memory             = 512
  cpu                = 256
  region             = "us-west-2"
}

variable "versioning" {
  type    = string
  default = "latest"
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
          valueFrom = data.aws_ssm_parameter.hot_wallet.arn
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
        }
      ]
      environment = [
        {
          name  = "PORT"
          value = local.container_port
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
          name  = "ECS_FARGATE"
          value = "true"
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
          name  = "DD_APM_ENABLED"
          value = "true"
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
      image     = "gcr.io/datadoghq/agent:latest"
      essential = false
      secrets = [{
        name      = "DD_API_KEY"
        valueFrom = data.aws_ssm_parameter.datadog.arn
      }],
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

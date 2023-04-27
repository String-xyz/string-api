data "aws_iam_policy_document" "ecs_task_policy" {
  statement {
    sid     = "AllowECSAndTaskAssumeRole"
    actions = ["sts:AssumeRole"]
    effect  = "Allow"
    principals {
      type        = "Service"
      identifiers = ["ecs.amazonaws.com", "ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "task_ecs_role" {
  name               = "${local.service_name}-task-ecs-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_policy.json
}

data "aws_iam_policy_document" "task_policy" {
  statement {
    sid    = "AllowReadToResourcesInListToTask"
    effect = "Allow"
    actions = [
      "ecs:*",
      "ecr:*"
    ]
    
    resources = ["*"]
  }

  statement {
    sid    = "AllowAccessToSSM"
    effect = "Allow"
    actions = [
      "ssm:GetParameters"
    ]
    resources = [
      data.aws_ssm_parameter.datadog.arn,
      data.aws_ssm_parameter.evm_private_key.arn,
      data.aws_ssm_parameter.string_encryption_secret.arn,
      data.aws_ssm_parameter.string_internal_id.arn,
      data.aws_ssm_parameter.string_wallet_id.arn,
      data.aws_ssm_parameter.string_bank_id.arn,
      data.aws_ssm_parameter.string_platform_id.arn,
      data.aws_ssm_parameter.unit21_api_key.arn,
      data.aws_ssm_parameter.checkout_public_key.arn,
      data.aws_ssm_parameter.checkout_private_key.arn,
      data.aws_ssm_parameter.owlracle_api_key.arn,
      data.aws_ssm_parameter.owlracle_api_secret.arn,
      data.aws_ssm_parameter.db_password.arn,
      data.aws_ssm_parameter.db_username.arn,
      data.aws_ssm_parameter.db_name.arn,
      data.aws_ssm_parameter.db_host.arn,
      data.aws_ssm_parameter.redis_host_url.arn,
      data.aws_ssm_parameter.redis_auth_token.arn,
      data.aws_ssm_parameter.fingerprint_api_key.arn,
      data.aws_ssm_parameter.sendgrid_api_key.arn,
      data.aws_ssm_parameter.twilio_sms_sid.arn,
      data.aws_ssm_parameter.twilio_account_sid.arn,
      data.aws_ssm_parameter.twilio_auth_token.arn
    ]
  }

  statement {
    sid    = "AllowDecrypt"
    effect = "Allow"
    actions = [
      "kms:Decrypt"
    ]
    resources = [data.aws_kms_key.kms_key.arn]
  }
}

data "aws_iam_policy_document" "ecs_service_scaling" {

  statement {
    effect = "Allow"

    actions = [
      "application-autoscaling:*",
      "ecs:DescribeServices",
      "ecs:UpdateService",
      "cloudwatch:DescribeAlarms",
      "cloudwatch:PutMetricAlarm",
      "cloudwatch:DeleteAlarms",
      "cloudwatch:DescribeAlarmHistory",
      "cloudwatch:DescribeAlarms",
      "cloudwatch:DescribeAlarmsForMetric",
      "cloudwatch:GetMetricStatistics",
      "cloudwatch:ListMetrics",
      "cloudwatch:PutMetricAlarm",
      "cloudwatch:DisableAlarmActions",
      "cloudwatch:EnableAlarmActions",
      "iam:CreateServiceLinkedRole",
      "sns:CreateTopic",
      "sns:Subscribe",
      "sns:Get*",
      "sns:List*"
    ]

    resources = [
      "*"
    ]
  }
}

resource "aws_iam_policy" "ecs_service_scaling" {
  name = "ecs-to-scaling"
  path = "/"
  description = "Allow ecs service scaling"
  policy = data.aws_iam_policy_document.ecs_service_scaling.json
}

resource "aws_iam_role_policy_attachment" "ecs_service_scaling" {
  role = aws_iam_role.task_ecs_role.name
  policy_arn = aws_iam_policy.ecs_service_scaling.arn
}

resource "aws_iam_role_policy" "task_ecs_policy" {
  name   = "${local.env}-${local.service_name}-task-ecs-policy"
  role   = aws_iam_role.task_ecs_role.id
  policy = data.aws_iam_policy_document.task_policy.json
}

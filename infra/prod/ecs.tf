resource "aws_ecs_cluster" "cluster" {
  name = local.cluster_name
}

resource "aws_ecs_task_definition" "task_definition" {
  container_definitions    = local.task_definition
  family                   = local.service_name
  cpu                      = local.cpu
  memory                   = local.memory
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.task_ecs_role.arn
  task_role_arn            = aws_iam_role.task_ecs_role.arn
}

resource "aws_ecr_repository" "repo" {
  name                 = local.service_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Environment = local.env
    Name        = local.service_name
  }
}

resource "aws_ecs_service" "ecs_service" {
  name            = local.service_name
  task_definition = local.service_name
  desired_count   = local.desired_task_count
  cluster         = aws_ecs_cluster.cluster.name
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = data.terraform_remote_state.vpc.outputs.public_subnets
    security_groups  = [aws_security_group.ecs_task_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    container_name   = local.service_name
    container_port   = local.container_port
    target_group_arn = aws_alb_target_group.ecs_task_target_group.arn
  }

  depends_on = [
    aws_alb_listener_rule.ecs_alb_listener_rule,
    aws_iam_role_policy.task_ecs_policy,
    aws_ecs_task_definition.task_definition
  ]

  tags = {
    Environment = local.env
    Name        = local.service_name
  }
}

resource "aws_appautoscaling_target" "ecs_target" {
  max_capacity       = 20
  min_capacity       = 1
  resource_id        = "service/${aws_ecs_cluster.cluster.name}/${aws_ecs_service.ecs_service.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  service_namespace  = "ecs"
}

resource "aws_appautoscaling_policy" "ecs_target_cpu" {
  name               = "${local.service_name}-scaling-policy-cpu"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.ecs_target.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs_target.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value = 60
  }
  depends_on = [aws_appautoscaling_target.ecs_target]
}

resource "aws_appautoscaling_policy" "ecs_target_memory" {
  name               = "${local.service_name}-scaling-policy-memory"
  policy_type        = "TargetTrackingScaling"
  resource_id        = aws_appautoscaling_target.ecs_target.resource_id
  scalable_dimension = aws_appautoscaling_target.ecs_target.scalable_dimension
  service_namespace  = aws_appautoscaling_target.ecs_target.service_namespace

  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
    target_value = 60
  }
  depends_on = [aws_appautoscaling_target.ecs_target]
}

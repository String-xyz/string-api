resource "aws_security_group" "ecs_alb_https_sg" {
  name        = "${local.service_name}-alb-https-sg"
  description = "Security group for ALB to cluster"
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
  
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${local.service_name}-alb-https-sg"
    Environment = local.env
  }
}

resource "aws_security_group" "ecs_task_sg" {
  name   = "${local.service_name}-task-sg"
  vpc_id = data.terraform_remote_state.vpc.outputs.id
  ingress {
    from_port   = local.container_port
    to_port     = local.container_port
    protocol    = "TCP"
    cidr_blocks = [data.terraform_remote_state.vpc.outputs.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${local.service_name}-task-sg"
    environment = local.env
  }
}

# Give access to DB through Security group rule
data "aws_security_group" "rds" {
  name   = "pg-cluster-20221103192516479600000004"
  vpc_id = data.terraform_remote_state.vpc.outputs.id
}

data "aws_security_group" "redis" {
  name        = "redis-client-redis"
  vpc_id      = data.terraform_remote_state.vpc.outputs.id
}

resource "aws_security_group_rule" "rds_to_ecs" {
  type                     = "ingress"
  protocol                 = "TCP"
  from_port                = local.db_port
  to_port                  = local.db_port
  source_security_group_id = aws_security_group.ecs_task_sg.id
  security_group_id        = data.aws_security_group.rds.id
}

resource "aws_security_group_rule" "redis_to_ecs" {
  type                     = "ingress"
  protocol                 = "TCP"
  from_port                = local.redis_port
  to_port                  = local.redis_port
  source_security_group_id = aws_security_group.ecs_task_sg.id
  security_group_id        = data.aws_security_group.redis.id
}

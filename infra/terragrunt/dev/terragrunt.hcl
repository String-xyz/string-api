include { 
    path = find_in_parent_folders()
}

terraform { 
    source = "../../.."
}

locals { 
  cluster_name       = "string-core"
  env                = "dev"
  service_name       = "string-api"
  root_domain        = "dev.string-api.xyz"
  container_port     = "3000"
  origin_id          = "string-api"
  desired_task_count = "1"
  db_port            = "5432"
  redis_port         = "6379"
  memory             = 512
  cpu                = 256
  region             = "us-west-2"
}

inputs { 
 cluster_name        = local.cluster_name
  env                = local.env
  service_name       = local.service_name
  root_domain        = local.root_domain
  container_port     = local.container_port
  origin_id          = local.origin_id
  desired_task_count = local.desired_task_count
  db_port            = local.db_port
  redis_port         = local.redis_port
  memory             = local.memory
  cpu                = local.cpu
  region             = local.region
}
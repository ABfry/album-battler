terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  # S3バックエンド設定
  # backend "s3" {
  #   bucket         = "album-battler-terraform-state"
  #   key            = "terraform.tfstate"
  #   region         = "ap-northeast-1"
  #   encrypt        = true
  #   dynamodb_table = "terraform-state-lock"
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project_name
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# VPC
resource "aws_vpc" "main" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.project_name}-vpc"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id

  tags = {
    Name = "${var.project_name}-igw"
  }
}

# Public Subnets
resource "aws_subnet" "public" {
  count                   = length(var.availability_zones)
  vpc_id                  = aws_vpc.main.id
  cidr_block              = var.public_subnet_cidrs[count.index]
  availability_zone       = var.availability_zones[count.index]
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.project_name}-public-subnet-${count.index + 1}"
  }
}

# Route Table
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }

  tags = {
    Name = "${var.project_name}-public-rt"
  }
}

# Route Table Association
resource "aws_route_table_association" "public" {
  count          = length(var.availability_zones)
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# S3ストレージモジュール
module "s3_storage" {
  source = "./modules/s3_storage"

  bucket_name        = var.s3_bucket_name
  region             = var.aws_region
  enable_versioning  = var.s3_enable_versioning
  allowed_origins    = var.s3_allowed_origins

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

# RDSモジュール
module "rds" {
  source = "./modules/rds"

  project_name       = var.project_name
  environment        = var.environment
  subnet_ids         = aws_subnet.public[*].id
  security_group_ids = [module.ecs.ecs_tasks_security_group_id]

  db_name     = var.db_name
  db_username = var.db_username
  db_password = var.db_password

  instance_class    = var.rds_instance_class
  allocated_storage = var.rds_allocated_storage
  multi_az          = var.rds_multi_az

  skip_final_snapshot     = var.rds_skip_final_snapshot
  backup_retention_period = var.rds_backup_retention_period

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}


# ECSモジュール
module "ecs" {
  source = "./modules/ecs"

  project_name       = var.project_name
  aws_region         = var.aws_region
  vpc_id             = aws_vpc.main.id
  public_subnet_ids  = aws_subnet.public[*].id
  s3_bucket_name     = var.s3_bucket_name

  log_retention_days = var.ecs_log_retention_days

  frontend_cpu           = var.frontend_cpu
  frontend_memory        = var.frontend_memory
  frontend_desired_count = var.frontend_desired_count

  backend_cpu           = var.backend_cpu
  backend_memory        = var.backend_memory
  backend_desired_count = var.backend_desired_count

  frontend_environment_variables = [
    {
      name  = "NEXT_PUBLIC_API_URL"
      value = "http://${module.ecs.alb_dns_name}"
    },
    {
      name  = "NEXT_PUBLIC_WEBSOCKET_URL"
      value = "ws://${module.ecs.alb_dns_name}/ws"
    }
  ]

  backend_environment_variables = [
    {
      name  = "PORT"
      value = "8080"
    },
    {
      name  = "DB_DRIVER"
      value = "mysql"
    },
    {
      name  = "DB_ENV"
      value = "production"
    },
    {
      name  = "DB_HOST"
      value = module.rds.endpoint
    },
    {
      name  = "DB_PORT"
      value = "3306"
    },
    {
      name  = "DB_NAME"
      value = var.db_name
    },
    {
      name  = "DB_USER"
      value = var.db_username
    },
    {
      name  = "DB_PASSWORD"
      value = var.db_password
    },
    {
      name  = "DB_PARAMS"
      value = "charset=utf8mb4&parseTime=true&loc=Local"
    },
    {
      name  = "DB_MAX_OPEN_CONNS"
      value = "32"
    },
    {
      name  = "DB_MAX_IDLE_CONNS"
      value = "16"
    },
    {
      name  = "DB_CONN_MAX_LIFETIME"
      value = "5400"
    },
    {
      name  = "DB_CONN_MAX_IDLE_TIME"
      value = "900"
    },
    {
      name  = "AWS_REGION"
      value = var.aws_region
    },
    {
      name  = "S3_BUCKET"
      value = var.s3_bucket_name
    },
    {
      name  = "GEMINI_API_KEY"
      value = var.gemini_api_key
    }
  ]

  tags = {
    Project     = var.project_name
    Environment = var.environment
  }
}

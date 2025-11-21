variable "aws_region" {
  description = "AWSリージョン"
  type        = string
  default     = "ap-northeast-1"
}

variable "project_name" {
  description = "プロジェクト名"
  type        = string
  default     = "album-battler"
}

variable "environment" {
  description = "環境名"
  type        = string
  default     = "production"
}

# S3設定
variable "s3_bucket_name" {
  description = "S3バケット名"
  type        = string
  default     = "album-battler-images"
}

variable "s3_enable_versioning" {
  description = "S3バージョニングを有効化"
  type        = bool
  default     = true
}

variable "s3_allowed_origins" {
  description = "CORS許可オリジン"
  type        = list(string)
  default     = ["*"]
}

# VPC設定
variable "vpc_cidr" {
  description = "VPC CIDR"
  type        = string
  default     = "10.0.0.0/16"
}

variable "availability_zones" {
  description = "アベイラビリティゾーン"
  type        = list(string)
  default     = ["ap-northeast-1a", "ap-northeast-1c"]
}

variable "public_subnet_cidrs" {
  description = "パブリックサブネットのCIDR"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

# RDS設定
variable "db_name" {
  description = "データベース名"
  type        = string
  default     = "album_battler"
}

variable "db_username" {
  description = "データベースユーザー名"
  type        = string
  default     = "album_user"
}

variable "db_password" {
  description = "データベースパスワード"
  type        = string
  sensitive   = true
}

variable "rds_instance_class" {
  description = "RDSインスタンスクラス"
  type        = string
  default     = "db.t4g.micro"
}

variable "rds_allocated_storage" {
  description = "割り当てストレージ (GB)"
  type        = number
  default     = 20
}

variable "rds_multi_az" {
  description = "Multi-AZ配置"
  type        = bool
  default     = false
}

variable "rds_skip_final_snapshot" {
  description = "削除時の最終スナップショットをスキップ"
  type        = bool
  default     = true
}

variable "rds_backup_retention_period" {
  description = "バックアップ保持期間 (日)"
  type        = number
  default     = 7
}

# ECS設定
variable "ecs_log_retention_days" {
  description = "CloudWatch Logsの保持期間 (日)"
  type        = number
  default     = 7
}

variable "frontend_cpu" {
  description = "Frontendタスクのvirtual CPU units"
  type        = string
  default     = "512"
}

variable "frontend_memory" {
  description = "Frontendタスクのメモリ (MB)"
  type        = string
  default     = "1024"
}

variable "frontend_desired_count" {
  description = "Frontendタスクの希望数"
  type        = number
  default     = 1
}

variable "backend_cpu" {
  description = "Backendタスクのvirtual CPU units"
  type        = string
  default     = "1024"
}

variable "backend_memory" {
  description = "Backendタスクのメモリ (MB)"
  type        = string
  default     = "2048"
}

variable "backend_desired_count" {
  description = "Backendタスクの希望数"
  type        = number
  default     = 1
}

# API Keys
variable "gemini_api_key" {
  description = "Gemini API Key"
  type        = string
  sensitive   = true
}

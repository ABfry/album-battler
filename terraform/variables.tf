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

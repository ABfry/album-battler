variable "project_name" {
  description = "プロジェクト名"
  type        = string
}

variable "aws_region" {
  description = "AWSリージョン"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "パブリックサブネットID"
  type        = list(string)
}

variable "s3_bucket_name" {
  description = "S3バケット名"
  type        = string
}

variable "log_retention_days" {
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

variable "frontend_environment_variables" {
  description = "Frontend環境変数"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
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

variable "backend_environment_variables" {
  description = "Backend環境変数"
  type = list(object({
    name  = string
    value = string
  }))
  default = []
}

variable "tags" {
  description = "タグ"
  type        = map(string)
  default     = {}
}

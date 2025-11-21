variable "project_name" {
  description = "プロジェクト名"
  type        = string
}

variable "environment" {
  description = "環境名"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "subnet_ids" {
  description = "サブネットID"
  type        = list(string)
}

variable "ecs_tasks_security_group_id" {
  description = "ECSタスクのセキュリティグループID（RDSへのアクセスを許可）"
  type        = string
}

variable "instance_class" {
  description = "RDSインスタンスクラス"
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "割り当てストレージ (GB)"
  type        = number
  default     = 20
}

variable "max_allocated_storage" {
  description = "最大割り当てストレージ (GB)"
  type        = number
  default     = 100
}

variable "db_name" {
  description = "データベース名"
  type        = string
}

variable "db_username" {
  description = "データベースユーザー名"
  type        = string
}

variable "db_password" {
  description = "データベースパスワード"
  type        = string
  sensitive   = true
}

variable "multi_az" {
  description = "Multi-AZ配置"
  type        = bool
  default     = false
}

variable "skip_final_snapshot" {
  description = "削除時の最終スナップショットをスキップ"
  type        = bool
  default     = true
}

variable "backup_retention_period" {
  description = "バックアップ保持期間 (日)"
  type        = number
  default     = 7
}

variable "tags" {
  description = "タグ"
  type        = map(string)
  default     = {}
}

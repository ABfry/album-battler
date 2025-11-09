variable "bucket_name" {
  description = "S3バケット名"
  type        = string
}

variable "region" {
  description = "AWSリージョン"
  type        = string
  default     = "ap-northeast-1"
}

variable "enable_versioning" {
  description = "バケットのバージョニングを有効化するか"
  type        = bool
  default     = true
}

variable "lifecycle_rules_enabled" {
  description = "ライフサイクルルールを有効化するか"
  type        = bool
  default     = false
}

variable "expiration_days" {
  description = "オブジェクトの自動削除日数（lifecycle_rules_enabled=trueの場合）"
  type        = number
  default     = 90
}

variable "tags" {
  description = "リソースに付与するタグ"
  type        = map(string)
  default = {
    Project = "album-battler"
    Service = "image-storage"
  }
}

variable "allowed_origins" {
  description = "CORS設定で許可するオリジン"
  type        = list(string)
  default     = ["*"]
}

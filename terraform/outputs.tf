output "vpc_id" {
  description = "VPC ID"
  value       = aws_vpc.main.id
}

output "vpc_cidr" {
  description = "VPC CIDR"
  value       = aws_vpc.main.cidr_block
}

output "public_subnet_ids" {
  description = "パブリックサブネットID"
  value       = aws_subnet.public[*].id
}

output "s3_bucket_name" {
  description = "S3バケット名"
  value       = module.s3_storage.bucket_name
}

output "s3_bucket_arn" {
  description = "S3バケットARN"
  value       = module.s3_storage.bucket_arn
}

output "s3_bucket_domain_name" {
  description = "S3バケットドメイン名"
  value       = module.s3_storage.bucket_domain_name
}

output "s3_bucket_url" {
  description = "S3バケットURL"
  value       = "https://${module.s3_storage.bucket_name}.s3.${var.aws_region}.amazonaws.com"
}

output "s3_bucket_region" {
  description = "S3バケットリージョン"
  value       = module.s3_storage.bucket_region
}

# RDS Outputs
output "rds_endpoint" {
  description = "RDSエンドポイント"
  value       = module.rds.endpoint
}

output "rds_address" {
  description = "RDSアドレス"
  value       = module.rds.address
}

output "rds_port" {
  description = "RDSポート"
  value       = module.rds.port
}

# ECS Outputs
output "ecs_cluster_name" {
  description = "ECSクラスター名"
  value       = module.ecs.cluster_name
}

output "alb_dns_name" {
  description = "ALB DNS名（アプリケーションURL）"
  value       = module.ecs.alb_dns_name
}

output "alb_url" {
  description = "アプリケーションURL"
  value       = "http://${module.ecs.alb_dns_name}"
}

output "frontend_ecr_repository_url" {
  description = "Frontend ECRリポジトリURL"
  value       = module.ecs.frontend_ecr_repository_url
}

output "backend_ecr_repository_url" {
  description = "Backend ECRリポジトリURL"
  value       = module.ecs.backend_ecr_repository_url
}

output "cluster_id" {
  description = "ECSクラスターID"
  value       = aws_ecs_cluster.main.id
}

output "cluster_name" {
  description = "ECSクラスター名"
  value       = aws_ecs_cluster.main.name
}

output "alb_dns_name" {
  description = "ALB DNS名"
  value       = aws_lb.main.dns_name
}

output "alb_zone_id" {
  description = "ALB Zone ID"
  value       = aws_lb.main.zone_id
}

output "alb_arn" {
  description = "ALB ARN"
  value       = aws_lb.main.arn
}

output "frontend_ecr_repository_url" {
  description = "Frontend ECRリポジトリURL"
  value       = aws_ecr_repository.frontend.repository_url
}

output "backend_ecr_repository_url" {
  description = "Backend ECRリポジトリURL"
  value       = aws_ecr_repository.backend.repository_url
}

output "frontend_service_name" {
  description = "Frontend ECSサービス名"
  value       = aws_ecs_service.frontend.name
}

output "backend_service_name" {
  description = "Backend ECSサービス名"
  value       = aws_ecs_service.backend.name
}

output "ecs_tasks_security_group_id" {
  description = "ECSタスクセキュリティグループID"
  value       = aws_security_group.ecs_tasks.id
}

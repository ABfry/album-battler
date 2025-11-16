output "endpoint" {
  description = "RDSエンドポイント（ホスト名）"
  value       = split(":", aws_db_instance.main.endpoint)[0]
}

output "address" {
  description = "RDSアドレス"
  value       = aws_db_instance.main.address
}

output "port" {
  description = "RDSポート"
  value       = aws_db_instance.main.port
}

output "db_name" {
  description = "データベース名"
  value       = aws_db_instance.main.db_name
}

output "arn" {
  description = "RDS ARN"
  value       = aws_db_instance.main.arn
}

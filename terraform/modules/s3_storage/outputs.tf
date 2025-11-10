output "bucket_name" {
  description = "S3バケット名"
  value       = aws_s3_bucket.images.id
}

output "bucket_arn" {
  description = "S3バケットARN"
  value       = aws_s3_bucket.images.arn
}

output "bucket_domain_name" {
  description = "S3バケットのドメイン名"
  value       = aws_s3_bucket.images.bucket_domain_name
}

output "bucket_regional_domain_name" {
  description = "S3バケットのリージョナルドメイン名"
  value       = aws_s3_bucket.images.bucket_regional_domain_name
}

output "bucket_region" {
  description = "S3バケットのリージョン"
  value       = aws_s3_bucket.images.region
}

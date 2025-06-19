output "cloudfront_domain_name" {
  description = "El nombre de dominio de la distribución de CloudFront."
  value       = aws_cloudfront_distribution.cdn.domain_name
}

output "s3_bucket_name" {
  description = "El nombre del bucket S3 para subir los archivos del frontend."
  value       = aws_s3_bucket.frontend_bucket.id
}

output "cloudfront_distribution_id" {
  description = "ID de la distribución CloudFront"
  value       = aws_cloudfront_distribution.cdn.id
}
# terraform-project/outputs.tf

output "instrucciones_despliegue" {
  value = <<-EOT
  Infraestructura creada. Siguientes pasos manuales:

  1. BACKEND: Construye, etiqueta y pushea tu imagen de Docker a ECR:
     aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${module.backend.ecr_repository_url}
     docker build -t app-backend .
     docker tag app-backend:latest ${module.backend.ecr_repository_url}:latest
     docker push ${module.backend.ecr_repository_url}:latest

     (Después del push, App Runner se desplegará automáticamente)

  2. FRONTEND: Sube tu carpeta 'dist' al bucket S3:
     aws s3 sync ./path/to/your/dist s3://${module.frontend.s3_bucket_name} --delete
  EOT
}

output "frontend_url" {
  description = "URL del Frontend (CloudFront)"
  value       = "https://${module.frontend.cloudfront_domain_name}"
}

output "backend_url" {
  description = "URL del Backend (App Runner)"
  value       = module.backend.app_runner_service_url
}
# terraform-project/outputs.tf

output "instrucciones_despliegue" {
  value = <<-EOT
  Infraestructura creada. Siguientes pasos manuales:

  1. BACKEND: Construye, etiqueta y pushea tu imagen de Docker a ECR:
     aws ecr get-login-password --region ${var.aws_region} | docker login --username AWS --password-stdin ${module.backend.ecr_repository_url}
     docker build -t app-backend ./backend
     docker tag app-backend:latest ${module.backend.ecr_repository_url}:latest
     docker push ${module.backend.ecr_repository_url}:latest

     (Después del push, App Runner se desplegará automáticamente)

  2. FRONTEND: Configura la URL del backend y construye:
     cd frontend
     export VITE_API_URL="${module.backend.app_runner_service_url}/api"
     npm run build
     aws s3 sync ./dist s3://${module.frontend.s3_bucket_name} --delete
     aws cloudfront create-invalidation --distribution-id ${module.frontend.cloudfront_distribution_id} --paths "/*"
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

output "aws_region" {
  description = "Región de AWS"
  value       = var.aws_region
}

# Outputs para el script de deploy
output "backend" {
  description = "Backend outputs"
  value = {
    ecr_repository_url = module.backend.ecr_repository_url
  }
}

output "frontend" {
  description = "Frontend outputs"
  value = {
    s3_bucket_name = module.frontend.s3_bucket_name
    cloudfront_distribution_id = module.frontend.cloudfront_distribution_id
    cloudfront_domain_name = module.frontend.cloudfront_domain_name
  }
}
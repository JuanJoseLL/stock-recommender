output "app_runner_service_url" {
  description = "La URL pública del servicio de backend."
  value       = aws_apprunner_service.backend_service.service_url
}

output "ecr_repository_url" {
  description = "La URL del repositorio ECR para pushear la imagen."
  value       = aws_ecr_repository.backend_repo.repository_url
}
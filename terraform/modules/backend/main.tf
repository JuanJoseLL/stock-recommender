resource "aws_ecr_repository" "backend_repo" {
  name                 = "${var.project_name}-backend"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_iam_role" "apprunner_ecr_role" {
  name = "${var.project_name}-apprunner-ecr-role"
  assume_role_policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "build.apprunner.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "apprunner_ecr_policy" {
  role       = aws_iam_role.apprunner_ecr_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}

# Instance role for App Runner to access secrets
resource "aws_iam_role" "apprunner_instance_role" {
  name = "${var.project_name}-apprunner-instance-role"
  assume_role_policy = jsonencode({
    Version   = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "tasks.apprunner.amazonaws.com"
      }
    }]
  })
}

# Policy to allow reading SSM parameters
resource "aws_iam_role_policy" "apprunner_secrets_policy" {
  name = "${var.project_name}-apprunner-secrets-policy"
  role = aws_iam_role.apprunner_instance_role.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "ssm:GetParameter",
        "ssm:GetParameters"
      ]
      Resource = var.secret_db_url_arn
    }]
  })
}

resource "aws_apprunner_service" "backend_service" {
  service_name = "${var.project_name}-backend-service"
  
  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.apprunner_ecr_role.arn
    }
    image_repository {
      image_identifier      = "${aws_ecr_repository.backend_repo.repository_url}:latest" 
      image_repository_type = "ECR"
      image_configuration {
        port = "8080" 
        
        runtime_environment_secrets = {
          DATABASE_URL = var.secret_db_url_arn
        }
      }
    }
    auto_deployments_enabled = true
  }

  instance_configuration {
    instance_role_arn = aws_iam_role.apprunner_instance_role.arn
  }

  tags = {
    Project = var.project_name
  }
}
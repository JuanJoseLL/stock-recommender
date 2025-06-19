terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# Obtener la VPC por defecto
data "aws_vpc" "default" {
  default = true
}

# Obtener la primera subnet disponible
data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Módulo de Security Groups
module "security" {
  source = "./modules/security"
  
  project_name = var.project_name
  vpc_id       = data.aws_vpc.default.id
}

# Clave SSH para la instancia EC2
resource "tls_private_key" "main" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "main" {
  key_name   = "${var.project_name}-key"
  public_key = tls_private_key.main.public_key_openssh
}

resource "local_file" "private_key" {
  content  = tls_private_key.main.private_key_pem
  filename = "${path.module}/${var.project_name}-key.pem"
  file_permission = "0400"
}

# Módulo de EC2
module "ec2" {
  source = "./modules/ec2"
  
  project_name       = var.project_name
  instance_type      = var.instance_type
  subnet_id          = data.aws_subnets.default.ids[0]
  security_group_ids = [module.security.web_security_group_id]
  key_name           = aws_key_pair.main.key_name
}

 
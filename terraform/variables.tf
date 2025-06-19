variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "stock-recommender"
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "admin_username" {
  description = "Admin username for EC2 instance"
  type        = string
  default     = "admin"
}

variable "admin_password" {
  description = "Admin password for EC2 instance"
  type        = string
  sensitive   = true
} 
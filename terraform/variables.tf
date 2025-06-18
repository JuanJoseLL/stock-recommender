# terraform-project/variables.tf

variable "project_name" {
  description = "El nombre del proyecto, usado para nombrar recursos."
  type        = string
  default     = "stock-recommender"
}

variable "aws_region" {
  description = "La región de AWS para desplegar los recursos."
  type        = string
  default     = "us-east-1"
}
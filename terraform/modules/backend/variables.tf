variable "project_name" {
  description = "Nombre del proyecto para las etiquetas y nombres."
  type        = string
}

variable "secret_db_url_arn" {
  description = "ARN del secreto en AWS Secrets Manager que contiene la URL de la BD."
  type        = string
}
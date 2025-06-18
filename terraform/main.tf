# terraform-project/main.tf

# Obtenemos el secreto de la BD que creamos manualmente
data "aws_ssm_parameter" "db_url" {
  name = "/stock-recommender/production/database/url" 
  with_decryption = true
}

module "backend" {
  source = "./modules/backend"

  project_name      = var.project_name
  secret_db_url_arn = data.aws_ssm_parameter.db_url.arn
}

module "frontend" {
  source = "./modules/frontend"

  project_name = var.project_name
}
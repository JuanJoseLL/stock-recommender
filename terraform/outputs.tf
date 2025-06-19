output "vpc_id" {
  description = "ID of the VPC"
  value       = data.aws_vpc.default.id
}

output "public_subnet_ids" {
  description = "IDs of the public subnets"
  value       = data.aws_subnets.default.ids
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = module.ec2.instance_id
}

output "instance_public_ip" {
  description = "Public IP address of the EC2 instance, to be used for SSH and in the Ansible inventory."
  value       = module.ec2.public_ip
}

output "private_ip" {
  description = "Private IP address of the EC2 instance"
  value       = module.ec2.private_ip
}

output "security_group_id" {
  description = "ID of the security group"
  value       = module.security.web_security_group_id
}

output "admin_username" {
  description = "Admin username for SSH access"
  value       = var.admin_username
}

output "ssh_connection" {
  description = "SSH connection command"
  value       = "ssh ${var.admin_username}@${module.ec2.public_ip}"
}

output "application_url" {
  description = "URL to access the application"
  value       = "http://${module.ec2.public_ip}"
}

output "frontend_url" {
  description = "URL to access the frontend directly"
  value       = "http://${module.ec2.public_ip}:3000"
}

output "backend_url" {
  description = "URL to access the backend API directly"
  value       = "http://${module.ec2.public_ip}:8080"
}

output "public_ip" {
  description = "Public IP of the EC2 instance"
  value       = module.ec2.public_ip
}

output "ssh_private_key_filename" {
  description = "Filename for the private key to SSH into the EC2 instance."
  value       = local_file.private_key.filename
}

output "ssh_connection_command" {
  description = "Example command to SSH into the EC2 instance."
  value       = "ssh -i ${local_file.private_key.filename} ubuntu@${module.ec2.public_ip}"
} 
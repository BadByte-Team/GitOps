variable "region" {
  description = "Region de AWS"
  type        = string
  default     = "us-east-1"
}

variable "bucket_name" {
  description = "Nombre del bucket S3 para el state de Terraform"
  type        = string
  default     = "curso-gitops-terraform-state"
}

variable "dynamodb_table_name" {
  description = "Nombre de la tabla DynamoDB para state locking"
  type        = string
  default     = "curso-gitops-terraform-locks"
}

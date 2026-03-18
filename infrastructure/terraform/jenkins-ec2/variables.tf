variable "region" {
  description = "Region de AWS"
  type        = string
  default     = "us-east-1"
}

variable "ami_id" {
  description = "AMI de Ubuntu 22.04 LTS (us-east-1)"
  type        = string
  default     = "ami-0c7217cdde317cfec"
}

variable "instance_type" {
  description = "Tipo de instancia para Jenkins"
  type        = string
  default     = "t2.medium"
}

variable "key_name" {
  description = "Nombre del Key Pair en AWS"
  type        = string
}

variable "allowed_cidr" {
  description = "CIDR permitido para SSH (tu IP)"
  type        = string
  default     = "0.0.0.0/0"
}

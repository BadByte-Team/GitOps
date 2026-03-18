variable "region" {
  description = "Region de AWS"
  type        = string
  default     = "us-east-1"
}

variable "cluster_name" {
  description = "Nombre del cluster EKS"
  type        = string
  default     = "curso-gitops-eks"
}

variable "kubernetes_version" {
  description = "Version de Kubernetes para EKS"
  type        = string
  default     = "1.29"
}

variable "node_instance_types" {
  description = "Tipos de instancia para los nodos"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "node_desired_size" {
  description = "Cantidad deseada de nodos"
  type        = number
  default     = 2
}

variable "node_max_size" {
  description = "Maximo de nodos"
  type        = number
  default     = 3
}

variable "node_min_size" {
  description = "Minimo de nodos"
  type        = number
  default     = 1
}

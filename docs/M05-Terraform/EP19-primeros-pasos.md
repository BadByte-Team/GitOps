# EP 19: Primeros Pasos — Provider, Resource, Variables

**Tipo:** PRACTICA
## Objetivo
Escribir tu primer archivo Terraform: crear una EC2 desde cero con provider, resource y variables.

## Pasos Detallados

### Estructura de Archivos
```
mi-primera-ec2/
├── main.tf          # Recursos
├── variables.tf     # Variables de entrada
└── outputs.tf       # Valores de salida
```

### `main.tf`
```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

resource "aws_instance" "ejemplo" {
  ami           = var.ami_id
  instance_type = var.instance_type

  tags = {
    Name = "Mi-Primera-EC2"
  }
}
```

### `variables.tf`
```hcl
variable "region" {
  default = "us-east-1"
}
variable "ami_id" {
  default = "ami-0c7217cdde317cfec"
}
variable "instance_type" {
  default = "t2.micro"
}
```

### `outputs.tf`
```hcl
output "instance_ip" {
  value = aws_instance.ejemplo.public_ip
}
```

### Ejecutar
```bash
terraform init      # Descargar providers
terraform plan      # Ver qué se va a crear
terraform apply     # Crear recursos (escribir "yes")
terraform destroy   # Eliminar todo
```

## Verificación
- [ ] `terraform init` descarga el provider de AWS
- [ ] `terraform plan` muestra 1 recurso a crear
- [ ] `terraform apply` crea la EC2 correctamente
- [ ] `terraform destroy` la elimina

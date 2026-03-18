# EP 20: Backend Remoto — State File en S3

**Tipo:** CONFIGURACION
## Objetivo
Configurar Terraform para guardar el state file en S3 con DynamoDB locking.

## Prerequisitos
- Bucket S3 y tabla DynamoDB creados (EP17)

## Agregar Backend al Proyecto

### `backend.tf`
```hcl
terraform {
  backend "s3" {
    bucket         = "curso-gitops-terraform-state"
    key            = "mi-proyecto/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "curso-gitops-terraform-locks"
    encrypt        = true
  }
}
```

### Migrar State Local a Remoto
```bash
terraform init
# Terraform detecta el cambio de backend y pregunta si migrar
# Responder: yes
```

## Archivos Involucrados
- Todos los directorios en `infrastructure/terraform/*/` ya tienen backend configurado
- Cada uno usa un `key` diferente para no sobreescribirse

## Verificación
- [ ] `terraform init` migra el state a S3
- [ ] En S3 aparece el archivo `.tfstate`
- [ ] Dos personas no pueden hacer `apply` al mismo tiempo (locking)

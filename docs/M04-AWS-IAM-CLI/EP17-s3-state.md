# EP 17: S3 — Crear Bucket para Terraform State

**Tipo:** CONFIGURACION
## Objetivo
Crear un bucket S3 y una tabla DynamoDB para almacenar y proteger el state file de Terraform.

## Pasos Detallados

### Opción A: Manual con AWS CLI
```bash
# Crear bucket S3
aws s3api create-bucket --bucket curso-gitops-terraform-state --region us-east-1

# Activar versionado
aws s3api put-bucket-versioning --bucket curso-gitops-terraform-state --versioning-configuration Status=Enabled

# Activar encriptación
aws s3api put-bucket-encryption --bucket curso-gitops-terraform-state --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

# Bloquear acceso público
aws s3api put-public-access-block --bucket curso-gitops-terraform-state --public-access-block-configuration BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true

# Crear tabla DynamoDB para locking
aws dynamodb create-table   --table-name curso-gitops-terraform-locks   --attribute-definitions AttributeName=LockID,AttributeType=S   --key-schema AttributeName=LockID,KeyType=HASH   --billing-mode PAY_PER_REQUEST
```

### Opción B: Con Terraform (recomendado)
```bash
cd infrastructure/terraform/backend
terraform init
terraform plan
terraform apply
```

## Archivos Involucrados
- `infrastructure/terraform/backend/main.tf`
- `infrastructure/terraform/backend/variables.tf`
- `infrastructure/terraform/backend/outputs.tf`

## Verificación
- [ ] El bucket S3 existe y tiene versionado activado
- [ ] La tabla DynamoDB `curso-gitops-terraform-locks` existe
- [ ] El acceso público está bloqueado

## Notas
> Este bucket es el **único** recurso que se crea antes de usar remote backend. Después, todos los demás Terraform configs usan este bucket.

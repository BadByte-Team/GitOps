# EP 28: Crear Cluster EKS con Terraform

**Tipo:** PRACTICA
## Objetivo
Usar los archivos Terraform del proyecto para crear un cluster EKS completo.

## Prerequisitos
- AWS CLI configurado (EP15)
- Backend S3 configurado (EP17)
- Terraform instalado (EP18)

## Pasos Detallados

### 1. Revisar los Archivos
```bash
ls infrastructure/terraform/eks-cluster/
# main.tf       — VPC, subnets, IAM roles, EKS cluster, node group
# variables.tf  — Configuración del cluster
# outputs.tf    — Endpoint, comando kubeconfig
```

### 2. Inicializar y Aplicar
```bash
cd infrastructure/terraform/eks-cluster

terraform init
terraform plan
# Revisar que va a crear: VPC, subnets, IGW, EKS cluster, node group

terraform apply
# Tarda ~15-20 minutos
```

### 3. Verificar en AWS
```bash
aws eks list-clusters
aws eks describe-cluster --name curso-gitops-eks --query "cluster.status"
# Output: "ACTIVE"
```

## Archivos Involucrados
- `infrastructure/terraform/eks-cluster/main.tf`
- `infrastructure/terraform/eks-cluster/variables.tf`
- `infrastructure/terraform/eks-cluster/outputs.tf`

## Verificación
- [ ] `terraform apply` completa sin errores
- [ ] El cluster aparece como ACTIVE en AWS
- [ ] Los nodos están en estado Ready

## Notas
> ⚠️ **COSTO**: El cluster EKS cuesta ~$0.10/hr. **Destruye al terminar**: `terraform destroy`

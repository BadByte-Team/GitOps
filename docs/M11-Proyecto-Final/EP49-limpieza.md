# EP 49: Limpieza de Recursos — Evitar Costos en AWS

**Tipo:** PRACTICA
## Objetivo
Eliminar todos los recursos de AWS para evitar cargos innecesarios.

## Orden de Eliminación

> ⚠️ **IMPORTANTE**: Seguir este orden exacto para evitar dependencias huérfanas.

### 1. Eliminar App en ArgoCD
```bash
kubectl delete application curso-gitops -n argocd
kubectl delete namespace curso-gitops
```

### 2. Desinstalar ArgoCD
```bash
kubectl delete namespace argocd
```

### 3. Destruir EKS con Terraform
```bash
cd infrastructure/terraform/eks-cluster
terraform destroy -auto-approve
# Tarda ~15 minutos
```

### 4. Destruir Jenkins EC2 con Terraform
```bash
cd infrastructure/terraform/jenkins-ec2
terraform destroy -auto-approve
```

### 5. (Opcional) Eliminar Backend S3
```bash
# Solo si ya no vas a usar Terraform
aws s3 rb s3://curso-gitops-terraform-state --force
aws dynamodb delete-table --table-name curso-gitops-terraform-locks
```

### O Usar el Script de Limpieza
```bash
chmod +x infrastructure/scripts/cleanup.sh
./infrastructure/scripts/cleanup.sh
```

## Verificar que No Quedan Recursos
```bash
# EC2
aws ec2 describe-instances --filters "Name=instance-state-name,Values=running" --query "Reservations[].Instances[].InstanceId"

# EKS
aws eks list-clusters

# Billing
# Ir a AWS Console → Billing → ver que no hay cargos activos
```

## Archivos Involucrados
- `infrastructure/scripts/cleanup.sh`
- `infrastructure/jenkins/Jenkinsfile-destroy`
- `infrastructure/terraform/eks-cluster/`
- `infrastructure/terraform/jenkins-ec2/`

## Verificación
- [ ] No hay instancias EC2 corriendo
- [ ] No hay clusters EKS
- [ ] No hay Load Balancers huérfanos
- [ ] AWS Billing muestra $0 cargos pendientes

## Notas
> 💡 **Buena práctica**: Al final de cada sesión de práctica, destruir los recursos que cuestan dinero (EKS, EC2). Volver a crearlos toma ~20 min con Terraform.

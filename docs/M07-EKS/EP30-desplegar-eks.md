# EP 30: Desplegar el Proyecto en EKS

**Tipo:** PRACTICA
## Objetivo
Aplicar los manifiestos YAML al cluster EKS y verificar que la app corre en la nube.

## Prerequisitos
- kubectl conectado a EKS (EP29)
- Imagen en Docker Hub (EP11)

## Pasos Detallados

### 1. Crear Secrets
```bash
# Crear secret para la base de datos
kubectl create secret generic db-credentials -n curso-gitops   --from-literal=username=curso_app   --from-literal=password='C4rs0_S3cur3_P@ss!'

# Crear secret para JWT
kubectl create secret generic app-secrets -n curso-gitops   --from-literal=jwt-secret='gk8s_pr0d_s3cr3t_ch4ng3_m3!'
```

### 2. Aplicar Manifiestos
```bash
kubectl apply -f infrastructure/kubernetes/app/namespace.yaml
kubectl apply -f infrastructure/kubernetes/app/deployment.yaml
kubectl apply -f infrastructure/kubernetes/app/service.yaml
```

### 3. Verificar
```bash
kubectl get pods -n curso-gitops
kubectl get svc -n curso-gitops
# Copiar el EXTERNAL-IP del LoadBalancer

# Esperar ~2 min a que el LB se provisionne
kubectl get svc -n curso-gitops -w
```

### 4. Acceder a la App
```bash
# Obtener la URL
EXTERNAL_IP=$(kubectl get svc curso-gitops-svc -n curso-gitops -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
echo "App URL: http://$EXTERNAL_IP"
```

## Archivos Involucrados
- `infrastructure/kubernetes/app/namespace.yaml`
- `infrastructure/kubernetes/app/deployment.yaml`
- `infrastructure/kubernetes/app/service.yaml`

## Verificación
- [ ] Pods en estado "Running"
- [ ] Service tiene EXTERNAL-IP asignado
- [ ] La app es accesible desde el navegador

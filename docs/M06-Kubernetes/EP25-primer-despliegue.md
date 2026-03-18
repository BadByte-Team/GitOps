# EP 25: Primer Despliegue con Manifiestos YAML

**Tipo:** PRACTICA
## Objetivo
Escribir un Deployment + Service y desplegar una app en el cluster Minikube.

## Prerequisitos
- Minikube corriendo (EP24)

## Pasos Detallados

### 1. Crear el Namespace
```bash
kubectl apply -f infrastructure/kubernetes/app/namespace.yaml
```

### 2. Aplicar el Deployment
```bash
kubectl apply -f infrastructure/kubernetes/app/deployment.yaml
```

### 3. Aplicar el Service
```bash
kubectl apply -f infrastructure/kubernetes/app/service.yaml
```

### 4. Verificar
```bash
# Ver pods corriendo
kubectl get pods -n curso-gitops

# Ver servicios
kubectl get svc -n curso-gitops

# Ver deployment
kubectl get deployments -n curso-gitops

# En Minikube, acceder al servicio
minikube service curso-gitops-svc -n curso-gitops
```

## Archivos Involucrados
- `infrastructure/kubernetes/app/namespace.yaml`
- `infrastructure/kubernetes/app/deployment.yaml`
- `infrastructure/kubernetes/app/service.yaml`

## Verificación
- [ ] Los pods están en estado "Running"
- [ ] El servicio tiene un endpoint asignado
- [ ] Puedes acceder a la app por el navegador

# EP 38: Instalar ArgoCD en EKS

**Tipo:** INSTALACION
## Objetivo
Instalar ArgoCD en el cluster EKS y verificar que los pods estén corriendo.

## Prerequisitos
- kubectl conectado a EKS (EP29)

## Pasos Detallados

### 1. Crear Namespace
```bash
kubectl create namespace argocd
```

### 2. Instalar ArgoCD
```bash
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### 3. Verificar Pods
```bash
kubectl get pods -n argocd
# Esperar a que todos estén "Running" (1-2 min)
# Debe haber ~7 pods
```

### Pods Esperados
| Pod | Función |
|---|---|
| `argocd-server` | API y UI web |
| `argocd-repo-server` | Clona repos y genera manifiestos |
| `argocd-application-controller` | Reconcilia estado |
| `argocd-dex-server` | Autenticación SSO |
| `argocd-redis` | Cache |

## Archivos Involucrados
- `infrastructure/kubernetes/argocd/install.yaml` (referencia)

## Verificación
- [ ] Todos los pods de ArgoCD están en "Running"
- [ ] `kubectl get svc -n argocd` muestra el argocd-server

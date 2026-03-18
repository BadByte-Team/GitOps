# EP 40: Crear Application en ArgoCD — Conectar al Repo

**Tipo:** CONFIGURACION / PRACTICA
## Objetivo
Configurar una Application en ArgoCD que conecte al repositorio de manifiestos y haga sync automático.

## Opción A: Por UI

1. Abrir ArgoCD → Applications → New App
2. Configurar:
   - **Name**: `curso-gitops`
   - **Project**: `default`
   - **Sync Policy**: Automatic (+ Prune + Self Heal)
   - **Repository URL**: `https://github.com/TU_USUARIO/curso-gitops-manifests.git`
   - **Path**: `infrastructure/kubernetes/app`
   - **Cluster URL**: `https://kubernetes.default.svc`
   - **Namespace**: `curso-gitops`
3. Create

## Opción B: Por YAML (recomendado)

```bash
kubectl apply -f infrastructure/kubernetes/argocd/application.yaml
```

### Verificar Sync
```bash
kubectl get application curso-gitops -n argocd
# STATUS debe ser "Synced" y HEALTH "Healthy"
```

## Archivos Involucrados
- `infrastructure/kubernetes/argocd/application.yaml`

## Verificación
- [ ] La app aparece en el dashboard de ArgoCD
- [ ] STATUS: Synced, HEALTH: Healthy
- [ ] Los pods de curso-gitops están corriendo

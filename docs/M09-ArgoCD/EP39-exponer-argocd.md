# EP 39: Exponer ArgoCD con LoadBalancer y Obtener Contraseña

**Tipo:** CONFIGURACION
## Objetivo
Exponer el UI de ArgoCD con un LoadBalancer y obtener la contraseña de admin.

## Pasos Detallados

### 1. Cambiar Service a LoadBalancer
```bash
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'
```

### 2. Obtener URL
```bash
kubectl get svc argocd-server -n argocd
# Copiar el EXTERNAL-IP (tarda ~2 min en asignarse)

# O en una línea:
ARGOCD_URL=$(kubectl get svc argocd-server -n argocd -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
echo "ArgoCD: https://$ARGOCD_URL"
```

### 3. Obtener Contraseña
```bash
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
echo  # Nueva línea
```

### 4. Login
- **URL**: `https://EXTERNAL-IP` (aceptar certificado self-signed)
- **Usuario**: `admin`
- **Contraseña**: (la del paso anterior)

### 5. Cambiar Contraseña (recomendado)
- En el UI: User Info → Update Password

## Verificación
- [ ] ArgoCD es accesible por el navegador
- [ ] Puedes loguearte como admin

# EP 41: Despliegue Automático — Push Dispara Sync en EKS

**Tipo:** PRACTICA
## Objetivo
Ver el flujo GitOps completo: un push al repo de manifiestos dispara un sync automático en EKS.

## Flujo Completo
```
1. Cambiar imagen tag en deployment.yaml
2. Push a GitHub (repo de manifiestos)
3. ArgoCD detecta el cambio (~3 min o manual sync)
4. ArgoCD aplica el cambio al cluster EKS
5. Rolling update — pods nuevos reemplazan a los viejos
```

## Pasos Detallados

### 1. Modificar el Tag de la Imagen
```bash
# En el repo de manifiestos
cd infrastructure/kubernetes/app

# Cambiar el tag de la imagen
sed -i 's|image: TU_USUARIO/curso-gitops:.*|image: TU_USUARIO/curso-gitops:v2|' deployment.yaml

git add deployment.yaml
git commit -m "ci: update image to v2"
git push origin main
```

### 2. Observar en ArgoCD
- Abrir el dashboard de ArgoCD
- En ~3 minutos, verás que la app pasa a "OutOfSync"
- ArgoCD automáticamente aplica el cambio (auto-sync)
- Los pods hacen rolling update

### 3. Verificar
```bash
kubectl get pods -n curso-gitops -w
# Ver cómo se crean pods nuevos y se eliminan los viejos

kubectl describe deployment curso-gitops -n curso-gitops | grep Image
# Debe mostrar la imagen v2
```

## Verificación
- [ ] El cambio en Git dispara un sync en ArgoCD
- [ ] Los pods se actualizan sin downtime (rolling update)
- [ ] La nueva versión es accesible

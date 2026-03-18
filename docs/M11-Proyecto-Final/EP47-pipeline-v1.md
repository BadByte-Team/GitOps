# EP 47: Pipeline V1 — El Proyecto en Producción

**Tipo:** PRACTICA
## Objetivo
Ejecutar el pipeline CI completo por primera vez y verificar que la app está en producción.

## Checklist Pre-Requisitos
- [ ] EC2 Jenkins corriendo con Docker, Trivy, SonarQube (EP31-EP44)
- [ ] Credenciales configuradas: Docker Hub, GitHub, SonarQube (EP34, EP44)
- [ ] Cluster EKS corriendo (EP28)
- [ ] ArgoCD instalado y configurado (EP38-EP40)

## Pasos Detallados

### 1. Ejecutar el Pipeline
1. Jenkins → `curso-gitops-ci` → Build Now
2. Observar los stages en tiempo real

### 2. Verificar cada Stage
- **Checkout**: ✅ Código clonado
- **SonarQube**: ✅ Análisis completado (ver reporte en `http://IP:9000`)
- **Quality Gate**: ✅ Passed
- **Docker Build**: ✅ Imagen construida
- **Trivy Scan**: ✅ Reporte generado
- **Docker Push**: ✅ Imagen en Docker Hub
- **Update Manifest**: ✅ deployment.yaml actualizado

### 3. Verificar ArgoCD
- Abrir ArgoCD → ver que la app se sincronizó
- STATUS: Synced, HEALTH: Healthy

### 4. Acceder a la App
```bash
EXTERNAL_IP=$(kubectl get svc curso-gitops-svc -n curso-gitops -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
echo "App: http://$EXTERNAL_IP"
```

## Verificación
- [ ] Pipeline completado en verde
- [ ] Imagen en Docker Hub con tag correcto
- [ ] ArgoCD muestra la app Synced y Healthy
- [ ] La app es accesible por el LoadBalancer

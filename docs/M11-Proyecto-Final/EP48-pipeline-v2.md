# EP 48: Actualizar a V2 — CI/CD en Acción

**Tipo:** PRACTICA
## Objetivo
Hacer un cambio en el código, pushear a GitHub, y ver cómo todo el flujo CI/CD se ejecuta automáticamente.

## Pasos Detallados

### 1. Hacer un Cambio en el Código
```bash
# Cambiar algo visible en la app, por ejemplo el footer
cd curso-gitops/frontend/dashboard.html
# Cambiar "Curso GitOps © 2026" por "Curso GitOps v2 © 2026"
```

### 2. Commit y Push
```bash
git add .
git commit -m "feat: update to v2"
git push origin main
```

### 3. Observar el Flujo Automático
1. **Jenkins** detecta el push (webhook o polling) → ejecuta pipeline
2. **Pipeline** construye nueva imagen, escanea, sube a Docker Hub
3. **Pipeline** actualiza el tag en el repo de manifiestos
4. **ArgoCD** detecta el cambio → sync automático
5. **EKS** hace rolling update → pods v2 reemplazan a v1

### 4. Verificar Zero-Downtime
```bash
# En otra terminal, hacer requests continuos
while true; do curl -s http://$EXTERNAL_IP | grep -o "v[0-9]"; sleep 1; done
# Verás la transición de v1 a v2 sin errores
```

## Verificación
- [ ] El push disparó el pipeline automáticamente
- [ ] La nueva imagen está en Docker Hub
- [ ] ArgoCD sincronizó el cambio
- [ ] La app muestra los cambios de v2
- [ ] No hubo downtime durante la actualización

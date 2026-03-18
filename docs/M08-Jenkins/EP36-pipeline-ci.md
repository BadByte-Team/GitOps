# EP 36: Pipeline CI Completo — Build, Tag y Push a Docker Hub

**Tipo:** PRACTICA
## Objetivo
Implementar el pipeline CI completo: build imagen Docker, tag con el número de build, y push a Docker Hub.

## El Pipeline Completo

El archivo `infrastructure/jenkins/Jenkinsfile` ya contiene el pipeline completo:

```
Checkout → SonarQube → Quality Gate → Docker Build → Trivy Scan → Docker Push → Update Manifest → Cleanup
```

### Stages Clave
1. **Checkout**: Clona el repositorio
2. **Docker Build**: Construye la imagen con tag `BUILD_NUMBER-GIT_HASH`
3. **Docker Push**: Sube la imagen a Docker Hub
4. **Update Manifest**: Actualiza el tag de la imagen en el deployment.yaml del repo de manifiestos

### Crear el Pipeline en Jenkins
1. New Item → Pipeline → nombre: `curso-gitops-ci`
2. Pipeline → "Pipeline script from SCM"
3. SCM: Git → URL del repo de la app
4. Credentials: seleccionar `github-token`
5. Branch: `*/main`
6. Script Path: `infrastructure/jenkins/Jenkinsfile`
7. Build Now

## Archivos Involucrados
- `infrastructure/jenkins/Jenkinsfile`

## Verificación
- [ ] El pipeline completa todos los stages en verde
- [ ] La imagen aparece en Docker Hub con el tag correcto
- [ ] El deployment.yaml fue actualizado con el nuevo tag

## Notas
> En este punto, si ArgoCD está configurado (EP40), el push al repo de manifiestos disparará automáticamente el despliegue en EKS.

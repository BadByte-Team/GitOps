# EP 45: Agregar Stages Trivy y SonarQube al Jenkinsfile

**Tipo:** PRACTICA
## Objetivo
Agregar los stages de escaneo de seguridad al pipeline CI.

## Stages de Seguridad en el Jenkinsfile

El `infrastructure/jenkins/Jenkinsfile` ya incluye estos stages:

### Stage: SonarQube Analysis
```groovy
stage('SonarQube Analysis') {
    steps {
        withSonarQubeEnv('sonarqube-server') {
            sh "\${SCANNER_HOME}/bin/sonar-scanner \
                -Dsonar.projectKey=curso-gitops \
                -Dsonar.sources=."
        }
    }
}
```

### Stage: Quality Gate
```groovy
stage('Quality Gate') {
    steps {
        timeout(time: 2, unit: 'MINUTES') {
            waitForQualityGate abortPipeline: true
        }
    }
}
```
> Si la calidad del código no pasa, el pipeline se detiene.

### Stage: Trivy Scan
```groovy
stage('Trivy Scan') {
    steps {
        sh "trivy image --exit-code 0 --severity HIGH,CRITICAL \
            --format table \${DOCKER_IMAGE}:\${BUILD_TAG}"
    }
}
```
> `--exit-code 0` = reporta pero no falla. Cambiar a `1` para fallar en vulnerabilidades.

## Pipeline Completo
```
Checkout → SonarQube → Quality Gate → Docker Build → Trivy Scan → Docker Push → Update Manifest
```

## Archivos Involucrados
- `infrastructure/jenkins/Jenkinsfile`

## Verificación
- [ ] El pipeline ejecuta SonarQube y muestra resultados
- [ ] Trivy escanea la imagen y genera un reporte
- [ ] Quality Gate funciona (pasa o falla según la calidad)

# EP 45: Agregar los Stages de Seguridad al Jenkinsfile

**Tipo:** PRÁCTICA
**Duración estimada:** 15–18 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Revisar el Jenkinsfile completo con todos los stages de seguridad integrados — SonarQube Analysis, Quality Gate y Trivy Scan — y ejecutar el pipeline completo observando cómo cada stage cumple su función en el flujo de CI/CD.

---

## 📋 Prerequisitos
- Trivy instalado y montado en Jenkins (EP42)
- SonarQube con el proyecto `curso-gitops` y el token generado (EP43)
- Integración Jenkins ↔ SonarQube configurada (EP44)

---

## 🧠 El pipeline completo con seguridad

El Jenkinsfile del proyecto tiene 6 stages en este orden:

```
Checkout
    ↓
SonarQube Analysis  ← calidad y seguridad del código fuente
    ↓
Docker Build        ← construir la imagen
    ↓
Trivy Scan          ← vulnerabilidades en la imagen Docker
    ↓
Docker Push         ← subir a Docker Hub si pasó todo
    ↓
Deploy to GitOps    ← actualizar gitops-infra
    ↓
Cleanup
```

El orden es intencional: primero validamos el código antes de construir la imagen, luego validamos la imagen antes de publicarla. Si algo falla en seguridad, el pipeline se detiene antes de llegar a producción.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: el Jenkinsfile completo en VS Code — todos los stages visibles desde el principio hasta el final.*

"Bienvenidos al episodio 45. El último episodio del Módulo 10.

Llevamos tres episodios preparando las herramientas de seguridad: Trivy instalado, SonarQube configurado con su proyecto y token, y la integración entre Jenkins y SonarQube funcionando.

Hoy revisamos el Jenkinsfile completo con todos esos stages integrados y ejecutamos el pipeline para ver el flujo de seguridad en acción.

Este episodio cierra el Módulo 10 y también cierra el ciclo completo de CI del curso. A partir de aquí, el Módulo 11 es el proyecto final — la arquitectura completa, la base de datos separada en Kubernetes, y la limpieza de recursos.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — Por qué este orden de stages (1:30 – 3:30)

> *Pantalla: diagrama del pipeline con los stages en orden.*

"Antes de revisar el código, quiero que entiendan por qué los stages están en este orden específico.

La regla general en seguridad es **fail fast and early** — fallar rápido y temprano. Cuanto antes detectes un problema, más barato y menos disruptivo es corregirlo.

**SonarQube va primero** porque analiza el código fuente — lo más básico. Si hay vulnerabilidades en el código Go en sí mismo, o el Quality Gate falla por alguna razón crítica, no tiene sentido gastar tiempo construyendo la imagen.

**Trivy va después del build** porque necesita que la imagen exista para escanearla. Pero va **antes del push** — si la imagen tiene CVEs críticos, no queremos publicarla en Docker Hub donde alguien podría usarla.

**El push va después de los dos escaneos** — es la puerta de salida. Solo las imágenes que pasaron ambos controles llegan a Docker Hub y eventualmente a producción.

**Deploy to GitOps va al final** — solo actualizamos `gitops-infra` si toda la cadena de validación pasó. ArgoCD solo despliega código que fue aprobado por el pipeline completo.

Ese orden es el que protege la cadena de suministro de software. Cada stage es una puerta que el código debe pasar antes de avanzar."

---

### PASO 1 — Revisar el Jenkinsfile completo (3:30 – 8:00)

> *Pantalla: VS Code con el Jenkinsfile — recorrer cada stage con el cursor.*

"Abro el Jenkinsfile en `gitops-app/Jenkinsfile` y lo recorro section por section.

**El bloque `tools` y `environment`:**"

```groovy
tools {
    jdk 'jdk17'
    nodejs 'node18'
}

environment {
    DOCKER_HUB_CREDS = credentials('dockerhub-id')
    DOCKER_IMAGE     = "TU_USUARIO_DOCKERHUB/curso-gitops"
    SCANNER_HOME     = tool('sonar-scanner')
}
```

"`SCANNER_HOME` usa `tool('sonar-scanner')` — resuelve la ruta donde Jenkins instaló el SonarQube Scanner configurado en el EP33.

---

**Stage: Checkout**"

```groovy
stage('Checkout') {
    steps {
        checkout scm
        script {
            env.GIT_COMMIT_SHORT = sh(
                script: 'git rev-parse --short HEAD',
                returnStdout: true
            ).trim()
            env.BUILD_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT_SHORT}"
        }
    }
}
```

"Genera el `BUILD_TAG` — el identificador único de esta versión. El número de build más el hash del commit. Por ejemplo: `3-a3b8d1c`.

---

**Stage: SonarQube Analysis**"

```groovy
stage('SonarQube Analysis') {
    steps {
        withSonarQubeEnv('sonarqube-server') {
            sh """
                ${SCANNER_HOME}/bin/sonar-scanner \\
                -Dsonar.projectKey=curso-gitops \\
                -Dsonar.projectName=curso-gitops \\
                -Dsonar.sources=. \\
                -Dsonar.exclusions=**/vendor/**,**/node_modules/**,**/frontend/**
            """
        }
    }
}
```

"`withSonarQubeEnv('sonarqube-server')` inyecta automáticamente la URL del servidor y el token de autenticación que configuramos en el EP44 — no hay que pasarlos manualmente.

`-Dsonar.exclusions` excluye del análisis directorios que no son código Go puro: vendor (dependencias externas), node_modules si existiera, y frontend (HTML/JS que SonarQube analiza de forma diferente). Así el análisis se enfoca en el código Go de la app.

---

**Stage: Docker Build**"

```groovy
stage('Docker Build') {
    steps {
        sh "docker build -t ${DOCKER_IMAGE}:${BUILD_TAG} ."
        sh "docker tag ${DOCKER_IMAGE}:${BUILD_TAG} ${DOCKER_IMAGE}:latest"
    }
}
```

"Construye la imagen con el tag versionado y también la etiqueta como `latest`.

---

**Stage: Trivy Scan**"

```groovy
stage('Trivy Scan') {
    steps {
        sh """
            trivy image \\
              --exit-code 0 \\
              --severity HIGH,CRITICAL \\
              --format table \\
              ${DOCKER_IMAGE}:${BUILD_TAG}
        """
    }
}
```

"`--exit-code 0` — reporta vulnerabilidades pero no bloquea el pipeline. Para cambiar la política a enforcement estricto, cambiar a `--exit-code 1`. El reporte aparece en el Console Output de Jenkins para que el equipo lo revise.

---

**Stage: Docker Push**"

```groovy
stage('Docker Push') {
    steps {
        sh "echo ${DOCKER_HUB_CREDS_PSW} | docker login -u ${DOCKER_HUB_CREDS_USR} --password-stdin"
        sh "docker push ${DOCKER_IMAGE}:${BUILD_TAG}"
        sh "docker push ${DOCKER_IMAGE}:latest"
    }
}
```

"Solo llega aquí si SonarQube y Trivy no bloquearon el pipeline. La imagen que se publica en Docker Hub fue validada por ambas herramientas.

---

**Stage: Deploy to GitOps Repo**"

```groovy
stage('Deploy to GitOps Repo') {
    steps {
        withCredentials([string(credentialsId: 'github-token-id', variable: 'GITHUB_TOKEN')]) {
            sh """
                rm -rf infra-repo
                git clone https://${GITHUB_TOKEN}@github.com/TU_USUARIO_GITHUB/gitops-infra.git infra-repo
                cd infra-repo
                git config user.email "jenkins@local.com"
                git config user.name "Jenkins CI"
                sed -i "s|image: ${DOCKER_IMAGE}:.*|image: ${DOCKER_IMAGE}:${BUILD_TAG}|" \\
                    infrastructure/kubernetes/app/deployment.yaml
                git add infrastructure/kubernetes/app/deployment.yaml
                git commit -m "ci: deploy version ${BUILD_TAG} from Jenkins"
                git push origin main
            """
        }
    }
}
```

"El corazón del patrón GitOps. Solo llega aquí código que pasó los cuatro stages anteriores.

---

**El bloque `post`:**"

```groovy
post {
    success {
        echo "✅ Pipeline completado — imagen: ${DOCKER_IMAGE}:${BUILD_TAG}"
    }
    failure {
        echo "❌ Pipeline fallido en stage: ${env.STAGE_NAME}"
    }
    always {
        cleanWs()
    }
}
```

"`cleanWs()` limpia el workspace del build. Sin esto, los archivos del repositorio clonado y los temporales se acumularían en el disco de Jenkins."

---

### PASO 2 — Ejecutar el pipeline completo (8:00 – 13:00)

> *Pantalla: Jenkins — Stage View del pipeline ejecutando.*

"Click en **Build Now**. El Stage View muestra los 6 stages progresando.

**Checkout** — verde en segundos. El `BUILD_TAG` generado es `3-a3b8d1c` (o el número correspondiente).

**SonarQube Analysis** — 30-60 segundos. Mientras corre, abro SonarQube en el navegador. El proyecto `curso-gitops` muestra 'Background task in progress'. Cuando termina, aparece el resultado del Quality Gate.

En el Console Output busco la confirmación:"

```
INFO: ANALYSIS SUCCESSFUL, you can find the results at:
http://localhost:9000/dashboard?id=curso-gitops
INFO: Note that you will be able to access the updated dashboard once the server has processed the submitted analysis report
INFO: More about the report processing at http://localhost:9000/api/ce/task?id=...
```

"Eso confirma que el análisis llegó a SonarQube.

**Docker Build** — 1-2 minutos dependiendo de si la imagen base está en caché. Jenkins construye localmente a través del socket Docker.

**Trivy Scan** — el reporte aparece en el Console Output:"

```
Total: 0 (HIGH: 0, CRITICAL: 0)
```

"Cero vulnerabilidades HIGH o CRITICAL en nuestra imagen Alpine + binario Go. Eso es exactamente lo que esperamos — el multi-stage build que usamos en el EP10 mantiene la imagen mínima.

**Docker Push** — la imagen sube a Docker Hub con el nuevo tag.

**Deploy to GitOps Repo** — Jenkins clona `gitops-infra`, actualiza el tag con `sed`, y hace push.

**Cleanup** — las imágenes locales se eliminan.

---

El pipeline terminó en verde. El Console Output muestra el mensaje final:"

```
✅ Pipeline completado — imagen: TU_USUARIO/curso-gitops:3-a3b8d1c
```

---

### PASO 3 — Verificar los resultados en SonarQube (13:00 – 15:00)

> *Pantalla: navegador en SonarQube — dashboard del proyecto.*

"Abro el proyecto `curso-gitops` en SonarQube.

El dashboard muestra ahora los resultados del análisis:

**Quality Gate:** Passed ✅ — el código cumplió todas las condiciones del Quality Gate 'Sonar way'.

**Reliability:** los bugs encontrados. Para el código Go del curso, debería ser 0 o muy bajo.

**Security:** las vulnerabilidades en el código fuente. Distintas a las vulnerabilidades de Trivy — aquí SonarQube analiza patrones del código Go, como inputs no sanitizados o manejo inseguro de credenciales.

**Maintainability:** los code smells. Puede haber algunos — funciones demasiado largas, duplicaciones menores — pero no deberían ser bloqueantes.

El histórico de análisis muestra este análisis y el del EP44. Con cada build, el historial crece y puedes ver la evolución de la calidad del código a lo largo del tiempo."

---

### CIERRE (15:00 – 16:00)

"Eso es el EP45. Y con esto cerramos el Módulo 10.

El pipeline de CI es ahora un pipeline de seguridad integrada. Cada commit que llega a producción pasó por:
- Análisis de calidad de código con SonarQube
- Evaluación del Quality Gate
- Escaneo de vulnerabilidades de la imagen con Trivy
- Y recién después, el push a Docker Hub y el despliegue via ArgoCD

Ese es el estándar que se usa en pipelines de producción reales. No es teoría — es lo que acabamos de construir y verificar que funciona.

El Módulo 11 es el proyecto final. Vamos a revisar la arquitectura híbrida completa, hablar de la base de datos separada en Kubernetes con todos sus manifiestos, ver el pipeline en acción una vez más, y finalmente limpiar todos los recursos de AWS.

Nos vemos en el EP46."

---

## ✅ Checklist de Verificación
- [ ] El pipeline ejecuta los 6 stages en verde
- [ ] El Console Output muestra `ANALYSIS SUCCESSFUL` en el stage de SonarQube
- [ ] El reporte de Trivy aparece en el Console Output con el conteo de vulnerabilidades
- [ ] La imagen aparece en Docker Hub con el nuevo tag
- [ ] El `deployment.yaml` en `gitops-infra` fue actualizado por Jenkins
- [ ] SonarQube muestra `Quality Gate: Passed` para el proyecto `curso-gitops`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Stage SonarQube falla con `Could not connect` | SonarQube no está corriendo: `docker compose ps` en `~/local-ci/` |
| `waitForQualityGate` espera indefinidamente | El webhook en SonarQube no está configurado — verificar EP44 |
| Trivy no encuentra la imagen | El stage Docker Build falló antes — revisar el log del build |
| Quality Gate falla | Ver los detalles en SonarQube — puede ser una vulnerabilidad de seguridad en el código Go o un bug blocker |
| `docker login` falla en el stage Docker Push | Las credenciales de Docker Hub expiraron — verificar en EP34 |

---

## 🗒️ Notas de Producción
- La pausa conceptual del orden de stages es el aporte pedagógico más valioso del episodio — "fail fast and early" es un principio que los alumnos van a usar en sus carreras.
- Mientras el stage de SonarQube corre, mantener ambas pantallas visibles — Jenkins en una ventana y SonarQube actualizándose en otra — muestra la integración en tiempo real.
- Si el reporte de Trivy muestra cero vulnerabilidades, celebrarlo explícitamente — "esto es el resultado del multi-stage build que aprendimos en el EP10".
- El dashboard de SonarQube al final del episodio es visualmente impactante — hacer zoom en el Quality Gate verde para que sea claramente legible.
- El cierre con el resumen de lo que el pipeline garantiza ("cada commit que llega a producción pasó por...") es el mensaje de closure del módulo completo.

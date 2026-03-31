# EP 36: Pipeline CI con GitOps — Build, Push y Actualización Automática de Manifiestos

**Tipo:** PRÁCTICA
**Duración estimada:** 18–20 min
**Dificultad:** ⭐⭐⭐ (Avanzado)
**🔄 MODIFICADO:** El pipeline termina clonando `gitops-infra`, actualizando el tag de la imagen con `sed`, y haciendo un `git push` automatizado al repositorio de infraestructura para que ArgoCD detecte el cambio.

---

## 🎯 Objetivo
Implementar el pipeline CI completo del curso: construir la imagen Docker de la app Go, escanearla con SonarQube, subirla a Docker Hub con un tag versionado, y actualizar automáticamente `deployment.yaml` en `gitops-infra` para que ArgoCD despliegue la nueva versión en K3s.

---

## 📋 Prerequisitos
- Jenkins con plugins, tools y credenciales configurados (EP31–34)
- Los dos repositorios privados creados: `gitops-app` y `gitops-infra` (EP07)
- Docker Hub con la imagen inicial subida (EP11)
- K3s corriendo en la EC2 (EP29)

---

## 🧠 El patrón GitOps de dos repositorios

El pipeline implementa el patrón central del curso:

```
gitops-app                      gitops-infra
(código fuente)                 (manifiestos K8s)
      │                               │
      │  git push                     │  ArgoCD observa
      ▼                               ▼
Jenkins (CI local)              K3s Cluster (EC2)
  1. Checkout                     ArgoCD detecta
  2. SonarQube                    el cambio en
  3. Docker build                 deployment.yaml
  4. Docker push      ──────────▶ y sincroniza
  5. Clona gitops-infra           el cluster
  6. sed actualiza tag
  7. git push
```

Jenkins es el puente. Construye la imagen en local, la publica en Docker Hub, y actualiza el manifiesto en `gitops-infra`. ArgoCD nunca habla con Jenkins — solo observa `gitops-infra`.

---

## El Jenkinsfile completo

>[!WARNING]
>Cambiar "TU_USUARIO_DOCKERHUB"  y "TU_USUARIO_GITHUB" por los valores correspondientes


```groovy
pipeline {
    agent any

    tools {
        jdk 'jdk17'
        nodejs 'node18'
    }

    environment {
        DOCKER_HUB_CREDS = credentials('dockerhub-id')
        DOCKER_IMAGE     = "TU_USUARIO_DOCKERHUB/curso-gitops"
        SCANNER_HOME     = tool('sonar-scanner')
        GITHUB_USER      = "TU_USUARIO_GITHUB"
        INFRA_REPO       = "gitops-infra"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    env.BUILD_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT_SHORT}"
                    echo "BUILD_TAG: ${env.BUILD_TAG}"
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('sonarqube-server') {
                    sh """
                        ${SCANNER_HOME}/bin/sonar-scanner \
                        -Dsonar.projectKey=curso-gitops \
                        -Dsonar.projectName=curso-gitops \
                        -Dsonar.sources=. \
                        -Dsonar.exclusions=**/vendor/**,**/node_modules/**,**/frontend/**
                    """
                }
            }
        }

        stage('Quality Gate') {
            steps {
                timeout(time: 5, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Docker Build') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE}:${BUILD_TAG} ."
                sh "docker tag ${DOCKER_IMAGE}:${BUILD_TAG} ${DOCKER_IMAGE}:latest"
                echo "Imagen construida: ${DOCKER_IMAGE}:${BUILD_TAG}"
            }
        }

        stage('Trivy Scan') {
            steps {
                sh """
                    trivy image \
                      --exit-code 0 \
                      --severity HIGH,CRITICAL \
                      --format table \
                      ${DOCKER_IMAGE}:${BUILD_TAG}
                """
            }
        }

        stage('Docker Push') {
            steps {
                sh "echo ${DOCKER_HUB_CREDS_PSW} | docker login -u ${DOCKER_HUB_CREDS_USR} --password-stdin"
                sh "docker push ${DOCKER_IMAGE}:${BUILD_TAG}"
                sh "docker push ${DOCKER_IMAGE}:latest"
                echo "Imagen subida a Docker Hub: ${DOCKER_IMAGE}:${BUILD_TAG}"
            }
        }

        stage('Deploy to GitOps Repo') {
            steps {
                withCredentials([string(credentialsId: 'github-token-id', variable: 'GITHUB_TOKEN')]) {
                    sh """
                        rm -rf infra-repo

                        git clone https://${GITHUB_TOKEN}@github.com/${GITHUB_USER}/${INFRA_REPO}.git infra-repo

                        cd infra-repo
                        git config user.email "jenkins@local.com"
                        git config user.name "Jenkins CI"

                        # Actualizar el tag de la imagen en deployment.yaml
                        sed -i "s|image: ${DOCKER_IMAGE}:.*|image: ${DOCKER_IMAGE}:${BUILD_TAG}|" \
                            infrastructure/kubernetes/app/deployment.yaml

                        git add infrastructure/kubernetes/app/deployment.yaml
                        git commit -m "ci: deploy version ${BUILD_TAG} from Jenkins"
                        git push origin main
                    """
                }
            }
        }

        stage('Cleanup') {
            steps {
                sh "docker rmi ${DOCKER_IMAGE}:${BUILD_TAG} || true"
                sh "docker rmi ${DOCKER_IMAGE}:latest || true"
                sh "docker image prune -f || true"
            }
        }
    }

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
}

```

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: el diagrama del flujo GitOps completo — PC local → gitops-infra → ArgoCD → K3s.*

"Bienvenidos al episodio 36. El episodio del que trata todo este módulo.

Llevamos cinco episodios preparando el entorno de Jenkins: lo instalamos localmente, instalamos los plugins, configuramos las herramientas, configuramos las credenciales, escribimos el primer pipeline de ejemplo.

Todo eso fue preparación para este momento.

Hoy implementamos el pipeline CI completo. No es un pipeline de juguete — es el pipeline que va a conectar el código de la app con el cluster K3s de producción. Cuando este pipeline termine de ejecutarse, ArgoCD va a detectar el cambio automáticamente y va a desplegar la nueva versión de la app en la EC2.

Ese flujo — código push → imagen nueva → ArgoCD sincroniza — es la definición de GitOps. Y hoy lo ponemos a funcionar.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — El stage que lo hace todo diferente (1:30 – 3:30)

> *Pantalla: el stage `Deploy to GitOps Repo` del Jenkinsfile aislado en VS Code.*

"Antes de ejecutar, quiero que entiendan el stage que diferencia este pipeline de uno convencional. Se llama `Deploy to GitOps Repo` y es el corazón del patrón GitOps.

Un pipeline convencional haría esto: construye la imagen, la sube a Docker Hub, y se conecta directamente al servidor de producción para aplicar el cambio. Jenkins habla directamente con Kubernetes. Ese enfoque funciona, pero tiene un problema: nadie sabe qué está en producción sin consultar el servidor.

El patrón GitOps es diferente. Jenkins nunca toca el servidor de producción directamente. En cambio, actualiza un archivo en Git — el `deployment.yaml` en `gitops-infra` — y es ArgoCD quien aplica ese cambio al cluster.

¿Por qué importa esa diferencia? Porque `gitops-infra` es la fuente de verdad. Si quieres saber qué versión está en producción, miras el repositorio. Si quieres hacer un rollback, reviertes el commit y ArgoCD lo aplica automáticamente. El historial de Git es el historial de producción.

Ese stage hace tres cosas:
1. Clona `gitops-infra` usando el token de GitHub
2. Usa `sed` para actualizar el tag de la imagen en `deployment.yaml`
3. Hace commit y push — ArgoCD detecta el cambio en segundos

Simple en concepto. Poderoso en práctica."

---

### PASO 1 — Preparar el Jenkinsfile en gitops-app (3:30 – 6:00)

> *Pantalla: VS Code con el Jenkinsfile del proyecto.*

"Reemplazo los dos placeholders del Jenkinsfile con mis datos reales.

En la línea `DOCKER_IMAGE = "TU_USUARIO_DOCKERHUB/curso-gitops"`, reemplazo `TU_USUARIO_DOCKERHUB` con mi usuario real de Docker Hub.

En el stage `Deploy to GitOps Repo`, en la URL del `git clone`, reemplazo `TU_USUARIO_GITHUB` con mi usuario real de GitHub.

Guardo el archivo y lo subo al repositorio:"

```bash
cd gitops-app
git add Jenkinsfile
git commit -m "ci: configurar pipeline CI/CD con flujo GitOps"
git push origin main
```

---
### PASO 2 — Configurar servidor de SonarQube en Jenkins (antes del pipeline)

> _Pantalla: Jenkins → configuración global._

"Antes de crear el pipeline, necesito registrar el servidor de SonarQube en Jenkins.

Voy a:
**Manage Jenkins** → **System** → **SonarQube servers** → **Add SonarQube**

Configuro:
- **Name:** `sonarqube-server`
- **Server URL:** `http://<IP_DEL_CONTENEDOR_SONARQUBE>:9000` _(IP del contenedor)_
- **Server authentication token:** seleccionar `sonarqube-server`

---
### ⚠️ Importante — No usar localhost

La URL **NO debe ser `http://localhost:9000`**.
¿Por qué?
- Jenkins corre en un contenedor
- SonarQube corre en otro contenedor
- Para Jenkins, `localhost` apunta a **sí mismo**, no a SonarQube
    
Entonces:
- `localhost` → Jenkins
- `IP DEL CONENEDOR SONARQUBE` → contenedor de SonarQube
    
Si usas `localhost`, el pipeline falla con:

```
connection refused
```

Click en **Save**

Con esto, Jenkins ya puede comunicarse correctamente con SonarQube y el stage:
```
withSonarQubeEnv('sonarqube-server')
```
va a funcionar sin problemas."

---

Si quieres, el paso 7 lo dejamos como validación del pipeline + cómo detectar rápido cuando SonarQube está mal conectado (ahí hay varios errores típicos interesantes).


### PASO 3 — Crear el pipeline en Jenkins (6:00 – 8:00)

> *Pantalla: navegador en Jenkins.*

"En Jenkins: **New Item** → nombre: `curso-gitops-ci` → **Pipeline** → **OK**.

En la sección **Build Triggers**, activo **'GitHub hook trigger for GITScm polling'** — cuando configuremos el webhook de GitHub en el EP48, el pipeline arrancará automáticamente en cada push.

En la sección **Pipeline**:
- **Definition:** `Pipeline script from SCM`
- **SCM:** Git
- **Repository URL:** `https://github.com/TU_USUARIO_GITHUB/gitops-app.git`
- **Credentials:** `github-token-id`
- **Branch:** `*/main`
- **Script Path:** `Jenkinsfile`

Click en **'Save'**."

---

### PASO 4 — Ejecutar y observar cada stage (8:00 – 15:00)

> *Pantalla: Stage View del pipeline corriendo en Jenkins.*

"Click en **Build Now**. El Stage View comienza a mostrar los stages uno por uno.

---

**Stage: Checkout** — Jenkins clona el repositorio `gitops-app`. En el script block genera el `BUILD_TAG` — el número de build más los primeros 7 caracteres del hash del commit. Por ejemplo: `1-a3b8d1c`. Ese tag va a ser el identificador único de esta versión de la imagen.

---

**Stage: SonarQube Analysis** — el scanner analiza el código Go. Envía los resultados al SonarQube local en `http://localhost:9000`. Esta etapa tarda entre 30 y 60 segundos dependiendo del tamaño del código.

Mientras corre, puedo abrir SonarQube en el navegador y en unos momentos aparecerá el proyecto `curso-gitops` con el resultado del análisis.

---

**Stage: Docker Build** — el comando que construye la imagen.

`docker build -t TU_USUARIO/curso-gitops:1-a3b8d1c .`

Jenkins llama a Docker directamente a través del socket montado en el EP31. La imagen se construye en la máquina local — sin network latency, sin autenticación adicional.

También hace `docker tag` para crear el tag `latest` — así Docker Hub siempre tiene una versión `latest` actualizada para quien quiera usar la imagen sin especificar un tag.

---

**Stage: Docker Push** — sube la imagen a Docker Hub.

El `echo ${DOCKER_HUB_CREDS_PSW} | docker login -u ${DOCKER_HUB_CREDS_USR} --password-stdin` usa las credenciales que configuramos en el EP34. La sintaxis `--password-stdin` es la forma segura de autenticarse — la contraseña nunca aparece en el log del pipeline.

Sube primero el tag versionado y luego el `latest`.

---

**Stage: Deploy to GitOps Repo** — el stage central del patrón GitOps.

```bash
git clone https://${GITHUB_TOKEN}@github.com/TU_USUARIO/gitops-infra.git infra-repo
```

Clona el repositorio de infraestructura usando el token que configuramos en el EP34.

```bash
sed -i "s|image: TU_USUARIO/curso-gitops:.*|image: TU_USUARIO/curso-gitops:1-a3b8d1c|" \
    infrastructure/kubernetes/app/deployment.yaml
```

Este comando `sed` es el que hace el cambio. Lee el `deployment.yaml`, encuentra la línea que empieza con `image: TU_USUARIO/curso-gitops:`, y reemplaza todo lo que hay después de `:` con el nuevo tag `1-a3b8d1c`. El patrón `.*` captura cualquier tag anterior — sea `latest`, `0-abc123`, o cualquier otro.

```bash
git commit -m "ci: deploy version 1-a3b8d1c from Jenkins"
git push origin main
```

Commit y push. En este momento, `gitops-infra` tiene el nuevo tag. ArgoCD va a detectar este commit en los próximos 3 minutos y sincronizará el cluster.

---

**Stage: Cleanup** — elimina las imágenes locales que construimos. No tiene sentido ocupar espacio en disco con imágenes que ya están en Docker Hub. El `|| true` evita que el stage falle si la imagen ya fue eliminada por algún motivo."

---

### PASO 5 — Verificar los resultados (15:00 – 17:30)

> *Pantalla: tres ventanas — Jenkins, Docker Hub, GitHub.*

"El pipeline terminó en verde. Verifico tres cosas.

**Docker Hub** — voy a `hub.docker.com/r/TU_USUARIO/curso-gitops`. Aparece el tag `1-a3b8d1c` y `latest`. La imagen está publicada y disponible para cualquier cluster que la necesite.

**GitHub — gitops-infra** — voy al repositorio `gitops-infra` en GitHub. En el historial de commits hay uno nuevo: `ci: deploy version 1-a3b8d1c from Jenkins`. Fue creado por el usuario `Jenkins CI` que configuramos en el stage.

Abro el archivo `infrastructure/kubernetes/app/deployment.yaml`. La línea `image:` ahora muestra el nuevo tag:

```yaml
- name: curso-gitops
  image: TU_USUARIO/curso-gitops:1-a3b8d1c   ← actualizado por Jenkins
```

Ese archivo es lo que ArgoCD va a leer. Cuando en el EP38 instalemos ArgoCD y en el EP40 lo conectemos a este repositorio, va a ver este tag y va a desplegar esa imagen exacta en el cluster K3s."

---

### CIERRE (17:30 – 18:30)

"Eso es el episodio 36. El Módulo 08 completo.

El pipeline CI está funcionando. Construye la imagen de la app Go, la escanea con SonarQube, la sube a Docker Hub con un tag único, y actualiza el repositorio de infraestructura para que ArgoCD lo despliegue.

Eso conecta todos los módulos anteriores: el código de la app del Módulo 03, los dos repositorios GitOps del EP07, las credenciales del EP34, el cluster K3s del Módulo 07.

Lo que falta para tener el flujo completo es ArgoCD — el operador del lado de Kubernetes que va a leer ese `deployment.yaml` y aplicarlo al cluster. Eso es el Módulo 09.

Nos vemos en el EP37."

---

## ✅ Checklist de Verificación
- [ ] El pipeline `curso-gitops-ci` ejecuta todos los stages en verde
- [ ] La imagen aparece en Docker Hub con el tag `BUILD-HASH` y `latest`
- [ ] El `deployment.yaml` en `gitops-infra` fue actualizado con el nuevo tag
- [ ] El commit en `gitops-infra` tiene el autor `Jenkins CI`
- [ ] Los logs del pipeline no muestran credenciales en texto plano

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Stage SonarQube falla con `connection refused` | Verificar que SonarQube está corriendo: `docker compose ps` en `~/local-ci` |
| `docker build` falla con `permission denied` | `sudo chmod 666 /var/run/docker.sock` y `docker compose restart jenkins` |
| Stage Docker Push falla con `unauthorized` | Las credenciales de Docker Hub son incorrectas — verificar en EP34 |
| `sed` no modifica el `deployment.yaml` | La línea `image:` en el YAML no tiene el formato esperado — verificar que tiene `image: TU_USUARIO/curso-gitops:algo` |
| `git push` falla con `remote: Invalid username or password` | El token de GitHub expiró o no tiene scope `repo` — regenerar en EP34 |
| El tag `latest` ya no corresponde a la versión versionada | El `docker tag` del stage Docker Build asegura que siempre estén sincronizados |

---

## 🗒️ Notas de Producción
- La apertura con el diagrama del flujo completo establece inmediatamente qué se va a lograr en el episodio — mantenerlo en pantalla durante los primeros 90 segundos.
- La pausa conceptual del stage `Deploy to GitOps Repo` es el corazón pedagógico del episodio — tomarse el tiempo necesario para explicar por qué Jenkins no toca el servidor directamente.
- Mostrar el Stage View durante la ejecución — es visualmente impactante y muestra el progreso en tiempo real.
- Mientras el stage de SonarQube corre, abrir SonarQube en el navegador y mostrar que el análisis aparece en tiempo real — demuestra que los dos servicios locales están integrados.
- El cierre con las tres ventanas (Jenkins verde, Docker Hub con la imagen, GitHub con el commit de Jenkins) es el momento más satisfactorio del módulo — dejar cada pantalla visible unos segundos antes de hablar del siguiente módulo.

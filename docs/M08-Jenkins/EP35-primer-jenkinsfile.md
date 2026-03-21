# EP 35: Creación del Primer Jenkinsfile

**Tipo:** PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Entender la estructura completa de un pipeline declarativo de Jenkins escribiendo un Jenkinsfile de ejemplo desde cero, para que cuando veamos el pipeline real del EP36 cada sección sea reconocible.

---

## 📋 Prerequisitos
- Jenkins corriendo con plugins, tools y credenciales configurados (EP31–34)
- El repositorio `gitops-app` en GitHub

---

## 🧠 Anatomía de un Jenkinsfile declarativo

```groovy
pipeline {
    agent any                    // dónde ejecutar el pipeline

    tools { }                    // herramientas: JDK, Node, etc.

    environment { }              // variables de entorno del pipeline

    stages {                     // las etapas — el corazón del pipeline
        stage('nombre') {
            steps {              // los pasos dentro de la etapa
                sh '...'
            }
        }
    }

    post {                       // qué hacer al terminar — siempre, en éxito, en fallo
        always  { }
        success { }
        failure { }
    }
}
```

Cada bloque tiene su propósito:
- **`agent`** — en qué nodo o contenedor corre el pipeline. `any` significa "usa el nodo disponible".
- **`tools`** — declara herramientas gestionadas por Jenkins (configuradas en el EP33).
- **`environment`** — variables accesibles en todos los stages. Las credenciales se inyectan aquí.
- **`stages`** — la lista de etapas en orden. Jenkins las ejecuta secuencialmente.
- **`post`** — acciones post-ejecución. `always` corre sin importar el resultado. `success` solo si todo fue bien. `failure` solo si algo falló.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: el Jenkinsfile completo del proyecto en VS Code — todas las líneas visibles.*

"Bienvenidos al episodio 35.

Este es el Jenkinsfile completo del proyecto — el pipeline que Jenkins va a ejecutar en el EP36. Tiene un bloque `tools`, un bloque `environment`, seis stages, y un bloque `post`.

Si lo miraran ahora mismo sin contexto previo, sería bastante denso. Hay referencias a credenciales, comandos de shell, variables de entorno con sintaxis específica de Groovy.

Hoy no vamos a escribir ese pipeline todavía. Vamos a escribir uno más simple desde cero, entendiendo cada sección. La idea es exactamente la misma que con Terraform en el EP19: construir los bloques fundamentales primero, para que el archivo real del EP36 no tenga ninguna sorpresa.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — ¿Por qué un Jenkinsfile y no configuración por UI? (1:00 – 2:30)

> *Pantalla: la interfaz de configuración de un job de Jenkins en el navegador.*

"Jenkins tiene dos formas de configurar pipelines. La primera es a través de la interfaz gráfica — clickear opciones, llenar formularios. La segunda es un Jenkinsfile en el repositorio.

¿Por qué preferimos el Jenkinsfile?

La razón es la misma que con Terraform: **el código es la fuente de verdad**. Si el Jenkinsfile vive en el repositorio `gitops-app`, cualquier cambio al pipeline es un commit de Git. Si algo se rompe, puedes ver exactamente qué cambió, cuándo, y quién lo hizo. Si necesitas recrear Jenkins desde cero, el pipeline se recupera automáticamente — Jenkins simplemente lee el Jenkinsfile del repositorio.

Con la configuración por UI, si pierdes el servidor de Jenkins, perdiste también toda la configuración del pipeline. Sin historial, sin posibilidad de recuperación automática.

Es el mismo principio de GitOps aplicado al CI: el estado deseado del pipeline vive en Git."

---

### PASO 1 — Crear el Jenkinsfile de ejemplo (2:30 – 5:30)

> *Pantalla: VS Code, creando un archivo `Jenkinsfile-ejemplo` en un directorio temporal.*

"Voy a crear un Jenkinsfile de ejemplo que muestre todos los bloques fundamentales. Lo escribimos en VS Code con la extensión de Jenkins instalada — tiene resaltado de sintaxis para Groovy."

```groovy
pipeline {
    agent any

    tools {
        jdk 'jdk17'
        nodejs 'node18'
    }

    environment {
        NOMBRE_APP = "mi-primera-app"
        VERSION    = "1.0.0"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
                sh 'echo "Código clonado en: $(pwd)"'
                sh 'ls -la'
            }
        }

        stage('Build') {
            steps {
                sh 'echo "Construyendo ${NOMBRE_APP} versión ${VERSION}"'
                sh 'java -version'
            }
        }

        stage('Test') {
            steps {
                sh 'echo "Ejecutando tests..."'
                sh 'echo "Tests completados"'
            }
        }

        stage('Deploy') {
            steps {
                sh 'echo "Desplegando ${NOMBRE_APP}:${VERSION}"'
            }
        }

    }

    post {
        success {
            echo "✅ Pipeline completado — ${NOMBRE_APP}:${VERSION} desplegado"
        }
        failure {
            echo "❌ Pipeline fallido en el stage: ${env.STAGE_NAME}"
        }
        always {
            cleanWs()
        }
    }
}
```

"Voy a explicar cada sección.

**`agent any`** — el pipeline puede correr en cualquier nodo disponible de Jenkins. En nuestro caso solo hay un nodo — el Jenkins local — así que siempre corre ahí.

**`tools`** — declara las herramientas que necesitamos. `jdk 'jdk17'` le dice a Jenkins que busque la configuración llamada `jdk17` que configuramos en el EP33 y ponga ese JDK disponible en el PATH para todos los stages.

**`environment`** — variables de entorno accesibles en todo el pipeline. Se usan con la sintaxis `${VARIABLE}` dentro de los steps. En el pipeline real del proyecto vamos a ver aquí la inyección de credenciales con `credentials('id')`.

**`stages`** — la lista de etapas en orden. Jenkins las ejecuta de forma secuencial. Si un stage falla, los siguientes no se ejecutan — a menos que uses configuraciones especiales para forzar la ejecución.

**`post`** — lo que pasa después. `success` solo corre si todos los stages terminaron bien. `failure` solo si alguno falló. `always` corre siempre — sin importar el resultado. El `cleanWs()` en `always` limpia el workspace del pipeline — buena práctica para no acumular archivos entre builds."

---

### PASO 2 — Crear el pipeline en Jenkins (5:30 – 8:30)

> *Pantalla: navegador en Jenkins.*

"Ahora conecto este Jenkinsfile a Jenkins para verlo ejecutar.

En el dashboard de Jenkins: **New Item** → nombre: `mi-primer-pipeline` → selecciono **Pipeline** → **OK**.

En la configuración del pipeline, bajo hasta la sección **Pipeline**:
- **Definition:** `Pipeline script from SCM`
- **SCM:** Git
- **Repository URL:** la URL SSH del repositorio `gitops-app`
- **Credentials:** selecciono `github-token-id` que configuramos en el EP34
- **Branch Specifier:** `*/main`
- **Script Path:** `Jenkinsfile` — el nombre del archivo en la raíz del repo

Click en **'Save'**.

Antes de ejecutar, subo el Jenkinsfile que escribimos al repositorio `gitops-app`:"

```bash
cd gitops-app
cp Jenkinsfile-ejemplo Jenkinsfile
git add Jenkinsfile
git commit -m "feat: agregar Jenkinsfile base"
git push origin main
```

"De vuelta en Jenkins: **Build Now**."

---

### PASO 3 — Ver el pipeline ejecutarse (8:30 – 11:00)

> *Pantalla: navegador mostrando el pipeline en ejecución — Stage View visible.*

"Jenkins muestra el Stage View — un diagrama visual de cada stage con su estado y tiempo de ejecución. Los stages en azul son los que están corriendo. Los en verde son los que terminaron bien.

Click en el número del build — `#1` — para ver los detalles. Luego click en **Console Output** para ver la salida completa del pipeline línea por línea.

Veo las salidas de los `echo` y los `sh`. El mensaje de `post.success` al final confirma que todo terminó bien.

Hay algo que quiero que noten en el Console Output: la variable `${NOMBRE_APP}` fue reemplazada por `mi-primera-app`. Así es como Jenkins interpola las variables de `environment` dentro de los steps."

---

### PASO 4 — Conectar con el Jenkinsfile real del proyecto (11:00 – 13:00)

> *Pantalla: VS Code con el Jenkinsfile real del proyecto.*

"Ahora abro el Jenkinsfile real del proyecto en `gitops-app/Jenkinsfile` — el que vamos a usar en el EP36.

Reconocen todo:

`tools { jdk 'jdk17' nodejs 'node18' }` — los mismos nombres que configuramos en el EP33.

`environment { DOCKER_HUB_CREDS = credentials('dockerhub-id') ... }` — las credenciales del EP34 inyectadas aquí.

Los stages: Checkout, SonarQube Analysis, Docker Build, Docker Push, Deploy to GitOps Repo, Cleanup. Cada uno hace exactamente lo que su nombre dice.

`post { success { ... } failure { ... } always { cleanWs() } }` — el mismo patrón de post que acabamos de usar en el ejemplo.

La única parte nueva que vamos a explicar en detalle en el EP36 es el stage `Deploy to GitOps Repo` — que es donde Jenkins clona `gitops-infra` y actualiza el tag de la imagen. Ese es el corazón del patrón GitOps."

---

### CIERRE (13:00 – 14:00)

"Eso es el EP35.

Escribieron su primer Jenkinsfile, lo conectaron a Jenkins, lo vieron ejecutar, y leyeron el pipeline real del proyecto reconociendo cada bloque.

En el siguiente episodio reemplazamos ese Jenkinsfile de ejemplo con el real — el que construye la imagen de la app Go, la sube a Docker Hub, y actualiza el repositorio `gitops-infra` para que ArgoCD la despliegue.

Ese es el episodio del que trata todo este módulo.

Nos vemos en el EP36."

---

## ✅ Checklist de Verificación
- [ ] El Jenkinsfile de ejemplo existe en el repositorio `gitops-app`
- [ ] El pipeline `mi-primer-pipeline` ejecuta en verde
- [ ] Puedes leer el Console Output y entender cada línea
- [ ] Puedes abrir el Jenkinsfile real del proyecto y reconocer todos los bloques

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Tool type "jdk" does not have an install of "jdk17"` | El nombre en Tools no coincide — verificar EP33 |
| `ERROR: Couldn't find any revision to build` | La rama en el pipeline no coincide — verificar `*/main` |
| `Failed to connect to repository` | El token de GitHub no tiene permisos — verificar scope `repo` en EP34 |
| El Console Output muestra `${NOMBRE_APP}` sin resolver | La variable no está en el bloque `environment` del Jenkinsfile |

---

## 🗒️ Notas de Producción
- Abrir el Jenkinsfile completo del proyecto al inicio — es el destino visual que motiva escribir el ejemplo simple primero.
- Escribir el Jenkinsfile de ejemplo en VS Code con la extensión de Jenkins para que tenga resaltado de sintaxis Groovy.
- El Stage View en la interfaz de Jenkins es visualmente muy atractivo — hacer zoom para que se vea bien en el video.
- Al conectar con el Jenkinsfile real al final, señalar con el cursor cada bloque mientras lo mencionas verbalmente — la conexión explícita entre el ejemplo y el real es el momento pedagógico más valioso.

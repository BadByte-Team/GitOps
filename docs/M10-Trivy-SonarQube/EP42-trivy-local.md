# EP 42: Instalación de Trivy Local

**Tipo:** INSTALACIÓN / PRÁCTICA
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Trivy se instala y ejecuta en el entorno CI local del alumno — la misma máquina donde corre Jenkins — no en la nube.

---

## 🎯 Objetivo
Instalar Trivy en tu PC local, montarlo dentro del contenedor de Jenkins actualizando el `docker-compose.yml`, y verificar que el stage de escaneo del pipeline puede ejecutarse correctamente.

---

## 📋 Prerequisitos
- Jenkins local corriendo con Docker Compose (EP31)
- Docker instalado localmente (EP08)

---

## 🧠 ¿Qué es Trivy y por qué lo necesitamos?

Trivy es un escáner de seguridad de código abierto creado por Aqua Security. Analiza imágenes Docker en busca de vulnerabilidades conocidas en las bibliotecas y paquetes que contiene.

El flujo en el pipeline:

```
docker build → imagen construida localmente
      ↓
trivy image → escanea la imagen antes del push
      ↓
reporte de vulnerabilidades HIGH y CRITICAL
      ↓
docker push → imagen sube a Docker Hub
```

Trivy actúa como un control de calidad de seguridad: si la imagen tiene vulnerabilidades críticas, las detectamos antes de que lleguen a producción.

En la arquitectura original, Trivy corría en la EC2 de Jenkins. En la arquitectura gratuita, corre en tu máquina local — donde Jenkins ya construye las imágenes. Sin latencia de red, sin configuración adicional.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: una imagen Docker siendo pusheada a Docker Hub. Al lado, el stage del Jenkinsfile con el comando `trivy image`.*

"Bienvenidos al episodio 42. Arrancamos el Módulo 10 — Seguridad.

En los módulos anteriores construimos el pipeline completo. Jenkins construye, SonarQube analiza la calidad del código, ArgoCD despliega. El flujo funciona.

Pero hay un eslabón que falta: ¿qué pasa con las vulnerabilidades de seguridad en las imágenes que estamos subiendo a producción? La imagen puede tener código impecable y aun así traer paquetes del sistema operativo con CVEs conocidos.

Trivy resuelve eso. Es un escáner que analiza la imagen Docker antes del push y reporta qué vulnerabilidades tiene, con su nivel de severidad y si existe una versión parcheada.

El episodio es corto: instalar Trivy en local, montarlo en Jenkins, y ver el reporte. Vamos."

---

### 🔍 PAUSA CONCEPTUAL — CVEs y niveles de severidad (1:00 – 2:30)

> *Pantalla: un reporte de Trivy mostrando las columnas Library, Vulnerability, Severity, Status.*

"Antes de instalar, el vocabulario que vamos a ver en los reportes.

Un **CVE** — Common Vulnerabilities and Exposures — es un identificador único para una vulnerabilidad de seguridad conocida públicamente. Cuando aparece `CVE-2024-12345` en el reporte de Trivy, ese número te permite buscar exactamente de qué se trata en la base de datos nacional de vulnerabilidades de NIST.

Los niveles de severidad que maneja Trivy son cinco: UNKNOWN, LOW, MEDIUM, HIGH y CRITICAL. En el pipeline del curso escaneamos con `--severity HIGH,CRITICAL` — nos interesan las vulnerabilidades que podrían ser explotables con impacto real.

La columna **Status** es la que más importa para decidir qué hacer:
- `fixed` — existe una versión del paquete que corrige la vulnerabilidad. Actualizar el paquete en el Dockerfile la resuelve.
- `affected` — la vulnerabilidad existe pero todavía no hay parche disponible. Evaluar el riesgo.
- `will_not_fix` — los mantenedores decidieron no parchearla. Evaluar si es crítico para el contexto.

Para nuestra imagen — Alpine + binario Go estático — el número de vulnerabilidades debería ser muy bajo. Alpine es minimalista por diseño y el binario Go compila sus dependencias estáticamente."

---

### PASO 1 — Instalar Trivy en la PC local (2:30 – 5:00)

> *Pantalla: terminal en la PC local.*

"Trivy tiene su propio repositorio de paquetes, igual que Terraform o HashiCorp. Lo instalo desde la fuente oficial:

**Ubuntu / Debian:**"

```bash
sudo apt-get install -y wget apt-transport-https gnupg lsb-release

wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key \
  | gpg --dearmor \
  | sudo tee /usr/share/keyrings/trivy.gpg > /dev/null

echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] \
  https://aquasecurity.github.io/trivy-repo/deb \
  $(lsb_release -sc) main" \
  | sudo tee /etc/apt/sources.list.d/trivy.list

sudo apt-get update -y && sudo apt-get install -y trivy
```

"**Arch Linux:**"

```bash
sudo pacman -S trivy
# o con yay: yay -S trivy
```

"**macOS:**"

```bash
brew install aquasecurity/trivy/trivy
```

"Verifico la instalación:"

```bash
trivy --version
# Version: 0.50.x
```

---

### PASO 2 — Primera prueba de escaneo manual (5:00 – 7:30)

> *Pantalla: terminal ejecutando trivy con la salida visible.*

"Antes de integrarlo con Jenkins, hago un escaneo manual para ver cómo funciona y familiarizarme con el output.

Escaneo la imagen de nginx como ejemplo:"

```bash
trivy image --severity HIGH,CRITICAL nginx:latest
```

"La primera ejecución tarda entre 30 y 60 segundos — Trivy descarga su base de datos de vulnerabilidades (~100 MB). Las siguientes son casi instantáneas porque la base de datos queda en caché.

El output muestra una tabla por cada imagen base encontrada en la imagen de Docker. Para nginx en Alpine, algo como:"

```
nginx:latest (alpine 3.19.1)
Total: 3 (HIGH: 2, CRITICAL: 1)

┌─────────────────┬────────────────┬──────────┬──────────┬──────────────────┐
│     Library     │ Vulnerability  │ Severity │  Status  │    Fixed In      │
├─────────────────┼────────────────┼──────────┼──────────┼──────────────────┤
│ libssl3         │ CVE-2024-XXXXX │ CRITICAL │  fixed   │ 3.1.5-r0         │
│ libcrypto3      │ CVE-2024-XXXXX │ HIGH     │  fixed   │ 3.1.5-r0         │
│ openssl         │ CVE-2024-XXXXX │ HIGH     │  fixed   │ 3.1.5-r0         │
└─────────────────┴────────────────┴──────────┴──────────┴──────────────────┘
```

"Los tres están en status `fixed` — existe una versión del paquete que los corrige. Si fuera tu imagen, actualizarías la versión base de Alpine en el Dockerfile.

Ahora escaneo nuestra propia imagen para saber con qué partimos:"

```bash
trivy image --severity HIGH,CRITICAL TU_USUARIO/curso-gitops:latest
```

"La imagen Alpine + Go estático generalmente tiene muy pocas vulnerabilidades — a veces ninguna en HIGH/CRITICAL. Ese es uno de los beneficios del multi-stage build que usamos en el Dockerfile: la imagen final solo contiene el binario compilado, sin compiladores ni herramientas de desarrollo que son fuentes comunes de CVEs."

---

### PASO 3 — Montar Trivy en el contenedor de Jenkins (7:30 – 9:30)

> *Pantalla: VS Code con `~/local-ci/docker-compose.yml`.*

"El contenedor de Jenkins necesita acceso al binario de Trivy para poder ejecutarlo en el pipeline. La forma más limpia es montarlo como volumen — igual que hicimos con el socket de Docker en el EP31.

Edito `~/local-ci/docker-compose.yml` y agrego una línea al servicio Jenkins:"

```yaml
jenkins:
  image: jenkins/jenkins:lts
  container_name: jenkins
  user: root
  ports:
    - "8080:8080"
    - "50000:50000"
  volumes:
    - jenkins_data:/var/jenkins_home
    - /var/run/docker.sock:/var/run/docker.sock
    - /usr/bin/docker:/usr/bin/docker
    - /usr/bin/trivy:/usr/bin/trivy        # ← línea nueva
  restart: unless-stopped
```

"El volumen `/usr/bin/trivy:/usr/bin/trivy` monta el binario de Trivy del sistema operativo dentro del contenedor de Jenkins en la misma ruta. Cuando el Jenkinsfile ejecute el comando `trivy`, el contenedor encontrará el binario exactamente donde lo busca.

Aplico el cambio reiniciando solo el servicio de Jenkins — SonarQube y su base de datos no necesitan reiniciarse:"

```bash
cd ~/local-ci
docker compose restart jenkins
```

"Espero que Jenkins reinicie — 20-30 segundos — y verifico que el binario está accesible desde dentro del contenedor:"

```bash
docker exec jenkins trivy --version
# Version: 0.50.x
```

"Perfecto. Jenkins puede ejecutar Trivy."

---

### PASO 4 — Agregar el stage de Trivy al Jenkinsfile (9:30 – 11:00)

> *Pantalla: VS Code con el Jenkinsfile del proyecto en `gitops-app`.*

"El Jenkinsfile del proyecto ya tiene el stage preparado. Solo hay que verificar que está en el lugar correcto — después del Docker Build y antes del Docker Push:

```groovy
stage('Docker Build') {
    steps {
        sh "docker build -t ${DOCKER_IMAGE}:${BUILD_TAG} ."
        sh "docker tag ${DOCKER_IMAGE}:${BUILD_TAG} ${DOCKER_IMAGE}:latest"
    }
}

stage('Trivy Scan') {          // ← este stage
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

stage('Docker Push') {
    steps {
        // ...
    }
}
```

"Dos flags que merecen una explicación.

`--exit-code 0` — el stage termina exitosamente sin importar si Trivy encontró vulnerabilidades. El reporte se genera pero el pipeline continúa. Es la configuración correcta para empezar: primero observas qué vulnerabilidades tienes, y cuando tienes un criterio de aceptación claro, cambias a `--exit-code 1`.

`--exit-code 1` — el pipeline falla si Trivy encuentra vulnerabilidades con la severidad especificada. Úsalo cuando quieras enforcement: ninguna imagen con CVEs críticos llega a Docker Hub.

Para el curso usamos `--exit-code 0` — así el pipeline no se rompe si nuestra imagen tiene algún CVE menor."

---

### CIERRE (11:00 – 12:00)

"Eso es el EP42.

Trivy instalado localmente, montado en Jenkins, y el stage de escaneo configurado en el Jenkinsfile. La próxima vez que corra el pipeline, después de construir la imagen aparecerá el reporte de vulnerabilidades en el Console Output antes de hacer el push a Docker Hub.

En el siguiente episodio confirmamos que SonarQube está correctamente configurado — ya corre desde el EP31, pero en el EP43 creamos el proyecto, generamos el token de análisis, y verificamos que el flujo está listo para el EP44 donde lo integramos con Jenkins.

Nos vemos en el EP43."

---

## ✅ Checklist de Verificación
- [ ] `trivy --version` funciona en la PC local
- [ ] `docker exec jenkins trivy --version` funciona dentro del contenedor
- [ ] El volumen `/usr/bin/trivy:/usr/bin/trivy` está en el `docker-compose.yml`
- [ ] `trivy image nginx:latest` genera un reporte sin errores
- [ ] El stage `Trivy Scan` existe en el Jenkinsfile entre Docker Build y Docker Push

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `trivy: command not found` dentro de Jenkins | El volumen no está en el compose o Jenkins no fue reiniciado — `docker compose restart jenkins` |
| `chmod +x` en el binario | `sudo chmod +x /usr/bin/trivy` en el sistema local |
| El escaneo tarda mucho la primera vez | Trivy descarga la base de datos (~100 MB) — es normal solo la primera vez |
| `permission denied` accediendo al cache de Trivy | Trivy guarda la DB en `/root/.cache/trivy/` dentro del contenedor — verificar permisos |
| Trivy no encuentra la imagen | La imagen debe existir localmente: `docker images | grep curso-gitops` |

---

## 🗒️ Notas de Producción
- La pausa conceptual de CVEs y severidades es breve pero necesaria — el alumno va a ver esos términos en el reporte y necesita entenderlos.
- El escaneo de `nginx:latest` es un buen ejemplo porque siempre tiene alguna vulnerabilidad — muestra el reporte real con datos concretos.
- El escaneo de `curso-gitops:latest` puede mostrar cero vulnerabilidades — si es así, explicar que es una señal positiva del multi-stage build, no un error.
- La diferencia `--exit-code 0` vs `--exit-code 1` merece énfasis verbal — es la decisión que los equipos toman según su política de seguridad.
- Mostrar `docker compose restart jenkins` (solo Jenkins, no todo el compose) — es más rápido y no interrumpe SonarQube.

# EP 31: Instalación de Jenkins Local con Docker Compose

**Tipo:** INSTALACIÓN
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Jenkins ya no se instala en AWS. Se levanta en la PC del alumno usando Docker Compose local, eliminando completamente ese costo de EC2.

---

## 🎯 Objetivo
Levantar Jenkins y SonarQube en tu máquina local usando Docker Compose, obtener la contraseña inicial de Jenkins, completar la configuración inicial, y verificar que Jenkins puede construir imágenes Docker.

---

## 📋 Prerequisitos
- Docker instalado y corriendo en local (EP08)
- `docker compose version` responde sin errores

---

## 🧠 ¿Por qué Jenkins local y no en AWS?

En la arquitectura original, Jenkins vivía en una EC2 t2.medium. En la arquitectura gratuita del curso:

| Componente | Dónde vive | Costo |
|---|---|---|
| Jenkins + SonarQube | Tu PC local (Docker Compose) | $0 |
| K3s + ArgoCD + App | EC2 t3.micro (AWS Free Tier) | $0 |
| **Total** | | **$0** |

Jenkins en tu PC tiene además una ventaja técnica real: acceso directo al Docker daemon local para construir imágenes. Sin configuración adicional, sin permisos especiales de red, sin tokens de autenticación entre Jenkins y Docker.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: la tabla de costos del módulo anterior — calculadora de AWS mostrando ~$33/mes por una EC2 t2.medium para Jenkins.*

"Bienvenidos al episodio 31. Arrancamos el Módulo 08 — Jenkins.

En el módulo anterior terminamos con el cluster K3s corriendo en la EC2 Free Tier. Esa instancia tiene 1 GB de RAM y la tenemos prácticamente al límite con K3s y ArgoCD. Si además instaláramos Jenkins ahí, colapsaría.

Y la alternativa obvia — una segunda instancia EC2 para Jenkins — cuesta alrededor de $33 al mes solo por el servidor, antes de considerar almacenamiento o tráfico.

Hay una solución mucho mejor y que además simplifica el flujo técnico: Jenkins en tu propia máquina local usando Docker Compose. Gratis, sin latencia de red para construir imágenes, y con acceso directo al Docker daemon.

En este episodio levantamos Jenkins y SonarQube juntos en un solo `docker compose up`. Hacia el final del episodio van a tener los dos corriendo y accesibles desde el navegador, listos para configurar en los siguientes episodios.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — La arquitectura CI/CD local (1:30 – 3:00)

> *Pantalla: diagrama de la arquitectura híbrida.*

"Antes de escribir cualquier archivo, quiero que vean exactamente dónde encaja Jenkins en el flujo completo del curso.

```
Tu PC Local
├── Jenkins  → construye imágenes, corre tests, actualiza gitops-infra
├── SonarQube → analiza la calidad del código
└── Docker   → construye y almacena las imágenes localmente antes del push

       ↓ git push a gitops-infra
       
EC2 t3.micro (AWS)
└── K3s → ArgoCD detecta el cambio → despliega la nueva versión
```

Jenkins vive en tu máquina. ArgoCD vive en AWS. El puente entre los dos es GitHub — Jenkins hace push a `gitops-infra`, ArgoCD lo detecta y sincroniza el cluster.

Esto es lo que se llama una arquitectura CI/CD híbrida: la integración continua corre localmente, el despliegue continuo corre en la nube. Es un patrón real que usan equipos con presupuesto limitado o que quieren control completo sobre el entorno de build."

---

### PASO 1 — Preparar el entorno local (3:00 – 4:30)

> *Pantalla: terminal en la PC local.*

"Primero un paso crítico que mucha gente olvida: dar permisos al socket de Docker.

El socket `/var/run/docker.sock` es el canal de comunicación con el Docker daemon. Cuando Jenkins construya una imagen, lo hará a través de este socket. Por defecto, solo root puede usarlo. Con este comando lo abrimos para que el contenedor de Jenkins pueda acceder:"

```bash
sudo chmod 666 /var/run/docker.sock
```

"Este permiso se resetea cuando Docker se reinicia. Si en algún momento Jenkins no puede construir imágenes, este es el primer lugar donde verificar.

Ahora creo el directorio donde va a vivir la configuración local de CI. Este directorio no se sube a GitHub — es solo tuyo, local:"

```bash
mkdir ~/local-ci
cd ~/local-ci
```

---

### PASO 2 — Crear el Docker Compose (4:30 – 8:00)

> *Pantalla: VS Code creando el archivo `docker-compose.yml`.*

"Ahora el corazón del episodio: el archivo que define el entorno completo de CI.

Voy a crear `~/local-ci/docker-compose.yml`. Son tres servicios: Jenkins, SonarQube y la base de datos de SonarQube."

```yaml
version: '3'
services:

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
    restart: unless-stopped

  sonarqube:
    image: sonarqube:lts-community
    container_name: sonarqube
    ports:
      - "9000:9000"
    environment:
      - SONAR_JDBC_URL=jdbc:postgresql://sonar-db:5432/sonar
      - SONAR_JDBC_USERNAME=sonar
      - SONAR_JDBC_PASSWORD=sonar_p4ssw0rd
    volumes:
      - sonarqube_data:/opt/sonarqube/data
      - sonarqube_extensions:/opt/sonarqube/extensions
      - sonarqube_logs:/opt/sonarqube/logs
    depends_on:
      - sonar-db
    restart: unless-stopped

  sonar-db:
    image: postgres:15-alpine
    container_name: sonar-db
    environment:
      - POSTGRES_USER=sonar
      - POSTGRES_PASSWORD=sonar_p4ssw0rd
      - POSTGRES_DB=sonar
    volumes:
      - sonar_postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  jenkins_data:
  sonarqube_data:
  sonarqube_extensions:
  sonarqube_logs:
  sonar_postgres_data:
```

"Voy a explicar las partes que más importan entender.

**El servicio Jenkins** tiene tres volúmenes. El primero — `jenkins_data:/var/jenkins_home` — es el más importante: guarda toda la configuración de Jenkins en un volumen nombrado. Plugins instalados, pipelines creados, credenciales configuradas — todo persiste aunque hagas `docker compose down` y vuelvas a levantarlo. Solo se pierde con `docker compose down -v`.

El segundo volumen — `/var/run/docker.sock:/var/run/docker.sock` — monta el socket de Docker del sistema operativo dentro del contenedor de Jenkins. Así cuando Jenkins ejecute `docker build`, el comando habla directamente con el Docker daemon de tu máquina.

El tercero — `/usr/bin/docker:/usr/bin/docker` — monta el binario de Docker dentro del contenedor para que esté disponible como comando.

**SonarQube** usa PostgreSQL como base de datos en lugar del H2 embebido — es más estable para uso prolongado. Los tres volúmenes de SonarQube guardan los análisis de código, los plugins y los logs respectivamente.

**Los volúmenes nombrados** al final son fundamentales. Sin ellos, cada `docker compose down` borraría toda la configuración acumulada. Con ellos, el estado persiste indefinidamente."

---

### PASO 3 — Ajuste para SonarQube (8:00 – 9:00)

> *Pantalla: terminal.*

"SonarQube usa Elasticsearch internamente y requiere que el parámetro `vm.max_map_count` del kernel sea al menos 262144. Sin este ajuste, SonarQube no arranca correctamente.

Lo configuro antes de levantar los servicios:"

```bash
sudo sysctl -w vm.max_map_count=262144

# Para hacerlo permanente entre reinicios:
echo 'vm.max_map_count=262144' | sudo tee -a /etc/sysctl.conf
```

---

### PASO 4 — Levantar los servicios (9:00 – 10:30)

> *Pantalla: terminal.*

```bash
cd ~/local-ci
docker compose up -d
```

"Docker descarga las imágenes la primera vez — puede tardar entre 2 y 5 minutos dependiendo de la conexión. Las siguientes veces arranca en segundos porque las imágenes están en caché.

Verifico que los tres contenedores están corriendo:"

```bash
docker compose ps
```

```
NAME         IMAGE                    COMMAND   STATUS    PORTS
jenkins      jenkins/jenkins:lts      ...       Up        0.0.0.0:8080->8080/tcp
sonarqube    sonarqube:lts-community  ...       Up        0.0.0.0:9000->9000/tcp
sonar-db     postgres:15-alpine       ...       Up        5432/tcp
```

"Los tres en estado `Up`. SonarQube tarda un par de minutos adicionales en inicializarse internamente — incluso cuando el contenedor está `Up`, el proceso de arranque de SonarQube continúa. Lo vemos en el siguiente paso."

---

### PASO 5 — Configuración inicial de Jenkins (10:30 – 12:30)

> *Pantalla: navegador en `http://localhost:8080`.*

"Obtengo la contraseña inicial:"

```bash
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```

"Copio ese hash. Abro el navegador en `http://localhost:8080` y lo pego.

La pantalla de bienvenida pregunta cómo instalar los plugins. Click en **'Install suggested plugins'**. Jenkins instalará automáticamente los plugins más comunes — Git, Pipeline, Credentials, etc. Tarda entre 2 y 4 minutos.

Cuando termina, me pide crear el usuario administrador. Lo lleno:
- **Username:** `admin`
- **Password:** una contraseña segura que recuerde
- **Full name:** mi nombre
- **Email:** mi correo

Confirmo la URL: `http://localhost:8080/` — sin cambios. Click en **'Start using Jenkins'**.

El dashboard de Jenkins aparece. Limpio, sin pipelines todavía. Lo vamos a poblar en los próximos episodios."

---

### PASO 6 — Verificar que Jenkins puede usar Docker (12:30 – 13:30)

> *Pantalla: terminal.*

"La verificación más importante: confirmar que el socket de Docker funciona correctamente desde dentro del contenedor de Jenkins. Si este paso falla, los pipelines no van a poder construir imágenes."

```bash
docker exec jenkins docker ps
```

"Si responde con la lista de contenedores corriendo — jenkins, sonarqube, sonar-db — el socket está correctamente montado y Jenkins puede hablar con Docker.

Si responde con un error de permisos, ejecutar `sudo chmod 666 /var/run/docker.sock` y reiniciar el contenedor de Jenkins con `docker compose restart jenkins`."

---

### PASO 7 — Verificar SonarQube (13:30 – 14:30)

> *Pantalla: navegador en `http://localhost:9000`.*

"SonarQube puede tardar entre 1 y 2 minutos en estar completamente listo después de que el contenedor arranca. Si la página no carga, esperar unos segundos y refrescar.

Login con `admin` / `admin`. SonarQube pide cambiar la contraseña en el primer login — lo hago con una contraseña segura y la anoto. La voy a necesitar en el EP43 para generar el token de integración con Jenkins."

---

### CIERRE (14:30 – 15:00)

"Eso es el episodio 31.

Jenkins y SonarQube corriendo localmente en Docker Compose. La configuración persiste en volúmenes nombrados — si cierras la laptop y la vuelves a abrir, `docker compose up -d` y todo está exactamente donde lo dejaste.

En el siguiente episodio instalamos los plugins que el pipeline de CI va a necesitar: Docker Pipeline, SonarQube Scanner, y algunos más. Eso es lo que convierte a este Jenkins básico en la herramienta completa para el curso.

Nos vemos en el EP32."

---

## ✅ Checklist de Verificación
- [ ] `docker compose ps` muestra 3 contenedores en estado `Up`
- [ ] Jenkins responde en `http://localhost:8080`
- [ ] Usuario `admin` creado en Jenkins
- [ ] SonarQube responde en `http://localhost:9000`
- [ ] Contraseña de SonarQube cambiada y anotada
- [ ] `docker exec jenkins docker ps` muestra los contenedores sin errores

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Puerto 8080 ya en uso | `sudo lsof -i :8080` para ver qué proceso lo usa — o cambiar el puerto en el compose |
| `Got permission denied: /var/run/docker.sock` | `sudo chmod 666 /var/run/docker.sock` y `docker compose restart jenkins` |
| SonarQube no arranca (contenedor se reinicia) | `sudo sysctl -w vm.max_map_count=262144` |
| SonarQube tarda en cargar | Normal — el proceso de arranque interno toma 1-2 minutos |
| Jenkins pierde la configuración al reiniciar | Usar `docker compose stop` / `start` — nunca `down -v` en producción |

---

## 🗒️ Notas de Producción
- La apertura con la tabla de costos y el $0 es el hook del episodio — enfatizar que Jenkins local elimina completamente el costo de una segunda EC2.
- Al explicar los volúmenes del compose, señalar con el cursor cada línea mientras se explica — el alumno necesita conectar la sintaxis con el comportamiento.
- Hacer el `docker compose up -d` con la terminal visible y esperar a que terminen las descargas — mostrar que la primera vez tarda y las siguientes son rápidas.
- La verificación con `docker exec jenkins docker ps` es el momento técnico más crítico del episodio — si falla aquí, el pipeline del EP36 nunca va a funcionar.
- Abrir el navegador para SonarQube y cambiar la contraseña en vivo — así el alumno ve que está completamente funcional desde el día uno.

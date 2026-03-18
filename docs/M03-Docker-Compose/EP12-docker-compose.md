# EP 12: Docker Compose — Stack Multi-Contenedor

**Tipo:** PRÁCTICA  
**Duración estimada:** 12–15 min  
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Usar Docker Compose para levantar el proyecto curso-gitops completo (app + MySQL) en un solo comando, entender la estructura del archivo `docker-compose.yml`, y dominar los comandos de gestión.

---

## 📋 Prerequisitos
- Docker con Compose plugin instalado (EP08)
- Proyecto `curso-gitops` con Dockerfile (EP10)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Abrir el Archivo `docker-compose.yml`

**Ejecutar en:** VS Code

```bash
cd ~/UTTT/Proyectos/GitOps/curso-gitops
code docker-compose.yml
```

**Archivo:** `curso-gitops/docker-compose.yml`

```yaml
services:
  api:
    build: .                           # Construye desde el Dockerfile local
    ports: ["8080:8080"]               # Mapea puerto
    environment:                       # Variables de entorno
      - DB_HOST=mysql-db               # Nombre del servicio MySQL
      - DB_USER=curso_app
      - DB_PASSWORD=C4rs0_S3cur3_P@ss!
      - DB_NAME=curso_db
      - JWT_SECRET=gk8s_pr0d_s3cr3t_ch4ng3_m3!
    depends_on:                        # Espera a que MySQL esté listo
      mysql-db: { condition: service_healthy }
    networks: [db-curso, web-curso]
    restart: unless-stopped

  mysql-db:
    image: mysql:8.0                   # Usa imagen oficial de MySQL
    environment:
      - MYSQL_ROOT_PASSWORD=r00t_S3cur3_P@ss!
      - MYSQL_DATABASE=curso_db
      - MYSQL_USER=curso_app
      - MYSQL_PASSWORD=C4rs0_S3cur3_P@ss!
    volumes:
      - mysql-data:/var/lib/mysql              # Datos persistentes
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql  # Script de inicialización
    networks: [db-curso]
    healthcheck:                       # Verifica que MySQL esté listo
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

networks:
  db-curso: { internal: true }         # Red interna (MySQL no accesible desde fuera)
  web-curso: { driver: bridge }        # Red externa (API accesible)

volumes: { mysql-data: }               # Volumen nombrado para persistencia
```

---

### Paso 2 — Explicar Conceptos Clave

| Concepto | Qué es | Dónde está |
|---|---|---|
| **services** | Cada contenedor que se va a levantar | `api` y `mysql-db` |
| **build** | Construye imagen desde Dockerfile | `api → build: .` |
| **image** | Usa imagen de Docker Hub | `mysql-db → image: mysql:8.0` |
| **depends_on** | Orden de arranque, espera healthcheck | `api → depends_on: mysql-db` |
| **networks** | Redes para comunicación entre contenedores | `db-curso` (interna), `web-curso` (expuesta) |
| **volumes** | Datos persistentes (sobreviven `docker compose down`) | `mysql-data` |
| **healthcheck** | Verifica que el servicio está listo | `mysqladmin ping` cada 10s |

> 📌 **Seguridad de redes:** MySQL solo es accesible dentro de la red `db-curso` (internal: true). No tiene puerto expuesto al exterior.

---

### Paso 3 — Levantar el Stack

**Ejecutar en:** Terminal, en la raíz de `curso-gitops/`

```bash
cd ~/UTTT/Proyectos/GitOps/curso-gitops

# Construir y levantar en background
docker compose up --build -d
```

**Resultado esperado:**
```
✔ Network curso-gitops_db-curso    Created
✔ Network curso-gitops_web-curso   Created
✔ Volume "curso-gitops_mysql-data" Created
✔ Container curso-gitops-mysql-db-1 Healthy
✔ Container curso-gitops-api-1      Started
```

---

### Paso 4 — Verificar el Stack

```bash
# Ver servicios corriendo
docker compose ps

# Ver logs de todos los servicios
docker compose logs

# Ver logs de un servicio específico en tiempo real
docker compose logs -f api

# Verificar que la app responde
curl http://localhost:8080
```

**Probar en el navegador:** Ir a **http://localhost:8080**
- Debe mostrar la página de login del curso
- Login con: **Leo** / **admin123**

---

### Paso 5 — Comandos de Gestión

```bash
# Detener (mantiene datos en volumen)
docker compose stop

# Levantar de nuevo (sin rebuild)
docker compose start

# Detener y eliminar contenedores + redes
docker compose down

# Detener, eliminar contenedores + redes + VOLÚMENES (reset total BD)
docker compose down -v

# Rebuild forzado (después de cambios en el código)
docker compose up --build -d

# Ver uso de recursos
docker compose top
```

---

### Paso 6 — El Archivo `init.sql`

**Archivo:** `curso-gitops/init.sql`

Este archivo se ejecuta automáticamente la **primera vez** que se crea el volumen de MySQL. Crea las tablas y los usuarios seed.

```bash
# Ver qué hay en init.sql
cat init.sql
```

> 📌 Si necesitas reiniciar la BD desde cero: `docker compose down -v && docker compose up --build -d`

---

## 🗂️ Archivos del Proyecto Involucrados

| Archivo | Dónde está | Qué hace |
|---|---|---|
| `docker-compose.yml` | `curso-gitops/docker-compose.yml` | Define los servicios (API + MySQL) |
| `Dockerfile` | `curso-gitops/Dockerfile` | Construye la imagen de la API |
| `init.sql` | `curso-gitops/init.sql` | Crea BD, tablas y usuarios seed |

---

## ✅ Checklist de Verificación
- [ ] `docker compose up --build -d` levanta ambos servicios sin error
- [ ] `docker compose ps` muestra 2 servicios running/healthy
- [ ] La app responde en **http://localhost:8080**
- [ ] Puedes loguearte con **Leo** / **admin123**
- [ ] `docker compose down -v` limpia todo correctamente

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `api` se reinicia constantemente | MySQL no ha terminado de inicializarse — esperar o chequear `docker compose logs mysql-db` |
| `Access denied for user` | Hacer `docker compose down -v` para resetear la BD |
| Puerto 8080 ya en uso | `sudo lsof -i :8080` para ver qué lo usa, o cambiar el puerto en compose |
| `init.sql` no se ejecuta | Solo se ejecuta la **primera vez**. Si ya existe el volumen: `docker compose down -v` |

---

## 📌 Notas para el Video
- Abrir `docker-compose.yml` y explicar cada sección
- Levantar todo con `docker compose up --build -d`
- Mostrar `docker compose ps` y `docker compose logs`
- Abrir la app en el navegador y hacer login
- Mostrar `docker compose down -v` para reset total
- Terminar explicando que esto es lo que Kubernetes hará en producción, pero a escala

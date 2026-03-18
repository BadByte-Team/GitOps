# EP 09: Imágenes y Contenedores — Comandos Esenciales

**Tipo:** PRÁCTICA  
**Duración estimada:** 12–15 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Dominar los comandos esenciales de Docker para descargar imágenes, ejecutar contenedores, ver logs, entrar a un contenedor y limpiar recursos.

---

## 📋 Prerequisitos
- Docker instalado y funcionando sin sudo (EP08)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Descargar una Imagen

**Ejecutar en:** Terminal con Docker

> 📌 **Docker Hub** es el registro por defecto: https://hub.docker.com

```bash
# Buscar imágenes en Docker Hub
# Ir a: https://hub.docker.com y buscar "nginx"

# Descargar la imagen oficial de Nginx
docker pull nginx:latest

# Descargar Ubuntu
docker pull ubuntu:22.04

# Ver las imágenes descargadas
docker images
```

**Resultado esperado de `docker images`:**
```
REPOSITORY   TAG       IMAGE ID       CREATED       SIZE
nginx        latest    abc123def      2 weeks ago   187MB
ubuntu       22.04     ghi456jkl      3 weeks ago   77.8MB
```

---

### Paso 2 — Ejecutar un Contenedor

```bash
# Ejecutar Nginx en modo detached (background)
docker run -d --name mi-nginx -p 8080:80 nginx:latest
```

**Desglose del comando:**
| Flag | Significado |
|---|---|
| `-d` | Detached — corre en background |
| `--name mi-nginx` | Nombre del contenedor |
| `-p 8080:80` | Mapear puerto local 8080 → puerto 80 del contenedor |
| `nginx:latest` | Imagen y tag a usar |

**Verificar:** Abrir en el navegador: **http://localhost:8080** → debe mostrar "Welcome to nginx!"

---

### Paso 3 — Gestionar Contenedores

```bash
# Listar contenedores activos
docker ps

# Listar TODOS (incluyendo detenidos)
docker ps -a

# Ver logs del contenedor
docker logs mi-nginx

# Ver logs en tiempo real (follow)
docker logs -f mi-nginx
# Presionar Ctrl+C para salir

# Ver uso de recursos (CPU, memoria)
docker stats mi-nginx
# Presionar Ctrl+C para salir
```

---

### Paso 4 — Entrar a un Contenedor

```bash
# Abrir una shell dentro del contenedor
docker exec -it mi-nginx bash

# Ahora estás DENTRO del contenedor:
cat /etc/nginx/nginx.conf    # Ver config de Nginx
ls /usr/share/nginx/html/    # Ver archivos web
exit                          # Salir del contenedor
```

**Desglose:**
| Flag | Significado |
|---|---|
| `-i` | Interactive — mantener STDIN abierto |
| `-t` | TTY — asignar terminal |
| `bash` | Comando a ejecutar dentro del contenedor |

> 📌 Algunos contenedores usan `sh` en vez de `bash` (ej: Alpine)

---

### Paso 5 — Detener y Eliminar

```bash
# Detener (no elimina)
docker stop mi-nginx

# Verificar que se detuvo
docker ps
# No aparece (está detenido)

docker ps -a
# Aparece con STATUS "Exited"

# Eliminar el contenedor
docker rm mi-nginx

# Eliminar una imagen
docker rmi nginx:latest
```

---

### Paso 6 — Limpieza Masiva

```bash
# Eliminar todos los contenedores detenidos
docker container prune -f

# Eliminar todas las imágenes no usadas
docker image prune -a -f

# Limpieza total (contenedores, imágenes, redes, cache)
docker system prune -a -f
```

> ⚠️ `docker system prune -a` elimina TODO lo no utilizado. Usar con cuidado.

---

### Resumen de Comandos

| Comando | Acción | Ejemplo |
|---|---|---|
| `docker pull` | Descargar imagen | `docker pull nginx` |
| `docker images` | Listar imágenes | `docker images` |
| `docker run -d` | Ejecutar en background | `docker run -d -p 8080:80 nginx` |
| `docker ps` | Listar contenedores activos | `docker ps` |
| `docker ps -a` | Listar todos los contenedores | `docker ps -a` |
| `docker logs` | Ver logs | `docker logs -f mi-nginx` |
| `docker exec -it` | Entrar al contenedor | `docker exec -it mi-nginx bash` |
| `docker stop` | Detener | `docker stop mi-nginx` |
| `docker rm` | Eliminar contenedor | `docker rm mi-nginx` |
| `docker rmi` | Eliminar imagen | `docker rmi nginx:latest` |

---

## ✅ Checklist de Verificación
- [ ] Puedes descargar una imagen con `docker pull`
- [ ] Puedes ejecutar un contenedor y acceder en `http://localhost:8080`
- [ ] Puedes ver los logs con `docker logs`
- [ ] Puedes entrar al contenedor con `docker exec`
- [ ] Puedes detener y eliminar contenedores

---

## 📌 Notas para el Video
- Abrir Docker Hub en el navegador y mostrar la imagen de Nginx
- Ejecutar el contenedor y abrir `localhost:8080` para ver la página de Nginx
- Entrar al contenedor con `exec` y explorar el filesystem
- Mostrar `docker stats` para el uso de recursos en vivo
- Limpiar todo al final con `docker system prune`

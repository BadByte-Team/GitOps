# EP 11: Docker Hub — Publicar y Versionar Imágenes

**Tipo:** PRÁCTICA  
**Duración estimada:** 10–12 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Publicar la imagen de curso-gitops en Docker Hub con tags versionados para que pueda ser descargada desde cualquier máquina o cluster de Kubernetes.

---

## 📋 Prerequisitos
- Docker instalado (EP08)
- Imagen construida (EP10)
- Cuenta en Docker Hub

---

## 📺 Paso a Paso para el Video

### Paso 1 — Crear Cuenta en Docker Hub

1. Ir a: **https://hub.docker.com/signup**
2. Registrarse con email y contraseña
3. Verificar el email
4. Recordar tu **username de Docker Hub** (lo usarás como prefijo de tus imágenes)

---

### Paso 2 — Login desde la Terminal

**Ejecutar en:** Terminal con Docker

```bash
docker login
# Username: TU_USUARIO_DOCKERHUB
# Password: tu contraseña
```

**Resultado esperado:**
```
Login Succeeded
```

> 📌 **Alternativa más segura:** Usar Access Tokens en vez de contraseña:
> 1. Docker Hub → Account Settings → Security → New Access Token
> 2. `docker login -u TU_USUARIO --password-stdin <<< "TU_TOKEN"`

---

### Paso 3 — Tag de la Imagen

La imagen necesita un tag con formato `usuario/repositorio:version`:

```bash
# Tag con versión
docker tag curso-gitops:v1 TU_USUARIO/curso-gitops:v1

# Tag latest (siempre apunta a la última versión)
docker tag curso-gitops:v1 TU_USUARIO/curso-gitops:latest

# Verificar
docker images | grep curso-gitops
```

**Resultado esperado:**
```
TU_USUARIO/curso-gitops   v1       abc123   5 min ago   25MB
TU_USUARIO/curso-gitops   latest   abc123   5 min ago   25MB
curso-gitops              v1       abc123   5 min ago   25MB
```

---

### Paso 4 — Push a Docker Hub

```bash
docker push TU_USUARIO/curso-gitops:v1
docker push TU_USUARIO/curso-gitops:latest
```

> El primer push sube todas las capas. Los siguientes solo las capas que cambiaron.

---

### Paso 5 — Verificar en Docker Hub

1. Ir a: **https://hub.docker.com/r/TU_USUARIO/curso-gitops**
2. Debe mostrar los tags `v1` y `latest`
3. Click en un tag para ver las capas y el tamaño

---

### Paso 6 — Probar el Pull desde Otro Lugar

```bash
# Eliminar la imagen local
docker rmi TU_USUARIO/curso-gitops:v1

# Descargar desde Docker Hub
docker pull TU_USUARIO/curso-gitops:v1

# Verificar
docker images | grep curso-gitops
```

---

### Convención de Tags — Buenas Prácticas

| Tag | Cuándo usarlo | Ejemplo |
|---|---|---|
| `latest` | Siempre apunta a la última versión | `user/app:latest` |
| `v1`, `v2` | Versión mayor | `user/app:v2` |
| `1.0.0` | Semantic versioning | `user/app:1.0.0` |
| `build-42` | Número de build CI | `user/app:build-42` |
| `42-abc1234` | Build + commit hash | `user/app:42-abc1234` |

> 📌 **En nuestro Jenkinsfile** usamos: `BUILD_NUMBER-GIT_COMMIT_SHORT`, por ejemplo: `15-a3b8d1c`

---

## 🗂️ Archivos del Proyecto Involucrados
| Archivo | Relevancia |
|---|---|
| `curso-gitops/Dockerfile` | La imagen que se sube parte de este archivo |
| `infrastructure/jenkins/Jenkinsfile` | Automatiza el build + push (stage "Docker Push", líneas 65-70) |
| `infrastructure/kubernetes/app/deployment.yaml` | Referencia la imagen de Docker Hub (línea `image:`) |

---

## ✅ Checklist de Verificación
- [ ] Tienes cuenta en Docker Hub
- [ ] `docker login` funciona
- [ ] La imagen aparece en `https://hub.docker.com/r/TU_USUARIO/curso-gitops`
- [ ] Puedes hacer `docker pull` de tu imagen desde otra máquina

---

## 📌 Notas para el Video
- Mostrar Docker Hub en el navegador (registro y dashboard)
- Hacer login desde terminal
- Tag y push en vivo
- Ir a Docker Hub y mostrar la imagen publicada
- Explicar la convención de tags que usaremos en Jenkins

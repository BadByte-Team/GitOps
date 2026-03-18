# EP 10: Escribir un Dockerfile desde Cero

**Tipo:** PRÁCTICA  
**Duración estimada:** 12–15 min  
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Entender cada instrucción de un Dockerfile, construir la imagen del proyecto curso-gitops y ejecutarla localmente.

---

## 📋 Prerequisitos
- Docker instalado (EP08)
- El proyecto `curso-gitops` clonado

---

## 📺 Paso a Paso para el Video

### Paso 1 — Abrir el Dockerfile del Proyecto

**Ejecutar en:** VS Code, en la raíz del proyecto

```bash
cd ~/UTTT/Proyectos/GitOps/curso-gitops
code Dockerfile
```

**Archivo:** `curso-gitops/Dockerfile`

```dockerfile
# ── Etapa 1: Build ──
FROM golang:1.25.8-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o curso-gitops ./cmd/api

# ── Etapa 2: Runtime ──
FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /home/appuser
COPY --from=builder --chown=appuser:appgroup /app/curso-gitops .
COPY --from=builder --chown=appuser:appgroup /app/frontend ./frontend
USER appuser
EXPOSE 8080
CMD ["./curso-gitops"]
```

---

### Paso 2 — Explicar Cada Instrucción

| Instrucción | Qué hace | Línea |
|---|---|---|
| `FROM golang:1.25.8-alpine AS builder` | Usa imagen Go para compilar (etapa 1) | 2 |
| `WORKDIR /app` | Establece directorio de trabajo | 3 |
| `COPY go.mod go.sum ./` | Copia dependencias PRIMERO (cache) | 4 |
| `RUN go mod download` | Descarga dependencias (se cachea si go.mod no cambia) | 5 |
| `COPY . .` | Copia todo el código fuente | 6 |
| `RUN CGO_ENABLED=0 ...` | Compila el binario | 7 |
| `FROM alpine:latest` | Imagen final mínima (5MB) — etapa 2 | 9 |
| `RUN addgroup... adduser...` | Crea usuario non-root (seguridad) | 10 |
| `COPY --from=builder` | Copia SOLO el binario de la etapa 1 | 12-13 |
| `USER appuser` | Ejecuta como non-root | 14 |
| `EXPOSE 8080` | Documenta el puerto (no lo abre) | 15 |
| `CMD [...]` | Comando al iniciar el contenedor | 16 |

> 📌 **Multi-stage build:** La imagen final NO tiene Go, compiladores ni código fuente — solo el binario y el frontend. Más segura y más pequeña.

---

### Paso 3 — Construir la Imagen

**Ejecutar en:** Terminal, en la raíz de `curso-gitops/`

```bash
cd ~/UTTT/Proyectos/GitOps/curso-gitops

# Construir con tag
docker build -t curso-gitops:v1 .
```

**Resultado esperado:** Muestra cada paso y termina con `Successfully tagged curso-gitops:v1`

```bash
# Verificar la imagen creada
docker images | grep curso-gitops
# curso-gitops   v1   abc123   5 seconds ago   25MB
```

> 📌 Notar que la imagen final es ~25MB (Alpine + binario Go estático)

---

### Paso 4 — Ejecutar el Contenedor

```bash
docker run -p 8080:8080 curso-gitops:v1
```

> ⚠️ Va a fallar porque no tiene conexión a MySQL. Es normal — la app necesita `docker-compose` para levantar con la BD. Pero puedes ver que el contenedor inicia e intenta conectar.

Para detener: `Ctrl+C`

---

### Paso 5 — Inspeccionar la Imagen

```bash
# Ver las capas de la imagen
docker history curso-gitops:v1

# Ver info detallada
docker inspect curso-gitops:v1 | head -50
```

---

## ��️ Archivos del Proyecto Involucrados
| Archivo | Dónde está | Qué hacer |
|---|---|---|
| `Dockerfile` | `curso-gitops/Dockerfile` | Revisar, entender cada línea |
| `.dockerignore` | `curso-gitops/.dockerignore` | Crear si no existe (excluir .git, docs, etc.) |

---

## ✅ Checklist de Verificación
- [ ] Entiendes cada instrucción del Dockerfile
- [ ] `docker build -t curso-gitops:v1 .` compila sin errores
- [ ] `docker images` muestra la imagen creada (~25MB)
- [ ] Entiendes el concepto de multi-stage build

---

## 📌 Notas para el Video
- Abrir el Dockerfile y explicar línea por línea
- Mostrar el concepto de multi-stage build con un diagrama mental
- Construir la imagen en vivo y mostrar el tamaño final
- Intentar ejecutar para mostrar que necesita la BD (setup completo en EP12)
- Comparar con una imagen sin multi-stage para ver la diferencia de tamaño

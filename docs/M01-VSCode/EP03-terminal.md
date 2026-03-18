# EP 03: Terminal Integrada y Atajos Clave

**Tipo:** PRÁCTICA  
**Duración estimada:** 10–12 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Dominar la terminal integrada de VS Code y los atajos de teclado que vas a usar durante todo el curso para navegar archivos, editar código y ejecutar comandos.

---

## 📋 Prerequisitos
- VS Code instalado con extensiones (EP01, EP02)
- Tener el proyecto `curso-gitops` clonado (o al menos una carpeta para practicar)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Abrir VS Code en el Proyecto

**Ejecutar en:** Terminal de tu máquina local

```bash
cd ~/UTTT/Proyectos/GitOps/curso-gitops
code .
```

> VS Code abre con el explorador de archivos mostrando la estructura del proyecto.

---

### Paso 2 — Terminal Integrada

#### Abrir la Terminal
- **Atajo:** `` Ctrl+` `` (la tecla backtick, al lado del 1)
- **Menú:** Terminal → New Terminal
- La terminal se abre **abajo** del editor, dentro de VS Code

#### Crear y Dividir Terminales
| Acción | Atajo |
|---|---|
| Abrir/cerrar terminal | `` Ctrl+` `` |
| Nueva terminal | `` Ctrl+Shift+` `` |
| Dividir terminal (lado a lado) | `` Ctrl+Shift+5 `` |
| Cambiar entre terminales | `Ctrl+PageUp` / `Ctrl+PageDown` |
| Cerrar terminal activa | Escribir `exit` o click en el ícono de basura |

#### Práctica en el Video
1. Abrir la terminal: `` Ctrl+` ``
2. Ejecutar: `ls -la` para ver los archivos del proyecto
3. Crear una segunda terminal: `` Ctrl+Shift+` ``
4. En la segunda terminal: `docker compose ps` (o cualquier comando)
5. Dividir terminal: `Ctrl+Shift+5`
6. Mostrar que puedes tener un terminal para Git y otro para Docker al mismo tiempo

---

### Paso 3 — Navegación de Archivos

| Atajo | Qué Hace | Cuándo Usarlo |
|---|---|---|
| `Ctrl+P` | **Buscar archivo por nombre** | Abrir rápidamente `docker-compose.yml`, `main.go`, etc. |
| `Ctrl+Shift+F` | **Buscar texto en todo el proyecto** | Buscar "8080" en todos los archivos |
| `Ctrl+Shift+E` | Abrir explorador de archivos | Ver la estructura del proyecto |
| `Ctrl+G` | Ir a una línea específica | Ir a la línea 45 de un archivo |
| `Ctrl+Tab` | Cambiar entre archivos abiertos | Alternar entre archivos rápidamente |

#### Práctica en el Video
1. Presionar `Ctrl+P` → escribir `docker` → seleccionar `docker-compose.yml`
2. Presionar `Ctrl+P` → escribir `main` → seleccionar `cmd/api/main.go`
3. Presionar `Ctrl+Shift+F` → buscar `8080` → ver todos los archivos donde aparece
4. En un archivo, presionar `Ctrl+G` → escribir `20` → salta a la línea 20

---

### Paso 4 — Atajos de Edición

| Atajo | Qué Hace |
|---|---|
| `Alt+↑` / `Alt+↓` | Mover línea arriba/abajo |
| `Alt+Shift+↑` / `Alt+Shift+↓` | Duplicar línea |
| `Ctrl+D` | Seleccionar siguiente ocurrencia |
| `Ctrl+Shift+L` | Seleccionar TODAS las ocurrencias |
| `Ctrl+/` | Comentar/descomentar línea |
| `Ctrl+Shift+K` | Eliminar línea completa |
| `Ctrl+Z` | Deshacer |
| `Ctrl+Shift+Z` | Rehacer |
| `Alt+Click` | Agregar cursor (multi-cursor) |
| `Ctrl+L` | Seleccionar línea completa |

#### Práctica en el Video
1. Abrir `docker-compose.yml`
2. Mover una línea con `Alt+↑`
3. Duplicar una línea con `Alt+Shift+↓`
4. Seleccionar todas las ocurrencias de `curso` con `Ctrl+Shift+L` → cambiar todos a la vez
5. Comentar varias líneas seleccionándolas y presionando `Ctrl+/`

---

### Paso 5 — Paleta de Comandos

El atajo más importante de VS Code:

| Atajo | Qué Hace |
|---|---|
| `Ctrl+Shift+P` | **Abre la paleta de comandos** — desde aquí puedes hacer CUALQUIER cosa |

#### Ejemplos Útiles
- Escribir `theme` → cambiar el tema del editor
- Escribir `terminal` → crear nueva terminal
- Escribir `format` → formatear el archivo actual
- Escribir `git` → ver todos los comandos de Git disponibles
- Escribir `preferences` → abrir settings

---

### Paso 6 — Zen Mode y Trucos Finales

| Atajo | Qué Hace |
|---|---|
| `Ctrl+B` | Mostrar/ocultar la barra lateral |
| `Ctrl+J` | Mostrar/ocultar el panel inferior (terminal) |
| `Ctrl+K Z` | **Zen Mode** — pantalla completa sin distracciones |
| `Ctrl+\` (backslash) | Dividir el editor en dos columnas |

---

## 🗂️ Archivos del Proyecto que Usarás para Practicar

| Archivo | Ruta desde la raíz del proyecto |
|---|---|
| Docker Compose | `curso-gitops/docker-compose.yml` |
| Main del API | `curso-gitops/cmd/api/main.go` |
| Dockerfile | `curso-gitops/Dockerfile` |
| Terraform backend | `infrastructure/terraform/backend/main.tf` |
| K8s deployment | `infrastructure/kubernetes/app/deployment.yaml` |

---

## ✅ Checklist de Verificación
- [ ] Puedes abrir la terminal con `` Ctrl+` ``
- [ ] Puedes dividir la terminal en dos paneles
- [ ] `Ctrl+P` abre el buscador de archivos y encuentras rápido cualquier archivo
- [ ] `Ctrl+Shift+F` busca texto en todo el proyecto
- [ ] Puedes mover y duplicar líneas con Alt+flechas
- [ ] Conoces la paleta de comandos `Ctrl+Shift+P`

---

## 📌 Notas para el Video
- Abrir el proyecto y mostrar la terminal integrada
- Dividir terminal y ejecutar dos comandos simultáneamente
- Hacer una demostración de `Ctrl+P` buscando archivos rápidamente
- Mostrar multi-cursor editando varias líneas a la vez
- Cerrar con la paleta de comandos como "el atajo que lo puede todo"

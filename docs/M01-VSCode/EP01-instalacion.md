# EP 01: Instalación de VS Code en Windows, Mac y Linux

**Tipo:** INSTALACIÓN  
**Duración estimada:** 10–15 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Descargar, instalar y configurar Visual Studio Code en tu sistema operativo. Al terminar, tendrás VS Code listo con los ajustes básicos para trabajar con Docker, Kubernetes, Terraform y archivos YAML.

---

## 📋 Prerequisitos
- Computadora con Windows 10+, macOS 10.15+, o Linux (Ubuntu 20.04+)
- Conexión a internet estable
- Permisos de administrador en tu máquina

---

## 📺 Paso a Paso para el Video

### Paso 1 — Ir a la Página de Descarga

1. Abrir el navegador e ir a: **https://code.visualstudio.com/**
2. El sitio detecta tu sistema operativo automáticamente
3. Hacer click en el botón azul grande **"Download for [tu SO]"**

> 📌 **URL directa por SO:**
> - Windows: https://code.visualstudio.com/sha/download?build=stable&os=win32-x64-user
> - Mac (Intel): https://code.visualstudio.com/sha/download?build=stable&os=darwin
> - Mac (Apple Silicon): https://code.visualstudio.com/sha/download?build=stable&os=darwin-arm64
> - Linux .deb: https://code.visualstudio.com/sha/download?build=stable&os=linux-deb-x64

---

### Paso 2 — Instalación por Sistema Operativo

#### 🐧 Linux (Ubuntu/Debian) — Método recomendado (repositorio APT)

**Ejecutar en:** Tu terminal local (`Ctrl+Alt+T`)

```bash
# 1. Instalar dependencias
sudo apt update
sudo apt install -y software-properties-common apt-transport-https wget

# 2. Importar la clave GPG de Microsoft
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
rm packages.microsoft.gpg

# 3. Agregar el repositorio oficial
echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" | sudo tee /etc/apt/sources.list.d/vscode.list

# 4. Instalar VS Code
sudo apt update
sudo apt install -y code
```

**Método alternativo (Snap):**
```bash
sudo snap install code --classic
```

#### 🪟 Windows
1. Ejecutar el archivo `.exe` descargado
2. Aceptar el acuerdo de licencia
3. **Marcar las opciones:**
   - ✅ Agregar "Open with Code" al menú contextual de archivos
   - ✅ Agregar "Open with Code" al menú contextual de directorios
   - ✅ Agregar a PATH (importante para usar `code` desde terminal)
4. Click en "Install" → esperar → "Finish"

#### 🍎 Mac
1. Abrir el archivo `.dmg` descargado
2. Arrastrar VS Code al folder **Applications**
3. Para usar desde terminal, abrir VS Code → `Cmd+Shift+P` → escribir "Shell Command" → click en **"Install 'code' command in PATH"**

---

### Paso 3 — Verificar la Instalación

**Ejecutar en:** Terminal de tu máquina local

```bash
code --version
```

**Resultado esperado:**
```
1.96.x
abc12345def
x64
```

Abrir VS Code en el directorio actual:
```bash
code .
```

---

### Paso 4 — Primeros Ajustes del Editor

Abrir VS Code y configurar los siguientes ajustes:

#### Cambiar el Tema
1. Ir a: **File → Preferences → Color Theme** (o `Ctrl+K Ctrl+T`)
2. Seleccionar: **"Dark Modern"** (viene incluido) o **"One Dark Pro"** (extensión)

#### Cambiar Tamaño de Fuente
1. Ir a: **File → Preferences → Settings** (o `Ctrl+,`)
2. Buscar: `Font Size`
3. Cambiar a **14** o **15**

#### Activar Auto-Save
1. En Settings, buscar: `Auto Save`
2. Cambiar a: **afterDelay** (guarda automáticamente cada segundo)

#### Terminal por Defecto
1. En Settings, buscar: `Terminal Default Profile`
2. **Linux/Mac:** Seleccionar `bash` o `zsh`
3. **Windows:** Seleccionar `Git Bash` o `PowerShell`

#### Mostrar Minimap (opcional)
1. En Settings, buscar: `Minimap`
2. Desactivar si prefieres más espacio

---

### Paso 5 — Abrir el Proyecto del Curso

```bash
# Si ya tienes el repo clonado:
cd ~/UTTT/Proyectos/GitOps/curso-gitops
code .
```

> 💡 VS Code abre con el explorador de archivos a la izquierda mostrando la estructura del proyecto.

---

## 🗂️ Archivos del Proyecto Involucrados
*Ninguno — este es un episodio de setup del entorno.*

---

## ✅ Checklist de Verificación
- [ ] VS Code instalado y abre correctamente
- [ ] `code --version` funciona desde la terminal
- [ ] `code .` abre VS Code en el directorio actual
- [ ] Tema oscuro configurado
- [ ] Tamaño de fuente ajustado
- [ ] Auto-Save activado

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `code: command not found` (Linux) | Cerrar y abrir la terminal, o ejecutar `source ~/.bashrc` |
| `code: command not found` (Mac) | Abrir VS Code → Cmd+Shift+P → "Install 'code' command" |
| VS Code se ve borroso (Linux) | Lanchar con `code --disable-gpu` |
| No aparece la opción de terminal | Reinstalar VS Code desde el repositorio oficial |

---

## 📌 Notas para el Video
- Mostrar la descarga desde la web oficial
- Hacer la instalación paso a paso en pantalla
- Abrir VS Code y cambiar el tema en vivo
- Abrir la terminal integrada con `Ctrl+`` para mostrar que funciona
- Terminar abriendo el proyecto del curso con `code .`

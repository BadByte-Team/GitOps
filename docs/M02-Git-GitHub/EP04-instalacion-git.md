e# EP 04: Instalación de Git y Configuración Inicial

**Tipo:** INSTALACIÓN / CONFIGURACIÓN  
**Duración estimada:** 12–15 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar Git, configurar tu identidad global, generar SSH keys y conectarlas a GitHub para que puedas clonar, pushear y pullear repositorios de forma segura.

---

## 📋 Prerequisitos
- Terminal (VS Code o nativa)
- Cuenta de GitHub — crearla en: **https://github.com/signup**

---

## 📺 Paso a Paso para el Video

### Paso 1 — Instalar Git

**Ejecutar en:** Terminal de tu máquina local

```bash
# Ubuntu/Debian
sudo apt update && sudo apt install -y git

# Verificar instalación
git --version
# Resultado: git version 2.43.x
```

> 📌 **Windows:** Descargar desde https://git-scm.com/download/win  
> 📌 **Mac:** `brew install git` o `xcode-select --install`

---

### Paso 2 — Configurar tu Identidad

Estos datos aparecen en cada commit que hagas:

```bash
git config --global user.name "Tu Nombre Completo"
git config --global user.email "tu-email@ejemplo.com"
```

Configuraciones adicionales recomendadas:
```bash
# Rama por defecto: main (no master)
git config --global init.defaultBranch main

# Editor por defecto: VS Code
git config --global core.editor "code --wait"

# Colores en la terminal
git config --global color.ui auto

# Verificar toda la configuración
git config --list
```

**Resultado esperado de `git config --list`:**
```
user.name=Tu Nombre
user.email=tu-email@ejemplo.com
init.defaultbranch=main
core.editor=code --wait
```

---

### Paso 3 — Generar SSH Key

Las SSH keys te permiten autenticarte con GitHub sin escribir contraseña cada vez.

```bash
# Generar clave Ed25519 (la más segura y rápida)
ssh-keygen -t ed25519 -C "tu-email@ejemplo.com"
```

Cuando pregunte:
- **Enter file:** Presionar Enter (acepta ruta por defecto `~/.ssh/id_ed25519`)
- **Enter passphrase:** Presionar Enter (sin passphrase) o escribir una para mayor seguridad

```bash
# Verificar que se crearon los archivos
ls -la ~/.ssh/
# Debe mostrar:  id_ed25519  (clave privada - NUNCA compartir)
#                id_ed25519.pub  (clave pública - esta se sube a GitHub)
```

Copiar la clave **pública**:
```bash
cat ~/.ssh/id_ed25519.pub
# Copiar TODO el output (empieza con ssh-ed25519...)
```

---

### Paso 4 — Agregar la SSH Key en GitHub

1. Ir a: **https://github.com/settings/keys**
2. Click en **"New SSH key"**
3. **Title:** `Mi laptop` (o un nombre descriptivo)
4. **Key type:** Authentication Key
5. **Key:** Pegar el contenido de `id_ed25519.pub`
6. Click **"Add SSH key"**
7. Confirmar con tu contraseña de GitHub

---

### Paso 5 — Verificar la Conexión

```bash
ssh -T git@github.com
```

La primera vez preguntará si confías en el host:
```
Are you sure you want to continue connecting (yes/no)? yes
```

**Resultado esperado:**
```
Hi TU_USUARIO! You've successfully authenticated, but GitHub does not provide shell access.
```

> ✅ Si ves ese mensaje, **la conexión SSH funciona correctamente.**

---

### Paso 6 — Iniciar el Agente SSH (si la key tiene passphrase)

Si le pusiste passphrase a tu key, para no escribirla cada vez:
```bash
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
```

---

## 🗂️ Archivos del Proyecto Involucrados
*Ninguno — este es un episodio de configuración del entorno.*

---

## ✅ Checklist de Verificación
- [ ] `git --version` muestra la versión instalada
- [ ] `git config user.name` muestra tu nombre
- [ ] `git config user.email` muestra tu email
- [ ] El archivo `~/.ssh/id_ed25519.pub` existe
- [ ] La SSH key aparece en https://github.com/settings/keys
- [ ] `ssh -T git@github.com` dice "successfully authenticated"

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Permission denied (publickey)` | Verificar que la key pública está en GitHub y que el archivo privado está en `~/.ssh/` |
| `Could not resolve hostname github.com` | Verificar conexión a internet |
| La key tiene passphrase y pide cada vez | Ejecutar `ssh-add` como se indica en el Paso 6 |
| `git config --list` no muestra nada | Verificar que ejecutaste los comandos con `--global` |

---

## 📌 Notas para el Video
- Mostrar la instalación de Git en pantalla
- Ejecutar `git config` y mostrar el resultado
- Generar la SSH key mostrando el proceso completo
- Ir a GitHub → Settings → SSH Keys → pegar la key en vivo
- Terminar con el test `ssh -T git@github.com` exitoso

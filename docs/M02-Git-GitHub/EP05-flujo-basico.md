# EP 05: Flujo Básico — clone, add, commit, push

**Tipo:** PRÁCTICA  
**Duración estimada:** 12–15 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Crear un repositorio en GitHub, clonarlo en tu máquina local y dominar el flujo básico de Git: modificar archivos, hacer stage, commit y push.

---

## 📋 Prerequisitos
- Git instalado y configurado (EP04)
- SSH key agregada en GitHub (EP04)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Crear Repositorio en GitHub

1. Ir a: **https://github.com/new**
2. Configurar:
   - **Repository name:** `curso-gitops`
   - **Description:** `Proyecto del curso de GitOps`
   - **Visibility:** Public
   - ✅ **Add a README file**
   - ✅ **Add .gitignore** → seleccionar template `Go`
3. Click **"Create repository"**

> 📌 Se crea el repo con un `README.md` y `.gitignore` iniciales.

---

### Paso 2 — Clonar el Repositorio

**Ejecutar en:** Terminal de tu máquina local

```bash
# Ir al directorio del proyecto
cd ~/UTTT/Proyectos/GitOps

# Clonar con SSH (el botón verde "Code" → SSH → copiar URL)
git clone git@github.com:TU_USUARIO/curso-gitops.git

# Entrar al directorio
cd curso-gitops

# Verificar
ls -la
# Debe mostrar: .git/  README.md  .gitignore
```

---

### Paso 3 — Entender el Estado con `git status`

```bash
git status
```

**Resultado cuando hay cambios:**
```
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
```

---

### Paso 4 — El Ciclo de Git (add → commit → push)

#### 4.1 Crear un archivo nuevo
```bash
echo "# Notas del Curso GitOps" > notas.md
echo "" >> notas.md
echo "## Módulo 1: VS Code" >> notas.md
echo "- Instalación completada" >> notas.md
```

#### 4.2 Ver el estado
```bash
git status
# notas.md aparece en rojo como "Untracked file"
```

#### 4.3 Stage (agregar al área de preparación)
```bash
# Un archivo específico
git add notas.md

# O todos los cambios de una vez
git add .

git status
# notas.md ahora aparece en verde como "new file"
```

#### 4.4 Commit (guardar snapshot)
```bash
git commit -m "feat: agregar archivo de notas del curso"
```

> 📌 **Convención de commits:**
> - `feat:` nueva funcionalidad
> - `fix:` corrección de bug
> - `docs:` cambios en documentación
> - `chore:` tareas de mantenimiento

#### 4.5 Push (subir a GitHub)
```bash
git push origin main
```

---

### Paso 5 — Verificar en GitHub

1. Ir a: `https://github.com/TU_USUARIO/curso-gitops`
2. El archivo `notas.md` debe aparecer en la lista
3. Click en el archivo para ver su contenido
4. Click en "commits" para ver el historial

---

### Paso 6 — Hacer Más Cambios (repetir el ciclo)

```bash
# Modificar el archivo
echo "## Módulo 2: Git + GitHub" >> notas.md
echo "- SSH keys configuradas" >> notas.md

# Ver qué cambió
git diff

# Stage + commit + push
git add .
git commit -m "docs: agregar notas del módulo 2"
git push origin main
```

---

### Paso 7 — Comandos Útiles

| Comando | Qué Hace | Cuándo Usarlo |
|---|---|---|
| `git status` | Ver archivos modificados/staged | Antes de cada commit |
| `git add <file>` | Stage un archivo específico | Cuando quieres control granular |
| `git add .` | Stage todos los cambios | Cuando todos los cambios van juntos |
| `git commit -m "msg"` | Guardar snapshot con mensaje | Después de hacer stage |
| `git push origin main` | Subir commits a GitHub | Después de commitear |
| `git pull origin main` | Bajar cambios de GitHub | Antes de empezar a trabajar |
| `git log --oneline` | Ver historial en una línea | Para revisar qué commits hay |
| `git diff` | Ver cambios no staged | Para revisar antes de hacer add |
| `git diff --staged` | Ver cambios staged | Para revisar antes de commitear |

---

## ✅ Checklist de Verificación
- [ ] El repositorio `curso-gitops` existe en tu GitHub
- [ ] Puedes clonar con `git clone git@github.com:...`
- [ ] Puedes hacer el ciclo completo: add → commit → push
- [ ] Los cambios se ven reflejados en GitHub
- [ ] `git log --oneline` muestra tus commits

---

## 📌 Notas para el Video
- Crear el repo en GitHub en vivo (mostrar la web)
- Clonar y mostrar la carpeta `.git/` que demuestra que es un repo
- Crear un archivo, mostrrar `git status` en rojo, luego en verde tras `add`
- Hacer commit y push, luego ir a GitHub y mostrar el archivo
- Mostrar `git log --oneline` con los commits acumulados

# EP 06: Ramas, Pull Requests y Merge Conflicts

**Tipo:** PRÁCTICA  
**Duración estimada:** 15–18 min  
**Dificultad:** ⭐⭐ (Intermedio)

---

## �� Objetivo
Trabajar con ramas de Git, abrir Pull Requests en GitHub, revisar y mergear cambios, y resolver conflictos de merge.

---

## 📋 Prerequisitos
- Repositorio clonado y flujo básico dominado (EP05)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Crear y Cambiar a una Rama

**Ejecutar en:** Terminal de tu máquina local, dentro del repo `curso-gitops`

```bash
# Verificar en qué rama estás
git branch
# * main  (el asterisco indica la rama actual)

# Crear rama y cambiar a ella (un solo comando)
git checkout -b feature/nueva-seccion
# Switched to a new branch 'feature/nueva-seccion'

# Verificar
git branch
#   main
# * feature/nueva-seccion
```

> 📌 **Convención de nombres de ramas:**
> - `feature/nombre` → nuevas funcionalidades
> - `fix/nombre` → corrección de bugs
> - `hotfix/nombre` → corrección urgente
> - `docs/nombre` → cambios en documentación

---

### Paso 2 — Hacer Cambios en la Rama

```bash
# Agregar contenido nuevo
echo "## Sección Nueva" >> notas.md
echo "Esta sección fue creada en una rama feature" >> notas.md

# Commit
git add .
git commit -m "feat: agregar nueva sección al documento"
```

---

### Paso 3 — Push de la Rama a GitHub

```bash
git push origin feature/nueva-seccion
```

> 📌 Es el mismo `git push` pero con el nombre de la rama en lugar de `main`.

---

### Paso 4 — Crear Pull Request en GitHub

1. Ir a: `https://github.com/TU_USUARIO/curso-gitops`
2. Aparece un banner amarillo: **"feature/nueva-seccion had recent pushes"** → Click **"Compare & pull request"**
3. Configurar el PR:
   - **Title:** `feat: agregar nueva sección`
   - **Description:** Explicar qué cambios se hicieron y por qué
   - **Base:** `main` ← **Compare:** `feature/nueva-seccion`
4. Click **"Create pull request"**

---

### Paso 5 — Revisar y Mergear el PR

1. En la página del PR, ir a la pestaña **"Files changed"**
2. Revisar los cambios (líneas verdes = agregadas, rojas = eliminadas)
3. Si todo está bien, click **"Merge pull request"**
4. Click **"Confirm merge"**
5. Opcional: **"Delete branch"** para limpiar

---

### Paso 6 — Actualizar tu Local

```bash
# Volver a main
git checkout main

# Bajar los cambios mergeados
git pull origin main

# Verificar que los cambios están
cat notas.md

# Eliminar la rama local (ya se mergeó)
git branch -d feature/nueva-seccion
```

---

### Paso 7 — Simular y Resolver un Conflicto de Merge

#### 7.1 Preparar el conflicto
```bash
# En main, editar la línea 1 de un archivo
git checkout main
echo "Contenido desde MAIN" > conflicto.txt
git add . && git commit -m "cambio en main"

# Crear rata con contenido diferente en la misma línea
git checkout -b feature/conflicto
echo "Contenido desde FEATURE" > conflicto.txt
git add . && git commit -m "cambio en feature"
```

#### 7.2 Intentar merge
```bash
git checkout main
git merge feature/conflicto
```

**Resultado:**
```
Auto-merging conflicto.txt
CONFLICT (content): Merge conflict in conflicto.txt
Automatic merge failed; fix conflicts and then commit the result.
```

#### 7.3 Resolver el conflicto

Abrir `conflicto.txt` en VS Code:
```
<<<<<<< HEAD
Contenido desde MAIN
=======
Contenido desde FEATURE
>>>>>>> feature/conflicto
```

**Opciones en VS Code (aparecen como botones arriba del conflicto):**
- "Accept Current Change" → mantener MAIN
- "Accept Incoming Change" → mantener FEATURE
- "Accept Both Changes" → mantener ambos
- Editar manualmente → escribir lo que quieras

Elegir una opción, luego:
```bash
git add conflicto.txt
git commit -m "fix: resolver conflicto en conflicto.txt"
```

> ✅ ¡Conflicto resuelto!

---

### Resumen de Comandos de Ramas

| Comando | Qué Hace |
|---|---|
| `git branch` | Ver ramas locales |
| `git branch -a` | Ver ramas locales y remotas |
| `git checkout -b nombre` | Crear rama y cambiar a ella |
| `git checkout main` | Cambiar a main |
| `git merge nombre` | Mergear rama a la actual |
| `git branch -d nombre` | Eliminar rama mergeada |
| `git push origin nombre` | Push de rama a GitHub |

---

## ✅ Checklist de Verificación
- [ ] Puedes crear ramas con `git checkout -b`
- [ ] Puedes pushear ramas y crear PRs en GitHub
- [ ] Puedes revisar y mergear un PR desde la web
- [ ] Puedes resolver un conflicto de merge en VS Code

---

## 📌 Notas para el Video
- Crear la rama en terminal y mostrar `git branch`
- Hacer cambios, push, y mostrar el banner de PR en GitHub
- Crear el PR mostrando la interfaz web completa
- Mergear y mostrar cómo `git pull` trae los cambios
- Simular un conflicto y resolverlo en VS Code mostrando los botones de resolución

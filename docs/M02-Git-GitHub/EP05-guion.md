# 🎬 Guión — EP05: Flujo Básico de Git (clone, add, commit, push)

**Duración estimada:** 12–15 min  
**Tono:** Directo, didáctico, con pausas visuales para que el espectador pueda seguir en su máquina.

---

## 🎙️ INTRO (0:00 – 0:45)

> *Pantalla: terminal limpia o slide de título.*

"Bienvenidos al episodio 5. Hasta aquí ya tienes Git instalado, configurado, y tu llave SSH funcionando en GitHub. Ahora viene el momento que lo une todo: el **flujo básico de trabajo con Git**.

En este episodio vas a aprender a crear un repositorio en GitHub, clonarlo en tu máquina, y repetir el ciclo fundamental que vas a usar todos los días como desarrollador: **add, commit, push**.

Vamos."

---

## 📌 PASO 1 — Crear el repositorio en GitHub (0:45 – 2:00)

> *Pantalla: navegador, abriendo github.com/new*

"Primero creamos el repositorio en GitHub. Voy a ir a **github.com/new** — también puedes llegar desde el botón verde que dice 'New' en tu página principal.

Aquí voy a llenar tres cosas importantes:

- **Nombre:** `curso-gitops` — sin espacios, todo en minúsculas, con guiones si hace falta.
- **Descripción:** `Proyecto del curso de GitOps` — esto es opcional pero es buena práctica.
- **Visibilidad:** Public.

Y voy a marcar dos casillas:
- **Add a README file** — esto inicializa el repo con un commit, lo que nos permite clonarlo de inmediato.
- **Add .gitignore** → voy a seleccionar el template de `Go`, porque ese será el lenguaje del curso.

Le doy a **Create repository**."

> *Pantalla: GitHub muestra el repositorio recién creado.*

"Perfecto. Ya tenemos el repositorio con su `README.md` y su `.gitignore`. Nótese que ya existe un commit — GitHub lo hizo por nosotros al inicializarlo."

---

## 📌 PASO 2 — Clonar el repositorio (2:00 – 3:30)

> *Pantalla: botón verde 'Code' en GitHub → pestaña SSH → copiar URL.*

"Para traer este repositorio a nuestra máquina, vamos a **clonarlo**. Click en el botón verde que dice 'Code', me voy a la pestaña **SSH**, y copio la URL. Se ve algo así: `git@github.com:TU_USUARIO/curso-gitops.git`.

Ahora en la terminal:"

```bash
cd ~/UTTT/Proyectos/GitOps

git clone git@github.com:TU_USUARIO/curso-gitops.git

cd curso-gitops

ls -la
```

"Tres cosas aparecen: `.git/`, `README.md`, y `.gitignore`.

El `README` y el `.gitignore` ya los conocemos. Pero ese directorio oculto `.git/` — ese es el corazón de todo. Y antes de seguir, quiero que lo entendamos."

---

## 🔍 PAUSA CONCEPTUAL — ¿Qué hay dentro de `.git/`? (3:30 – 6:00)

> *Pantalla: terminal, explorar .git/*

"Vamos a abrir ese directorio y ver qué hay adentro."

```bash
ls -la .git/
```

"Verás algo así:"

```
.git/
├── HEAD
├── config
├── description
├── hooks/
├── info/
├── objects/
└── refs/
```

"Esto puede parecer intimidante, pero cada cosa tiene su propósito. Te lo explico rápido:

---

**`HEAD`** — Es un puntero. Le dice a Git *en qué rama estás parado ahora mismo*. Si lo abres, dice algo como `ref: refs/heads/main`. Es básicamente un post-it que dice 'estás aquí'.

---

**`config`** — La configuración local de este repositorio. Aquí vive la URL del `origin` — es decir, a dónde apunta `git push` cuando lo ejecutas. También puede tener ajustes específicos del repo que sobreescriben tu configuración global.

---

**`objects/`** — Aquí Git guarda *todo*. Cada archivo, cada commit, cada árbol de directorios. Git los convierte en objetos con un hash SHA-1 y los guarda aquí comprimidos. Cuando haces un commit, Git no guarda un 'diff' — guarda una **foto completa** del estado de tus archivos en ese momento.

---

**`refs/`** — Referencias. Aquí viven los punteros a los commits. Dentro de `refs/heads/` están tus ramas locales. En `refs/remotes/` están las referencias a las ramas remotas, como `origin/main`.

---

**`hooks/`** — Scripts opcionales que se ejecutan automáticamente en ciertos momentos del flujo de Git: antes de un commit, antes de un push, después de un merge. Por defecto están todos desactivados (terminan en `.sample`). Los exploraremos en episodios más avanzados.

---

**`COMMIT_EDITMSG`** *(aparece después del primer commit)* — Simplemente guarda el mensaje de tu último commit. Git lo usa internamente.

---

**La regla de oro: nunca modifiques nada dentro de `.git/` a mano.** Git gestiona todo eso. Si algo se rompe ahí, el repositorio puede corromperse. Déjaselo a Git."

> *Pantalla: regresar a la raíz del proyecto.*

"Dicho eso, volvamos a lo práctico."

---

## 📌 PASO 3 — Entender `git status` (6:00 – 7:00)

> *Pantalla: terminal en la raíz del proyecto.*

```bash
git status
```

"Git te responde: `nothing to commit, working tree clean`. Eso significa que tu copia local está **idéntica** a lo que está en GitHub. No hay cambios pendientes.

`git status` es el comando que más vas a usar. Úsalo constantemente. Antes de hacer add, después de hacer add, antes del commit — siempre."

---

## 📌 PASO 4 — El ciclo: add → commit → push (7:00 – 10:30)

> *Pantalla: terminal.*

"Ahora sí, el ciclo principal. Vamos a crear un archivo nuevo:"

```bash
echo "# Notas del Curso GitOps" > notas.md
echo "" >> notas.md
echo "## Módulo 1: VS Code" >> notas.md
echo "- Instalación completada" >> notas.md
```

"Veamos qué dice Git:"

```bash
git status
```

"Aparece `notas.md` en **rojo**, como 'Untracked file'. Eso significa que Git lo detecta, pero todavía no lo está siguiendo. Todavía no existe para Git.

---

### Stage: `git add`

```bash
git add notas.md
```

"Ejecuto `git status` de nuevo:"

```bash
git status
```

"Ahora aparece en **verde**, como 'new file'. Esto significa que `notas.md` está en el **staging area** — el área de preparación. Git ya lo conoce y está listo para incluirlo en el próximo commit.

Piénsalo así: el staging area es como preparar los ingredientes antes de cocinar. Todavía no hiciste el plato — solo pusiste todo en la mesa.

---

### Commit: guardar el snapshot

```bash
git commit -m "feat: agregar archivo de notas del curso"
```

"El flag `-m` es el mensaje. Este mensaje es tu firma — describe *qué cambiaste y por qué*. Fíjate en la convención que estamos usando:

- `feat:` → nueva funcionalidad
- `fix:` → corrección de bug
- `docs:` → cambios en documentación
- `chore:` → mantenimiento

Esto se llama **Conventional Commits** y es un estándar muy adoptado en la industria.

---

### Push: subir a GitHub

```bash
git push origin main
```

"`origin` es el nombre que Git le da a tu repositorio remoto por defecto — el que está en GitHub. `main` es la rama. Estamos diciendo: 'sube mis commits locales a la rama `main` del repositorio remoto'."

---

## 📌 PASO 5 — Verificar en GitHub (10:30 – 11:15)

> *Pantalla: navegador, ir al repo en GitHub.*

"Vamos a GitHub. Refresco la página... y ahí está: `notas.md` aparece en la lista. Click en él para ver el contenido — perfecto.

Ahora hago click en 'commits', arriba a la derecha. Vemos el historial: el commit inicial que hizo GitHub, y nuestro commit con el mensaje `feat: agregar archivo de notas del curso`.

Esto es fundamental: Git lleva el historial completo de todo lo que pasó."

---

## 📌 PASO 6 — Repetir el ciclo (11:15 – 12:30)

> *Pantalla: terminal.*

"El flujo siempre es el mismo. Lo repetimos:"

```bash
# Modificar el archivo
echo "## Módulo 2: Git + GitHub" >> notas.md
echo "- SSH keys configuradas" >> notas.md

# Ver qué cambió exactamente
git diff

# Stage + commit + push
git add .
git commit -m "docs: agregar notas del módulo 2"
git push origin main
```

"Nota el `git diff` antes de hacer add — esto te muestra línea por línea qué cambió. Las líneas en verde son adiciones, en rojo son eliminaciones. Es tu última oportunidad de revisar antes de hacer stage.

Y usé `git add .` en lugar de `git add notas.md` — el punto agrega **todos los archivos modificados** de una vez. Úsalo cuando todos los cambios van juntos en el mismo commit."

---

## 📌 PASO 7 — Referencia rápida de comandos (12:30 – 13:30)

> *Pantalla: tabla de comandos o slide.*

"Para cerrar, aquí está el cheatsheet del episodio:

| Comando | Para qué sirve |
|---|---|
| `git status` | Ver el estado actual — úsalo siempre |
| `git add <archivo>` | Stage un archivo específico |
| `git add .` | Stage todos los cambios |
| `git diff` | Ver cambios antes del stage |
| `git diff --staged` | Ver cambios después del stage, antes del commit |
| `git commit -m "msg"` | Crear el snapshot con mensaje |
| `git push origin main` | Subir commits a GitHub |
| `git pull origin main` | Bajar cambios de GitHub |
| `git log --oneline` | Ver el historial de commits en formato corto |

`git log --oneline` lo ejecutamos una vez más para ver todo lo que hemos hecho hoy."

```bash
git log --oneline
```

---

## 🎙️ CIERRE (13:30 – 14:00)

"Eso es todo para el episodio 5. Ahora ya conoces el flujo que vas a repetir miles de veces en tu carrera: **clonar, modificar, add, commit, push**.

En el siguiente episodio vamos a hablar de **ramas** — cómo trabajar en paralelo sin romper lo que ya funciona.

Si tienes dudas, déjalas en los comentarios. Nos vemos en el EP06."

---

## 🗒️ Notas de producción

- Mostrar `.git/` con `ls -la` y hacer zoom en la terminal para que se lea bien.
- Al hacer `git status` por primera vez (limpio) y luego con el archivo nuevo: detener la grabación si es necesario para crear el archivo, o escribirlo en vivo despacio.
- El `git diff` es visual — considera usar una fuente grande o un tema de terminal con buen contraste para rojo/verde.
- Terminar en el navegador mostrando `git log` del repo en GitHub ("commits" count) para reforzar la relación local ↔ remoto.

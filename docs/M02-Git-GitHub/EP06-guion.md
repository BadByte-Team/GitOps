# 🎬 Guión — EP06: Ramas, Pull Requests y Merge Conflicts

**Duración estimada:** 15–18 min  
**Tono:** Directo, didáctico, con pausas conceptuales antes de cada operación nueva.

---

## 🎙️ INTRO (0:00 – 0:50)

> *Pantalla: terminal con el repo `gitops-app` abierto, mostrando el `tree` del proyecto.*

"En el episodio anterior aprendiste el flujo básico: add, commit, push. Todo en una sola rama llamada `main`.

Ahora tenemos nuestra aplicación Go en `main` — handlers, modelos, autenticación, el frontend. Todo ahí. Pero le falta algo fundamental para correrla de forma reproducible: **los archivos de Docker**.

Y aquí es donde entra el concepto central de hoy. En el mundo real, nadie agrega cosas directamente a `main`. Creas una **rama separada**, propones el cambio, alguien lo revisa, y solo entonces se integra.

Hoy vas a aprender exactamente eso: **ramas, Pull Requests y cómo resolver conflictos de merge**. Esto es lo que diferencia a alguien que 'sabe Git' de alguien que trabaja en equipo con Git.

Empecemos."

---

## 🔍 PAUSA CONCEPTUAL — ¿Qué es una rama? (0:50 – 2:00)

> *Pantalla: diagrama o terminal limpia.*

"Antes de crear cualquier rama, quiero que entiendas qué es una.

Imagina el historial de commits como una línea del tiempo. Cada commit es un punto en esa línea. `main` es esa línea principal.

Cuando creas una **rama**, estás diciendo: 'quiero probar algo sin tocar esta línea'. Git crea un puntero nuevo que parte desde donde estás ahora, y todos los commits que hagas a partir de ahí van por ese camino separado.

Lo importante: `main` no se mueve. No se altera. Tus cambios viven en su propia línea hasta que decidas integrarlos, o descartarlos.

Eso es una rama. Un camino paralelo. Barato de crear, seguro de usar."

---

## 📌 PASO 1 — Crear y cambiar a una rama (2:00 – 3:30)

> *Pantalla: terminal dentro de `gitops-app`.*

"Primero, siempre verifico en qué rama estoy:"

```bash
git branch
```

"El asterisco marca la rama actual — en este caso, `main`. Bien.

Ahora creo una rama nueva y me cambio a ella con un solo comando:"

```bash
git checkout -b feature/docker
```

"`-b` significa 'crear'. Este comando equivale a hacer `git branch feature/docker` y luego `git checkout feature/docker` — todo junto.

Verifico:"

```bash
git branch
```

"El asterisco ahora está en `feature/docker`. Estoy en la rama nueva. Cualquier commit que haga a partir de aquí **no toca `main`**.

---

Una nota sobre nombres de ramas. La convención más usada es:

- `feature/nombre` → cuando agregas algo nuevo
- `fix/nombre` → cuando corriges un bug
- `hotfix/nombre` → corrección urgente, directo a producción
- `docs/nombre` → cambios solo de documentación

En nuestro caso usamos `feature/docker` porque estamos agregando capacidad nueva al proyecto: containerización. El prefijo le dice a cualquiera que lea el repo de qué tipo de cambio se trata, sin siquiera abrir el código."

---

## 📌 PASO 2 — Agregar los archivos de Docker (3:30 – 4:45)

> *Pantalla: terminal. Tenemos en paralelo el directorio `app-docker/` con sus tres archivos.*

"Tengo los archivos de Docker en una carpeta separada llamada `app-docker`. Voy a copiar los que necesito a la raíz del proyecto.

El `init.sql` no lo copio — ya existe en `gitops-app`. Sería duplicar algo que ya está y eso generaría confusión. Solo traigo lo que falta:"

```bash
cp ../app-docker/Dockerfile .
cp ../app-docker/docker-compose.yml .
```

"Verifico:"

```bash
ls -la
# Debe mostrar Dockerfile y docker-compose.yml junto al resto del proyecto

git status
# Aparecen los dos archivos en rojo como 'Untracked files'
```

"Y hago el commit. Aquí voy a ser granular — en lugar de `git add .`, listo explícitamente los dos archivos que quiero stagear:"

```bash
git add Dockerfile docker-compose.yml
git commit -m "feat: agregar Dockerfile y docker-compose para contenedores"
```

"El flujo es idéntico al EP05. Lo único que cambió es la rama en la que estamos parados."

---

## 📌 PASO 3 — Push de la rama a GitHub (4:45 – 5:30)

> *Pantalla: terminal.*

"Ahora subo esta rama a GitHub:"

```bash
git push origin feature/docker
```

"Es exactamente el mismo comando que ya conoces, pero en lugar de `main`, escribes el nombre de tu rama. GitHub recibe la rama y la guarda.

Voy al navegador."

---

## 📌 PASO 4 — Crear el Pull Request en GitHub (5:30 – 7:30)

> *Pantalla: navegador, repo en GitHub.*

"GitHub ya detectó que hice un push reciente. Aparece un banner amarillo que dice: **'feature/docker had recent pushes'** con un botón que dice **'Compare & pull request'**. Click ahí.

Se abre el formulario del Pull Request. Esto es lo que hay que entender:

Un **Pull Request** — o PR — es una propuesta formal de cambios. Le estás diciendo al equipo: 'oigan, agregué los archivos de Docker en esta rama, ¿los revisamos antes de meterlos a `main`?'.

Lleno el formulario:

- **Title:** `feat: agregar Dockerfile y docker-compose` — claro y descriptivo.
- **Description:** Aquí explico qué hice y por qué. En este caso podría decir: 'Se agregan el Dockerfile y el docker-compose.yml para levantar la app junto con la base de datos en local. El init.sql ya existía en el repo, no se duplica.' En proyectos reales, esto puede incluir el comando para probar localmente, o el ticket que lo originó. No lo dejes vacío.
- Verifico que la base sea `main` y el compare sea `feature/docker`. Eso significa: 'quiero mezclar `feature/docker` hacia `main`'.

Click en **'Create pull request'**."

> *Pantalla: página del PR ya creado.*

"El PR está abierto. Ahora cualquier miembro del equipo puede revisarlo antes de que llegue a `main`."

---

## 📌 PASO 5 — Revisar y mergear el PR (7:30 – 9:00)

> *Pantalla: pestaña 'Files changed' en el PR.*

"Click en la pestaña **'Files changed'**. Aquí está el **diff** completo — los dos archivos nuevos, `Dockerfile` y `docker-compose.yml`, línea por línea en verde.

Este es el momento de revisión. En un equipo real, tu compañero de infra haría este review: revisaría si el puerto del `docker-compose` es el correcto, si la imagen base del `Dockerfile` está actualizada, si las variables de entorno tienen sentido. Dejaría comentarios si algo no está bien.

Hoy lo hacemos solos, pero el proceso es el mismo.

Todo se ve bien. Vuelvo a la pestaña **'Conversation'** y hago click en **'Merge pull request'** → **'Confirm merge'**.

Hecho. `Dockerfile` y `docker-compose.yml` ya están en `main` en GitHub.

GitHub me ofrece eliminar la rama. La acepto: **'Delete branch'**. La rama ya cumplió su propósito."

---

## 📌 PASO 6 — Actualizar tu copia local (9:00 – 10:15)

> *Pantalla: terminal.*

"Pero ojo — el merge ocurrió en GitHub. Mi máquina local todavía no sabe que `main` se actualizó. Tengo que bajar esos cambios:"

```bash
git checkout main

git pull origin main

ls -la
# Aparecen Dockerfile y docker-compose.yml en la raíz del proyecto
```

"Ahí están. El `Dockerfile` y el `docker-compose.yml` llegaron a mi `main` local. El ciclo está completo.

Y limpio la rama local también:"

```bash
git branch -d feature/docker
```

"`-d` solo elimina la rama si ya fue mergeada. Es seguro. Si trataras de borrar una rama con cambios sin mergear, Git te lo impediría."

---

## 🔍 PAUSA CONCEPTUAL — ¿Qué es un conflicto de merge? (10:15 – 11:00)

> *Pantalla: terminal limpia o slide.*

"Hasta aquí todo fue limpio porque nadie más estaba tocando los mismos archivos. Pero qué pasa cuando dos personas modifican la misma línea del mismo archivo en ramas distintas.

Piénsalo en el contexto de este proyecto: tú estás ajustando el puerto de la base de datos en el `docker-compose.yml` en tu rama, y tu compañero también lo cambió en la suya, pero a un valor distinto. Git no puede saber cuál es el correcto. Se detiene y te dice: **'hay un conflicto, resuélvelo tú'**.

Eso es un **merge conflict**. No es un error. No es que algo salió mal. Es Git siendo honesto: 'no sé cuál versión quieres conservar, necesito que me lo digas.'

Vamos a simularlo para que lo veas en vivo."

---

## 📌 PASO 7 — Simular y resolver un conflicto (11:00 – 15:30)

### 7.1 Preparar el escenario

> *Pantalla: terminal.*

"Voy a simular el escenario real: dos personas editando la misma configuración en ramas distintas. Uso un archivo de ejemplo para no tocar el proyecto:"

```bash
git checkout main
echo "DB_PORT=5432" > config.env
git add . && git commit -m "chore: agregar config de base de datos"
```

"Ahora creo una rama que parte desde este mismo punto y edita ese mismo archivo con un valor diferente:"

```bash
git checkout -b feature/conflicto
echo "DB_PORT=5433" > config.env
git add . && git commit -m "chore: cambiar puerto de base de datos"
```

"Dos ramas. Mismo archivo. Misma línea. Contenido diferente. El conflicto está listo."

---

### 7.2 Intentar el merge

```bash
git checkout main
git merge feature/conflicto
```

"Git responde:"

```
CONFLICT (content): Merge conflict in config.env
Automatic merge failed; fix conflicts and then commit the result.
```

"No pudo resolver solo. Nos toca a nosotros."

---

### 7.3 Abrir el archivo en VS Code

```bash
code config.env
```

> *Pantalla: VS Code con el archivo abierto mostrando los marcadores de conflicto.*

"VS Code muestra esto:"

```
<<<<<<< HEAD
DB_PORT=5432
=======
DB_PORT=5433
>>>>>>> feature/conflicto
```

"Los tres bloques significan:

- **`<<<<<<< HEAD`** — lo que tiene tu rama actual (`main`): puerto 5432.
- **`=======`** — el separador.
- **`>>>>>>> feature/conflicto`** — lo que viene de la rama que estás mergeando: puerto 5433.

VS Code pone botones arriba de cada conflicto para que no tengas que editar el texto a mano:

- **'Accept Current Change'** → se queda con `main` (5432).
- **'Accept Incoming Change'** → se queda con `feature/conflicto` (5433).
- **'Accept Both Changes'** → conserva ambas líneas — en este caso no tiene sentido, no puedes tener dos puertos.
- O editas manualmente — borras los marcadores y escribes el valor correcto directamente.

En la vida real, aquí hablarías con tu compañero para saber cuál puerto es el correcto. Hoy elijo **'Accept Current Change'** — me quedo con el 5432. Guardo el archivo."

---

### 7.4 Completar el merge

> *Pantalla: terminal.*

```bash
git add config.env
git commit -m "fix: resolver conflicto de puerto en config.env"
```

"Y listo. El conflicto está resuelto. Git guardó el merge con la decisión que tomaste.

La regla práctica para evitar conflictos frecuentes: **hacer `git pull` antes de empezar a trabajar**, y mantener las ramas cortas — que vivan días, no semanas. Mientras más tiempo pasa sin mergear, más divergen las ramas y más conflictos acumulan."

---

## 📌 Resumen de comandos (15:30 – 16:15)

> *Pantalla: tabla de referencia.*

"El cheatsheet del episodio:

| Comando | Para qué sirve |
|---|---|
| `git branch` | Ver ramas locales (asterisco = actual) |
| `git branch -a` | Ver ramas locales y remotas |
| `git checkout -b nombre` | Crear rama y cambiar a ella |
| `git checkout main` | Volver a main |
| `git push origin nombre` | Subir una rama a GitHub |
| `git merge nombre` | Mergear rama a la actual |
| `git pull origin main` | Bajar cambios de GitHub |
| `git branch -d nombre` | Eliminar rama ya mergeada |"

---

## 🎙️ CIERRE (16:15 – 17:00)

"Eso es EP06.

Ya sabes trabajar con ramas, abrir Pull Requests, integrar cambios de forma controlada y resolver conflictos cuando aparecen. Nuestra app Go ahora tiene sus archivos de Docker viviendo en `main`, integrados de forma limpia y revisada.

En el siguiente episodio vamos a entrar a territorio de GitOps propiamente dicho: vamos a ver cómo Git deja de ser solo una herramienta de control de versiones para convertirse en la **fuente de verdad** de tu infraestructura.

Nos vemos en el EP07."

---

## 🗒️ Notas de producción

- Hacer `tree gitops-app` al inicio para que el espectador recuerde la estructura del proyecto antes de arrancar.
- Al mostrar `git branch` con el asterisco, hacer zoom para que se lea claro.
- Al copiar los archivos de `app-docker/`, mostrar ambos directorios lado a lado si la pantalla lo permite — refuerza visualmente que son dos carpetas distintas que se unen.
- En el PR, detente en 'Files changed' y menciona en voz alta qué revisarías tú en un equipo real (puerto, imagen base, variables de entorno).
- Para el conflicto: abrir VS Code en pantalla completa, con fuente grande. Los colores de resaltado son el gancho visual del segmento.
- Usar `config.env` en lugar de editar el `docker-compose.yml` real evita romper el proyecto durante la grabación.

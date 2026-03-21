# 🎬 Guión — EP07: Gitflow y Estructura de Repositorio DevOps

**Duración estimada:** 10–12 min  
**Tono:** Más teórico que los anteriores, pero anclado siempre en el proyecto real. Pausas visuales con el diagrama de ramas.

---

## 🎙️ INTRO (0:00 – 0:45)

> *Pantalla: terminal con `git log --oneline` del repo `gitops-app`, mostrando los commits de los episodios anteriores.*

"En los episodios anteriores aprendiste el flujo básico de Git y a trabajar con ramas. Ya tienes la app Go en `main`, con los archivos de Docker integrados via Pull Request.

Pero hasta ahora hemos estado tomando decisiones sobre la marcha: ¿cómo nombro esta rama? ¿A qué rama hago el PR? ¿Qué escribo en el mensaje del commit?

Hoy vamos a resolver eso de una vez. Vamos a hablar de **Gitflow** — la convención de ramas que usaremos durante todo el curso — y vamos a configurar el `.gitignore` para que el repositorio esté listo para Docker y Terraform.

Esto es teoría y configuración. Corto pero importante."

---

## 🔍 PAUSA CONCEPTUAL — ¿Qué es Gitflow? (0:45 – 3:00)

> *Pantalla: diagrama de Gitflow — puede ser un slide, una imagen, o dibujarlo en la terminal con texto.*

"Gitflow es una **convención**. No es una herramienta ni un plugin — es un acuerdo sobre cómo nombrar las ramas y cuál es el flujo entre ellas. Fue propuesto por Vincent Driessen en 2010 y sigue siendo el estándar en la mayoría de equipos de desarrollo.

La idea central es esta: tienes dos ramas permanentes y varias temporales.

```
main ────────────────────────────────────────────▶  producción, siempre estable
  │
  └── dev ─────────────────────────────────────▶  integración, donde se acumula el trabajo
        │
        ├── feature/login ──── PR ──▶ dev
        ├── feature/api ────── PR ──▶ dev
        │
        └── release/v1.0 ──── PR ──▶ main + dev

hotfix/fix-critico ────────── PR ──▶ main + dev
```

Veamos qué hace cada una:

**`main`** es producción. Lo que vive aquí es lo que está corriendo para los usuarios. Nadie hace commits directamente aquí — solo entran PRs aprobados, desde `dev` o desde `hotfix`.

**`dev`** es la rama de integración. Aquí llegan todos los features terminados. Es el lugar donde se combinan y se prueba que todo funcione junto antes de mandarlo a producción.

**`feature/*`** son temporales. Una por funcionalidad. Parten de `dev`, viven mientras se desarrolla esa feature, y mueren cuando el PR se mergea de vuelta a `dev`.

**`release/*`** se crean cuando `dev` está listo para una nueva versión. Se hace QA final, se corrigen últimos detalles, y se mergea a `main` y a `dev`.

**`hotfix/*`** es para emergencias. Un bug crítico en producción que no puede esperar al ciclo normal. Parte de `main`, se arregla, y se mergea de vuelta a `main` y a `dev` para que ninguno se quede desactualizado.

---

Ahora, una aclaración importante para el curso.

Gitflow es poderoso, pero tiene peso. Para un proyecto de equipo grande tiene todo el sentido. Para nuestro proyecto de curso — donde trabajamos solos y queremos que cada episodio sea claro — vamos a usar una versión simplificada:

- **`main`** → rama principal, equivale a producción.
- **`feature/*`** → una rama por episodio o por tema, PR directo a `main`.

Sin rama `dev`, sin `release`. Lo mencionamos en el curso para que conozcas el estándar completo, pero practicamos la versión liviana. Cuando trabajes en un equipo real, ya sabrás cómo escalar."

---

## 📌 PASO 1 — Convención de mensajes de commit (3:00 – 5:30)

> *Pantalla: terminal o editor de texto con ejemplos de commits.*

"Ya usamos esta convención en episodios anteriores sin explicarla formalmente. Se llama **Conventional Commits** y es un estándar adoptado por miles de proyectos open source y empresas.

La estructura es simple:

```
<tipo>: <descripción en imperativo, minúsculas>
```

Los tipos más comunes:

| Tipo | Cuándo usarlo |
|---|---|
| `feat` | Nueva funcionalidad |
| `fix` | Corrección de bug |
| `docs` | Cambios en documentación |
| `chore` | Mantenimiento, configs, dependencias |
| `ci` | Cambios en pipelines de CI/CD |
| `refactor` | Reestructurar código sin cambiar lo que hace |
| `style` | Formato, espacios — sin cambio de lógica |
| `test` | Agregar o modificar tests |

Ejemplos concretos con nuestro proyecto:

```
feat: agregar endpoint de health check
fix: corregir validación de token JWT en handlers
docs: actualizar README con instrucciones de Docker
chore: actualizar dependencias de Go
ci: agregar stage de build al Jenkinsfile
refactor: separar lógica de auth en paquete propio
```

---

La regla que más se rompe: el mensaje debe estar en **imperativo** y **minúsculas**.

No 'Agregué el endpoint', no 'Agregando el endpoint' — sino 'agregar el endpoint'. Como si le dijeras a Git qué hacer, no qué hiciste.

¿Por qué importa esto? Porque cuando usas `git log --oneline` en un proyecto con cientos de commits, o cuando GitHub genera un changelog automático, los mensajes consistentes hacen la diferencia entre un historial legible y uno inútil."

---

## 📌 PASO 2 — Crear el `.gitignore` para DevOps (5:30 – 8:30)

> *Pantalla: terminal en la raíz de `gitops-app`. Primero mostrar el `.gitignore` actual que vino del template de Go.*

"Cuando creamos el repositorio en GitHub, seleccionamos el template de `.gitignore` para Go. Eso nos dio una base para ignorar binarios y caches de Go. Bien.

Pero este proyecto va a crecer — vamos a agregar Terraform, Kubernetes, archivos de configuración con secretos. El `.gitignore` actual no cubre nada de eso.

Lo voy a reemplazar con uno completo para un proyecto DevOps:"

```bash
cat > .gitignore << 'EOF'
# ══════════════════════════════
#  .gitignore — Proyecto DevOps
# ══════════════════════════════

# ── Terraform ──
*.tfstate
*.tfstate.backup
*.tfstate.lock.info
.terraform/
.terraform.lock.hcl
*.tfvars
!example.tfvars
crash.log
override.tf
override.tf.json
*_override.tf
*_override.tf.json

# ── Docker ──
*.log
docker-compose.override.yml

# ── Kubernetes ──
kubeconfig
*.kubeconfig

# ── IDE ──
.vscode/settings.json
.vscode/*.code-workspace
.idea/
*.swp
*.swo
*~

# ── OS ──
.DS_Store
Thumbs.db
desktop.ini

# ── Secrets / Credenciales ──
*.pem
*.key
*.crt
.env
.env.*
!.env.example
aws-credentials
credentials

# ── Binarios Go ──
*.exe
*.dll
*.so
*.dylib
gitops-app

# ── Temporary ──
tmp/
temp/
*.tmp
EOF
```

"Voy a pausar en las partes más importantes.

La sección de **Terraform** ignora los `.tfstate` — esos archivos contienen el estado de tu infraestructura y pueden incluir IPs, passwords, claves de acceso. Nunca deben estar en el repo. Los `.tfvars` también, excepto el `example.tfvars` que sirve como plantilla pública.

La sección de **Secrets** es crítica. El `.env` nunca va al repo — ahí viven las variables de entorno con contraseñas y tokens. El `!.env.example` es la excepción: ese sí lo commiteamos porque es la plantilla vacía que le dice al equipo qué variables configurar.

Ahora verifico que funciona:"

```bash
# Ver qué ignora antes de commitear
git status

# Crear un archivo de prueba que debería ser ignorado
echo "DB_PASSWORD=supersecret" > .env
git status
# .env NO debe aparecer — está ignorado

# Limpiar
rm .env
```

"Perfecto. El `.env` no aparece en el status. Ahora sí commiteo:"

```bash
git add .gitignore
git commit -m "chore: configurar .gitignore para proyecto DevOps"
git push origin main
```

---

## 📌 PASO 3 — Estructura ideal del repositorio (8:30 – 10:30)

> *Pantalla: VS Code con el explorador de archivos abierto, o `tree` en terminal.*

"Antes de cerrar el episodio, quiero mostrarte hacia dónde va la estructura de este repositorio a lo largo del curso.

Hoy tenemos esto:"

```bash
tree gitops-app -L 2
```

```
gitops-app/
├── cmd/api/main.go
├── frontend/
├── internal/
├── go.mod / go.sum
├── init.sql
├── Dockerfile
├── docker-compose.yml
└── README.md
```

"A lo largo del curso, le vamos a agregar capas. La estructura objetivo es esta:"

```
gitops-app/
├── .github/              ← GitHub Actions (CI automático)
├── .gitignore
├── README.md
├── Dockerfile
├── docker-compose.yml
├── Jenkinsfile           ← Pipeline de CI/CD
├── cmd/                  ← Código Go
├── internal/
├── frontend/
├── infrastructure/       ← Todo lo de infra
│   ├── terraform/        ← Archivos .tf para provisionar infra
│   ├── kubernetes/       ← Manifiestos YAML (Deployment, Service...)
│   ├── jenkins/          ← Configuración de Jenkins
│   └── scripts/          ← Scripts de instalación y utilidad
└── docs/                 ← Documentación por módulo
```

"Cada carpeta nueva va a llegar en su episodio correspondiente, siempre en su propia rama `feature/*` con su PR. Así que el árbol de ramas va a contar la historia del proyecto commit a commit.

Eso es exactamente GitOps: **el repositorio como fuente de verdad**. No solo del código, sino de toda la infraestructura."

---

## 🎙️ CIERRE (10:30 – 11:30)

"Eso es EP07.

Ahora tienes las reglas del juego para el resto del curso: Gitflow simplificado con `feature/*` hacia `main`, Conventional Commits para mensajes consistentes, y un `.gitignore` que protege secretos y archivos de estado desde el día uno.

En el siguiente episodio arrancamos con CI/CD: vamos a crear el `Jenkinsfile` que va a construir y testear la app automáticamente cada vez que hagas un push.

Nos vemos en el EP08."

---

## 🗒️ Notas de producción

- El diagrama de Gitflow es el momento más visual del episodio — vale la pena tenerlo preparado como imagen o slide en lugar de dibujarlo en texto en la terminal. Si lo haces en terminal, usa `echo` con caracteres de box drawing y haz zoom.
- Al explicar el `.gitignore`, hacer la demo del `.env` ignorado es el momento práctico clave — asegúrate de que `git status` se lea claramente.
- La estructura del repositorio final puede mostrarse en VS Code (explorador lateral) — es más visual que el `tree` en terminal.
- Mencionar explícitamente "esto lo veremos en el EP08/09/..." al mostrar cada carpeta de `infrastructure/` — le da al espectador un mapa mental del curso.

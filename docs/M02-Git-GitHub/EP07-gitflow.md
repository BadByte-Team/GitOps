# EP 07: Gitflow y Estructura de Repositorio DevOps

**Tipo:** TEORÍA / CONFIGURACIÓN  
**Duración estimada:** 10–12 min  
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Entender la convención de ramas Gitflow que usaremos durante todo el curso, los patrones de mensajes de commit, y configurar el `.gitignore` para proyectos que usan Docker y Terraform.

---

## 📋 Prerequisitos
- Ramas y PRs dominados (EP06)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Gitflow: Estructura de Ramas

```
main ────────────────────────────────────────────▶  (PRODUCCIÓN — siempre estable)
  │
  └── dev ───────────────────────────────────────▶  (DESARROLLO — integración)
        │
        ├── feature/login ───── PR ──── merge ▶ dev
        ├── feature/api ─────── PR ──── merge ▶ dev
        │
        └── release/v1.0 ────── PR ──── merge ▶ main + dev
                                                  │
        hotfix/fix-critical ──── PR ──── merge ──▶ main + dev
```

### Ramas y Su Propósito

| Rama | Propósito | Quién la toca | Vida útil |
|---|---|---|---|
| `main` | Código en producción, siempre deployable | Nadie directamente (solo PRs) | Permanente |
| `dev` | Integración de features antes de release | Features mergeadas aquí | Permanente |
| `feature/*` | Una funcionalidad nueva específica | Un desarrollador | Temporal |
| `hotfix/*` | Arreglo urgente directo a producción | Quien detecte el bug | Temporal |
| `release/*` | Preparar nueva versión (testing final) | QA / lead | Temporal |

### Flujo en Nuestro Curso
Para simplificar, en el curso usaremos:
- **`main`** → rama principal de producción
- **`feature/*`** → para desarrollar cada episodio
- PRs directos a `main` (sin rama `dev` por ser un proyecto pequeño)

---

### Paso 2 — Convención de Mensajes de Commit

Usamos **Conventional Commits** (https://www.conventionalcommits.org):

```
<tipo>: <descripción corta>
```

| Tipo | Cuándo usarlo | Ejemplo |
|---|---|---|
| `feat` | Nueva funcionalidad | `feat: agregar endpoint de login` |
| `fix` | Corrección de bug | `fix: corregir error en autenticación JWT` |
| `docs` | Documentación | `docs: actualizar README con instrucciones` |
| `chore` | Mantenimiento, configs | `chore: actualizar dependencias de Go` |
| `ci` | Cambios en CI/CD | `ci: agregar stage de Trivy al Jenkinsfile` |
| `refactor` | Reestructurar código sin cambiar funcionalidad | `refactor: separar handlers en archivos` |
| `style` | Formato, espacios, comas | `style: formatear archivos Terraform` |
| `test` | Tests | `test: agregar test para login handler` |

> 📌 **Regla:** El mensaje debe ser en **imperativo** y **minúsculas**: "agregar", no "Agregué" ni "Agregando".

---

### Paso 3 — Crear el `.gitignore` para DevOps

**Ejecutar en:** Raíz del proyecto `curso-gitops`

**Archivo a crear/modificar:** `curso-gitops/.gitignore`

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

# ── Binarios ──
*.exe
*.dll
*.so
*.dylib
curso-gitops

# ── Temporary ──
tmp/
temp/
*.tmp
EOF

git add .gitignore
git commit -m "chore: configurar .gitignore para proyecto DevOps"
git push origin main
```

---

### Paso 4 — Estructura Ideal de un Repo DevOps

```
proyecto/
├── .github/              # GitHub Actions workflows (si usas GA)
├── .gitignore            # Archivo que acabamos de crear
├── README.md             # Documentación principal
├── Dockerfile            # Imagen de la app
├── docker-compose.yml    # Stack local
├── Jenkinsfile           # Pipeline CI/CD
├── cmd/                  # Código fuente (Go)
├── internal/             # Lógica interna
├── frontend/             # UI web
├── infrastructure/       # Todo lo de infra
│   ├── terraform/        # Archivos .tf
│   ├── kubernetes/       # Manifiestos YAML
│   ├── jenkins/          # Jenkinsfiles
│   └── scripts/          # Scripts de instalación
└── docs/                 # Documentación por módulo
```

> 📌 **Nuestro proyecto ya tiene esta estructura.** Puedes mostrarlo con `tree -L 2` o `ls -R`.

---

## 🗂️ Archivos del Proyecto Involucrados
| Archivo | Qué hacer con él |
|---|---|
| `curso-gitops/.gitignore` | Crear o actualizar con el contenido de arriba |

---

## ✅ Checklist de Verificación
- [ ] Entiendes la diferencia entre `main`, `dev` y `feature/*`
- [ ] Conoces la convención de commits (feat, fix, docs, ci, etc.)
- [ ] El `.gitignore` está configurado y committeado
- [ ] `git status` no muestra archivos que deberían estar ignorados (.tfstate, .env, etc.)

---

## 📌 Notas para el Video
- Dibujar o mostrar el diagrama de Gitflow en pantalla
- Mostrar ejemplos de mensajes de commit buenos vs malos
- Crear el `.gitignore` en vivo y mostrar qué archivos ignora
- Hacer `git status` antes y después del `.gitignore` para ver la diferencia
- Mostrar la estructura del proyecto con `tree` o VS Code

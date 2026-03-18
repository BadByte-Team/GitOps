# EP 02: Extensiones Indispensables para DevOps

**Tipo:** CONFIGURACIÓN  
**Duración estimada:** 8–12 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar y configurar las extensiones de VS Code que vas a necesitar a lo largo de todo el curso: Docker, Kubernetes, Terraform, Git, SSH remoto y soporte YAML.

---

## 📋 Prerequisitos
- VS Code instalado (EP01)

---

## 📺 Paso a Paso para el Video

### Paso 1 — Abrir el Panel de Extensiones

1. Abrir VS Code
2. Click en el ícono de **Extensiones** en la barra lateral izquierda (ícono de cuadritos)
3. O usar el atajo: `Ctrl+Shift+X`

---

### Paso 2 — Instalar las Extensiones Esenciales

Buscar cada extensión por su nombre y hacer click en **"Install"**:

#### 🐳 Docker
- **Nombre:** Docker
- **Editor:** Microsoft
- **ID:** `ms-azuretools.vscode-docker`
- **Qué hace:** Syntax highlighting para Dockerfiles, manage contenedores e imágenes desde VS Code, autocompletado en docker-compose.yml
- **Lo usarás en:** Módulo 3 (Docker + Compose)

#### ☸️ Kubernetes
- **Nombre:** Kubernetes
- **Editor:** Microsoft
- **ID:** `ms-kubernetes-tools.vscode-kubernetes-tools`
- **Qué hace:** Explorar clusters, ver pods/services/deployments, editar manifiestos con autocompletado
- **Lo usarás en:** Módulo 6 (Kubernetes), Módulo 7 (EKS)

#### 🏗️ HashiCorp Terraform
- **Nombre:** HashiCorp Terraform
- **Editor:** HashiCorp
- **ID:** `hashicorp.terraform`
- **Qué hace:** Syntax highlighting para `.tf`, autocompletado de recursos/providers, format on save
- **Lo usarás en:** Módulo 5 (Terraform), Módulo 7 (EKS)

#### 🔍 GitLens
- **Nombre:** GitLens — Git supercharged
- **Editor:** GitKraken
- **ID:** `eamodio.gitlens`
- **Qué hace:** Ver quién modificó cada línea (blame), historial de archivos, comparar ramas visualmente
- **Lo usarás en:** Módulo 2 (Git + GitHub), todo el curso

#### 🔗 Remote - SSH
- **Nombre:** Remote - SSH
- **Editor:** Microsoft
- **ID:** `ms-vscode-remote.remote-ssh`
- **Qué hace:** Conectar a servidores remotos (EC2) y editar archivos como si fueran locales
- **Lo usarás en:** Módulo 4 (AWS), Módulo 8 (Jenkins)

#### 📄 YAML
- **Nombre:** YAML
- **Editor:** Red Hat
- **ID:** `redhat.vscode-yaml`
- **Qué hace:** Validación de YAML, autocompletado para schemas de Kubernetes, Docker Compose, etc.
- **Lo usarás en:** Todo el curso (K8s manifests, docker-compose, etc.)

---

### Paso 3 — Instalación Rápida por Terminal

Si prefieres instalar todas de golpe:

**Ejecutar en:** Terminal de tu máquina local

```bash
code --install-extension ms-azuretools.vscode-docker
code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
code --install-extension hashicorp.terraform
code --install-extension eamodio.gitlens
code --install-extension ms-vscode-remote.remote-ssh
code --install-extension redhat.vscode-yaml
```

**Verificar que quedaron instaladas:**
```bash
code --list-extensions
```

**Resultado esperado:** Debe mostrar las 6 extensiones en la lista.

---

### Paso 4 — Configurar Extensiones (Opcional pero Recomendado)

#### Terraform — Format on Save
1. `Ctrl+,` → buscar `terraform`
2. Activar: **Terraform: Format On Save** → ✅

#### YAML — Schema de Kubernetes
La extensión YAML detecta automáticamente archivos de Kubernetes. Verificar abriendo cualquier `.yaml` del proyecto — debe mostrar autocompletado.

#### GitLens — Mostrar Blame
- Al abrir cualquier archivo, GitLens muestra al final de cada línea quién y cuándo la modificó
- Para desactivar si molesta: Settings → buscar `GitLens: Current Line Blame` → desactivar

---

### Paso 5 — Verificar que Todo Funciona

1. **Docker:** Abrir `curso-gitops/Dockerfile` → debe tener syntax highlighting en colores
2. **Terraform:** Abrir `infrastructure/terraform/backend/main.tf` → debe tener autocompletado
3. **YAML:** Abrir `infrastructure/kubernetes/app/deployment.yaml` → debe validar el schema
4. **GitLens:** Al lado de cada línea debe aparecer el blame en gris

---

## 🗂️ Archivos del Proyecto para Probar las Extensiones
| Archivo | Extensión que lo usa |
|---|---|
| `curso-gitops/Dockerfile` | Docker |
| `curso-gitops/docker-compose.yml` | Docker + YAML |
| `infrastructure/terraform/backend/main.tf` | Terraform |
| `infrastructure/kubernetes/app/deployment.yaml` | Kubernetes + YAML |

---

## ✅ Checklist de Verificación
- [ ] Las 6 extensiones aparecen en `code --list-extensions`
- [ ] El ícono de Docker aparece en la barra lateral
- [ ] Los archivos `.tf` tienen colores de syntax highlighting
- [ ] Los archivos `.yaml` de Kubernetes tienen autocompletado
- [ ] GitLens muestra el blame en las líneas

---

## 💡 Extensiones Extras Opcionales
| Extensión | Para qué |
|---|---|
| **Material Icon Theme** | Íconos bonitos para archivos |
| **Error Lens** | Muestra errores directamente en la línea |
| **Indent Rainbow** | Colorea la indentación (útil para YAML) |
| **Thunder Client** | Cliente REST integrado (como Postman) |

---

## 📌 Notas para el Video
- Abrir el panel de extensiones y mostrar cómo buscar
- Instalar Docker y mostrar el ícono nuevo en la barra lateral
- Abrir un Dockerfile y mostrar el syntax highlighting
- Instalar Terraform y abrir un `.tf` para mostrar el autocompletado
- Mostrar GitLens blame en un archivo con historial de Git
- Mencionar que se pueden instalar todas por terminal de una sola vez

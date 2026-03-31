 EP 34: Configuración de Credenciales — Docker Hub y GitHub Token

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Agregar en Jenkins las credenciales de Docker Hub y el Personal Access Token de GitHub con los IDs exactos que el Jenkinsfile espera. Sin estas credenciales, el pipeline no puede hacer push de imágenes ni actualizar el repositorio `gitops-infra`.

---

## 📋 Prerequisitos
- Jenkins corriendo (EP31)
- Cuenta en Docker Hub (EP11)
- Personal Access Token de GitHub con scope `repo` (EP07)

---

## 🧠 ¿Por qué no hardcodear credenciales en el Jenkinsfile?

El Jenkinsfile vive en el repositorio `gitops-app` — es código que GitHub puede ver. Si pusiéramos el token de GitHub o la contraseña de Docker Hub directamente en el archivo, cualquiera con acceso al repositorio tendría esas credenciales.

Jenkins tiene un sistema de credenciales que las guarda cifradas y las inyecta en el pipeline en tiempo de ejecución. El Jenkinsfile solo referencia el ID de la credencial — nunca el valor real.

```groovy
// En el Jenkinsfile — solo el ID, nunca el valor real
DOCKER_HUB_CREDS = credentials('dockerhub-id')
```

| Credencial | ID en Jenkins | Tipo |
|---|---|---|
| Docker Hub | `dockerhub-id` | Username with password |
| GitHub Token | `github-token-id` | Secret text |

---

## 🎬 Guión del Video

### INTRO (0:00 – 0:50)

> *Pantalla: el bloque `environment` del Jenkinsfile con `credentials('dockerhub-id')` visible.*

"Bienvenidos al episodio 34.

Esta línea del Jenkinsfile — `credentials('dockerhub-id')` — está intentando leer una credencial almacenada en Jenkins con el ID `dockerhub-id`. Si esa credencial no existe, el pipeline falla en el primer momento que intenta autenticarse con Docker Hub.

Lo mismo pasa con el token de GitHub que el pipeline usa para hacer push a `gitops-infra`.

Hoy configuramos las dos credenciales. Es uno de los episodios más cortos del módulo pero es crítico — sin este paso, ningún pipeline puede llegar a producción.

Empecemos."

---

### PASO 1 — Ir a la gestión de credenciales (0:50 – 1:30)

> *Pantalla: navegador en Jenkins.*

"Desde el dashboard de Jenkins:

**Manage Jenkins** → **Credentials** → **System** → **Global credentials (unrestricted)** → **Add Credentials**

Estas credenciales globales estarán disponibles para todos los pipelines. Para el curso es lo correcto."

---

### PASO 2 — Credencial de Docker Hub (1:30 – 4:30)

> *Pantalla: navegador en el formulario de nueva credencial.*

"Primera credencial: Docker Hub.

- **Kind:** `Username with password`
- **Scope:** Global
- **Username:** tu usuario de Docker Hub — el mismo con el que haces `docker login`
- **Password:** tu contraseña de Docker Hub o un Access Token

Una nota sobre la contraseña: Docker Hub permite generar Access Tokens en la sección de seguridad de la cuenta — en `hub.docker.com` → tu perfil → **Account Settings → Security → New Access Token**. Son preferibles a la contraseña real porque puedes revocarlos individualmente si se comprometen.

- **ID:** `dockerhub-id` — este campo es el más importante. Tiene que ser exactamente `dockerhub-id` — con guion, en minúsculas — porque es el nombre que el Jenkinsfile va a buscar. Si escriben `dockerhubId` con camelCase o `docker-hub-id` con dos guiones, el pipeline no encontrará la credencial.
- **Description:** `Docker Hub credentials` — opcional, pero ayuda a identificarla después

Click en **'Create'**.

La credencial aparece en la lista. El password no es visible — Jenkins lo guarda cifrado. Solo aparece el ID y el username."

---

### PASO 3 — Generar el Personal Access Token en GitHub (4:30 – 6:30)

> *Pantalla: navegador en GitHub → Settings.*

"Para la segunda credencial necesito primero generar el token en GitHub. Si ya lo generé en el EP07, puedo usar ese mismo — solo verifico que tenga el scope `repo`.

En GitHub: **Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)** → **Generate new token (classic)**.

Configuro:
- **Note:** `jenkins-gitops-token` — descripción para recordar para qué es
- **Expiration:** Sin expiración o la fecha que prefieran — para el curso, sin expiración es lo más cómodo
- **Scopes:** marco `repo` — acceso completo a repositorios privados

Click en **'Generate token'**.

El token aparece una sola vez — empieza con `ghp_`. Lo copio inmediatamente. Si cierro esta página sin copiarlo, tendré que generar uno nuevo."

---

### PASO 4 — Credencial del GitHub Token en Jenkins (6:30 – 9:00)

> *Pantalla: navegador en Jenkins → Credentials.*

"De vuelta en Jenkins, **Add Credentials** de nuevo.

- **Kind:** `Usermame with Password` — porque el token de GitHub es una cadena de texto, no usuario + contraseña
- **Scope:** Global
- **Username:**  Tu usuario de GitHub con acceso a los repositorios privados.
- **Password:** pego el token que acabo de copiar — `ghp_...`
- **ID:** `github-token-id` — exactamente así. El Jenkinsfile lo referencia como `credentialsId: 'github-token-id'` en el stage de Deploy to GitOps Repo
- **Description:** `GitHub Personal Access Token`

Click en **'Create'**.

---
### PASO 5 — Configuración inicial de SonarQube (contenedor)

> _Pantalla: navegador accediendo a SonarQube._

"Antes de generar el token, verifico que SonarQube esté listo.

Accedo desde el navegador a: **http://<IP_DEL_SERVIDOR>:9000**

Inicio sesión con:

- **Username:** `admin`
- **Password:** `admin`

El sistema solicita cambiar la contraseña.

Después verifico:
- **Administration** → **System** → **Status** → Estado en **UP**
    
Con esto, SonarQube queda listo para usarse con Jenkins."

---

### PASO 6 — Generar el Access Token en SonarQube

> _Pantalla: SonarQube → perfil._

"Ahora genero el token.
Voy a: **My Account** → **Security**

Configuro:
- **Name:** `jenkins-sonarqube-token`
    
Click en **Generate**

Se muestra el token (ejemplo: `sqp_XXXX...`) — lo copio inmediatamente porque solo aparece una vez.

"De vuelta en Jenkins, **Add Credentials** de nuevo.

- **Kind:** `Secret text` — porque el token de SonarQube es una cadena de texto, no usuario + contraseña
- **Scope:** Global
- **Secret:** pego el token que acabo de copiar — `ghp_...`
- **ID:** `sonarqube-server` — exactamente así.
- **Description:** `SonarQube Server Access Token`

Click en **'Create'**.

Aparece en la lista como Secret text — el valor cifrado no es visible."

---

### PASO 7 — Verificar las credenciales (9:00 – 10:00)

> *Pantalla: navegador mostrando la lista de credenciales.*

"La lista de credenciales globales ahora muestra las dos:

```
ID               Tipo                    Descripción
dockerhub-id     Username with password  Docker Hub credentials
github-token-id  Username with password  GitHub Personal Access Token
sonarqube-server Secret text             SonarQube Credential (Token)
```

Verifico que los IDs coinciden exactamente con lo que usa el Jenkinsfile:"

```groovy
// En el Jenkinsfile del proyecto
DOCKER_HUB_CREDS = credentials('dockerhub-id')     ← debe coincidir
// ...
withCredentials([string(credentialsId: 'github-token-id', ...)]) ← debe coincidir
```

"Ambos coinciden. El pipeline va a encontrar estas credenciales cuando las necesite."

---

### CIERRE (10:00 – 10:30)

"Eso es el EP34.

Dos credenciales configuradas: Docker Hub para hacer push de imágenes y GitHub para actualizar el repositorio de infraestructura. El pipeline tiene ahora todo lo que necesita para autenticarse con los servicios externos.

En el siguiente episodio escribimos el primer Jenkinsfile — la estructura base de un pipeline declarativo. Antes de llegar al pipeline completo del EP36, quiero que la anatomía de un Jenkinsfile sea familiar: agent, tools, stages, steps, post.

Nos vemos en el EP35."

---

## ✅ Checklist de Verificación
- [ ] Credencial `dockerhub-id` existe como Username with password
- [ ] Credencial `github-token-id` existe como Secret text
- [ ] Los IDs coinciden carácter por carácter con el Jenkinsfile
- [ ] El token de GitHub tiene el scope `repo` activado

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| El pipeline falla con `Invalid credentials` | El ID de la credencial no coincide con el del Jenkinsfile — verificar mayúsculas y guiones |
| `docker login` falla en el pipeline | La contraseña de Docker Hub es incorrecta — usar un Access Token en lugar de la contraseña de la cuenta |
| `remote: Invalid username or password` en el git push | El token de GitHub expiró o no tiene el scope `repo` — generar uno nuevo |
| La credencial de GitHub no puede hacer push | El token no tiene permisos de escritura — verificar que el scope `repo` (no `public_repo`) está activado |

---

## 🗒️ Notas de Producción
- Mostrar la pantalla de GitHub al generar el token — el proceso paso a paso en la interfaz real es más claro que describirlo verbalmente.
- Al copiar el token de GitHub, enfatizar que solo aparece una vez. Hacer una pausa dramática antes de navegar fuera de esa página.
- La tabla de verificación final comparando IDs en Jenkins vs IDs en el Jenkinsfile es el cierre técnico más valioso del episodio.
- Tapar o pixelar el token real en el video antes de publicar — o usar un token de ejemplo que sea obviamente falso.

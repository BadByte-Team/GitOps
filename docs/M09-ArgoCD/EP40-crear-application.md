# EP 40: Crear Application — Conectar ArgoCD al Repositorio Privado

**Tipo:** CONFIGURACIÓN / PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)
**🔄 MODIFICADO:** ArgoCD se vincula con el repositorio privado `gitops-infra` usando un Personal Access Token de GitHub.

---

## 🎯 Objetivo
Conectar ArgoCD al repositorio privado `gitops-infra` usando autenticación por token, crear la Application `curso-gitops` que observará los manifiestos de Kubernetes, y verificar que el primer sync desplega correctamente MySQL y la app Go en el cluster K3s.

---

## 📋 Prerequisitos
- ArgoCD accesible en `http://<IP_EC2>:30080` (EP39)
- Repositorio privado `gitops-infra` con todos los manifiestos en `infrastructure/kubernetes/app/` (EP07, EP47)
- Token de GitHub con scope `repo` (el mismo del EP34)
- La imagen de la app en Docker Hub con el tag actualizado por Jenkins (EP36)

---

## 🧠 Lo que vamos a desplegar

ArgoCD va a aplicar estos 7 archivos del directorio `infrastructure/kubernetes/app/` en `gitops-infra`:

| Archivo | Qué crea |
|---|---|
| `namespace.yaml` | El namespace `curso-gitops` |
| `secrets.yaml` | Credenciales de MySQL y JWT en Base64 |
| `mysql-configmap.yaml` | Script SQL de inicialización |
| `mysql-deployment.yaml` | El Pod de MySQL |
| `mysql-service.yaml` | El Service interno `mysql-svc` |
| `deployment.yaml` | El Pod de la app Go |
| `service.yaml` | El NodePort 30081 para la app |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: el dashboard de ArgoCD abierto en el navegador — sin applications todavía. Al lado, `gitops-infra` en GitHub con los 7 archivos YAML visibles.*

"Bienvenidos al episodio 40.

En la pantalla de la izquierda está ArgoCD: vacío, esperando instrucciones. En la de la derecha está `gitops-infra`: los siete manifiestos de Kubernetes que definen toda la infraestructura de la app — MySQL, la app Go, sus services, sus secrets.

En este episodio conectamos los dos. Le decimos a ArgoCD: 'observa este repositorio, este directorio, y mantén el cluster sincronizado con lo que encuentres ahí'.

Cuando terminemos este episodio, ArgoCD va a desplegar automáticamente MySQL y la app Go en el cluster K3s. Ese es el momento del que trata todo el módulo.

Empecemos."

---

### PASO 1 — Verificar los manifiestos en gitops-infra (1:00 – 2:30)

> *Pantalla: VS Code o GitHub con los archivos del directorio.*

"Antes de conectar ArgoCD, verifico que todos los archivos necesarios están en `gitops-infra/infrastructure/kubernetes/app/`.

Una cosa importante: el `deployment.yaml` debe tener el tag de imagen actualizado por Jenkins en el EP36. Si Jenkins aún no ha corrido, el tag puede ser `latest` — eso funciona, pero no tendrás el control de versiones que da el tag específico."

```bash
ls gitops-infra/infrastructure/kubernetes/app/
# namespace.yaml
# secrets.yaml
# mysql-configmap.yaml
# mysql-deployment.yaml
# mysql-service.yaml
# deployment.yaml
# service.yaml
```

"Verifico también que el `deployment.yaml` tiene la imagen correcta:"

```bash
grep "image:" gitops-infra/infrastructure/kubernetes/app/deployment.yaml
# image: TU_USUARIO/curso-gitops:1-a3b8d1c
```

"El tag versionado está ahí. ArgoCD va a usar exactamente esa imagen."

---

### PASO 2 — Conectar el repositorio privado (2:30 – 5:30)

> *Pantalla: dashboard de ArgoCD en el navegador.*

"En ArgoCD: **Settings** (el engrane en el menú lateral) → **Repositories** → **'+ Connect Repo'**.

El formulario tiene varias opciones de método de conexión. Selecciono **HTTPS** — es el más simple para repositorios de GitHub con token.

Completo los campos:

**Repository URL:** `https://github.com/TU_USUARIO_GITHUB/gitops-infra.git`
— La URL HTTPS del repositorio. No la URL SSH — el token de GitHub funciona con HTTPS.

**Username:** mi usuario de GitHub
— El nombre de usuario de la cuenta que generó el token.

**Password:** el TOKEN_GITHUB
— Aquí va el Personal Access Token que generamos en el EP34, no la contraseña de la cuenta de GitHub. El token empieza con `ghp_...`.

Click en **'Connect'**.

---

Hay dos resultados posibles:

Si aparece un ✅ verde con el mensaje `Successful` — la conexión funciona. ArgoCD puede leer el repositorio.

Si aparece un ❌ con `401 Unauthorized` — el token es incorrecto o expiró. Volver al EP34 para generar uno nuevo.

Si aparece un ❌ con `Repository not found` — la URL del repositorio tiene un typo o el repositorio es privado y el token no tiene el scope `repo`.

Veo el ✅. ArgoCD puede leer `gitops-infra`."

---

### PASO 3 — Crear la Application (5:30 – 9:30)

> *Pantalla: dashboard de ArgoCD.*

"Ahora creo la Application — el objeto de ArgoCD que define qué observar y dónde aplicarlo.

Hay dos formas. Primero muestro la interfaz web para que vean todos los campos, y luego aplico con YAML porque es más rápido y reproducible.

**Por la interfaz web:**

**Applications** → **'+ New App'**.

El formulario tiene dos secciones principales:

**GENERAL:**
- **Application Name:** `curso-gitops`
- **Project:** `default`
- **Sync Policy:** `Automatic` — ArgoCD sincroniza solo sin que yo tenga que hacer click cada vez
- Activo también **Prune Resources** y **Self Heal**

**SOURCE:**
- **Repository URL:** selecciono el repositorio que acabo de conectar
- **Revision:** `HEAD` — siempre el commit más reciente de la rama
- **Path:** `infrastructure/kubernetes/app` — el directorio con los 7 manifiestos

**DESTINATION:**
- **Cluster URL:** `https://kubernetes.default.svc` — el propio cluster donde corre ArgoCD
- **Namespace:** `curso-gitops`

Click en **'Create'**.

---

**Por YAML (alternativa reproducible):**

El archivo `gitops-infra/infrastructure/kubernetes/argocd/application.yaml` ya tiene todo esto configurado. Solo necesito aplicarlo:"

```bash
kubectl apply -f gitops-infra/infrastructure/kubernetes/argocd/application.yaml
```

```
application.argoproj.io/curso-gitops created
```

"Esta es la forma que recomiendo para el curso — si alguna vez tienes que recrear el entorno desde cero, un solo comando restaura la Application completa."

---

### PASO 4 — Observar el primer sync (9:30 – 12:00)

> *Pantalla: dashboard de ArgoCD mostrando la sincronización en progreso.*

"Después de crear la Application, ArgoCD comienza el primer sync inmediatamente. El dashboard muestra el estado: `Syncing`.

Hago click en la Application `curso-gitops` para ver el detalle. Aparece el árbol de recursos — cada manifiesto del directorio representado como un nodo:

```
curso-gitops
├── Namespace/curso-gitops        ✅ Synced
├── Secret/db-credentials         ✅ Synced
├── Secret/app-secrets            ✅ Synced
├── ConfigMap/mysql-init-config   ✅ Synced
├── Deployment/mysql              ⟳ Progressing
├── Service/mysql-svc             ✅ Synced
├── Deployment/curso-gitops       ⟳ Progressing
└── Service/curso-gitops-svc      ✅ Synced
```

"Los Deployments están en `Progressing` — Kubernetes está descargando las imágenes y creando los pods. En 1-2 minutos deberían pasar a `Healthy`.

Verifico desde la terminal también:"

```bash
kubectl get pods -n curso-gitops -w
```

```
NAME                            READY   STATUS              RESTARTS
mysql-...                       0/1     ContainerCreating   0
curso-gitops-...                0/1     ContainerCreating   0
mysql-...                       1/1     Running             0
curso-gitops-...                1/1     Running             0
```

"Ambos pods en `Running`. ArgoCD ha completado el primer sync.

De vuelta en el dashboard, el estado ahora muestra:
- **Sync Status:** `Synced` ✅
- **Health Status:** `Healthy` ✅

Ese par de checkboxes verdes es la confirmación. El estado en Git y el estado en el cluster coinciden perfectamente."

---

### PASO 5 — Verificar que la app responde (12:00 – 13:30)

> *Pantalla: navegador y terminal.*

"Verifico que la app Go es accesible en el NodePort 30081:"

```bash
IP_EC2=$(cd gitops-infra/infrastructure/terraform/jenkins-ec2 && terraform output -raw prod_public_ip)
curl -s -o /dev/null -w "%{http_code}" http://$IP_EC2:30081
# 200
```

"HTTP 200. Abro en el navegador: `http://<IP_EC2>:30081`. Aparece la pantalla de login de la plataforma GitOps del curso.

La app está corriendo en K3s, desplegada automáticamente por ArgoCD, a partir de los manifiestos en `gitops-infra`, con la imagen que Jenkins subió a Docker Hub. El loop completo de GitOps está funcionando."

---

### PASO 6 — Probar el selfHeal (13:30 – 14:30)

> *Pantalla: terminal.*

"Una demostración rápida del poder del `selfHeal: true`. Voy a borrar uno de los pods manualmente y observar que ArgoCD lo recrea:"

```bash
# Borrar el pod de la app
kubectl delete pod -l app=curso-gitops -n curso-gitops

# Ver cómo ArgoCD lo recrea automáticamente
kubectl get pods -n curso-gitops -w
```

"En menos de 30 segundos, ArgoCD detecta que el estado del cluster no coincide con el estado deseado en Git — Git dice que debe haber una réplica, y ahora hay cero — y crea un pod nuevo automáticamente.

Eso es el selfHeal. El cluster siempre converge al estado definido en Git, sin intervención manual."

---

### CIERRE (14:30 – 15:00)

"Eso es el EP40.

ArgoCD conectado a `gitops-infra`, la Application `curso-gitops` creada, y el primer despliegue exitoso en K3s. MySQL y la app Go corriendo en producción, gestionados automáticamente por ArgoCD.

En el siguiente episodio vemos el flujo completo de extremo a extremo: hacemos un cambio en el código de la app, Jenkins construye la imagen, actualiza `gitops-infra`, y ArgoCD despliega la nueva versión en K3s automáticamente. Sin intervención manual.

Nos vemos en el EP41."

---

## ✅ Checklist de Verificación
- [ ] El repositorio `gitops-infra` conectado con ✅ en Settings → Repositories
- [ ] La Application `curso-gitops` muestra `Synced` y `Healthy`
- [ ] `kubectl get pods -n curso-gitops` muestra los pods de MySQL y la app en `Running`
- [ ] La app responde en `http://<IP_EC2>:30081`
- [ ] Al borrar un pod manualmente, ArgoCD lo recrea (selfHeal)

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `401 Unauthorized` al conectar el repo | El token de GitHub expiró o no tiene scope `repo` — regenerar en EP34 |
| Pod de MySQL en `CrashLoopBackOff` | `kubectl logs mysql-POD -n curso-gitops` — probablemente el Secret de credenciales tiene valores incorrectos |
| Pod de la app en `CrashLoopBackOff` | `kubectl logs curso-gitops-POD -n curso-gitops` — verificar que `DB_HOST=mysql-svc` resuelve |
| Application en `OutOfSync` permanente | Algún manifiesto tiene un error de validación — `kubectl describe application curso-gitops -n argocd` |
| `ImagePullBackOff` | La imagen en Docker Hub no existe o el nombre en `deployment.yaml` tiene un typo |

---

## 🗒️ Notas de Producción
- La apertura con las dos pantallas — ArgoCD vacío y gitops-infra con los archivos — establece perfectamente la tensión que el episodio va a resolver.
- Al conectar el repo, mostrar el proceso completo incluyendo un intento fallido con `401` si es posible — enseña qué pasa cuando el token es incorrecto.
- El árbol de recursos del primer sync es el momento más visual del episodio — hacer zoom para que cada nodo sea legible.
- La demo del selfHeal borrando un pod manualmente es el cierre más poderoso — ver cómo el sistema "se cura solo" hace tangible la promesa de GitOps.
- Anunciar el EP41 como "el flujo completo de extremo a extremo" — es el episodio de payoff de los tres módulos de CI/CD.

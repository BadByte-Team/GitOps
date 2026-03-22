# EP 48: Pipeline CI/CD en Acción

**Tipo:** PRÁCTICA
**Duración estimada:** 18–20 min
**Dificultad:** ⭐⭐ (Intermedio)
**🔄 MODIFICADO:** Jenkins trabaja localmente, sube el nuevo tag a gitops-infra, y ArgoCD en AWS se actualiza solo — sin que Jenkins toque directamente los servidores de Amazon.

---

## 🎯 Objetivo
Ejecutar el pipeline CI/CD completo con todos los stages de seguridad activos y observar el flujo GitOps de extremo a extremo: un cambio de código en la PC local → imagen nueva en Docker Hub → despliegue automático en K3s → nueva versión visible en producción.

---

## 📋 Prerequisitos
- Todo el stack funcionando:
  - Jenkins con pipeline `curso-gitops-ci` y todos los stages activos (EP45)
  - ArgoCD observando `gitops-infra` con la Application en `Synced / Healthy` (EP40)
  - La app accesible en `http://<IP_EC2>:30081`

---

## 🧠 El run final — qué vamos a observar

```
CAMBIO EN CÓDIGO           PIPELINE CI (local)         DESPLIEGUE CD (AWS)
────────────────           ───────────────────         ───────────────────
git push                → Checkout                   
                        → SonarQube + Quality Gate   
                        → Docker Build               
                        → Trivy Scan                 
                        → Docker Push ──────────────▶ Docker Hub
                        → Deploy GitOps ────────────▶ gitops-infra
                                                     ▼
                                               ArgoCD detecta
                                               K3s rolling update
                                                     ▼
                                              Nueva versión live
```

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: cuatro ventanas abiertas simultáneamente. Arriba izquierda: VS Code con el código. Abajo izquierda: Jenkins Stage View. Arriba derecha: ArgoCD dashboard. Abajo derecha: la app en el navegador en `http://<IP_EC2>:30081`.*

"Bienvenidos al episodio 48.

Cuarenta y siete episodios para llegar aquí. Jenkins, SonarQube, Trivy, K3s, ArgoCD, Terraform, Docker, Kubernetes — cada pieza construida, entendida, verificada.

Hoy las ponemos todas a trabajar al mismo tiempo.

El plan es simple: voy a hacer un cambio en el código de la app — un cambio visible en la interfaz — voy a hacer push, y voy a observar las cuatro pantallas mientras el sistema se mueve solo. Sin tocar el servidor. Sin SSH. Sin `kubectl apply` manual. Sin ningún comando apuntando a producción.

Si todo funciona como debe, en unos minutos la app en el navegador va a mostrar el cambio. Automáticamente.

Ese momento — cuando la pantalla del navegador cambia sin que hayas tocado ningún servidor — es lo que GitOps se siente en la práctica.

Empecemos."

---

### PASO 1 — Verificar el estado inicial de todo el stack (1:30 – 4:00)

> *Pantalla: comandos de verificación en la terminal.*

"Antes de hacer cualquier cambio, capturo el estado actual del sistema para poder comparar después.

**La versión corriendo en K3s:**"

```bash
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
echo
# TU_USUARIO/curso-gitops:3-a3b8d1c   ← versión en producción ahora
```

"**ArgoCD sincronizado:**"

```bash
kubectl get application curso-gitops -n argocd \
  -o jsonpath='{.status.sync.status}'
# Synced
```

"**Los pods corriendo:**"

```bash
kubectl get pods -n curso-gitops
# NAME                         READY   STATUS    RESTARTS
# mysql-...                    1/1     Running   0
# curso-gitops-...             1/1     Running   0
```

"**Jenkins disponible:**"

```bash
docker compose ps -q jenkins | xargs docker inspect --format='{{.State.Status}}'
# running
```

"**La app en el navegador.** Abro `http://<IP_EC2>:30081`. El footer dice 'GitOps © 2026'. Ese texto es lo que va a cambiar.

Todo el stack está en verde. Hago el cambio."

---

### PASO 2 — Hacer el cambio en el código (4:00 – 5:30)

> *Pantalla: VS Code con los archivos HTML del frontend.*

"Un cambio pequeño pero visible — el footer de todos los HTML del proyecto:"

```bash
cd gitops-app

sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/index.html
sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/dashboard.html
sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/admin.html

grep "GitOps" frontend/index.html
# GitOps v2 &copy; 2026   ← cambio confirmado
```

"Commit y push:"

```bash
git add frontend/
git commit -m "feat: actualizar footer a v2 — run final del curso"
git push origin main
```

"El push está en GitHub. A partir de aquí, el sistema trabaja solo."

---

### PASO 3 — Observar el pipeline de Jenkins (5:30 – 11:00)

> *Pantalla: Jenkins Stage View — seguir la ejecución en tiempo real.*

"Abro Jenkins en `http://localhost:8080`. El pipeline `curso-gitops-ci` — si hay webhook configurado, arranca solo. Si no, click en **Build Now**.

El Stage View comienza a mostrar el progreso.

---

**Checkout** — 10 segundos. Jenkins clona `gitops-app` del commit que acabo de hacer. El `BUILD_TAG` generado es `4-b7c9e2f` — el cuarto build, más el hash del commit del footer v2.

---

**SonarQube Analysis** — entre 30 y 60 segundos. El scanner analiza el código Go de la app. Mientras corre, abro SonarQube en `http://localhost:9000`. El proyecto `curso-gitops` muestra 'Background task in progress'. En el Console Output de Jenkins busco:

```
INFO: ANALYSIS SUCCESSFUL, you can find the results at:
http://localhost:9000/dashboard?id=curso-gitops
```

El análisis llegó a SonarQube. Espero la notificación del Quality Gate.

---

**Docker Build** — 1-2 minutos. Jenkins construye `TU_USUARIO/curso-gitops:4-b7c9e2f`. La imagen incluye el footer con 'v2' que acabo de editar. El build se hace localmente a través del socket Docker — sin latencia de red, sin autenticación adicional.

---

**Trivy Scan** — 30 segundos. El reporte aparece en el Console Output:

```
Total: 0 (HIGH: 0, CRITICAL: 0)
```

Cero vulnerabilidades. El multi-stage build Alpine + Go sigue siendo limpio.

---

**Docker Push** — 30-60 segundos. La imagen sube a Docker Hub. Sube el tag versionado `4-b7c9e2f` y `latest`. Mientras corre, puedo ir a `hub.docker.com/r/TU_USUARIO/curso-gitops` y ver cómo aparece el nuevo tag en la lista.

---

**Deploy to GitOps Repo** — el stage que activa todo lo demás.

Jenkins clona `gitops-infra`, ejecuta `sed` sobre `deployment.yaml`:

```
ANTES: image: TU_USUARIO/curso-gitops:3-a3b8d1c
DESPUÉS: image: TU_USUARIO/curso-gitops:4-b7c9e2f
```

Hace commit y push. El Console Output muestra:

```
ci: deploy version 4-b7c9e2f from Jenkins
```

Ese commit está en `gitops-infra`. ArgoCD lo va a detectar.

---

**Cleanup** — limpia las imágenes locales. El pipeline termina en verde.

```
✅ Pipeline completado — imagen: TU_USUARIO/curso-gitops:4-b7c9e2f
```"

---

### PASO 4 — Observar ArgoCD sincronizar (11:00 – 14:00)

> *Pantalla: ArgoCD dashboard en el navegador.*

"Abro ArgoCD en `http://<IP_EC2>:30080`.

La Application `curso-gitops` muestra **`OutOfSync`**. ArgoCD ya detectó el commit de Jenkins en `gitops-infra`. El `deployment.yaml` tiene el tag `4-b7c9e2f` pero en el cluster todavía corre `3-a3b8d1c`. Hay una diferencia.

Con `automated: true` en la syncPolicy, ArgoCD sincroniza automáticamente. Si quiero acelerar, click en **Sync** en la UI.

El estado cambia a `Syncing`. En la vista del árbol de recursos, el Deployment `curso-gitops` muestra que está actualizando.

Observo en la terminal al mismo tiempo:"

```bash
kubectl get pods -n curso-gitops -w
```

```
NAME                          READY   STATUS              RESTARTS
curso-gitops-OLD-aaa          1/1     Running             0
curso-gitops-NEW-bbb          0/1     ContainerCreating   0   ← pod nuevo arrancando
curso-gitops-NEW-bbb          1/1     Running             0   ← pod nuevo listo
curso-gitops-OLD-aaa          1/1     Terminating         0   ← pod viejo eliminado
curso-gitops-OLD-aaa          0/1     Terminating         0
```

"El rolling update sin downtime. El pod nuevo está completamente listo antes de que el viejo se elimine. En ningún momento la app estuvo inaccesible para un usuario.

ArgoCD vuelve a `Synced / Healthy`. El cluster refleja exactamente lo que dice `gitops-infra`."

---

### PASO 5 — Confirmar la nueva versión en producción (14:00 – 16:00)

> *Pantalla: navegador con la app y verificaciones en terminal.*

"La verificación final.

**La imagen en producción:**"

```bash
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
# TU_USUARIO/curso-gitops:4-b7c9e2f   ← nueva versión
```

"**ArgoCD sincronizado:**"

```bash
kubectl get application curso-gitops -n argocd \
  -o jsonpath='{.status.sync.status}'
# Synced
```

"**La app en el navegador.** Abro `http://<IP_EC2>:30081`. Si el navegador muestra la versión anterior, limpio el caché con `Ctrl+Shift+R`.

El footer dice **'GitOps v2 © 2026'**.

---

El cambio que hice en VS Code hace unos minutos está en producción. Pasó por:
- Un análisis de calidad de código con SonarQube
- Un escaneo de vulnerabilidades con Trivy
- Un push a Docker Hub con un tag único y trazable
- Un commit en Git que documenta exactamente qué versión está en producción
- Un rolling update en Kubernetes sin downtime

Y yo no toqué el servidor en ningún momento."

---

### PASO 6 — La verificación de los tres estados (16:00 – 17:30)

> *Pantalla: tres ventanas — gitops-infra en GitHub, kubectl, ArgoCD.*

"La verificación que cierra el ciclo GitOps. Los tres estados deben coincidir:

**Estado deseado en Git:**"

```bash
grep "image:" gitops-infra/infrastructure/kubernetes/app/deployment.yaml
# image: TU_USUARIO/curso-gitops:4-b7c9e2f
```

"**Estado real en K3s:**"

```bash
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
# TU_USUARIO/curso-gitops:4-b7c9e2f
```

"**Estado según ArgoCD:**"

```bash
kubectl get application curso-gitops -n argocd \
  -o jsonpath='{.status.sync.status}'
# Synced
```

"Los tres coinciden. Git, K3s y ArgoCD están en perfecta sincronía. El sistema está en el estado exacto que el repositorio dice que debe estar. Eso es la promesa de GitOps cumplida."

---

### CIERRE (17:30 – 18:30)

"Eso es el EP48.

El pipeline completo corrió. Todos los stages de seguridad pasaron. La imagen nueva llegó a Docker Hub, el commit llegó a `gitops-infra`, ArgoCD sincronizó el cluster, y la nueva versión está visible en producción.

Sin tocar el servidor. Sin SSH. Sin `kubectl` manual.

Queda un episodio. El EP49 es la limpieza: `terraform destroy` para eliminar la EC2, verificar que la cuenta de AWS quedó en cero, y el cierre del curso.

Nos vemos en el EP49."

---

## ✅ Checklist de Verificación
- [ ] El pipeline ejecuta los 6 stages en verde con todos los stages de seguridad
- [ ] El reporte de Trivy muestra 0 vulnerabilidades HIGH/CRITICAL
- [ ] La imagen `4-b7c9e2f` aparece en Docker Hub
- [ ] El commit de Jenkins aparece en el historial de `gitops-infra`
- [ ] El rolling update ocurre sin que la app deje de responder
- [ ] Los tres estados (Git, K3s, ArgoCD) coinciden con el mismo tag
- [ ] El footer de la app muestra 'v2' en el navegador

---

## 🗒️ Notas de Producción
- La apertura con las cuatro ventanas simultáneas es el setup más ambicioso del curso — prepararlo de antemano para que no haya demoras en el momento de grabar.
- Narrar cada stage en voz alta mientras corre — no quedarse en silencio esperando. El alumno necesita escuchar qué está pasando mientras mira la pantalla.
- El momento del `kubectl get pods -w` mostrando el pod nuevo aparecer y el viejo desaparecer merece silencio breve y luego comentario — es el momento más impactante visualmente.
- La verificación de los tres estados al final es el cierre más poderoso del episodio — dejar las tres líneas en pantalla unos segundos mientras se dice "los tres coinciden".
- Anunciar el EP49 como el último episodio — el alumno sabe que está llegando al final y merece ese reconocimiento explícito.

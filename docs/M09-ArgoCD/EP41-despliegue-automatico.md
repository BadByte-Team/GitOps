# EP 41: Despliegue Automático y Sincronización

**Tipo:** PRÁCTICA
**Duración estimada:** 15–18 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Ver el flujo GitOps completo funcionando de extremo a extremo: un cambio de código en la PC local desencadena automáticamente un build en Jenkins, una imagen nueva en Docker Hub, una actualización en `gitops-infra`, y un despliegue automático en K3s sin intervención manual.

---

## 📋 Prerequisitos
- Pipeline CI funcionando (EP36)
- ArgoCD con la Application `curso-gitops` en estado `Synced` (EP40)
- La app accesible en `http://<IP_EC2>:30081`

---

## 🧠 El flujo completo de un push a producción

```
1. git push a gitops-app
        ↓
2. Jenkins ejecuta el pipeline (CI local)
   ├── SonarQube analiza el código
   ├── docker build → imagen nueva con tag 2-xyz789
   ├── docker push → Docker Hub
   └── git push → gitops-infra (deployment.yaml actualizado)
        ↓
3. ArgoCD detecta el cambio en gitops-infra
        ↓
4. K3s aplica el deployment.yaml nuevo
   ├── descarga imagen 2-xyz789 de Docker Hub
   ├── crea pod nuevo con la imagen nueva
   ├── espera a que el pod pase el readiness check
   └── elimina el pod viejo (rolling update sin downtime)
        ↓
5. La nueva versión está en producción
```

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: tres ventanas. Izquierda: el código de la app en VS Code. Centro: ArgoCD en el navegador mostrando `Synced / Healthy`. Derecha: la app corriendo en `http://<IP_EC2>:30081`.*

"Bienvenidos al episodio 41. El episodio de payoff del módulo.

Llevamos cuatro episodios construyendo este momento: instalamos ArgoCD, lo expusimos con NodePort, lo conectamos al repositorio privado, y creamos la Application que observa `gitops-infra`.

Hoy vamos a ver si todo funciona junto.

El plan es simple: voy a hacer un cambio visible en la app — cambiaré el footer para que muestre 'v2' — haré push a `gitops-app`, y luego observaremos qué pasa. Jenkins, ArgoCD, K3s — todos deberían moverse solos, sin que yo toque ningún servidor directamente.

Si el flujo funciona como debe, en unos minutos la app en producción va a mostrar el footer actualizado. Sin SSH, sin `kubectl apply` manual, sin ningún comando apuntando al servidor.

Eso es GitOps.

Empecemos."

---

### PASO 1 — Verificar el estado inicial (1:30 – 3:00)

> *Pantalla: tres terminales y el navegador.*

"Antes de hacer cualquier cambio, capturo el estado inicial para poder comparar después.

La versión actual corriendo en K3s:"

```bash
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
# TU_USUARIO/curso-gitops:1-a3b8d1c   ← versión actual
```

"El estado de ArgoCD — debe estar `Synced`:"

```bash
kubectl get application curso-gitops -n argocd \
  -o jsonpath='{.status.sync.status}'
# Synced
```

"Y abro la app en el navegador para ver el footer actual:"

```bash
open http://<IP_EC2>:30081
```

"Footer dice 'GitOps © 2026'. Eso es lo que va a cambiar."

---

### PASO 2 — Hacer el cambio en el código (3:00 – 5:00)

> *Pantalla: VS Code con los archivos HTML del frontend.*

"Hago un cambio pequeño pero visible. Edito el footer en los tres archivos HTML del frontend:"

```bash
cd gitops-app

# Cambiar el footer en los tres archivos HTML
sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/index.html
sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/dashboard.html
sed -i 's/GitOps &copy; 2026/GitOps v2 &copy; 2026/g' frontend/admin.html
```

"Verifico el cambio:"

```bash
grep "GitOps" frontend/index.html
# GitOps v2 &copy; 2026   ← 'v2' agregado
```

"Hago commit y push:"

```bash
git add frontend/
git commit -m "feat: actualizar app a v2"
git push origin main
```

"El commit está en `gitops-app`. A partir de aquí, el flujo es automático."

---

### PASO 3 — Observar Jenkins ejecutar el pipeline (5:00 – 9:00)

> *Pantalla: navegador en Jenkins — `http://localhost:8080`.*

"Abro Jenkins. Si el webhook de GitHub está configurado, el pipeline arranca automáticamente. Si no, click en **Build Now** en el pipeline `curso-gitops-ci`.

El Stage View muestra el progreso en tiempo real:

**Checkout** — Jenkins clona `gitops-app`. Genera el nuevo `BUILD_TAG`: `2-def5678` — el número 2 porque es el segundo build, más el hash del commit que acabo de hacer.

**SonarQube Analysis** — mientras corre, abro SonarQube en `http://localhost:9000`. El proyecto `curso-gitops` muestra un nuevo análisis en progreso.

**Docker Build** — construye `TU_USUARIO/curso-gitops:2-def5678`. La imagen incluye el cambio del footer que acabo de hacer.

**Docker Push** — sube la imagen a Docker Hub. En `hub.docker.com/r/TU_USUARIO/curso-gitops`, el tag `2-def5678` va a aparecer en unos momentos.

**Deploy to GitOps Repo** — este es el stage que activa todo lo demás.

Jenkins clona `gitops-infra`, ejecuta:
```bash
sed -i "s|image: TU_USUARIO/curso-gitops:.*|image: TU_USUARIO/curso-gitops:2-def5678|" \
    infrastructure/kubernetes/app/deployment.yaml
git commit -m "ci: deploy version 2-def5678 from Jenkins"
git push origin main
```

Ese push es lo que ArgoCD va a detectar.

**Cleanup** — elimina las imágenes locales.

El pipeline terminó en verde. Todo en menos de 3 minutos."

---

### PASO 4 — Observar ArgoCD sincronizar (9:00 – 12:00)

> *Pantalla: dashboard de ArgoCD.*

"Voy al dashboard de ArgoCD. La Application `curso-gitops` muestra **`OutOfSync`** — ArgoCD ya detectó que el `deployment.yaml` en `gitops-infra` tiene un tag diferente al que está corriendo en el cluster.

Con `automated sync` activo, ArgoCD sincroniza automáticamente en segundos. El estado cambia de `OutOfSync` a `Syncing`.

Si quiero acelerar el proceso en lugar de esperar el ciclo de 3 minutos, puedo forzar el sync manualmente:"

```bash
# Desde la terminal, sync manual
kubectl argo app sync curso-gitops -n argocd 2>/dev/null || \
  echo "Usar el botón Sync en la UI de ArgoCD"
```

"O simplemente hago click en el botón **'Sync'** en la UI.

ArgoCD aplica el `deployment.yaml` con el nuevo tag. Kubernetes inicia el rolling update:

1. Crea un pod nuevo con `TU_USUARIO/curso-gitops:2-def5678`
2. Descarga la imagen de Docker Hub
3. Arranca el contenedor
4. Espera a que pase el readiness check (`GET /` → HTTP 200)
5. Elimina el pod viejo con `TU_USUARIO/curso-gitops:1-a3b8d1c`

Observo en la terminal:"

```bash
kubectl get pods -n curso-gitops -w
```

```
NAME                         READY   STATUS              RESTARTS
curso-gitops-OLD-xxx         1/1     Running             0
curso-gitops-NEW-yyy         0/1     ContainerCreating   0   ← pod nuevo
curso-gitops-NEW-yyy         1/1     Running             0   ← pod nuevo listo
curso-gitops-OLD-xxx         1/1     Terminating         0   ← pod viejo eliminado
curso-gitops-OLD-xxx         0/1     Terminating         0
```

"El rolling update sin downtime. El pod nuevo está listo antes de que el viejo se elimine. En ningún momento la app estuvo inaccesible."

---

### PASO 5 — Confirmar la nueva versión en producción (12:00 – 14:00)

> *Pantalla: navegador, terminal.*

"Verifico que el pod nuevo tiene la imagen correcta:"

```bash
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
# TU_USUARIO/curso-gitops:2-def5678   ← nueva versión
```

"Abro la app en el navegador: `http://<IP_EC2>:30081`.

El footer ahora dice **'GitOps v2 © 2026'**. El cambio que hice en el código hace unos minutos está en producción.

---

Hago la verificación final del ciclo completo:"

```bash
# ¿Qué hay en gitops-infra?
grep "image:" gitops-infra/infrastructure/kubernetes/app/deployment.yaml
# image: TU_USUARIO/curso-gitops:2-def5678

# ¿Qué corre en K3s?
kubectl get deployment curso-gitops -n curso-gitops \
  -o jsonpath='{.spec.template.spec.containers[0].image}'
# TU_USUARIO/curso-gitops:2-def5678

# ¿ArgoCD está sincronizado?
kubectl get application curso-gitops -n argocd \
  -o jsonpath='{.status.sync.status}'
# Synced
```

"Los tres coinciden. Git, K3s y ArgoCD están en perfecta sincronía. Eso es exactamente lo que el patrón GitOps garantiza."

---

### PASO 6 — Probar el rollback (14:00 – 16:00)

> *Pantalla: GitHub → gitops-infra → historial de commits.*

"Para cerrar el módulo, demuestro cómo funciona el rollback en GitOps.

Imagina que la versión v2 tiene un bug crítico y necesitas volver a la versión anterior inmediatamente. En GitOps, el rollback es simplemente revertir el commit en `gitops-infra`:"

```bash
cd gitops-infra

# Ver el historial de commits
git log --oneline -3
# abc1234  ci: deploy version 2-def5678 from Jenkins   ← último commit (v2)
# def5678  ci: deploy version 1-a3b8d1c from Jenkins   ← commit anterior (v1)

# Revertir el último commit
git revert HEAD --no-edit

# Push — ArgoCD detectará el cambio y volverá a v1
git push origin main
```

"ArgoCD detecta el nuevo commit en `gitops-infra` y aplica la versión anterior del `deployment.yaml`. K3s hace un rolling update de vuelta a `1-a3b8d1c`.

El footer vuelve a mostrar 'GitOps © 2026'.

Ese rollback fue:
- Sin conectarse al servidor
- Sin ejecutar `kubectl` manualmente
- Con trazabilidad completa en el historial de Git — hay un commit que dice 'revert' con fecha y hora

En producción, esa trazabilidad puede ser la diferencia entre resolver un incidente en minutos o en horas."

---

### CIERRE (16:00 – 17:00)

"Eso es el episodio 41. Y con esto cerramos el Módulo 09.

El flujo GitOps completo está funcionando. Un commit en `gitops-app` se convierte en una imagen nueva en Docker Hub, en un commit en `gitops-infra`, en un despliegue automático en K3s. Sin intervención manual, sin SSH al servidor, sin `kubectl apply` a mano.

¿Qué falta? Los dos módulos que completan el curso:

El **Módulo 10** — Seguridad. Integraremos Trivy para escanear vulnerabilidades en las imágenes y SonarQube formalmente en el pipeline. El mismo pipeline que vimos hoy, reforzado.

El **Módulo 11** — El proyecto final. La arquitectura completa explicada, la base de datos separada en Kubernetes, el pipeline en acción con todo integrado, y la limpieza de recursos en el EP49.

Nos vemos en el EP42."

---

## ✅ Checklist de Verificación
- [ ] El pipeline de Jenkins se ejecuta completo en verde con el nuevo commit
- [ ] La imagen `2-xyz789` aparece en Docker Hub
- [ ] El `deployment.yaml` en `gitops-infra` tiene el nuevo tag
- [ ] ArgoCD muestra `Synced / Healthy` después del despliegue
- [ ] La app muestra el cambio visual en `http://<IP_EC2>:30081`
- [ ] El rollback con `git revert` + `git push` revierte la versión en producción

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Jenkins no detecta el push automáticamente | Usar **Build Now** manual o configurar el webhook de GitHub |
| ArgoCD no detecta el cambio de `gitops-infra` | El sync automático tiene un delay de hasta 3 min — forzar con el botón **Sync** |
| Pod nuevo en `ImagePullBackOff` | La imagen no se subió a Docker Hub — verificar el stage Docker Push en Jenkins |
| Rolling update no termina (`Pending` durante mucho tiempo) | La EC2 no tiene suficiente RAM — verificar `kubectl top nodes` y el Swap |
| La app muestra la versión anterior en el navegador | Limpiar caché del navegador (`Ctrl+Shift+R` o Cmd+Shift+R en Mac) |

---

## 🗒️ Notas de Producción
- La apertura con las tres ventanas simultáneas (código, ArgoCD, app en el navegador) establece el escenario perfectamente — todo en pantalla antes de tocar una sola línea de código.
- Hablar en voz alta mientras el pipeline corre — no quedarse en silencio esperando. Explicar qué hace cada stage mientras ocurre, como un narrador deportivo.
- El momento del rolling update con `kubectl get pods -w` mostrando el pod nuevo aparecer y el viejo desaparecer es el más satisfactorio del episodio — dejar la terminal en pantalla grande y silencio breve mientras ocurre.
- El rollback es el cierre más poderoso del módulo — demuestra que GitOps no es solo para deployments sino también para recovery. Enfatizar la trazabilidad del `git revert` como ventaja sobre el rollback manual con kubectl.
- El cierre con el mapa de los módulos restantes (10 y 11) da al alumno perspectiva de cuánto queda y motiva a continuar.

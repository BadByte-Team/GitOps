# EP 38: Instalar ArgoCD en el Cluster K3s

**Tipo:** INSTALACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar ArgoCD en el cluster K3s aplicando el manifiesto oficial, verificar que todos sus pods están en estado `Running`, y entender qué hace cada componente del sistema.

---

## 📋 Prerequisitos
- K3s corriendo en la EC2 con kubectl configurado localmente (EP29, EP30)
- Swap de 2 GB activo en la EC2 (EP28) — ArgoCD ocupa ~400 MB de RAM

---

## 🧠 Los pods de ArgoCD

ArgoCD se instala como un conjunto de pods dentro del namespace `argocd`. Cada uno tiene un rol específico:

| Pod | Función |
|---|---|
| `argocd-server` | API server y UI web — el que exponemos en el EP39 |
| `argocd-repo-server` | Clona los repositorios y genera los manifiestos |
| `argocd-application-controller` | El cerebro — compara Git vs cluster y reconcilia |
| `argocd-dex-server` | Autenticación SSO (no lo usamos directamente) |
| `argocd-redis` | Cache interno del application controller |
| `argocd-notifications-controller` | Notificaciones (Slack, email, etc.) |

---

## 🎬 Guión del Video

### INTRO (0:00 – 0:45)

> *Pantalla: terminal local con kubectl apuntando al cluster K3s — `kubectl get nodes` mostrando el nodo Ready.*

"Bienvenidos al episodio 38.

En el episodio anterior entendimos qué es GitOps y qué hace ArgoCD conceptualmente. Hoy lo instalamos.

El proceso es directo: un namespace, un `kubectl apply` con el manifiesto oficial, y esperar a que los pods levanten. En menos de cinco minutos ArgoCD está corriendo en el cluster K3s de la EC2.

Empecemos."

---

### PASO 1 — Verificar el estado del cluster (0:45 – 2:00)

> *Pantalla: terminal local.*

"Primero verifico que el cluster está en buen estado y que el kubectl local apunta a K3s:"

```bash
kubectl config current-context
# default  ← el contexto de K3s

kubectl get nodes
# NAME        STATUS   ROLES                  AGE
# ip-10-...   Ready    control-plane,master   Xd
```

"Y verifico la memoria disponible antes de instalar — ArgoCD va a necesitar ~400 MB más:"

```bash
kubectl top nodes
# NAME        CPU(cores)   MEMORY(bytes)   MEMORY%
# ip-10-...   150m         450Mi           46%
```

"450 MB usados de 981 MB de RAM física — el resto lo cubre el Swap del EP28. Tenemos margen suficiente."

---

### PASO 2 — Crear el namespace de ArgoCD (2:00 – 3:00)

> *Pantalla: terminal.*

"ArgoCD vive en su propio namespace, separado de los recursos de la aplicación:"

```bash
kubectl create namespace argocd
```

```
namespace/argocd created
```

"Verifico:"

```bash
kubectl get namespaces | grep argocd
# argocd   Active   5s
```

---

### PASO 3 — Instalar ArgoCD con el manifiesto oficial (3:00 – 6:00)

> *Pantalla: terminal.*

"La instalación se hace aplicando el manifiesto oficial de ArgoCD directamente desde su repositorio de GitHub. Este manifiesto contiene todos los recursos que ArgoCD necesita: Deployments, Services, ConfigMaps, RBAC, CRDs."

```bash
kubectl apply -n argocd -f \
  https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

"El output es largo — ArgoCD crea decenas de recursos. Roles, ServiceAccounts, Deployments, Services, ConfigMaps, CustomResourceDefinitions. Las CRDs son especialmente importantes — son las que le enseñan a Kubernetes el nuevo tipo de objeto `Application` que usaremos en el EP40.

Al terminar, el output dice:"

```
customresourcedefinition.apiextensions.k8s.io/applications.argoproj.io created
customresourcedefinition.apiextensions.k8s.io/appprojects.argoproj.io created
serviceaccount/argocd-application-controller created
...
deployment.apps/argocd-server created
service/argocd-server created
```

---

### PASO 4 — Esperar a que los pods estén listos (6:00 – 8:30)

> *Pantalla: terminal con el watch de pods.*

"Los pods de ArgoCD pueden tardar entre 1 y 3 minutos en pasar de `ContainerCreating` a `Running` — necesitan descargar las imágenes si no están en caché y arrancar sus procesos internos.

Monitoreo en tiempo real:"

```bash
kubectl get pods -n argocd -w
```

"La flag `-w` mantiene el output actualizado automáticamente. Veo cómo los pods van cambiando de estado:

```
NAME                                                READY   STATUS              RESTARTS
argocd-redis-...                                    0/1     ContainerCreating   0
argocd-dex-server-...                               0/1     ContainerCreating   0
argocd-repo-server-...                              0/1     Init:0/1            0
argocd-application-controller-...                  0/1     ContainerCreating   0
argocd-server-...                                   0/1     ContainerCreating   0
argocd-redis-...                                    1/1     Running             0
argocd-dex-server-...                              1/1     Running             0
argocd-repo-server-...                             1/1     Running             0
argocd-application-controller-...                  1/1     Running             0
argocd-server-...                                  1/1     Running             0
```

"Cuando todos los pods muestran `1/1 Running`, ArgoCD está listo. Presiono `Ctrl+C` para salir del watch.

---

Un momento importante: si algún pod muestra `OOMKilled` en los RESTARTS — Out of Memory Killed — significa que el Swap no está activo o es insuficiente. Si ocurre, verificar con `ssh -i aws-key.pem ubuntu@IP free -h` que el Swap de 2 GB está activo."

---

### PASO 5 — Verificar el estado final (8:30 – 10:00)

> *Pantalla: terminal.*

"Verificación completa:"

```bash
kubectl get pods -n argocd
```

```
NAME                                               READY   STATUS    RESTARTS   AGE
argocd-application-controller-0                    1/1     Running   0          2m
argocd-dex-server-...                              1/1     Running   0          2m
argocd-redis-...                                   1/1     Running   0          2m
argocd-repo-server-...                             1/1     Running   0          2m
argocd-server-...                                  1/1     Running   0          2m
argocd-notifications-controller-...               1/1     Running   0          2m
```

"Todos en `Running`. Verifico también los servicios creados:"

```bash
kubectl get svc -n argocd
```

```
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
argocd-dex-server       ClusterIP   10.43.x.x       <none>        5556/TCP,5557/TCP   2m
argocd-metrics          ClusterIP   10.43.x.x       <none>        8082/TCP            2m
argocd-redis            ClusterIP   10.43.x.x       <none>        6379/TCP            2m
argocd-repo-server      ClusterIP   10.43.x.x       <none>        8081/TCP,8084/TCP   2m
argocd-server           ClusterIP   10.43.x.x       <none>        80/TCP,443/TCP      2m
argocd-server-metrics   ClusterIP   10.43.x.x       <none>        8083/TCP            2m
```

"El `argocd-server` es el Service que vamos a modificar en el EP39 para exponerlo. Actualmente es `ClusterIP` — solo accesible desde dentro del cluster. En el próximo episodio lo convertimos en `NodePort` para acceder desde el navegador."

---

### CIERRE (10:00 – 10:30)

"Eso es el EP38.

ArgoCD instalado en el cluster K3s. Seis pods corriendo en el namespace `argocd`, cada uno con su rol específico en el sistema de reconciliación.

En el siguiente episodio exponemos la interfaz web de ArgoCD usando NodePort en el puerto 30080 — sin LoadBalancer, sin costo adicional. Ese puerto ya está abierto en el Security Group de la EC2 desde el EP22.

Nos vemos en el EP39."

---

## ✅ Checklist de Verificación
- [ ] `kubectl get pods -n argocd` muestra todos los pods en `Running`
- [ ] `kubectl get svc -n argocd` muestra el `argocd-server` como `ClusterIP`
- [ ] No hay pods con `OOMKilled` en los RESTARTS
- [ ] Entiendes el rol de cada pod de ArgoCD

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Pod en `OOMKilled` | El Swap no está activo — `ssh ubuntu@IP "free -h"` para verificar |
| Pod en `Pending` por mucho tiempo | Verificar recursos: `kubectl describe pod NOMBRE -n argocd` → sección Events |
| `ImagePullBackOff` | La EC2 no tiene acceso a internet — verificar el Security Group (egress rule) |
| `kubectl apply` falla con `connection refused` | kubectl no apunta al cluster K3s — `kubectl config current-context` |

---

## 🗒️ Notas de Producción
- El `kubectl get pods -n argocd -w` es uno de los momentos más visuales del módulo — ver los pods cambiar de estado en tiempo real. Mantener la terminal en pantalla completa con fuente grande.
- Mencionar cada pod mientras aparece en Running y recordar brevemente su función — conecta la teoría del EP37 con lo que se ve en pantalla.
- Si algún pod tarda más de 3 minutos, explicar que es normal en una t3.micro — la descarga de imágenes puede ser lenta la primera vez.
- El `kubectl get svc -n argocd` al final anticipa el EP39 de forma natural — "el `argocd-server` está en ClusterIP, en el próximo episodio lo exponemos".

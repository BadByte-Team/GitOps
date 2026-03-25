# EP 26: Comandos Esenciales de kubectl

**Tipo:** PRÁCTICA
**Duración estimada:** 15–18 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Dominar los comandos de operación diaria de kubectl — `logs`, `exec`, `describe`, `scale`, `rollout` — con el cluster de Minikube. Al terminar este episodio, cuando haya un problema en producción con K3s, sabrás exactamente qué comando ejecutar para diagnosticarlo.

---

## 📋 Prerequisitos
- Minikube corriendo (EP24)
- El Deployment y el Service del EP25 aplicados, o levantarlos de nuevo al inicio de este episodio

---

## 🧠 La diferencia entre `get`, `describe`, `logs` y `exec`

Es la distinción que más confunde al principio:

| Comando | Qué muestra | Cuándo usarlo |
|---|---|---|
| `kubectl get` | Lista de recursos con estado básico | Para ver qué existe y si está `Running` |
| `kubectl describe` | Todo el detalle de un recurso — eventos, condiciones, configuración | Para entender por qué algo no funciona |
| `kubectl logs` | La salida estándar del contenedor | Para ver errores de la aplicación |
| `kubectl exec` | Abre una shell dentro del contenedor | Para explorar o depurar desde dentro |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: un Pod en estado `CrashLoopBackOff` — puede ser simulado o una captura de pantalla.*

"Bienvenidos al episodio 26.

Esto es lo que ven cuando algo sale mal en Kubernetes. Un Pod en estado `CrashLoopBackOff`. El contenedor arranca, falla, Kubernetes intenta reiniciarlo, falla de nuevo, espera un poco más, vuelve a intentar. Un ciclo que se repite hasta que alguien lo soluciona.

¿Qué haces cuando ves esto? ¿Por dónde empiezas?

Si no conoces los comandos correctos, Kubernetes parece una caja negra que simplemente 'no funciona'. Si los conoces, en 60 segundos puedes saber exactamente qué está pasando: qué error tiene la aplicación, qué recursos le faltan, qué evento disparó el problema.

Ese es el propósito de este episodio. Vamos a cubrir los comandos de operación que usarás constantemente en los módulos de GitOps — especialmente cuando el pipeline despliegue una nueva versión y algo no arranque como esperabas.

Levanto el entorno del EP25 si no lo tienen activo:"

```bash
kubectl apply -f namespace.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

"Empecemos."

---

### `kubectl get` — el comando de estado (1:00 – 3:30)

> *Pantalla: terminal.*

"El comando que más vas a ejecutar. `get` lista los recursos y muestra su estado básico."

```bash
# Pods en el namespace practica-k8s
kubectl get pods -n practica-k8s
```

```
NAME                              READY   STATUS    RESTARTS   AGE
nginx-practica-6d8f7c5b9-xk8j2   1/1     Running   0          2m
```

"Cuatro columnas que siempre leer.

**READY** — cuántos contenedores dentro del Pod están listos sobre cuántos tiene en total. `1/1` es lo que quieres ver. `0/1` significa que el contenedor no pasó el health check.

**STATUS** — el estado del Pod. `Running` es bueno. `Pending` significa que Kubernetes todavía no pudo asignarlo a un nodo. `CrashLoopBackOff` significa que se está cayendo repetidamente. `ImagePullBackOff` significa que no pudo descargar la imagen.

**RESTARTS** — cuántas veces se reinició el contenedor. Un número creciente es señal de problema. Cero es lo que quieres.

**AGE** — cuánto tiempo tiene el Pod. Un Pod muy joven en estado `Running` después de un deploy reciente es buena señal.

---

Algunos flags útiles de `get`:"

```bash
# Con más información — muestra en qué nodo está y la IP del Pod
kubectl get pods -n practica-k8s -o wide

# En tiempo real — actualiza automáticamente cuando algo cambia
kubectl get pods -n practica-k8s -w

# Todo en el namespace — Pods, Services, Deployments
kubectl get all -n practica-k8s

# En formato YAML — ver la configuración completa del objeto
kubectl get deployment nginx-practica -n practica-k8s -o yaml
```

"El `-o yaml` es especialmente útil cuando quieres ver la diferencia entre lo que describiste en tu YAML y lo que Kubernetes realmente aplicó — Kubernetes agrega campos adicionales como el `uid`, la fecha de creación, el estado actual."

---

### `kubectl describe` — el diagnóstico profundo (3:30 – 6:00)

> *Pantalla: terminal, ejecutando `describe` y mostrando el output completo.*

"Cuando `get` muestra algo que no está bien — un Pod en `Pending`, un `0/1` en READY — el siguiente paso es `describe`.

`describe` muestra todo el detalle de un objeto: su configuración completa, las condiciones actuales, y lo más importante: los **eventos**."

```bash
kubectl describe pod nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s
```

"El output es largo. Pero hay una sección específica que hay que buscar primero: la sección `Events` al final."

```
Events:
  Type    Reason     Age   From               Message
  ----    ------     ----  ----               -------
  Normal  Scheduled  2m    default-scheduler  Successfully assigned practica-k8s/nginx-practica to minikube
  Normal  Pulled     2m    kubelet            Container image "nginx:alpine" already present on machine
  Normal  Created    2m    kubelet            Created container nginx
  Normal  Started    2m    kubelet            Started container nginx
```

"Cuando un Pod no arranca, los eventos te dicen exactamente por qué. Los mensajes más comunes que verás:

`Failed to pull image "...": not found` — la imagen no existe. Revisa el nombre y el tag en el YAML.

`Insufficient memory` — el nodo no tiene suficiente RAM. En nuestro t3.micro esto puede ocurrir si hay demasiados Pods corriendo sin límites de memoria.

`Back-off restarting failed container` — el contenedor se cayó varias veces. Mira los logs para ver qué error tiene la aplicación.

`0/1 nodes are available: 1 Insufficient cpu` — el Pod pide más CPU de la que hay disponible.

Cada uno de esos mensajes te dice exactamente qué ajustar."

```bash
# También funciona con Deployments y Services
kubectl describe deployment nginx-practica -n practica-k8s
kubectl describe service nginx-svc -n practica-k8s
```

---

### `kubectl logs` — ver lo que dice la aplicación (6:00 – 8:30)

> *Pantalla: terminal.*

"`describe` te dice qué pasó a nivel de Kubernetes. `logs` te dice qué pasó a nivel de tu aplicación."

```bash
# Ver los logs de un Pod
kubectl logs nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s
```

"Para nginx verás las líneas de acceso — cada petición HTTP que llegó al servidor.

---

El flag más importante de `logs` es `-f` — follow, igual que `tail -f`:"

```bash
kubectl logs -f nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s
```

"La terminal queda en vivo mostrando los logs en tiempo real. Mientras accedes a la app en el navegador, ves las líneas aparecer aquí. Para salir, `Ctrl+C`.

---

Para nuestro proyecto — la app Go y MySQL — los logs son el primer lugar donde buscar errores. Por ejemplo, si la app no puede conectarse a la base de datos:"

```
Error conectando a la base de datos: dial tcp mysql-svc:3306: connection refused
```

"Eso te dice inmediatamente que el Pod de MySQL no está listo, o que el Service `mysql-svc` no resuelve correctamente.

---

Dos flags adicionales útiles:"

```bash
# Ver los últimos 50 líneas
kubectl logs nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s --tail=50

# Ver logs del Pod anterior (si se cayó y fue recreado)
kubectl logs nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s --previous
```

"El `--previous` es muy útil cuando tienes un `CrashLoopBackOff`. El Pod se cae, se recrea, y el Pod nuevo no tiene los logs del Pod anterior. Con `--previous`, los recuperas."

---

### `kubectl exec` — entrar al contenedor (8:30 – 11:00)

> *Pantalla: terminal.*

"El comando que te da acceso completo al interior del contenedor — como hacer SSH a un servidor, pero dentro de un Pod.

Esto es muy poderoso para diagnóstico: puedes verificar qué archivos existen, qué variables de entorno están seteadas, si la red funciona, si el proceso está corriendo."

```bash
kubectl exec -it nginx-practica-6d8f7c5b9-xk8j2 -n practica-k8s -- sh
```

"El `-it` combina `-i` (interactive, mantiene STDIN abierto) con `-t` (TTY, asigna una terminal). El `--` separa los flags de kubectl del comando que se ejecuta dentro del contenedor. `sh` en lugar de `bash` porque nginx usa una imagen Alpine que no tiene bash.

Estoy dentro del contenedor:"

```bash
# Ver las variables de entorno
env | grep DB

# Ver si hay conectividad de red
wget -q -O- http://localhost:80

# Ver los archivos de configuración
cat /etc/nginx/nginx.conf

# Salir
exit
```

"En el contexto de nuestro proyecto, los casos de uso más comunes del `exec` son:

Verificar que las variables de entorno de la base de datos llegaron correctamente al contenedor de la app:"

```bash
kubectl exec -it POD_APP -n curso-gitops -- sh
env | grep DB_
# DB_HOST=mysql-svc
# DB_USER=curso_app
# DB_NAME=curso_db
```

"Verificar conectividad con MySQL desde dentro del Pod de la app:"

```bash
kubectl exec -it POD_APP -n curso-gitops -- sh
wget -q -O- http://mysql-svc:3306
# Si responde, el Service resuelve correctamente
exit
```

"Esas dos verificaciones resuelven el 80% de los problemas de conectividad que pueden aparecer en el EP40 cuando ArgoCD despliegue el stack completo."

---

### `kubectl scale` — cambiar el número de réplicas (11:00 – 12:30)

> *Pantalla: terminal.*

"Escalar un Deployment — cambiar el número de réplicas — se puede hacer directamente con `kubectl scale` sin editar el YAML:"

```bash
# Escalar a 3 réplicas
kubectl scale deployment nginx-practica -n practica-k8s --replicas=3

# Verificar
kubectl get pods -n practica-k8s
```

```
NAME                              READY   STATUS    RESTARTS   AGE
nginx-practica-6d8f7c5b9-xk8j2   1/1     Running   0          10m
nginx-practica-6d8f7c5b9-ab3k4   1/1     Running   0          5s
nginx-practica-6d8f7c5b9-mn7p1   1/1     Running   0          5s
```

"Tres Pods. Kubernetes creó los dos adicionales en segundos.

---

Una advertencia importante: si usas GitOps, **no escales con `kubectl scale`**. El estado deseado vive en Git. Si escalas a 3 con kubectl pero el `deployment.yaml` en el repositorio dice `replicas: 1`, ArgoCD va a detectar la discrepancia y revertirlo a 1 en la próxima sincronización.

La forma correcta en GitOps es editar el `deployment.yaml`, hacer commit, hacer push, y dejar que ArgoCD aplique el cambio. Así el estado en Git y el estado en el cluster siempre coinciden.

Por ahora, en este cluster de práctica donde no hay ArgoCD, escalar con kubectl está bien. Pero tengan ese principio en mente para cuando lleguemos a los módulos de GitOps."

```bash
# Volver a 1 réplica
kubectl scale deployment nginx-practica -n practica-k8s --replicas=1
```

---

### `kubectl rollout` — gestionar actualizaciones (12:30 – 14:00)

> *Pantalla: terminal.*

"El último conjunto de comandos: `rollout`. Sirve para gestionar las actualizaciones de Deployments.

**Ver el estado de un rollout** — útil cuando acabas de aplicar una nueva versión y quieres saber si terminó:"

```bash
kubectl rollout status deployment/nginx-practica -n practica-k8s
# deployment "nginx-practica" successfully rolled out
```

"**Ver el historial de versiones:**"

```bash
kubectl rollout history deployment/nginx-practica -n practica-k8s
```

"**Hacer rollback a la versión anterior** — el comando de emergencia cuando una nueva versión tiene un bug crítico:"

```bash
kubectl rollout undo deployment/nginx-practica -n practica-k8s
```

"Kubernetes revierte el Deployment a la versión anterior inmediatamente. Sin downtime, sin editar YAMLs, sin esperar.

En el flujo GitOps del curso, el rollback equivale a revertir el commit que actualizó la imagen en `deployment.yaml` y hacer push. ArgoCD detecta el cambio y aplica la versión anterior. La diferencia entre el rollback con `kubectl` y el rollback con Git es que el segundo deja trazabilidad — hay un commit que dice 'revertí al EP48 porque la versión v3 tenía un bug'."

---

### Resumen de comandos (14:00 – 14:30)

> *Pantalla: tabla completa de comandos.*

"El cheatsheet del episodio:

| Comando | Para qué |
|---|---|
| `kubectl get pods -n NS` | Ver el estado de los Pods |
| `kubectl get pods -n NS -w` | Monitorear en tiempo real |
| `kubectl get all -n NS` | Ver todos los recursos del namespace |
| `kubectl describe pod NOMBRE -n NS` | Diagnóstico profundo — busca la sección Events |
| `kubectl logs NOMBRE -n NS` | Ver output de la aplicación |
| `kubectl logs -f NOMBRE -n NS` | Logs en tiempo real |
| `kubectl logs --previous NOMBRE -n NS` | Logs del Pod anterior (si se cayó) |
| `kubectl exec -it NOMBRE -n NS -- sh` | Entrar al contenedor |
| `kubectl scale deployment NOMBRE -n NS --replicas=N` | Cambiar réplicas |
| `kubectl rollout status deployment/NOMBRE -n NS` | Estado de un rollout |
| `kubectl rollout undo deployment/NOMBRE -n NS` | Rollback a versión anterior |"

---

### CIERRE (14:30 – 15:30)

"Eso es el episodio 26.

Ahora tienen el vocabulario completo para operar Kubernetes. Cuando en el EP40 ArgoCD despliegue el stack completo en K3s y algo no funcione como esperan, saben exactamente qué hacer: `get` para ver el estado, `describe` para leer los eventos, `logs` para ver el error de la aplicación, `exec` para entrar y verificar desde dentro.

Antes de cerrar, detengo Minikube — hemos terminado con él por ahora:"

```bash
minikube stop
```

"El cluster queda en pausa, con todos sus datos intactos. Si en algún momento quieren repasar estos comandos, `minikube start` y está de vuelta.

En el siguiente episodio arrancamos el Módulo 07 — K3s en AWS. Volvemos a la EC2 del EP22, configuramos el Swap, e instalamos K3s. El cluster de producción del curso.

Nos vemos en el EP27."

---

## ✅ Checklist de Verificación
- [ ] `kubectl get pods -n practica-k8s -o wide` muestra la IP del Pod
- [ ] `kubectl describe pod` muestra la sección Events sin errores
- [ ] `kubectl logs -f` muestra los logs en tiempo real
- [ ] `kubectl exec -it ... -- sh` abre una shell dentro del contenedor
- [ ] Puedes escalar a 3 réplicas y volver a 1 con `kubectl scale`
- [ ] `kubectl rollout undo` revierte el Deployment a la versión anterior
- [ ] Entiendes por qué en GitOps se prefiere el rollback por Git sobre `kubectl rollout undo`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `kubectl logs` devuelve `previous terminated container ... not found` | El Pod no se ha reiniciado todavía — no hay logs previos |
| `kubectl exec` devuelve `OCI runtime exec failed: exec failed: ... executable file not found` | El contenedor no tiene `bash` — usar `sh` en su lugar |
| `kubectl scale` no persiste después de un rato | ArgoCD o el Deployment controller están revirtiendo al valor del YAML — editar el YAML directamente |
| `kubectl get pods` no muestra nada | Verificar el namespace con `-n NOMBRE` o usar `-A` para ver todos |
| `error: you must be logged in to the server` | El contexto de kubectl no apunta al cluster correcto — `kubectl config get-contexts` |

---

## 🗒️ Notas de Producción
- La apertura con el Pod en `CrashLoopBackOff` es el gancho más efectivo del episodio — genera la sensación de "eso me va a pasar a mí, necesito saber qué hacer".
- Al mostrar el output de `kubectl describe`, hacer scroll lento hasta la sección `Events` y detenerse ahí — es la información más valiosa y la que la gente tiende a ignorar porque el output es largo.
- El `kubectl exec` dentro del Pod es el momento más interactivo — explorar brevemente con `env` y `wget` hace que la audiencia entienda que están literalmente dentro del contenedor.
- La advertencia sobre `kubectl scale` en GitOps necesita énfasis verbal: "esto funciona ahora, pero cuando tengamos ArgoCD, no hagas esto". Es un anti-patrón que la gente aprende mal y luego cuesta desaprender.
- Hacer `minikube stop` al final del video en vivo — cierra el ciclo de los EP24-26 y comunica que la siguiente etapa es diferente.

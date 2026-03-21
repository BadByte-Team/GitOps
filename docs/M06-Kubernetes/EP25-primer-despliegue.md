# EP 25: Primer Despliegue con Manifiestos YAML

**Tipo:** PRÁCTICA
**Duración estimada:** 15–18 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Escribir y aplicar los tres manifiestos fundamentales — Namespace, Deployment y Service — a un cluster de Minikube, y verificar que el ciclo completo de despliegue funciona antes de llegar a K3s en producción.

---

## 📋 Prerequisitos
- Minikube corriendo (EP24)
- Conceptos de Pod, Deployment y Service claros (EP23)

---

## 🧠 Anatomía de un manifiesto YAML

Todo manifiesto de Kubernetes sigue la misma estructura de cuatro campos obligatorios:

```yaml
apiVersion: apps/v1     # versión de la API de Kubernetes para este objeto
kind: Deployment        # tipo de objeto
metadata:               # datos de identificación
  name: mi-app
  namespace: mi-namespace
spec:                   # la especificación — qué quieres que haga
  replicas: 1
  ...
```

El `apiVersion` cambia según el tipo de objeto:

| Objeto | apiVersion |
|---|---|
| Pod, Service, Namespace, ConfigMap, Secret | `v1` |
| Deployment, ReplicaSet | `apps/v1` |
| Ingress | `networking.k8s.io/v1` |
| Application (ArgoCD) | `argoproj.io/v1alpha1` |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: los archivos YAML del proyecto en `gitops-infra/infrastructure/kubernetes/app/` visibles en VS Code.*

"Bienvenidos al episodio 25.

Estos son los manifiestos de Kubernetes del proyecto — los archivos YAML que ArgoCD va a aplicar al cluster K3s en el EP40. Hay siete archivos: namespace, secrets, configmap, dos deployments, dos services.

Hoy no vamos a aplicar todos esos. Vamos a escribir versiones simplificadas de los más importantes — el Namespace, un Deployment, y un Service — y los vamos a aplicar a Minikube para ver el ciclo completo en acción.

La razón de usar versiones simplificadas es didáctica: los manifiestos completos del proyecto incluyen secretos, variables de entorno, health checks, límites de recursos. Tiene mucho ruido para un primer contacto. Hoy nos quedamos con lo esencial para entender la estructura. El próximo episodio cubrimos el `kubectl` para operar lo que creamos.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — La estructura de un manifiesto (1:00 – 3:00)

> *Pantalla: un YAML en blanco con los cuatro campos anotados.*

"Antes de escribir cualquier cosa, la estructura que todo manifiesto de Kubernetes comparte.

Cuatro campos siempre presentes:

**`apiVersion`** — qué versión de la API de Kubernetes usamos para este objeto. No es la versión de Kubernetes en sí — es la versión del grupo de APIs donde vive este tipo de objeto. Los Pods, Services y Namespaces son tan fundamentales que están en `v1`, la API original. Los Deployments, que llegaron después, están en `apps/v1`.

**`kind`** — el tipo de objeto. `Pod`, `Deployment`, `Service`, `Namespace`. Este campo le dice a Kubernetes qué clase de recurso estás describiendo.

**`metadata`** — los datos de identificación. El `name` es obligatorio — es el nombre con el que kubectl y otros objetos referencian este recurso. El `namespace` agrupa los recursos — si no lo pones, va al namespace `default`.

**`spec`** — la especificación. Aquí va la definición completa de lo que quieres. El contenido de `spec` cambia completamente según el `kind` — la spec de un Deployment es muy diferente a la de un Service.

Con esa estructura clara, podemos leer cualquier manifiesto YAML de Kubernetes."

---

### PASO 1 — Crear el directorio de práctica (3:00 – 3:30)

> *Pantalla: terminal.*

```bash
mkdir ~/k8s-practica && cd ~/k8s-practica
```

---

### PASO 2 — El Namespace (3:30 – 5:30)

> *Pantalla: VS Code creando `namespace.yaml`.*

"Empiezo por el Namespace. Es el más simple de los tres — crea el espacio de nombres donde van a vivir todos nuestros recursos."

```bash
cat > namespace.yaml << 'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: practica-k8s
EOF
```

"`apiVersion: v1` porque Namespace es un objeto fundamental. `kind: Namespace`. Y en `metadata`, solo el `name`.

Lo aplico:"

```bash
kubectl apply -f namespace.yaml
# namespace/practica-k8s created
```

"Esa palabra al final — `created` — es la confirmación. Kubernetes creó el Namespace.

Verifico:"

```bash
kubectl get namespaces
```

"Aparece `practica-k8s` en la lista junto con los namespaces del sistema como `kube-system` y `default`.

El `kubectl apply -f` es el comando que más vamos a usar. Le dice a Kubernetes: 'toma este archivo YAML y hazlo realidad'. Si el objeto no existe, lo crea. Si existe y el YAML cambió, lo actualiza. Si el YAML está igual que el estado actual, no hace nada. Esa idempotencia es fundamental para el flujo de GitOps."

---

### PASO 3 — El Deployment (5:30 – 9:30)

> *Pantalla: VS Code creando `deployment.yaml`.*

"Ahora el Deployment. Este es más complejo — tiene tres niveles de anidamiento que hay que entender."

```bash
cat > deployment.yaml << 'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-practica
  namespace: practica-k8s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-practica
  template:
    metadata:
      labels:
        app: nginx-practica
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
EOF
```

"Voy a explicar cada sección porque estas son exactamente las mismas secciones que tiene el `deployment.yaml` del proyecto.

`spec.replicas: 1` — cuántos Pods quiero que siempre estén corriendo.

`spec.selector.matchLabels` — cómo el Deployment identifica qué Pods gestiona. Los Pods que tengan el label `app: nginx-practica` son 'suyos'. Si uno muere, el Deployment crea uno nuevo con ese mismo label.

`spec.template` — la plantilla para crear los Pods. Todo lo que está dentro de `template` es la definición del Pod que el Deployment va a crear y mantener.

`spec.template.metadata.labels` — los labels del Pod. Tienen que coincidir con el `selector.matchLabels` de arriba. Si no coinciden, el Deployment no reconoce sus propios Pods — es un error común.

`spec.template.spec.containers` — la lista de contenedores dentro del Pod. En nuestro caso uno solo: nginx en su versión Alpine, expuesto en el puerto 80.

---

Lo aplico:"

```bash
kubectl apply -f deployment.yaml
# deployment.apps/nginx-practica created
```

"Verifico que el Deployment existe:"

```bash
kubectl get deployments -n practica-k8s
```

```
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
nginx-practica   1/1     1            1           30s
```

"El `1/1` en la columna READY significa: quiero 1 réplica, tengo 1 réplica lista. El Deployment está sano.

Verifico los Pods que creó:"

```bash
kubectl get pods -n practica-k8s
```

```
NAME                              READY   STATUS    RESTARTS   AGE
nginx-practica-6d8f7c5b9-xk8j2   1/1     Running   0          35s
```

"El Pod tiene un nombre generado automáticamente — el nombre del Deployment más un hash. Está en estado `Running`. El `0` en RESTARTS es buena señal — no ha tenido que reiniciarse.

Fíjense que el Pod está en el namespace `practica-k8s` — sin el flag `-n practica-k8s` en el `kubectl get pods`, no aparecería. Cuando se trabaja con múltiples namespaces, el `-n` es fundamental."

---

### PASO 4 — El Service (9:30 – 12:30)

> *Pantalla: VS Code creando `service.yaml`.*

"Por último el Service. Voy a crear un NodePort para poder acceder a nginx desde el navegador."

```bash
cat > service.yaml << 'EOF'
apiVersion: v1
kind: Service
metadata:
  name: nginx-svc
  namespace: practica-k8s
spec:
  type: NodePort
  selector:
    app: nginx-practica
  ports:
  - name: http
    protocol: TCP
    port: 80
    targetPort: 80
    nodePort: 30090
EOF
```

"El campo más importante aquí es `selector`. El Service usa `app: nginx-practica` para identificar a qué Pods debe enviar el tráfico. Tiene que coincidir exactamente con el label que definimos en el Deployment. Es el mecanismo de conexión entre el Service y sus Pods.

`port: 80` — el puerto en el que el Service escucha dentro del cluster.
`targetPort: 80` — el puerto del contenedor al que reenvía el tráfico.
`nodePort: 30090` — el puerto externo en la máquina para acceder desde fuera.

Lo aplico:"

```bash
kubectl apply -f service.yaml
# service/nginx-svc created
```

"Verifico:"

```bash
kubectl get services -n practica-k8s
```

```
NAME        TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)        AGE
nginx-svc   NodePort   10.96.xxx.xxx   <none>        80:30090/TCP   10s
```

"NodePort activo. El `80:30090/TCP` confirma que el puerto 80 interno está mapeado al 30090 externo.

Para acceder desde el navegador en Minikube, necesito la IP del nodo:"

```bash
minikube service nginx-svc -n practica-k8s --url
# http://127.0.0.1:PORT  ← abre directamente en el navegador
```

"O con la IP del nodo directamente:"

```bash
minikube ip
# 192.168.49.2

# En el navegador: http://192.168.49.2:30090
```

"Aparece la página de bienvenida de nginx. El flujo completo funciona: Deployment creó el Pod, Service expone el Pod, y desde fuera podemos acceder a la app."

---

### PASO 5 — Aplicar todo de una vez (12:30 – 13:30)

> *Pantalla: terminal.*

"Una práctica importante que usaremos mucho: aplicar un directorio completo en lugar de archivo por archivo.

`kubectl apply` puede apuntar a un directorio y aplicar todos los YAML que encuentre:"

```bash
# Primero borro todo lo que creé
kubectl delete namespace practica-k8s

# Espero a que se elimine
kubectl get namespaces -w

# Ahora lo aplico todo de una vez
kubectl apply -f ~/k8s-practica/
```

"Kubernetes aplica los archivos en orden y crea todos los objetos. Esto es exactamente lo que hace ArgoCD en el EP40 cuando sincroniza el directorio `gitops-infra/infrastructure/kubernetes/app/` — aplica todos los YAMLs del directorio de una sola vez.

La única precaución: el Namespace debe existir antes de que los recursos que viven en él. `kubectl apply` no garantiza el orden. Por eso en el EP40 verán que el `namespace.yaml` se aplica primero, separado del resto."

---

### PASO 6 — Limpiar (13:30 – 14:00)

> *Pantalla: terminal.*

"Limpio el directorio de práctica y el namespace del cluster:"

```bash
kubectl delete namespace practica-k8s
rm -rf ~/k8s-practica
```

"Eliminar el Namespace borra automáticamente todos los recursos que contiene — el Deployment, los Pods, el Service. Una sola operación limpia todo."

---

### CIERRE (14:00 – 15:00)

"Eso es el episodio 25.

Escribieron tres manifiestos YAML desde cero, los aplicaron a un cluster real, verificaron que los Pods corrían, y accedieron a la app desde el navegador. El ciclo completo de un despliegue en Kubernetes.

En el siguiente episodio cubrimos los comandos de `kubectl` que más vas a usar: `logs`, `exec`, `describe`, `scale`. Porque desplegar es solo la mitad del trabajo — la otra mitad es poder operar y depurar lo que está corriendo.

Nos vemos en el EP26."

---

## 📋 Referencia: estructura del Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nombre-del-deployment    # nombre único en el namespace
  namespace: mi-namespace
spec:
  replicas: 1                    # cuántos Pods mantener vivos
  selector:
    matchLabels:
      app: mi-app                # debe coincidir con los labels del template
  template:                      # plantilla para crear Pods
    metadata:
      labels:
        app: mi-app              # label que identifica los Pods de este Deployment
    spec:
      containers:
      - name: contenedor
        image: imagen:tag        # Jenkins actualiza esta línea en el flujo GitOps
        ports:
        - containerPort: 8080
```

---

## ✅ Checklist de Verificación
- [ ] `kubectl apply -f namespace.yaml` crea el Namespace sin errores
- [ ] `kubectl get deployments -n practica-k8s` muestra `1/1 READY`
- [ ] `kubectl get pods -n practica-k8s` muestra el Pod en estado `Running`
- [ ] `kubectl get services -n practica-k8s` muestra el NodePort activo
- [ ] La página de nginx es accesible desde el navegador
- [ ] `kubectl apply -f directorio/` aplica todos los YAMLs de una vez
- [ ] `kubectl delete namespace` borra todos los recursos del namespace

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Pod en estado `Pending` | El nodo no tiene recursos suficientes — `kubectl describe pod NOMBRE -n NAMESPACE` para ver el motivo |
| Pod en estado `ImagePullBackOff` | La imagen no existe o el nombre está mal — verificar `image:` en el YAML |
| Pod en `CrashLoopBackOff` | El contenedor arranca y se cae — `kubectl logs NOMBRE -n NAMESPACE` para ver el error |
| Service no conecta con los Pods | El `selector` del Service no coincide con los `labels` del Pod — verificar que son idénticos |
| `error: the namespace ... not found` | El Namespace no existe todavía — aplicar `namespace.yaml` primero |

---

## 🗒️ Notas de Producción
- La apertura con los YAMLs del proyecto reales es el ancla narrativa — "hoy hacemos la versión simple de esto" establece la motivación sin abrumar.
- Al escribir el Deployment, detenerse en el `selector.matchLabels` y el `template.metadata.labels` — son la misma cadena y la relación entre ambas confunde a mucha gente. Señalar con el cursor las dos aparecencias de `app: nginx-practica` mientras se explica.
- El momento de abrir nginx en el navegador es el más satisfactorio del episodio — esperar a que cargue con la pantalla en primer plano antes de continuar.
- La sección de `kubectl apply -f directorio/` tiene que conectarse explícitamente con lo que hace ArgoCD — "esto es exactamente lo que ArgoCD hace cuando sincroniza el repo de infraestructura".
- Limpiar al final en vivo — `kubectl delete namespace practica-k8s` — establece el hábito de no dejar recursos innecesarios en el cluster.

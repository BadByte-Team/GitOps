# EP 23: Conceptos Básicos — Pod, Deployment y Service

**Tipo:** TEORÍA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Entender los tres objetos fundamentales de Kubernetes — Pod, Deployment y Service — y la arquitectura del cluster, para que cuando en el EP25 apliquemos manifiestos YAML, cada línea tenga sentido desde el primer momento.

---

## 📋 Prerequisitos
- EC2 t3.micro creada (EP22) — el servidor donde instalaremos K3s en el EP28
- Ningún prerequisito técnico para este episodio — es teoría pura

---

## 🧠 Arquitectura de un cluster Kubernetes

```
┌─────────────────────────────────────────────────────┐
│                    CLUSTER K3s                       │
│                                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │                   NODE                        │   │
│  │   (nuestra EC2 t3.micro)                     │   │
│  │                                               │   │
│  │   ┌─────────────┐   ┌─────────────┐         │   │
│  │   │     Pod     │   │     Pod     │         │   │
│  │   │  ┌───────┐  │   │  ┌───────┐  │         │   │
│  │   │  │  App  │  │   │  │ MySQL │  │         │   │
│  │   │  │  Go   │  │   │  │       │  │         │   │
│  │   │  └───────┘  │   │  └───────┘  │         │   │
│  │   └─────────────┘   └─────────────┘         │   │
│  └──────────────────────────────────────────────┘   │
│                        ▲                             │
│                        │ Service (NodePort)          │
│                        │ puerto 30081               │
└────────────────────────┼────────────────────────────┘
                         │
                      Internet
```

En K3s (nuestro caso), el control plane y el nodo de trabajo viven en la misma máquina — la EC2 t3.micro del EP22. En EKS o GKE, estarían separados. La API es exactamente la misma.

---

	## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: el `docker-compose.yml` del proyecto abierto en VS Code — los servicios `api` y `mysql-db` visibles.*

"Bienvenidos al episodio 23.

Este es el `docker-compose.yml` del proyecto. Define dos servicios: la app Go y MySQL. Con `docker compose up`, ambos levantan en tu máquina local en segundos. Simple, claro, funcional.

¿Entonces para qué necesitamos Kubernetes?

La respuesta tiene varias partes, pero la más práctica para este curso es esta: Docker Compose vive en tu máquina. Kubernetes vive en la nube. Cuando decimos 'desplegar en producción', lo que queremos es que la app corra en un servidor, no en tu laptop. Y el estándar de la industria para gestionar contenedores en producción es Kubernetes.

Lo que hiciste con Docker Compose — definir servicios, configurar variables de entorno, manejar dependencias entre contenedores — vas a hacer exactamente lo mismo en Kubernetes. La diferencia es el formato y la escala.

Hoy es un episodio de conceptos. No hay comandos, no hay código. Solo los tres objetos de Kubernetes que necesitas entender antes de poder trabajar con él.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — La arquitectura del cluster (1:00 – 3:30)

> *Pantalla: diagrama del cluster — puede ser el del encabezado de este documento o un slide.*

"Antes de los objetos, la arquitectura.

Un cluster de Kubernetes tiene dos tipos de componentes: el **control plane** y los **nodes**.

El **control plane** es el cerebro. Contiene el API server — que recibe todas las instrucciones — el scheduler — que decide en qué nodo correr cada Pod — y el etcd — la base de datos que guarda el estado del cluster. Cuando ejecutas `kubectl apply -f deployment.yaml`, tu instrucción va al API server del control plane.

Los **nodes** son los músculos. Son las máquinas donde realmente corren tus aplicaciones. Cada node tiene un agente llamado `kubelet` que recibe instrucciones del control plane y las ejecuta.

En la arquitectura original del curso, el control plane vivía en EKS — el servicio administrado de AWS — y los nodes eran instancias EC2 separadas. Eso costaba ~$72 al mes solo por el control plane.

En nuestra arquitectura gratuita, usamos **K3s**: una distribución certificada de Kubernetes donde el control plane y el node corren en la misma máquina — nuestra EC2 t3.micro del EP22. Sin costo adicional, con la misma API, con los mismos manifiestos YAML.

Cuando lleguen al EP27, instalar K3s va a ser un solo comando `curl`. Pero antes de eso, necesitamos entender qué va a correr ahí."

---

### El Pod (3:30 – 6:00)

> *Pantalla: diagrama del Pod con un contenedor dentro.*

"El **Pod** es la unidad mínima en Kubernetes. Es lo que Kubernetes crea, gestiona y destruye.

Un Pod contiene uno o más contenedores Docker que comparten la misma red y el mismo sistema de archivos. En la práctica, la gran mayoría de los Pods tienen un solo contenedor.

La relación con Docker es directa: el Pod es el envoltorio de Kubernetes alrededor de un contenedor. Si en `docker-compose.yml` tenías un servicio `api` con su imagen, en Kubernetes tienes un Pod con la misma imagen.

Pero hay una diferencia crítica: **los Pods son efímeros**. Si un Pod muere — por un error, por falta de memoria, por lo que sea — Kubernetes no lo reinicia automáticamente. El Pod simplemente desaparece.

¿Eso no es un problema? Sí, lo sería. Por eso nadie crea Pods directamente en producción. Se crea un Deployment, y el Deployment se encarga de mantener los Pods vivos. Eso nos lleva al segundo objeto."

---

### El Deployment (6:00 – 8:30)

> *Pantalla: diagrama mostrando un Deployment gestionando múltiples Pods.*

"El **Deployment** es el gestor de Pods. Le dices a Kubernetes: 'quiero que siempre haya exactamente 1 réplica de este Pod corriendo'. Y el Deployment se encarga de que eso sea siempre verdad.

Si el Pod muere, el Deployment crea uno nuevo. Si hay un error en la imagen nueva y el Pod no arranca, el Deployment hace rollback a la versión anterior automáticamente. Si quieres escalar a 3 réplicas, cambias el número y el Deployment crea los dos Pods adicionales.

Cuando Jenkins actualiza el tag de la imagen en el `deployment.yaml` — que es el corazón del flujo GitOps que construiremos — lo que está haciendo es decirle al Deployment: 'ahora quiero que los Pods usen esta nueva imagen'. El Deployment hace el **rolling update**: crea el Pod nuevo con la imagen nueva, espera a que esté sano, y recién entonces elimina el Pod viejo. Así no hay downtime.

En nuestro proyecto tenemos dos Deployments:
- `mysql` — gestiona el Pod de MySQL
- `curso-gitops` — gestiona el Pod de la app Go

Miran el `deployment.yaml` en `gitops-infra/infrastructure/kubernetes/app/` — tiene `replicas: 1` porque estamos en una EC2 con 1 GB de RAM. Con más recursos, podríamos poner `replicas: 3` y Kubernetes distribuiría el tráfico entre las tres instancias."

---

### El Service (8:30 – 11:00)

> *Pantalla: diagrama mostrando el Service como punto de acceso estable a los Pods.*

"El tercer objeto, y el que más confusión genera al principio: el **Service**.

El problema que resuelve es este: los Pods tienen IPs, pero esas IPs cambian. Cada vez que Kubernetes recrea un Pod — por una actualización, por un fallo, por cualquier razón — la IP es diferente. Si tu app Go tuviera la IP de MySQL hardcodeada, se rompería cada vez que MySQL se reinicie.

El Service es una IP estable — o más precisamente, un DNS estable — que apunta siempre a los Pods correctos sin importar cuántas veces se hayan recreado. En lugar de `DB_HOST=10.x.x.x`, tienes `DB_HOST=mysql-svc`. Ese nombre siempre resuelve a los Pods de MySQL, sea cual sea su IP actual.

En el proyecto tenemos tres Services:

**`mysql-svc`** — de tipo `ClusterIP` (el default). Solo accesible desde dentro del cluster. La app Go se conecta a MySQL con `mysql-svc:3306`, y funciona sin importar cuántas veces se haya recreado el Pod de MySQL.

**`curso-gitops-svc`** — de tipo `NodePort`. Este es el que expone la app al mundo exterior. Cuando dices `http://IP_EC2:30081`, el tráfico entra al Node, el Service lo redirige al Pod de la app.

**`argocd-server`** — también `NodePort`, en el puerto 30080. Ya lo configuramos en el EP22 cuando abrimos ese puerto en el Security Group de Terraform.

---

Los tres tipos de Service que hay que conocer:

**`ClusterIP`** — solo dentro del cluster. Para la comunicación entre servicios internos, como la app hablando con la base de datos. No necesita puerto externo.

**`NodePort`** — expone el servicio en un puerto de la máquina. Accesible desde fuera con `IP_NODO:PUERTO`. Es lo que usamos en el curso porque no cuesta nada.

**`LoadBalancer`** — crea un balanceador de carga en AWS o GCP automáticamente. Más profesional, pero cuesta ~$20 al mes. En el curso lo reemplazamos con NodePort para mantener el costo en cero."

---

### La conexión con el proyecto (11:00 – 13:00)

> *Pantalla: VS Code con los archivos YAML del proyecto en `gitops-infra/infrastructure/kubernetes/app/`.*

"Antes de cerrar, quiero conectar todo esto con los archivos que ya existen en el proyecto.

Abro el directorio `gitops-infra/infrastructure/kubernetes/app/`. Ahí están todos los manifiestos:

- `namespace.yaml` — crea el namespace `curso-gitops` que agrupa todos nuestros recursos
- `secrets.yaml` — las credenciales de MySQL y el JWT secret, en Base64
- `mysql-configmap.yaml` — el script SQL de inicialización de la base de datos
- `mysql-deployment.yaml` — el Deployment de MySQL, el gestor del Pod
- `mysql-service.yaml` — el Service interno de MySQL (`mysql-svc`)
- `deployment.yaml` — el Deployment de la app Go, con `replicas: 1`
- `service.yaml` — el NodePort 30081 que expone la app al exterior

Son exactamente los tres objetos que acabamos de ver — Pods (gestionados por los Deployments), Services — más los objetos de configuración que los alimentan.

En el EP25 vamos a aplicar algunos de estos manifiestos a un cluster local de Minikube para ver el flujo completo. En el EP27 arrancamos con K3s y ArgoCD los desplegará todos juntos en la EC2 de producción."

---

### CIERRE (13:00 – 14:00)

"Eso es el episodio 23.

Pod, Deployment, Service. Los tres objetos sobre los que descansa todo lo que construiremos en los módulos de Kubernetes y GitOps.

Un Pod es la unidad mínima — el contenedor con su envoltorio de Kubernetes. Un Deployment es el gestor — mantiene los Pods vivos y hace rolling updates sin downtime. Y un Service es la dirección estable — permite que los componentes se encuentren entre sí sin depender de IPs que cambian.

En el siguiente episodio instalamos Minikube y kubectl para tener un cluster local donde practicar antes de llegar a K3s en AWS.

Nos vemos en el EP24."

---

## 📋 Referencia de Objetos

| Objeto | Qué es | Para qué se usa |
|---|---|---|
| **Pod** | Unidad mínima — 1 o más contenedores | Rara vez se crea directamente en producción |
| **Deployment** | Gestor de Pods | Mantiene N réplicas vivas, hace rolling updates |
| **Service ClusterIP** | Dirección interna estable | Comunicación entre servicios dentro del cluster |
| **Service NodePort** | Exposición externa en un puerto | Acceso desde fuera del cluster sin LoadBalancer |
| **Namespace** | Aislamiento lógico | Agrupa recursos relacionados |
| **ConfigMap** | Configuración no sensible | Variables de entorno, scripts de inicialización |
| **Secret** | Configuración sensible | Contraseñas, tokens, claves en Base64 |

---

## ✅ Checklist de Verificación
- [ ] Entiendes la diferencia entre Pod y Deployment — por qué no se crean Pods directamente
- [ ] Sabes por qué existe el Service — el problema de las IPs cambiantes
- [ ] Entiendes la diferencia entre ClusterIP y NodePort
- [ ] Puedes identificar en `gitops-infra/infrastructure/kubernetes/app/` qué archivo es qué objeto

---

## 🗒️ Notas de Producción
- Abrir el `docker-compose.yml` del proyecto al inicio — la comparación con lo que ya conocen es el gancho más efectivo para introducir Kubernetes.
- El diagrama de arquitectura del cluster puede presentarse como slide o dibujarse en vivo. Si se dibuja en vivo, hacerlo despacio y verbalizando cada componente mientras aparece.
- Al hablar del rolling update del Deployment, describir el proceso paso a paso: "Pod nuevo → espera sano → elimina viejo". Ese proceso es exactamente lo que verán en el EP48 cuando el pipeline despliegue una nueva versión.
- Al mostrar los archivos YAML al final, no leer el contenido — solo señalar cada archivo y nombrar qué objeto de Kubernetes representa. El contenido se ve en el EP25.
- Mencionar explícitamente que el `replicas: 1` en el `deployment.yaml` es una decisión de arquitectura consciente para el t3.micro — no un error ni una limitación técnica de Kubernetes.

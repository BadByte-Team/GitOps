# EP 24: Minikube y kubectl — Cluster Local para Practicar

**Tipo:** INSTALACIÓN
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar kubectl y Minikube en la máquina local, levantar un cluster de un solo nodo, y verificar que podemos comunicarnos con él. Este cluster local es el entorno de práctica para los EP25 y EP26 — antes de usar K3s en la EC2 de producción.

---

## 📋 Prerequisitos
- Docker instalado y corriendo en local (EP08)
- Terminal con acceso a internet

---

## 🧠 Minikube vs K3s — ¿por qué dos clusters?

Una pregunta válida: si en el EP28 vamos a instalar K3s en la EC2, ¿para qué instalar Minikube ahora?

| | Minikube (EP24–26) | K3s en EC2 (EP27–30) |
|---|---|---|
| **Dónde corre** | Tu máquina local | AWS EC2 t2.micro |
| **Propósito** | Práctica y aprendizaje | Producción del curso |
| **Costo** | $0 | $0 (Free Tier) |
| **Velocidad de setup** | 2 minutos | ~5 minutos |
| **Persistencia** | Solo mientras practicas | Permanente hasta EP49 |

La razón es pedagógica: practicar los comandos de `kubectl` localmente, donde si algo sale mal no afecta el servidor de producción, y donde el ciclo de prueba es más rápido. Una vez que los comandos son familiares, el salto a K3s en AWS es natural.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: terminal mostrando la IP de la EC2 creada en el EP22 — `terraform output prod_public_ip`.*

"Bienvenidos al episodio 24.

Aquí está la IP del servidor que creamos en el EP22. Esa EC2 ya está esperando que lleguemos al EP27 para instalar K3s encima. Pero antes de llegar ahí, necesitamos aprender a usar Kubernetes.

Y para aprender a usar Kubernetes — especialmente los primeros comandos, los primeros manifiestos, los primeros errores — es mucho mejor tener un cluster local que un servidor en la nube. Si rompes algo localmente, no pasa nada. Si rompes algo en la EC2 a mitad del curso, tienes que destruirla y recrearla.

Por eso instalamos Minikube hoy. Es un cluster de Kubernetes de un solo nodo que corre en Docker, en tu máquina. Setup en 2 minutos, práctica sin riesgo.

Los comandos de `kubectl` que aprendes aquí son exactamente los mismos que usarás con K3s, con EKS, con GKE, con cualquier cluster. La herramienta es universal.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — ¿Qué es kubectl? (1:00 – 2:30)

> *Pantalla: diagrama simple mostrando kubectl → API Server → Cluster.*

"Antes de instalar nada, una distinción importante que a mucha gente confunde al principio.

**Minikube** y **kubectl** son dos cosas diferentes.

**Minikube** crea y gestiona un cluster de Kubernetes local. Es la herramienta de infraestructura — crea el nodo, levanta los componentes del control plane, expone la API. Sin Minikube, no hay cluster con quien hablar.

**kubectl** es el cliente de línea de comandos que se comunica con cualquier cluster de Kubernetes — Minikube, K3s, EKS, GKE, lo que sea. Envía instrucciones al API server del cluster. Es como la AWS CLI para AWS, pero para cualquier cluster de Kubernetes.

La relación es la misma que en el módulo anterior: Terraform es la herramienta que crea la infraestructura, y la AWS CLI es el cliente que habla con esa infraestructura. Aquí, Minikube crea el cluster y kubectl habla con él.

Cuando en el EP27 tengamos K3s en la EC2, seguiremos usando el mismo kubectl. Solo cambia el cluster al que apunta — de Minikube local a K3s remoto."

---

### PASO 1 — Instalar kubectl (2:30 – 5:00)

> *Pantalla: terminal en la PC local.*

"Empiezo por kubectl — el cliente que usaremos siempre, independientemente del cluster que tengamos.

El instalador oficial descarga el binario de la versión estable más reciente:"

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s \
  https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
```

"El `$(curl -L -s https://dl.k8s.io/release/stable.txt)` detecta automáticamente cuál es la versión estable más reciente — así el mismo comando funciona sin importar cuándo lo ejecutes.

Instalo el binario en `/usr/local/bin` — el directorio estándar para herramientas que deben estar en el PATH:"

```bash
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl
```

"El `install` con esos permisos garantiza que el archivo sea ejecutable por todos los usuarios pero solo editable por root. Limpio el archivo temporal descargado.

Verifico:"

```bash
kubectl version --client
```

"Verás algo como `Client Version: v1.29.x`. Solo la versión del cliente — el servidor no existe todavía porque no hemos levantado ningún cluster.

---

Para **macOS**:"

```bash
brew install kubectl
```

"Para **Windows**, el instalador oficial está en la documentación de Kubernetes en `kubernetes.io/docs/tasks/tools/install-kubectl-windows`."

---

### PASO 2 — Instalar Minikube (5:00 – 7:00)

> *Pantalla: terminal.*

"Ahora Minikube. Similar proceso — descargo el binario y lo instalo:"

```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64

sudo install minikube-linux-amd64 /usr/local/bin/minikube
rm minikube-linux-amd64
```

"Verifico:"

```bash
minikube version
# minikube version: v1.32.x
```

"Para macOS:"

```bash
brew install minikube
```

---

### PASO 3 — Levantar el cluster (7:00 – 9:30)

> *Pantalla: terminal.*

"Ahora la parte interesante: levantar el cluster. Un solo comando:"

```bash
minikube start --driver=docker
```

"El flag `--driver=docker` le dice a Minikube que use Docker para crear el nodo. Como ya tenemos Docker instalado desde el EP08, es la opción más directa y no requiere instalar nada adicional.

Minikube va a hacer varias cosas en secuencia: descargar la imagen del nodo si no la tiene, crear un contenedor Docker con todos los componentes de Kubernetes, y configurar kubectl para que apunte a ese cluster.

Puede tardar entre 1 y 3 minutos la primera vez — está descargando la imagen del nodo si no la tiene en caché. Las siguientes veces es mucho más rápido.

Cuando termina, el output dice algo así:"

```
😄  minikube v1.32.x
✨  Using the docker driver based on user configuration
🔥  Creating docker container (CPUs=2, Memory=2200MB)
🐳  Preparing Kubernetes v1.29.x on Docker ...
🔗  Configuring bridge CNI (Container Networking Interface)
🔎  Verifying Kubernetes components...
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster
```

"Esa última línea es la clave: `kubectl is now configured to use "minikube" cluster`. Minikube configuró automáticamente kubectl para apuntar al nuevo cluster. No tuve que hacer nada manual."

---

### PASO 4 — Verificar el cluster (9:30 – 11:30)

> *Pantalla: terminal.*

"Verifico que el cluster está funcionando correctamente.

Primero, los nodos:"

```bash
kubectl get nodes
```

```
NAME       STATUS   ROLES           AGE   VERSION
minikube   Ready    control-plane   1m    v1.29.x
```

"Un nodo, estado `Ready`. En K3s la salida será muy similar — el nombre del nodo será el hostname de la EC2, pero el estado y los roles son los mismos.

Fíjense en el rol: `control-plane`. Eso confirma que en Minikube, igual que en K3s, el control plane y el nodo de trabajo son la misma máquina.

---

Los pods del sistema:"

```bash
kubectl get pods -A
```

"La flag `-A` significa `--all-namespaces` — muestra los pods de todos los namespaces, no solo del namespace por defecto.

Verás pods con nombres como `coredns`, `kube-apiserver`, `kube-scheduler`, `etcd`. Esos son los componentes del control plane corriendo como pods dentro del propio cluster. Es la arquitectura de Kubernetes: hasta el control plane corre en contenedores.

---

Información general del cluster:"

```bash
kubectl cluster-info
```

```
Kubernetes control plane is running at https://127.0.0.1:PORT
CoreDNS is running at https://127.0.0.1:PORT/api/v1/namespaces/...
```

"El API server está en `127.0.0.1` — nuestro localhost. En K3s, esta dirección será la IP pública de la EC2. El comando es idéntico, solo cambia la dirección."

---

### PASO 5 — El kubeconfig — cómo kubectl sabe a qué cluster hablar (11:30 – 13:00)

> *Pantalla: terminal y VS Code mostrando el archivo kubeconfig.*

"Un concepto que vale la pena entender ahora porque lo vamos a usar activamente en el EP30: el **kubeconfig**.

kubectl sabe a qué cluster conectarse porque lee un archivo de configuración en `~/.kube/config`. Ese archivo contiene la dirección del API server, las credenciales para autenticarse, y el nombre del cluster."

```bash
cat ~/.kube/config
```

"Ven que tiene una estructura con tres secciones: `clusters` — la dirección del API server — `users` — las credenciales — y `contexts` — que combina un cluster con un usuario y un namespace por defecto.

Cuando Minikube hizo el `start`, escribió automáticamente su entrada en este archivo y configuró el contexto activo para apuntar a él.

Cuando en el EP30 tengamos K3s en la EC2, vamos a descargar el kubeconfig de K3s, modificar la IP para que apunte a la EC2 en lugar de `127.0.0.1`, y agregarlo a este mismo archivo. A partir de ahí, kubectl podrá hablar tanto con Minikube local como con K3s en AWS — solo hay que cambiar el contexto activo."

---

### PASO 6 — Comandos de gestión de Minikube (13:00 – 14:00)

> *Pantalla: terminal.*

"Para cerrar, los comandos de Minikube que van a usar durante los próximos dos episodios:"

```bash
# Ver el estado del cluster
minikube status

# Abrir el dashboard web de Kubernetes en el navegador
minikube dashboard

# Detener el cluster (los datos se conservan)
minikube stop

# Reiniciar el cluster
minikube start

# Eliminar el cluster completamente
minikube delete
```

"El `minikube stop` detiene el cluster sin borrar nada. Si cierras la laptop y la vuelves a abrir, ejecutas `minikube start` y el cluster está exactamente como lo dejaste.

El `minikube delete` borra todo — el nodo, la configuración, el contexto en kubectl. Útil cuando quieres empezar desde cero.

Para este curso, después del EP26 vamos a hacer `minikube stop` para liberar recursos. El cluster de producción será K3s en la EC2 — Minikube ya cumplió su propósito."

---

### CIERRE (14:00 – 15:00)

"Eso es el episodio 24.

kubectl instalado, Minikube corriendo, el cluster verificado. Tienen un entorno de práctica completo en su máquina local.

En el siguiente episodio aplicamos los primeros manifiestos YAML a este cluster — un Namespace, un Deployment y un Service. Verán el ciclo completo: escribir el YAML, aplicarlo con `kubectl apply`, y verificar que los pods están corriendo.

Nos vemos en el EP25."

---

## ✅ Checklist de Verificación
- [ ] `kubectl version --client` muestra la versión del cliente
- [ ] `minikube version` muestra la versión instalada
- [ ] `minikube start --driver=docker` completa sin errores
- [ ] `kubectl get nodes` muestra el nodo `minikube` en estado `Ready`
- [ ] `kubectl get pods -A` muestra los pods del sistema corriendo
- [ ] Entiendes la diferencia entre kubectl (cliente) y Minikube (cluster)

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `kubectl: command not found` | Verificar que `/usr/local/bin` está en el PATH: `echo $PATH` |
| `minikube start` falla con `driver=docker` | Verificar que Docker está corriendo: `docker ps` |
| El nodo está en `NotReady` | Esperar 2-3 minutos más — los componentes del sistema toman tiempo en inicializarse |
| `Exiting due to PROVIDER_DOCKER_NOT_RUNNING` | Docker no está activo: `sudo systemctl start docker` |
| `kubectl get nodes` muestra cluster incorrecto | Verificar el contexto activo: `kubectl config current-context` — debe decir `minikube` |

---

## 🗒️ Notas de Producción
- La apertura mostrando la IP de la EC2 con `terraform output prod_public_ip` es el ancla con el módulo anterior — refuerza que estamos avanzando hacia ese servidor, no alejándonos.
- El diagrama de kubectl vs Minikube puede hacerse verbal con un gesto de las manos — Minikube a la izquierda (el cluster) y kubectl a la derecha (el cliente que habla con él).
- Al hacer `minikube start`, mantener la terminal en pantalla completa mientras carga — es uno de los momentos más satisfactorios del episodio visualmente (los emojis y el progreso).
- La sección del kubeconfig puede mostrarse en VS Code con la extensión YAML para que tenga colores — hace la estructura más legible visualmente.
- Al final del episodio, abrir el dashboard con `minikube dashboard` brevemente para mostrar la interfaz gráfica — muchos alumnos no saben que existe y es muy visual.

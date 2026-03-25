# EP 29: Instalación de K3s y kubectl en la EC2

**Tipo:** INSTALACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** En lugar de conectar kubectl a EKS, instalamos K3s directamente en la EC2 con el script oficial y configuramos los permisos del kubeconfig.

---

## 🎯 Objetivo
Instalar K3s en la EC2 con el script oficial de instalación, configurar los permisos del kubeconfig para que kubectl funcione sin `sudo`, y verificar que el cluster está corriendo y el nodo está `Ready`.

---

## 📋 Prerequisitos
- EC2 t3.micro con Swap de 2 GB configurado (EP28)
- Conectividad SSH a la instancia

---

## 🧠 K3s vs instalación estándar de Kubernetes

| | Kubernetes estándar (`kubeadm`) | K3s |
|---|---|---|
| Pasos de instalación | ~15 comandos: `kubeadm`, `kubelet`, `kubectl`, CNI, etcd... | **1 comando `curl`** |
| Tiempo | 15–30 minutos | ~2 minutos |
| Componentes separados | API server, scheduler, controller, etcd | Todo en un solo binario |
| RAM mínima | ~2 GB | ~512 MB |
| Resultado | Kubernetes certificado | Kubernetes certificado |

El script de instalación de K3s detecta automáticamente el sistema operativo, descarga la versión correcta del binario, configura el servicio de systemd, y levanta el cluster. Sin decisiones manuales, sin pasos intermedios.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: terminal local. Se ve la IP de la EC2.*

"Bienvenidos al episodio 29.

En el episodio anterior configuramos el Swap. La EC2 tiene ahora 1 GB de RAM física más 2 GB de memoria virtual. El servidor está listo para recibir K3s.

Si alguna vez han instalado Kubernetes con `kubeadm` — el instalador estándar — saben que son quince comandos, configuración de CNI, inicialización del cluster, unir los nodos... fácilmente media hora de trabajo.

K3s es un comando. Uno solo. El script oficial detecta el sistema operativo, descarga el binario correcto, configura el servicio de systemd, levanta el cluster, y te deja listo para usar kubectl. Todo en menos de dos minutos.

Hoy instalamos K3s y verificamos que el cluster está corriendo. Vamos."

---

### PASO 1 — Conectar a la EC2 (1:00 – 1:30)

> *Pantalla: terminal local.*

```bash
ssh -i aws-key.pem ubuntu@<IP_PUBLICA_EC2>
```

"Dentro del servidor. Todos los comandos que siguen se ejecutan aquí."

---

### PASO 2 — Instalar K3s (1:30 – 4:30)

> *Pantalla: terminal dentro de la EC2.*

"El comando de instalación. Un solo `curl` que descarga y ejecuta el script oficial de K3s desde get.k3s.io:"

```bash
curl -sfL https://get.k3s.io | sh -
```

"Voy a explicar las flags mientras descarga, porque hay gente que prefiere no ejecutar scripts de internet sin entenderlos.

`-s` silencia el progreso de curl — sin barras de progreso ni porcentajes.
`-f` hace que curl falle en silencio si hay un error HTTP — sin descargar páginas de error.
`-L` sigue redirecciones — algunos URLs redirigen a la versión final.

El `| sh -` ejecuta el script que descargó directamente. El script en sí está disponible en GitHub — `github.com/k3s-io/k3s/blob/master/install.sh` — cualquiera puede leerlo antes de ejecutarlo.

---

El proceso tarda entre 30 segundos y 2 minutos dependiendo de la velocidad de la conexión. El output va apareciendo progresivamente:"

```
[INFO]  Finding release for channel stable
[INFO]  Using v1.29.x+k3s1 as release
[INFO]  Downloading hash: https://github.com/k3s-io/k3s/releases/...
[INFO]  Downloading binary: https://github.com/k3s-io/k3s/releases/...
[INFO]  Verifying binary download
[INFO]  Installing K3s to /usr/local/bin/k3s
[INFO]  Creating /usr/local/bin/kubectl symlink to k3s
[INFO]  Creating /usr/local/bin/crictl symlink to k3s
[INFO]  Creating /usr/local/bin/ctr symlink to k3s
[INFO]  Creating killall script /usr/local/bin/k3s-killall.sh
[INFO]  Creating uninstall script /usr/local/bin/k3s-uninstall.sh
[INFO]  env: Creating environment file /etc/systemd/system/k3s.service.env
[INFO]  systemd: Creating service file /etc/systemd/system/k3s.service
[INFO]  systemd: Enabling k3s unit
[INFO]  systemd: Starting k3s
```

"Cada línea del output tiene su significado.

Descargó el binario de K3s y lo instaló en `/usr/local/bin/k3s`. Creó symlinks para que `kubectl`, `crictl` y `ctr` también funcionen como comandos directos. Configuró un servicio de systemd para que K3s arranque automáticamente con el sistema. Y finalmente arrancó el servicio.

El cluster ya está corriendo. Pero todavía no podemos usar kubectl directamente — necesitamos configurar el acceso."

---

### PASO 3 — Verificar que K3s está corriendo (4:30 – 6:00)

> *Pantalla: terminal dentro de la EC2.*

"Verifico el estado del servicio de systemd:"

```bash
sudo systemctl status k3s
```

```
● k3s.service - Lightweight Kubernetes
     Loaded: loaded (/etc/systemd/system/k3s.service; enabled; vendor preset: enabled)
     Active: active (running) since ...
```

"Dos cosas importantes: `active (running)` confirma que el proceso está corriendo, y `enabled` confirma que arrancará automáticamente en cada inicio del servidor.

Ahora intento usar kubectl con sudo para ver el cluster:"

```bash
sudo kubectl get nodes
```

```
NAME        STATUS   ROLES                  AGE   VERSION
ip-10-...   Ready    control-plane,master   1m    v1.29.x+k3s1
```

"Un nodo, estado `Ready`, rol `control-plane,master`. Exactamente lo que esperamos — en K3s, el control plane y el nodo de trabajo son la misma máquina.

La versión termina en `+k3s1` — eso es el identificador de la build de K3s sobre la versión base de Kubernetes. Es la forma de distinguir un cluster K3s de uno estándar."

---

### PASO 4 — Configurar kubectl sin sudo (6:00 – 8:30)

> *Pantalla: terminal dentro de la EC2.*

"Ahora el paso más importante de este episodio para la experiencia de uso. Fíjense que en el paso anterior tuve que usar `sudo kubectl` — con privilegios de root. Eso es incómodo y potencialmente peligroso para el trabajo diario.

El motivo es que el kubeconfig de K3s — el archivo que contiene las credenciales y la dirección del API server — se guarda en `/etc/rancher/k3s/k3s.yaml` con permisos de solo lectura para root.

La solución estándar es copiarlo al directorio personal del usuario y asignarle la propiedad correcta:"

```bash
# Crear el directorio .kube en el home del usuario ubuntu
mkdir -p ~/.kube

# Copiar el kubeconfig de K3s al lugar estándar
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config

# Asignar la propiedad al usuario ubuntu
sudo chown ubuntu:ubuntu ~/.kube/config

# Verificar los permisos resultantes
ls -la ~/.kube/config
# -rw------- 1 ubuntu ubuntu ... /home/ubuntu/.kube/config
```

"Los permisos resultantes son `600` — el usuario ubuntu puede leer y escribir, nadie más. kubectl requiere que el kubeconfig no sea legible por otros usuarios por razones de seguridad.

Ahora pruebo sin sudo:"

```bash
kubectl get nodes
```

```
NAME        STATUS   ROLES                  AGE   VERSION
ip-10-...   Ready    control-plane,master   3m    v1.29.x+k3s1
```

"Perfecto — sin sudo, sin errores. kubectl puede leer el kubeconfig correctamente."

---

### PASO 5 — Explorar el cluster recién instalado (8:30 – 10:30)

> *Pantalla: terminal dentro de la EC2.*

"Ahora que kubectl funciona correctamente, exploro el estado inicial del cluster.

Los pods del sistema:"

```bash
kubectl get pods -A
```

```
NAMESPACE     NAME                                      READY   STATUS    RESTARTS
kube-system   coredns-...                               1/1     Running   0
kube-system   local-path-provisioner-...                1/1     Running   0
kube-system   metrics-server-...                        1/1     Running   0
kube-system   traefik-...                               1/1     Running   0
```

"K3s instala por defecto cuatro componentes que merece la pena conocer:

**CoreDNS** — el servidor DNS interno del cluster. Cuando la app Go se conecta a `mysql-svc:3306`, CoreDNS resuelve `mysql-svc` a la IP del Pod de MySQL. Sin CoreDNS, los nombres de los Services no resolverían.

**local-path-provisioner** — el proveedor de almacenamiento local. Cuando un Pod pide un volumen persistente, este componente lo crea en el disco local de la instancia.

**metrics-server** — expone métricas de CPU y RAM de los pods. Lo usaremos cuando queramos ver `kubectl top pods` para entender el consumo de recursos.

**Traefik** — un Ingress Controller que K3s instala por defecto. No lo vamos a usar directamente en el curso — usamos NodePort — pero está ahí y consume un poco de memoria.

---

Verifico los recursos disponibles:"

```bash
kubectl top nodes
```

```
NAME        CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
ip-10-...   150m         15%    450Mi           46%
```

"450 MB de memoria usados — el 46% de los 981 MB de RAM física. Eso es K3s con sus componentes del sistema. Cuando ArgoCD, MySQL y la app entren en juego, la memoria subirá y el Swap empezará a trabajar. Todo dentro de lo esperado."

---

### PASO 6 — Salir del servidor (10:30 – 11:00)

> *Pantalla: terminal dentro de la EC2.*

"Todo listo. Salgo del servidor:"

```bash
exit
```

"En el siguiente episodio vamos a configurar el acceso **remoto** — desde la máquina local, sin tener que hacer SSH cada vez que queremos ejecutar un kubectl."

---

### CIERRE (11:00 – 11:30)

"Eso es el episodio 29.

K3s instalado en la EC2, kubectl funcionando sin sudo, el nodo en estado `Ready`. El cluster de producción del curso está corriendo en la capa gratuita de AWS.

En el siguiente episodio configuramos el kubeconfig en tu máquina local para que puedas ejecutar `kubectl get pods` directamente desde tu terminal, sin conectarte por SSH al servidor cada vez. Eso es lo que hace el flujo de trabajo diario con Kubernetes práctico y cómodo.

Nos vemos en el EP30."

---

## ✅ Checklist de Verificación
- [ ] `sudo systemctl status k3s` muestra `active (running)` y `enabled`
- [ ] `kubectl get nodes` (sin sudo) muestra el nodo en estado `Ready`
- [ ] `kubectl get pods -A` muestra los 4 pods del sistema en `Running`
- [ ] El archivo `~/.kube/config` existe y pertenece al usuario `ubuntu`
- [ ] `kubectl top nodes` muestra las métricas de CPU y RAM

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| El nodo está en `NotReady` | Esperar 2-3 minutos más — K3s tarda en inicializar todos los componentes |
| `error: no configuration file provided` | Verificar que `~/.kube/config` existe: `ls -la ~/.kube/` |
| K3s no inicia (OOM killer activo) | El Swap no está configurado — volver al EP28 |
| `connection refused` al hacer kubectl | K3s todavía está arrancando — esperar y volver a intentar |
| `WARN[...] Unable to read /etc/rancher/k3s/k3s.yaml` | El archivo tiene permisos incorrectos — `sudo chmod 644 /etc/rancher/k3s/k3s.yaml` |

---

## 🗒️ Notas de Producción
- Comparar brevemente el `curl | sh` de K3s con los 15 pasos de `kubeadm` — hace que el alumno aprecie la simplicidad.
- Mientras el script de instalación corre y el output aparece, leer las líneas en voz alta y explicar brevemente qué hace cada una.
- El momento donde `kubectl get nodes` (sin sudo) muestra `Ready` es el cierre técnico del episodio — hacer una pausa para que se registre.
- Al mostrar `kubectl get pods -A`, explicar brevemente cada componente del sistema — CoreDNS, Traefik, metrics-server. Son nombres que el alumno va a ver constantemente y vale la pena conocerlos desde el principio.
- `kubectl top nodes` es un buen momento para conectar con el EP28 — "aquí está el 46% de RAM siendo usado, y ya queda el Swap como colchón para cuando ArgoCD entre en escena".

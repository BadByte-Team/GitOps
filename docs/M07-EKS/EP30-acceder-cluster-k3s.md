# EP 30: Acceder al Cluster K3s desde tu PC Local

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐⭐ (Intermedio)
**🔄 MODIFICADO:** En lugar de conectar kubectl a EKS con `aws eks update-kubeconfig`, se configura el acceso remoto al cluster K3s modificando el kubeconfig para apuntar a la IP pública de la EC2.

---

## 🎯 Objetivo
Configurar kubectl en la máquina local para comunicarse con el cluster K3s de la EC2, de modo que puedas ejecutar cualquier comando kubectl directamente desde tu terminal sin necesidad de hacer SSH al servidor.

---

## 📋 Prerequisitos
- K3s instalado y kubectl configurado en la EC2 (EP29)
- kubectl instalado en tu máquina local (EP24)
- IP pública de la EC2 disponible

---

## 🧠 El problema del kubeconfig local

El kubeconfig que K3s crea en la EC2 apunta al servidor local:

```yaml
server: https://127.0.0.1:6443
```

Ese `127.0.0.1` funciona perfectamente desde dentro del servidor. Pero si copias ese archivo a tu laptop, kubectl intentará conectarse a `127.0.0.1:6443` en tu propia máquina, donde no hay ningún cluster. La solución es reemplazar `127.0.0.1` con la IP pública de la EC2 antes de descargar el archivo.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: dos terminales. Izquierda: kubectl funcionando dentro de la EC2 por SSH. Derecha: terminal local mostrando `connection refused`.*

"Bienvenidos al episodio 30.

Aquí está el problema que resolvemos hoy. En la terminal de la izquierda, estoy dentro de la EC2 por SSH y kubectl funciona perfectamente — el cluster responde, los pods están visibles. En la terminal de la derecha, desde mi máquina local, kubectl intenta conectarse a `127.0.0.1:6443` — que es mi laptop — y obviamente falla.

Para trabajar cómodamente durante el resto del curso — ver los logs de un pod, aplicar un manifiesto, verificar un despliegue de ArgoCD — necesito poder ejecutar kubectl desde mi terminal directamente, sin pasar por SSH. Eso es lo que configuramos hoy.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — El kubeconfig y los contextos (1:00 – 2:30)

> *Pantalla: el archivo `~/.kube/config` abierto en VS Code.*

"Antes de hacer cualquier cambio, entendamos qué es el kubeconfig y cómo funciona.

kubectl nunca habla directamente con un cluster Kubernetes. Primero lee un archivo de configuración llamado kubeconfig — que por defecto vive en `~/.kube/config` — y ahí encuentra la dirección del servidor, las credenciales para autenticarse, y el nombre del cluster.

El archivo tiene tres secciones:

**`clusters`** — lista los clusters que kubectl conoce. Cada uno tiene un nombre y la URL del API server.

**`users`** — las credenciales para autenticarse. En K3s son certificados de cliente que se generaron automáticamente durante la instalación.

**`contexts`** — el puente entre los dos. Un contexto combina un cluster con un usuario. Cuando cambias el contexto activo, kubectl cambia a qué cluster apunta.

Cuando en el EP24 instalamos Minikube, ese proceso escribió automáticamente un contexto en `~/.kube/config`. Hoy vamos a agregar un segundo contexto que apunta a K3s en la EC2.

Con los dos contextos disponibles, cambiar entre Minikube y K3s es un solo comando:"

```bash
kubectl config use-context minikube   # apunta a Minikube local
kubectl config use-context default    # apunta a K3s en AWS
```

---

### PASO 1 — Preparar el kubeconfig en la EC2 (2:30 – 5:00)

> *Pantalla: terminal conectado a la EC2 por SSH.*

"Me conecto a la EC2:"

```bash
ssh -i aws-key.pem ubuntu@<IP_PUBLICA_EC2>
```

"Veo el kubeconfig actual para confirmar el problema:"

```bash
cat ~/.kube/config | grep server
# server: https://127.0.0.1:6443
```

"Ahí está el `127.0.0.1` que causa el problema. Voy a crear una copia de este archivo donde reemplazo esa IP con la IP pública de la EC2.

Uso `sed` para hacer el reemplazo en una sola línea:"

```bash
IP_PUBLICA="<IP_PUBLICA_EC2>"   # reemplazar con la IP real

sed "s/127.0.0.1/$IP_PUBLICA/g" ~/.kube/config > ~/k3s-remote.yaml
```

"Verifico que el reemplazo fue correcto:"

```bash
grep server ~/k3s-remote.yaml
# server: https://54.X.X.X:6443   ← ahora tiene la IP pública
```

"Ese archivo contiene exactamente lo mismo que el kubeconfig original — el mismo certificado de la autoridad, las mismas credenciales de cliente — pero con la dirección correcta para conectarse desde afuera.

Salgo del servidor:"

```bash
exit
```

---

### PASO 2 — Verificar el puerto 6443 en el Security Group (5:00 – 6:30)

> *Pantalla: terminal local.*

"El puerto 6443 es el API server de Kubernetes. Para que kubectl pueda conectarse desde mi laptop, ese puerto tiene que estar abierto en el Security Group de la EC2.

Verifico si ya está abierto:"

```bash
aws ec2 describe-security-groups \
  --filters "Name=group-name,Values=prod-sg" \
  --query "SecurityGroups[0].IpPermissions[*].FromPort" \
  --output table
```

"Si el 6443 no aparece en la lista, lo agrego con la CLI:"

```bash
SG_ID=$(aws ec2 describe-security-groups \
  --filters "Name=group-name,Values=prod-sg" \
  --query "SecurityGroups[0].GroupId" \
  --output text)

aws ec2 authorize-security-group-ingress \
  --group-id $SG_ID \
  --protocol tcp \
  --port 6443 \
  --cidr 0.0.0.0/0
```

"Y si quieren que sea permanente — que sobreviva un `terraform destroy` seguido de un `terraform apply` — agregan este bloque al `main.tf` de Terraform:"

```hcl
ingress {
  description = "K3s API Server"
  from_port   = 6443
  to_port     = 6443
  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
}
```

---

### PASO 3 — Descargar el kubeconfig a la máquina local (6:30 – 8:00)

> *Pantalla: terminal local.*

"Con el puerto abierto, descargo el archivo que preparé en la EC2 usando `scp` — secure copy. Funciona igual que SSH pero para transferir archivos:"

```bash
scp -i aws-key.pem ubuntu@<IP_PUBLICA_EC2>:~/k3s-remote.yaml ~/.kube/k3s-config
```

"El formato de `scp` es siempre `scp origen destino`. El origen tiene el formato `usuario@host:ruta`. Lo guardo como `~/.kube/k3s-config` — un nombre distinto al `~/.kube/config` de Minikube para no sobreescribir nada.

Asigno permisos seguros — kubectl requiere que el kubeconfig no sea legible por otros usuarios:"

```bash
chmod 600 ~/.kube/k3s-config
```

---

### PASO 4 — Configurar kubectl local para usar K3s (8:00 – 10:00)

> *Pantalla: terminal local.*

"Ahora le digo a kubectl que use este archivo. La variable de entorno `KUBECONFIG` puede apuntar a múltiples archivos separados por `:`. kubectl los combina y me da acceso a todos sus contextos al mismo tiempo:"

```bash
export KUBECONFIG=~/.kube/config:~/.kube/k3s-config
```

"Verifico los contextos disponibles:"

```bash
kubectl config get-contexts
```

```
CURRENT   NAME      CLUSTER   AUTHINFO   NAMESPACE
          default   default   default
*         minikube  minikube  minikube
```

"Los dos están ahí. El asterisco está en Minikube porque fue el último que se configuró. Cambio al contexto de K3s:"

```bash
kubectl config use-context default
```

"Agregamos la IP del servidor adicionales al certificado"

```shell
nano /etc/rancher/k3s/config.yaml
```

"En el `config.yaml`"

```yml
tls-san:
 - IP_PUBLICA_EC2
```

"Reiniciamos del servicio"

```bash
sudo systemctl restart k3s 
```

"Ahora la prueba que importa:"

```bash
kubectl get nodes
```

```
NAME        STATUS   ROLES                  AGE   VERSION
ip-10-...   Ready    control-plane,master   20m   v1.29.x+k3s1
```

"Ahí está. Desde mi terminal local, sin SSH, sin estar dentro del servidor. El cluster K3s de la EC2 respondiendo a mis comandos.

---

Para que la variable sea permanente entre sesiones, la agrego al archivo de configuración del shell:"

```bash
echo 'export KUBECONFIG=~/.kube/config:~/.kube/k3s-config' >> ~/.bashrc
source ~/.bashrc
```

---

### PASO 5 — Prueba completa de conectividad (10:00 – 11:00)

> *Pantalla: terminal local.*

"Verifico con los comandos que ya conocen de los episodios anteriores, ahora apuntando al cluster remoto:"

```bash
# Información del cluster — debe mostrar la IP pública de la EC2
kubectl cluster-info
```

```
Kubernetes control plane is running at https://54.X.X.X:6443
CoreDNS is running at https://54.X.X.X:6443/api/v1/namespaces/kube-system/services/...
```

"La URL del control plane muestra la IP pública de la EC2 — confirmación de que kubectl habla con el cluster remoto.

```bash
# Todos los pods del sistema
kubectl get pods -A
```

"Los mismos cuatro pods que vimos en el EP29 — CoreDNS, Traefik, metrics-server, local-path-provisioner. Pero ahora los estoy viendo desde mi laptop.

```bash
# Métricas de recursos
kubectl top nodes
```

"El cluster K3s en la EC2, visto completamente desde la terminal local. El flujo de trabajo para el resto del módulo — ArgoCD, Jenkins, los despliegues — va a usar exactamente esta configuración."

---

### Una nota sobre la IP cambiante (11:00 – 11:30)

> *Pantalla: terminal local.*

"Una cosa práctica importante para el resto del curso. La IP pública de una instancia t3.micro **cambia cuando la instancia se reinicia**. Si en algún momento kubectl empieza a fallar con `connection refused` después de que la EC2 se reinició, el kubeconfig tiene una IP desactualizada.

La solución es simple — actualizar el kubeconfig con la nueva IP:"

```bash
NEW_IP=$(cd gitops-infra/infrastructure/terraform/jenkins-ec2 && terraform output -raw prod_public_ip)
sed -i "s/[0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+/$NEW_IP/g" ~/.kube/k3s-config
kubectl get nodes   # verificar que funciona
```

"Tarda 10 segundos. Si quieren evitar este problema completamente, pueden asociar una Elastic IP a la instancia — una IP fija que no cambia entre reinicios. Cuesta unos centavos al mes cuando la instancia está apagada. Para el curso, actualizar manualmente cuando sea necesario es perfectamente suficiente."

---

### CIERRE (11:30 – 12:00)

"Eso es el episodio 30. Y con esto cerramos el Módulo 07.

El cluster K3s está instalado en la EC2 del Free Tier, con Swap configurado para que no colapse, y con acceso remoto desde tu máquina local para trabajar cómodamente sin SSH.

Tenemos el 100% del stack de infraestructura listo. Lo que viene ahora son las herramientas de CI/CD: Jenkins en el EP31, SonarQube y Trivy, y finalmente ArgoCD en los episodios 38 al 40.

Cuando ArgoCD esté configurado apuntando a `gitops-infra`, y el pipeline de Jenkins empuje una nueva imagen, ese cluster K3s que configuramos hoy va a recibir el despliegue automáticamente. Ese es el momento del que se trata todo el curso.

Nos vemos en el EP31."

---

## ✅ Checklist de Verificación
- [ ] `kubectl get nodes` desde la terminal **local** muestra el nodo K3s en `Ready`
- [ ] `kubectl cluster-info` muestra la IP pública de la EC2 (no `127.0.0.1`)
- [ ] `kubectl get pods -A` desde local muestra los 4 pods del sistema
- [ ] La variable `KUBECONFIG` está en `~/.bashrc` para que persista
- [ ] Entiendes qué hacer si la IP cambia al reiniciar la instancia

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Unable to connect: dial tcp X.X.X.X:6443` | Puerto 6443 no está abierto en el Security Group — ver Paso 2 |
| `certificate signed by unknown authority` | El certificado del kubeconfig no coincide — rehacer el Paso 1 desde la EC2 |
| `context "default" not found` | `KUBECONFIG` no apunta al archivo correcto — `echo $KUBECONFIG` para verificar |
| kubectl muestra el cluster equivocado | `kubectl config current-context` — luego `kubectl config use-context default` para K3s |
| `connection refused` después de un reinicio | La IP de la EC2 cambió — actualizar con el comando `sed` de la sección "IP cambiante" |

---

## 🗒️ Notas de Producción
- La apertura con las dos terminales lado a lado es el gancho visual perfecto — preparar ese setup de antemano para que no haya demoras en la grabación.
- Al mostrar el `grep server` antes y después del `sed`, hacer zoom en la terminal para que el cambio de IP sea claramente visible.
- El momento donde `kubectl get nodes` funciona desde la terminal local es el cierre técnico más satisfactorio del episodio — hacer pausa verbal para que se registre.
- La nota sobre la IP cambiante es una de las cosas más prácticas del episodio — enfatizarla verbalmente para que el alumno sepa qué hacer si algo falla días después.
- Cerrar mencionando explícitamente que el Módulo 07 está completo y qué viene después — ArgoCD en los EP38-40 va a usar exactamente este cluster. Crea la anticipación del momento final del curso.

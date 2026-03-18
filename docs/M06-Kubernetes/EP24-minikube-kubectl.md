# EP 24: Instalación de Minikube y kubectl

**Tipo:** INSTALACION
## Objetivo
Instalar Minikube en local, kubectl, iniciar un cluster y verificar nodos.

## Pasos Detallados

### 1. Instalar kubectl
```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl
kubectl version --client
```

### 2. Instalar Minikube
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
rm minikube-linux-amd64
```

### 3. Iniciar Cluster
```bash
minikube start --driver=docker
# Esperar 2-3 minutos

# Verificar
kubectl get nodes
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   1m    v1.28.x
```

### Comandos Útiles de Minikube
```bash
minikube status          # Ver estado
minikube dashboard       # Abrir dashboard web
minikube stop            # Detener cluster
minikube delete          # Eliminar cluster
minikube service <svc>   # Abrir servicio en navegador
```

## Verificación
- [ ] `kubectl get nodes` muestra el nodo en estado "Ready"
- [ ] `minikube status` muestra todo "Running"

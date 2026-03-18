# EP 29: Conectar kubectl a EKS

**Tipo:** CONFIGURACION
## Objetivo
Configurar kubectl para conectarse al cluster EKS en la nube.

## Prerequisitos
- Cluster EKS creado (EP28)
- AWS CLI configurado (EP15)

## Pasos Detallados

### 1. Actualizar kubeconfig
```bash
aws eks update-kubeconfig --region us-east-1 --name curso-gitops-eks
# Output: Added new context arn:aws:eks:us-east-1:XXXX:cluster/curso-gitops-eks
```

### 2. Verificar Conexión
```bash
kubectl get nodes
# NAME                          STATUS   ROLES    AGE   VERSION
# ip-10-0-1-xxx.ec2.internal    Ready    <none>   5m    v1.29.x
# ip-10-0-2-xxx.ec2.internal    Ready    <none>   5m    v1.29.x

kubectl cluster-info
```

### 3. Ver Contextos (si tienes múltiples clusters)
```bash
kubectl config get-contexts
kubectl config use-context arn:aws:eks:us-east-1:XXXX:cluster/curso-gitops-eks
```

## Verificación
- [ ] `kubectl get nodes` muestra los nodos del cluster EKS
- [ ] El contexto apunta al cluster correcto

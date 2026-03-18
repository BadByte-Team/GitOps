# EP 26: Comandos Esenciales de kubectl

**Tipo:** PRACTICA
## Objetivo
Dominar los comandos de operación diaria con kubectl.

## Comandos

### Obtener Información
```bash
kubectl get pods -n curso-gitops               # Listar pods
kubectl get pods -o wide                        # Con más detalle
kubectl get all -n curso-gitops                 # Todo: pods, svc, deploy
kubectl describe pod <pod-name> -n curso-gitops # Detalle completo
```

### Logs y Debug
```bash
kubectl logs <pod-name> -n curso-gitops       # Ver logs
kubectl logs -f <pod-name>                     # Follow (tiempo real)
kubectl exec -it <pod-name> -n curso-gitops -- sh  # Entrar al contenedor
```

### Escalar
```bash
kubectl scale deployment curso-gitops -n curso-gitops --replicas=3
kubectl get pods -n curso-gitops  # Ahora hay 3 pods
```

### Eliminar
```bash
kubectl delete pod <pod-name> -n curso-gitops     # Eliminar pod (se recrea)
kubectl delete -f deployment.yaml                  # Eliminar por archivo
kubectl delete namespace curso-gitops              # Eliminar todo el namespace
```

### Resumen
| Comando | Acción |
|---|---|
| `get` | Listar recursos |
| `describe` | Detalle completo de un recurso |
| `logs` | Ver logs de un pod |
| `exec` | Ejecutar comando en un pod |
| `scale` | Cambiar número de réplicas |
| `delete` | Eliminar recursos |
| `apply -f` | Crear/actualizar desde YAML |

## Verificación
- [ ] Puedes escalar de 2 a 3 réplicas y ver los pods nuevos
- [ ] Puedes entrar a un pod con `exec`
- [ ] Puedes ver los logs en tiempo real

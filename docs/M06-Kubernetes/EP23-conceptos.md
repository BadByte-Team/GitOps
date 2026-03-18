# EP 23: Conceptos Clave — Pod, Deployment, Service

**Tipo:** TEORIA
## Objetivo
Entender la arquitectura básica de Kubernetes y la diferencia entre los objetos principales.

## Conceptos

### Arquitectura
```
┌─────────────────────────────────────────┐
│              CLUSTER                     │
│  ┌──────────┐  ┌──────────┐             │
│  │  Node 1  │  │  Node 2  │             │
│  │ ┌──────┐ │  │ ┌──────┐ │             │
│  │ │ Pod  │ │  │ │ Pod  │ │             │
│  │ │ ┌──┐ │ │  │ │ ┌──┐ │ │             │
│  │ │ │C │ │ │  │ │ │C │ │ │   C = Container
│  │ │ └──┘ │ │  │ │ └──┘ │ │             │
│  │ └──────┘ │  │ └──────┘ │             │
│  └──────────┘  └──────────┘             │
│          ▲                               │
│          │ Service (LoadBalancer)         │
│          │                               │
└──────────│───────────────────────────────┘
           │
        Internet
```

### Objetos Principales
| Objeto | Qué es |
|---|---|
| **Pod** | Unidad mínima — 1 o más contenedores |
| **Deployment** | Gestiona réplicas de Pods, rolling updates |
| **Service** | Expone Pods a la red (ClusterIP, NodePort, LoadBalancer) |
| **Namespace** | Aislamiento lógico dentro del cluster |
| **ConfigMap** | Variables de configuración |
| **Secret** | Datos sensibles (base64) |

### Pod vs Deployment
- **Pod**: Se muere y no se reinicia solo
- **Deployment**: Mantiene N réplicas vivas, autoheal si un Pod muere

### Tipos de Service
| Tipo | Acceso |
|---|---|
| `ClusterIP` | Solo dentro del cluster |
| `NodePort` | IP del nodo + puerto (30000-32767) |
| `LoadBalancer` | IP pública (cloud only) |

## Verificación
- [ ] Entiendes la jerarquía: Cluster → Node → Pod → Container
- [ ] Entiendes la diferencia entre Pod y Deployment
- [ ] Entiendes los tipos de Service

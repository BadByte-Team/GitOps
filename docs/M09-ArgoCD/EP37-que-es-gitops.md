# EP 37: Qué es GitOps y Cómo Funciona ArgoCD

**Tipo:** TEORIA
## Objetivo
Entender el concepto de GitOps y el loop de reconciliación de ArgoCD.

## GitOps — El Concepto

### La Regla de Oro
> **Git es la fuente de verdad.** El estado deseado de la infraestructura y las aplicaciones se define en un repositorio Git.

### Flujo GitOps
```
Dev hace push → Pipeline CI → Build imagen → Push a Docker Hub
                                     ↓
              Push tag nuevo al repo de manifiestos
                                     ↓
              ArgoCD detecta cambio → Sync a Kubernetes
```

### Dos Repositorios
| Repo | Contenido |
|---|---|
| **App repo** | Código fuente, Dockerfile, Jenkinsfile |
| **Manifests repo** | Kubernetes YAML (deployment, service, etc.) |

### ArgoCD — Loop de Reconciliación
```
┌──────────┐     compara     ┌──────────┐
│   Git    │ ◄──────────────►│   K8s    │
│ (desired)│                 │ (actual) │
└──────────┘                 └──────────┘
      ↑                           │
      │      si difiere           │
      │      ArgoCD aplica        │
      └───────────────────────────┘
```

Cada 3 minutos, ArgoCD compara el estado en Git con el estado real en el cluster. Si hay diferencia, **reconcilia automáticamente**.

## Verificación
- [ ] Entiendes que Git es la fuente de verdad
- [ ] Entiendes la separación de repos (app vs manifests)
- [ ] Entiendes el loop de reconciliación

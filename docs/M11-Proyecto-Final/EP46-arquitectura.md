# EP 46: Arquitectura Final — Visión Completa del Flujo

**Tipo:** TEORIA
## Objetivo
Recorrer todos los componentes del flujo GitOps y entender cómo se conectan.

## Arquitectura

```
┌─────────────┐     push      ┌─────────────┐
│  Developer  │──────────────▶│   GitHub     │
│  (VS Code)  │               │  (App Repo)  │
└─────────────┘               └──────┬───────┘
                                     │ webhook
                              ┌──────▼───────┐
                              │   Jenkins    │
                              │   (EC2)      │
                              │              │
                              │ 1. Checkout  │
                              │ 2. SonarQube │
                              │ 3. Build     │
                              │ 4. Trivy     │
                              │ 5. Push      │
                              └──┬───────┬───┘
                                 │       │
                    push tag     │       │ push image
                                 │       │
                    ┌────────────▼┐   ┌──▼──────────┐
                    │   GitHub    │   │  Docker Hub  │
                    │ (Manifests) │   │  (Registry)  │
                    └──────┬──────┘   └──────────────┘
                           │ detect change
                    ┌──────▼──────┐
                    │   ArgoCD    │
                    │  (in EKS)   │
                    │             │
                    │  reconcile  │
                    └──────┬──────┘
                           │ apply manifests
                    ┌──────▼──────┐
                    │    EKS      │
                    │  (K8s)      │
                    │             │
                    │ ┌─────────┐ │
                    │ │  Pods   │ │
                    │ │curso-   │ │
                    │ │gitops   │ │
                    │ └─────────┘ │
                    └─────────────┘
```

### Componentes y Quién los Provisionó
| Componente | Creado con | Episodio |
|---|---|---|
| EC2 Jenkins | Terraform | EP22 |
| SonarQube | Docker Compose | EP43 |
| Trivy | Script bash | EP42 |
| EKS Cluster | Terraform | EP28 |
| ArgoCD | kubectl apply | EP38 |
| App en EKS | ArgoCD (auto) | EP40 |

## Verificación
- [ ] Puedes explicar el flujo completo de un push a producción
- [ ] Entiendes qué hace cada componente

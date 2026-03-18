# EP 27: Qué es EKS y Cuándo Usarlo

**Tipo:** TEORIA
## Objetivo
Entender la diferencia entre un cluster Kubernetes self-managed y Amazon EKS.

## EKS — Elastic Kubernetes Service

### ¿Qué es?
- Servicio administrado de AWS para correr Kubernetes
- AWS gestiona el **Control Plane** (API server, etcd, scheduler)
- Tú solo gestionas los **Worker Nodes** (donde corren tus pods)

### Self-managed vs EKS
| Aspecto | Self-managed | EKS |
|---|---|---|
| Control Plane | Tú lo instalas y mantienes | AWS lo gestiona |
| Actualizaciones | Manuales | AWS las ofrece |
| Alta disponibilidad | Tú la configuras | Incluida (multi-AZ) |
| Costo | Solo EC2 | $0.10/hr (~$72/mes) + nodos |
| Complejidad | Alta | Media |

### ¿Cuándo usar EKS?
- ✅ Producción con alta disponibilidad
- ✅ Equipos que no quieren administrar etcd
- ✅ Integración nativa con AWS (ALB, IAM, CloudWatch)
- ❌ Para aprender (usa Minikube)
- ❌ Proyectos con presupuesto muy limitado

### Costos Aproximados del Curso
| Recurso | Costo/hora | Si lo usas 4 horas |
|---|---|---|
| EKS Control Plane | $0.10 | $0.40 |
| 2x t3.medium nodes | $0.0416 × 2 | $0.33 |
| **Total** | | **~$0.73** |

## Verificación
- [ ] Entiendes que EKS gestiona el control plane
- [ ] Entiendes los costos involucrados
- [ ] Sabes que debes destruir el cluster al terminar cada sesión

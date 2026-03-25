# EP 46: Arquitectura Final Híbrida

**Tipo:** TEORÍA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Explicación del diagrama completo: PC Local (Jenkins/SonarQube/Trivy) → Docker Hub / GitHub gitops-infra → AWS EC2 Gratuita con K3s (ArgoCD) → Pods.

---

## 🎯 Objetivo
Recorrer todos los componentes del flujo GitOps híbrido que construimos durante el curso, entender cómo se conectan entre sí, y contrastar la arquitectura gratuita con la arquitectura convencional de pago — cerrando el círculo de lo aprendido antes del proyecto final.

---

## 📋 Prerequisitos
- Pipeline CI/CD completamente funcional (EP36, EP41, EP45)
- ArgoCD sincronizando `gitops-infra` con K3s (EP40)

---

## 🧠 El diagrama de arquitectura completo

```
┌─────────────────────────────────────────────────────────────┐
│                    TU PC LOCAL                               │
│                                                              │
│  ┌──────────┐   ┌────────────┐   ┌──────────┐   ┌───────┐ │
│  │ VS Code  │   │  Jenkins   │   │SonarQube │   │ Trivy │ │
│  │(código)  │   │  (Docker)  │   │ (Docker) │   │(local)│ │
│  └────┬─────┘   └─────┬──────┘   └──────────┘   └───────┘ │
└───────┼───────────────┼──────────────────────────────────────┘
        │ git push      │ 1. checkout
        ▼               │ 2. sonarqube
┌───────────────┐       │ 3. docker build
│    GitHub     │       │ 4. trivy scan
│  ┌──────────┐ │       │ 5. docker push ──────▶ Docker Hub
│  │gitops-app│ │       │ 6. git push ──────────▶ gitops-infra
│  └──────────┘ │       │
│  ┌──────────┐ │◀──────┘
│  │gitops-   │ │
│  │infra     │ │
│  └─────┬────┘ │
└────────┼──────┘
         │ ArgoCD detecta cambio (cada 3 min)
         ▼
┌─────────────────────────────────────────────────────────────┐
│         AWS EC2 t3.micro — Free Tier ($0/mes)               │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                  K3s Cluster                          │   │
│  │                                                       │   │
│  │  [argocd] ArgoCD ─── observa gitops-infra            │   │
│  │                │                                      │   │
│  │                ▼ aplica manifiestos                   │   │
│  │  [curso-gitops]                                       │   │
│  │  ┌──────────────────┐   ┌──────────────────────────┐ │   │
│  │  │  MySQL Pod        │   │  App Go Pod              │ │   │
│  │  │  mysql:8.0        │   │  curso-gitops:N-xxxxxx   │ │   │
│  │  │  :3306 interno    │   │  NodePort: 30081         │ │   │
│  │  └──────────────────┘   └──────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────┘   │
│  ArgoCD UI → :30080    App → :30081                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 📦 Inventario completo: componentes y costos

| Componente | Dónde vive | Tecnología | Costo |
|---|---|---|---|
| Edición de código | PC Local | VS Code | $0 |
| CI — Build, Scan | PC Local | Jenkins (Docker Compose) | $0 |
| Calidad de código | PC Local | SonarQube (Docker Compose) | $0 |
| Seguridad de imagen | PC Local | Trivy (binario) | $0 |
| Repositorio app | GitHub | Git | $0 |
| Repositorio infra | GitHub | Git | $0 |
| Registro de imágenes | Docker Hub | Docker | $0 |
| IaC | Local → AWS | Terraform | $0 |
| Servidor nube | AWS | EC2 t3.micro | $0 (Free Tier) |
| Cluster Kubernetes | EC2 t3.micro | K3s | $0 |
| Operador GitOps | K3s | ArgoCD | $0 |
| Base de datos | K3s | MySQL Pod | $0 |
| **Total mensual** | | | **$0** |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: la tabla de costos de la arquitectura original — Jenkins EC2 t2.medium, EKS, ALB — sumando más de $125/mes. Al lado, la tabla de la arquitectura del curso con una sola columna de ceros.*

"Bienvenidos al episodio 46. Bienvenidos al Módulo 11 — el último módulo del curso.

Esta es la arquitectura que la mayoría de los cursos de GitOps proponen: Jenkins en una EC2 dedicada, EKS para el cluster, un Application Load Balancer para exponer los servicios. El total supera los $125 al mes. Para aprender.

Y esta es nuestra arquitectura: doce componentes, doce ceros. No porque hayamos simplificado el flujo ni sacrificado conceptos — el flujo GitOps es exactamente el mismo. Sino porque tomamos decisiones de diseño inteligentes en cada punto donde había una alternativa gratuita válida.

Hoy recorremos esas decisiones una por una. Por qué Jenkins local en lugar de EC2. Por qué K3s en lugar de EKS. Por qué NodePort en lugar de LoadBalancer. Cuál es el costo real de esas elecciones en términos de limitaciones. Y qué aprendiste que es directamente transferible a cualquier entorno profesional.

Es el episodio de cierre antes del proyecto final. Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — Las tres decisiones clave (1:30 – 4:30)

> *Pantalla: diagrama de la arquitectura con cada decisión destacada.*

"La arquitectura gratuita del curso descansa en tres decisiones de diseño. Quiero que las entiendan porque no son atajos — son compensaciones conscientes que tienen sentido en el contexto de aprendizaje.

---

**Decisión 1: Jenkins local en lugar de EC2**

La decisión más impactante en costo: $0 vs $33/mes.

La compensación: Jenkins solo corre cuando tu laptop está encendida. En producción real, el servidor de CI necesita estar disponible 24/7 para detectar pushes y ejecutar pipelines automáticamente. En un entorno de aprendizaje, ejecutar el pipeline manualmente o con la laptop abierta es completamente aceptable.

Lo que no cambia: el Jenkinsfile es idéntico. Los plugins, las credenciales, los stages, el patrón GitOps — todo igual. Si mañana llevas este conocimiento a una empresa con Jenkins en AWS, el archivo no cambia.

---

**Decisión 2: K3s en lugar de EKS**

$0 vs $72/mes solo por el control plane.

K3s es Kubernetes certificado por la CNCF. La misma API, los mismos manifiestos YAML, la misma compatibilidad con ArgoCD. La compensación: control plane y nodo de trabajo en la misma instancia t3.micro, sin multi-AZ, sin alta disponibilidad automática. Para aprender GitOps, eso no importa.

Lo que no cambia: los manifiestos YAML de Kubernetes que escribiste funcionan en EKS, GKE, o AKS sin cambiar una sola línea. El conocimiento es completamente portátil.

---

**Decisión 3: NodePort en lugar de LoadBalancer**

$0 vs $20/mes por el ALB.

NodePort expone el servicio directamente en un puerto de la instancia EC2. La compensación: sin balanceo de carga externo, sin failover automático, sin SSL termination gestionada. Para un entorno de aprendizaje con una instancia, NodePort es suficiente.

Lo que no cambia: el archivo `service.yaml`. Cambiar de NodePort a LoadBalancer en EKS es editar una línea — `type: LoadBalancer` — y ArgoCD lo aplica en el siguiente sync."

---

### El flujo completo, paso por paso (4:30 – 8:30)

> *Pantalla: diagrama animado o slide con las flechas del flujo.*

"Recorramos el flujo completo que construimos durante el curso, de principio a fin.

**Paso 1: el desarrollador hace un cambio.**

VS Code en local. Edita el código Go, hace `git commit` y `git push` a `gitops-app`. Ese push es el disparador de todo lo que sigue.

---

**Paso 2: Jenkins ejecuta el pipeline de CI.**

Jenkins detecta el push — manualmente con Build Now, o automáticamente si hay webhook configurado. El pipeline corre los seis stages en orden:

Checkout clona `gitops-app` y genera el `BUILD_TAG` — el número de build más el hash del commit.

SonarQube analiza el código Go y evalúa el Quality Gate. Si el código no cumple los criterios de calidad, el pipeline se detiene aquí.

Docker Build construye la imagen con el `BUILD_TAG`. La imagen incluye exactamente el código del commit que disparó el pipeline.

Trivy escanea la imagen en busca de CVEs de severidad HIGH y CRITICAL. El reporte aparece en el Console Output para revisión.

Docker Push sube la imagen a Docker Hub. Solo llega aquí código que pasó SonarQube y Trivy.

Deploy to GitOps Repo clona `gitops-infra`, actualiza el tag en `deployment.yaml` con `sed`, y hace push. Ese commit en `gitops-infra` es lo que activa la siguiente etapa.

---

**Paso 3: ArgoCD detecta el cambio.**

ArgoCD corre en el cluster K3s de la EC2. Cada 3 minutos compara el estado de `gitops-infra` con el estado del cluster. Cuando detecta que el `deployment.yaml` tiene un tag nuevo, sincroniza.

---

**Paso 4: K3s aplica el rolling update.**

Kubernetes crea un pod nuevo con la imagen actualizada. Espera a que pase el readiness check. Elimina el pod viejo. Sin downtime. La app nunca dejó de estar disponible durante el despliegue."

---

### Comparativa final: qué es igual, qué es diferente (8:30 – 10:30)

> *Pantalla: tabla comparativa.*

"Para cerrar con claridad absoluta:

| Aspecto | Arquitectura de pago | Arquitectura del curso |
|---|---|---|
| Jenkins | EC2 t2.medium — $33/mes | PC Local — $0 |
| Kubernetes | EKS — $72/mes | K3s en t3.micro — $0 |
| LoadBalancer | AWS ALB — $20/mes | NodePort — $0 |
| **Costo total** | **~$125+/mes** | **$0** |
| Jenkinsfile | Idéntico | Idéntico |
| Manifiestos YAML | Idénticos | Idénticos |
| kubectl | Idéntico | Idéntico |
| ArgoCD | Idéntico | Idéntico |
| Patrón GitOps | Idéntico | Idéntico |

La diferencia está en el costo y en las limitaciones operacionales — disponibilidad 24/7, alta disponibilidad, balanceo de carga. La diferencia no está en el conocimiento, en los archivos, ni en los comandos.

Todo lo que aprendiste en este curso funciona en producción. El Jenkinsfile que escribiste puede correr en Jenkins en AWS. Los manifiestos YAML que escribiste pueden desplegarse en EKS. ArgoCD funciona igual en cualquier cluster certificado por la CNCF.

Eso es exactamente lo que queríamos lograr."

---

### CIERRE (10:30 – 12:00)

"Eso es el EP46.

Hemos recorrido la arquitectura completa: doce componentes, doce ceros, un flujo GitOps de nivel profesional.

Los últimos tres episodios del curso son prácticos:

**EP47** — La base de datos separada y los manifiestos de Kubernetes. Vamos a revisar en detalle por qué MySQL no va en el mismo Dockerfile que la app Go, y a recorrer los siete manifiestos que ArgoCD despliega en el cluster.

**EP48** — El pipeline CI/CD en acción. El run final del pipeline completo con todos los stages de seguridad activos, viendo el flujo completo de extremo a extremo una última vez.

**EP49** — La limpieza de recursos. `terraform destroy` para dejar la cuenta de AWS en cero, y el cierre del curso.

Nos vemos en el EP47."

---

## ✅ Checklist de Verificación
- [ ] Puedes explicar las tres decisiones de diseño y sus compensaciones
- [ ] Puedes recorrer el flujo completo de un push a producción de memoria
- [ ] Entiendes qué cambia y qué no cambia entre la arquitectura del curso y EKS
- [ ] Sabes que los manifiestos YAML y el Jenkinsfile son directamente portables

---

## 🗒️ Notas de Producción
- Abrir con las dos tablas de costos lado a lado — el contraste visual es el gancho más fuerte del episodio.
- Las tres decisiones de diseño pueden presentarse como slides separados — una por una, con tiempo para que el alumno procese la compensación antes de pasar a la siguiente.
- Al recorrer el flujo completo, señalar con el cursor en el diagrama cada componente mientras lo describes — la conexión visual refuerza la comprensión.
- La tabla comparativa final merece quedarse en pantalla durante 15-20 segundos mientras se lee en voz alta — es el mensaje de cierre del módulo más importante.
- Anunciar los tres episodios finales con sus nombres da al alumno el mapa completo del final del curso.

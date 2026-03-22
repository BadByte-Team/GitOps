# EP 37: ¿Qué es GitOps y Cómo Funciona ArgoCD?

**Tipo:** TEORÍA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Entender el concepto de GitOps, cómo ArgoCD implementa el loop de reconciliación, y qué rol cumple en la arquitectura del curso — cerrando el ciclo que comenzó con el pipeline de Jenkins en el EP36.

---

## 📋 Prerequisitos
- Pipeline CI funcionando (EP36)
- K3s corriendo en la EC2 (EP29)
- El repositorio `gitops-infra` con el `deployment.yaml` actualizado por Jenkins

---

## 🧠 GitOps: la idea central

**Git como fuente de verdad para la infraestructura.** El estado deseado del cluster vive en un repositorio. Cualquier cambio al cluster pasa por un commit.

```
Estado deseado → gitops-infra (Git)
Estado actual  → K3s (Kubernetes)

ArgoCD         → compara los dos → reconcilia las diferencias
```

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: el flujo completo del EP36 en el diagrama — Jenkins actualizó `gitops-infra`, el commit es visible en GitHub. Pero el cluster K3s todavía tiene la versión antigua.*

"Bienvenidos al episodio 37. Bienvenidos al Módulo 09 — ArgoCD.

En el episodio anterior terminamos con algo interesante. Jenkins construyó la imagen, la subió a Docker Hub, y actualizó el `deployment.yaml` en `gitops-infra` con el nuevo tag. Ese commit está en GitHub ahora mismo.

Pero si voy al cluster K3s y pregunto qué imagen está corriendo... todavía es la versión anterior. El commit existe en Git, pero nadie lo ha aplicado al cluster todavía.

Esa es exactamente la brecha que ArgoCD cierra.

ArgoCD es un operador de Kubernetes que vive dentro del cluster y tiene un trabajo muy específico: observar un repositorio de Git, y cuando detecta que el estado deseado en Git difiere del estado actual del cluster, aplicar los cambios automáticamente.

Es la pieza que conecta el pipeline de CI con el cluster de producción. Sin ArgoCD, el flujo GitOps no está completo. Con ArgoCD, un `git push` se convierte en un despliegue automático.

Hoy entendemos el concepto. Mañana — los siguientes episodios — lo instalamos y lo configuramos.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — ¿Por qué 'GitOps'? (1:30 – 4:00)

> *Pantalla: diagrama comparando el flujo tradicional CI/CD vs el flujo GitOps.*

"El término 'GitOps' fue acuñado por Alexis Richardson, el fundador de Weaveworks, en 2017. La idea central es simple pero tiene implicaciones profundas.

En un pipeline CI/CD tradicional, el proceso es:
```
código → tests → build → [Jenkins aplica directamente al servidor] → producción
```
Jenkins o el pipeline de CI se conecta directamente al servidor de producción, ejecuta `kubectl apply`, y el cambio ocurre. Funciona, pero tiene un problema: el estado real del servidor puede divergir de lo que cualquiera esperaría. Alguien puede hacer un cambio manual en el cluster que nadie registró. Un deployment puede fallar a mitad y quedar en un estado inconsistente.

En GitOps, el proceso es diferente:
```
código → tests → build → [pipeline actualiza Git] → [operador aplica desde Git] → producción
```
El pipeline nunca toca el servidor directamente. Solo actualiza un archivo en Git. Un operador — ArgoCD en nuestro caso — observa ese repositorio y aplica los cambios.

¿Por qué importa esa diferencia?

**Auditoría completa.** Cada cambio en producción tiene un commit de Git. Puedes ver quién cambió qué, cuándo, y por qué. El `git log` de `gitops-infra` es el historial de producción.

**Rollback trivial.** Si una nueva versión tiene un bug crítico, el rollback es `git revert` + `git push`. ArgoCD aplica automáticamente la versión anterior. Sin comandos especiales de Kubernetes, sin scripts de emergencia.

**Drift detection.** Si alguien hace un cambio manual en el cluster — `kubectl edit deployment` directamente — ArgoCD lo detecta como una discrepancia entre el estado deseado (Git) y el estado real (cluster). Con `selfHeal: true`, lo revierte automáticamente.

Esos tres beneficios juntos son los que hacen que GitOps sea el estándar de la industria para despliegues en entornos críticos."

---

### El loop de reconciliación de ArgoCD (4:00 – 7:00)

> *Pantalla: diagrama del loop de reconciliación.*

"ArgoCD implementa GitOps a través de lo que se llama el **loop de reconciliación**. Es un proceso que corre constantemente dentro del cluster.

```
                    ┌─────────────────────┐
                    │   gitops-infra      │
                    │   (GitHub)          │
                    │   deployment.yaml   │
                    │   image: app:1-abc  │  ← estado deseado
                    └──────────┬──────────┘
                               │ cada 3 minutos
                               │ ArgoCD compara
                               ▼
                    ┌─────────────────────┐
                    │   K3s Cluster       │
                    │   (EC2 t2.micro)    │
                    │   Pod corriendo:    │
                    │   image: app:0-xyz  │  ← estado actual
                    └──────────┬──────────┘
                               │
                    ¿Son iguales? NO
                               │
                               ▼
                    ArgoCD aplica el cambio:
                    kubectl apply -f deployment.yaml
                               │
                               ▼
                    Pod nuevo: image: app:1-abc ✅
```

El loop ocurre cada 3 minutos por defecto. También se puede disparar manualmente desde la UI o la CLI de ArgoCD, o automáticamente con un webhook de GitHub.

Hay tres configuraciones del sync policy que usaremos:

**`automated`** — ArgoCD sincroniza automáticamente cuando detecta diferencias. Sin intervención manual.

**`prune: true`** — Si un recurso existe en el cluster pero no en el repositorio, ArgoCD lo elimina. Garantiza que el cluster refleja exactamente lo que está en Git.

**`selfHeal: true`** — Si alguien modifica un recurso manualmente en el cluster, ArgoCD lo revierte al estado de Git. El cluster nunca puede divergir del repositorio."

---

### ArgoCD en la arquitectura del curso (7:00 – 9:30)

> *Pantalla: diagrama completo de la arquitectura híbrida del curso.*

"Vamos a ver exactamente cómo encaja ArgoCD en la arquitectura que hemos construido.

```
Tu PC Local                          GitHub
├── gitops-app ──── git push ──────▶ gitops-app
└── Jenkins                          gitops-infra ◀─────────────────┐
    1. Checkout gitops-app                 │                         │
    2. SonarQube                           │ ArgoCD observa          │
    3. docker build                        │ cada 3 min              │
    4. docker push ──────▶ Docker Hub      │                         │
    5. clona gitops-infra                  │                         │
    6. sed actualiza tag                   │                         │
    7. git push ────────────────────────────────────────────────────┘

                          AWS EC2 t2.micro
                          K3s Cluster
                          └── namespace: argocd
                              └── ArgoCD
                                  └── observa gitops-infra
                                      └── aplica cambios a:
                                          namespace: curso-gitops
                                          ├── MySQL Pod
                                          └── App Go Pod ◀── imagen de Docker Hub
```

El flujo completo:
1. Developer hace push a `gitops-app`
2. Jenkins construye la imagen y la sube a Docker Hub
3. Jenkins actualiza `deployment.yaml` en `gitops-infra` con el nuevo tag
4. ArgoCD detecta el cambio en `gitops-infra`
5. ArgoCD aplica el `deployment.yaml` actualizado al cluster K3s
6. K3s descarga la nueva imagen de Docker Hub y actualiza los pods

En ningún momento Jenkins habla con Kubernetes directamente. Esa separación es el corazón del patrón GitOps."

---

### Los objetos de ArgoCD (9:30 – 11:30)

> *Pantalla: VS Code con el archivo `application.yaml` del proyecto en `gitops-infra/infrastructure/kubernetes/argocd/`.*

"ArgoCD extiende Kubernetes con sus propios tipos de objetos. El más importante es la **Application**.

Una Application de ArgoCD define:
- **Qué repositorio observar** — `gitops-infra`
- **Qué directorio** dentro del repositorio — `infrastructure/kubernetes/app/`
- **A qué cluster aplicar** — el K3s local del cluster
- **En qué namespace** — `curso-gitops`
- **Cómo sincronizar** — automático, con prune y selfHeal

Abrimos el `application.yaml` que ya existe en `gitops-infra`:"

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: curso-gitops
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/TU_USUARIO_GITHUB/gitops-infra.git
    targetRevision: main
    path: infrastructure/kubernetes/app
  destination:
    server: https://kubernetes.default.svc
    namespace: curso-gitops
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

"Este archivo es lo que crearemos en el EP40 para conectar ArgoCD con el repositorio. El `https://kubernetes.default.svc` es la dirección interna del API server de Kubernetes — ArgoCD lo usa para aplicar los manifiestos al mismo cluster donde él mismo corre."

---

### CIERRE (11:30 – 12:30)

"Eso es el episodio 37.

GitOps: Git como fuente de verdad. ArgoCD: el operador que cierra el loop entre el repositorio y el cluster. El loop de reconciliación: comparar el estado deseado en Git con el estado actual cada 3 minutos y aplicar las diferencias.

En los próximos cuatro episodios convertimos esta teoría en realidad:
- EP38: instalar ArgoCD en K3s
- EP39: exponerlo con NodePort en el puerto 30080
- EP40: conectarlo a `gitops-infra` y crear la Application
- EP41: ver el flujo completo funcionando — un push en `gitops-app` → imagen nueva → pod actualizado

Nos vemos en el EP38."

---

## ✅ Checklist de Verificación
- [ ] Entiendes la diferencia entre CI/CD tradicional y el patrón GitOps
- [ ] Puedes explicar el loop de reconciliación de ArgoCD
- [ ] Entiendes qué hacen `automated`, `prune` y `selfHeal`
- [ ] Puedes identificar el rol de ArgoCD en el diagrama de arquitectura del curso

---

## 🗒️ Notas de Producción
- La apertura con el commit de Jenkins en GitHub pero el cluster todavía desactualizado es el gancho perfecto — crea la tensión que ArgoCD va a resolver.
- El diagrama comparando CI/CD tradicional vs GitOps es el momento más pedagógico del episodio — dibujarlo lentamente mientras se describe verbalmente.
- El diagrama de la arquitectura completa del curso con ArgoCD integrado puede presentarse como slide — dejar tiempo para que el alumno lo procese.
- Al mostrar el `application.yaml`, señalar con el cursor cada campo mientras lo describes — es un YAML relativamente corto pero denso en significado.
- Anunciar explícitamente los cuatro episodios que vienen — el alumno necesita saber el mapa completo de lo que falta para ver el flujo completo.

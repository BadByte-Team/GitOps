# EP 27: Alternativas a EKS — ¿Qué es K3s?

**Tipo:** TEORÍA
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Este módulo ya no cubre EKS. Explica por qué K3s es la alternativa profesional gratuita para el curso.

---

## 🎯 Objetivo
Entender por qué EKS no es viable para este curso, qué es K3s y por qué es una alternativa técnicamente válida y usada en producción real, y qué vamos a construir en los próximos cuatro episodios sobre la EC2 del EP22.

---

## 📋 Prerequisitos
- EC2 t3.micro corriendo (EP22)
- Conceptos de Kubernetes dominados (EP23)
- Minikube detenido: `minikube stop`

---

## 🧠 El problema con EKS

EKS cobra por el Control Plane independientemente del uso:

| Recurso | Costo/hora | Costo/mes (24x7) |
|---|---|---|
| EKS Control Plane | $0.10 | **~$72 USD** |
| 1× t3.medium (nodo) | $0.0416 | ~$30 USD |
| **Total mínimo** | | **~$102 USD/mes** |

Con la capa gratuita de AWS y K3s, el total es **$0**.

---

## 🧠 La solución: K3s

| Característica | Kubernetes estándar / EKS | K3s |
|---|---|---|
| Binario | ~100 MB (múltiples componentes) | ~70 MB (todo en uno) |
| RAM mínima | ~2 GB | ~512 MB |
| Instalación | `kubeadm` + etcd + múltiples pasos | Un solo `curl` |
| Certificación CNCF | ✅ | ✅ |
| API de Kubernetes | Completa | Completa |
| Compatible con kubectl | ✅ | ✅ |
| Compatible con ArgoCD | ✅ | ✅ |
| Costo del control plane | ~$72/mes (EKS) | **$0** |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: calculadora de precios de AWS — aws.amazon.com/eks/pricing — con el costo del control plane visible. $0.10 por hora.*

"Bienvenidos al episodio 27. Bienvenidos al Módulo 07.

Acabamos de terminar seis episodios de Kubernetes. Aprendimos qué son los Pods, los Deployments y los Services. Practicamos con Minikube localmente. Dominamos los comandos de kubectl. Tenemos la base completa.

Ahora llega el momento de llevar todo eso a un cluster real en la nube. Y aquí es donde la mayoría de los cursos de DevOps te dicen: 'vamos a crear un cluster en EKS'.

Miren esto. $0.10 por hora, solo por el control plane de EKS. Sin haber levantado un solo nodo todavía. Eso son $72 al mes solo por tener el cluster existiendo — antes de correr ninguna aplicación. Un nodo t3.medium son otros $30. Estamos en más de cien dólares al mes como mínimo, y eso si destruyes el cluster al terminar cada sesión de práctica y lo vuelves a crear desde cero cada vez.

Para alguien que está aprendiendo, eso es una barrera completamente innecesaria. No deberías tener que pagar $100 al mes para aprender cómo funciona GitOps.

Por eso en este curso no usamos EKS. Usamos **K3s** — y hoy vamos a entender exactamente qué es, por qué es técnicamente serio, y qué vamos a construir en los próximos episodios sobre el servidor que ya tenemos corriendo en AWS.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — ¿Qué es K3s y quién lo usa? (1:30 – 4:30)

> *Pantalla: página oficial de K3s en k3s.io, mostrando el badge de certificación CNCF.*

"K3s es una distribución de Kubernetes creada por Rancher Labs — ahora parte de SUSE — y certificada por la **CNCF**: la Cloud Native Computing Foundation. Es la misma organización que certifica Kubernetes oficial, Prometheus, Helm, y la mayoría de las herramientas cloud-native que se usan en producción.

Esa palabra, certificada, importa mucho. No es un fork no oficial. No es una versión simplificada que sacrifica funcionalidad. Es Kubernetes completo, empaquetado de forma más eficiente para entornos con recursos limitados.

La API es la misma. Los manifiestos YAML son los mismos. kubectl funciona exactamente igual con K3s que con EKS. ArgoCD funciona exactamente igual. Si escribes un `deployment.yaml` hoy para K3s y mañana lo quieres desplegar en GKE o EKS, no cambias ni una línea.

---

La diferencia técnica está en cómo está construido internamente. Kubernetes estándar distribuye sus componentes en múltiples binarios: el API server corre por un lado, el scheduler por otro, el controller manager por otro, etcd separado. Cada uno es un proceso independiente.

K3s empaqueta todo eso en un único binario de ~70 MB. El resultado es que puede correr en una máquina con 512 MB de RAM, frente a los 2 GB que necesita Kubernetes estándar. Para nuestra EC2 t3.micro con 1 GB de RAM física más 2 GB de Swap, eso es exactamente lo que necesitamos.

---

Una pregunta que siempre aparece en este punto: '¿K3s es de juguete? ¿Lo usa alguien en producción de verdad?'

La respuesta es no, y sí respectivamente.

K3s es la distribución elegida para **edge computing** — los dispositivos que necesitan correr Kubernetes en condiciones de hardware limitado. Chick-fil-A usa K3s en los sistemas de punto de venta de sus restaurantes. Tesla lo usa en sus vehículos. Organismos como la NASA lo usan en sistemas embebidos de misiones espaciales.

Para nuestro curso, K3s corre en una EC2 t3.micro y nos da exactamente lo mismo que EKS: pods, deployments, services, namespaces, y soporte completo para ArgoCD. La diferencia es que no pagamos nada."

---

### La arquitectura de lo que vamos a construir (4:30 – 7:30)

> *Pantalla: diagrama ASCII o slide de la arquitectura completa del módulo.*

"Antes de arrancar con los comandos, quiero que vean el mapa completo de lo que construiremos en este módulo y cómo encaja en el stack final del curso.

Todo vive en la misma EC2 t3.micro que creamos en el EP22:"

```
EC2 t3.micro — Ubuntu 22.04 — Free Tier AWS
│
├── 2 GB Swap               ← EP28 — sin esto K3s colapsa
│
└── K3s Cluster             ← EP29 — un solo comando curl
    │
    ├── namespace: kube-system
    │   └── CoreDNS, Traefik   (componentes del sistema — K3s los instala)
    │
    ├── namespace: argocd      ← EP38 — lo instalaremos ahí
    │   └── ArgoCD (~7 pods)   ← el operador que observa gitops-infra
    │
    └── namespace: curso-gitops   ← EP47 — ArgoCD lo despliega
        ├── MySQL Pod             ← base de datos
        └── App Go Pod            ← nuestra aplicación
```

"Tres episodios para tener todo esto funcionando: EP28 para el Swap, EP29 para instalar K3s, EP30 para configurar el acceso remoto desde tu máquina local.

Después de eso, los episodios de Jenkins y ArgoCD construyen encima de este cluster. Cuando llegues al EP40, ArgoCD va a observar el repositorio `gitops-infra` y va a desplegar automáticamente el namespace `curso-gitops` con todos sus recursos. Ese momento es el payoff del curso completo — y empieza aquí, en este módulo."

---

### Comparativa completa: EKS vs K3s (7:30 – 9:30)

> *Pantalla: tabla comparativa en slide o terminal.*

"Para que quede absolutamente claro qué cambia y qué no entre la arquitectura original con EKS y la nuestra con K3s:"

| Aspecto | EKS (original) | K3s (nuestro) |
|---|---|---|
| Cluster | EKS en AWS (~$72/mes) | K3s en EC2 t3.micro ($0) |
| Control Plane | Gestionado por AWS | En la misma EC2 |
| Nodos | 2× t3.medium separados | La misma EC2 t3.micro |
| LoadBalancer | AWS ALB (~$20/mes) | NodePort (gratis) |
| Configurar kubectl | `aws eks update-kubeconfig` | Copiar y editar `k3s.yaml` |
| kubectl | Idéntico | Idéntico |
| ArgoCD | Idéntico | Idéntico |
| Manifiestos YAML | Idénticos | Idénticos |
| Tiempo de setup | ~20 min (Terraform) | ~3 min (un `curl`) |
| **Costo total/mes** | **~$102** | **$0** |

"Las dos últimas filas son las que importan para el alumno. Todo lo intermedio — kubectl, ArgoCD, los YAMLs — es exactamente igual. El conocimiento es completamente transferible. Si algún día trabajan en una empresa que usa EKS, todo lo que aprenden aquí aplica directamente. Solo cambia el comando para configurar el kubeconfig."

---

### Lo que NO cambia (9:30 – 10:30)

> *Pantalla: VS Code con `deployment.yaml` del proyecto abierto.*

"Quiero mostrar algo concreto para que esto no sea solo teoría.

Este es el `deployment.yaml` de nuestra app Go en `gitops-infra/infrastructure/kubernetes/app/`. Lo escribimos siguiendo exactamente el estándar de Kubernetes.

Este archivo es el que ArgoCD va a aplicar al cluster K3s en el EP40. Si mañana deciden llevar este proyecto a EKS en su trabajo, este mismo archivo funciona sin cambiar absolutamente nada. Solo le dices a ArgoCD que el cluster destino es diferente — una línea en la configuración de ArgoCD.

El campo que van a ver cambiar durante el curso es este:"

```yaml
containers:
- name: curso-gitops
  image: TU_USUARIO/curso-gitops:latest   ← Jenkins actualiza esta línea
```

"Cuando el pipeline de Jenkins termine en el EP36, va a clonar `gitops-infra`, cambiar ese tag con `sed`, y hacer push. ArgoCD va a detectar el cambio y aplicar el nuevo deployment al cluster K3s. Ese flujo — el corazón de GitOps — empieza a construirse en los episodios de este módulo."

---

### CIERRE (10:30 – 11:30)

"Eso es el episodio 27.

K3s es Kubernetes certificado, ligero, gratuito, y usado en producción en empresas y proyectos reales. No es una alternativa de segunda clase — es la elección correcta cuando los recursos son limitados y el costo importa.

En el siguiente episodio hacemos el primer paso práctico sobre el servidor: configurar 2 GB de Swap en la EC2. Sin este paso, cuando K3s, ArgoCD y MySQL comiencen a correr simultáneamente, el kernel de Linux va a matar procesos porque se queda sin memoria. Es un paso de Linux puro — sin Kubernetes todavía — pero es completamente necesario para que todo lo demás funcione.

Nos vemos en el EP28."

---

## ✅ Checklist de Verificación
- [ ] Entiendes por qué EKS cuesta ~$72/mes solo por el Control Plane
- [ ] Entiendes que K3s está certificado por la CNCF — no es un fork no oficial
- [ ] Sabes que kubectl y los manifiestos YAML son idénticos en K3s y EKS
- [ ] Puedes describir la arquitectura de lo que construiremos en EP28–30
- [ ] La EC2 t3.micro del EP22 está corriendo (`running`)

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| La EC2 no responde | `aws ec2 describe-instances --filters "Name=tag:Name,Values=Produccion-K3s" --query "Reservations[0].Instances[0].State.Name"` |
| Olvidé la IP de la EC2 | `cd gitops-infra/infrastructure/terraform/jenkins-ec2 && terraform output prod_public_ip` |
| Minikube sigue consumiendo RAM | `minikube stop` |

---

## 🗒️ Notas de Producción
- Abrir la calculadora de precios de AWS en vivo en el navegador — no usar un screenshot. Ver el número en tiempo real tiene más impacto.
- La página de k3s.io muestra el badge de certificación CNCF prominentemente — señalarlo con el cursor mientras se menciona.
- La tabla comparativa EKS vs K3s puede presentarse como slide — dejarla en pantalla 15 segundos mientras se lee en voz alta.
- Al abrir el `deployment.yaml`, señalar con el cursor el campo `image:` y mencionar que Jenkins lo va a modificar en el EP36 — crear la continuidad narrativa hacia los módulos siguientes.
- Este es un episodio de teoría y motivación — el ritmo debe ser más pausado que los episodios técnicos. Los números y las comparativas necesitan tiempo para que el alumno los procese.

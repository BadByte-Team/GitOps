# EP 49: Limpieza de Recursos — Terraform Destroy

**Tipo:** PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Usa `terraform destroy` para eliminar la EC2 y opcionalmente el Backend S3, dejando la cuenta de AWS en $0.

---

## 🎯 Objetivo
Eliminar todos los recursos de AWS creados durante el curso — la EC2 con K3s y opcionalmente el backend de Terraform — verificar que la cuenta queda en cero, y cerrar el curso con el resumen de lo que construiste y lo que te llevas.

---

## 📋 Prerequisitos
- Terraform instalado localmente (EP18)
- AWS CLI configurado (EP15)
- `gitops-infra` con los archivos de Terraform del EP22

---

## ⚠️ Orden de eliminación

```
1. Detener Jenkins y SonarQube locales (opcional — no generan costo)
2. Destruir la EC2 con K3s (Terraform) ← el único recurso que cuesta dinero
3. Opcional: destruir el backend S3 + DynamoDB (casi sin costo pero limpia todo)
```

> El backend S3 tiene `prevent_destroy = true`. Para eliminarlo hay que quitar esa protección primero.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: `aws ec2 describe-instances` mostrando la instancia `Produccion-K3s` en estado `running`. Al lado, la consola de billing de AWS mostrando el consumo del Free Tier.*

"Bienvenidos al episodio 49. El último episodio del curso.

Esta instancia EC2 ha estado corriendo durante todo el módulo de Kubernetes y GitOps. Ha alojado K3s, ArgoCD, MySQL y la app Go. Ha recibido despliegues automáticos de ArgoCD. Ha demostrado que el patrón GitOps funciona completamente en la capa gratuita de AWS.

Ahora la vamos a destruir.

No porque dejó de ser útil — sino porque en un entorno de aprendizaje, saber cuándo y cómo limpiar los recursos es tan importante como saber crearlos. El Free Tier de AWS tiene 750 horas mensuales de t2.micro. Si no destruyes la instancia al terminar el curso, esas horas se siguen consumiendo. Y cuando el primer año termina, los cargos empiezan.

Terraform Destroy elimina todo lo que Terraform creó, en el orden correcto, sin dejar recursos huérfanos. Un solo comando.

Empecemos."

---

### PASO 1 — Detener los servicios locales (opcional) (1:30 – 2:30)

> *Pantalla: terminal local.*

"Jenkins y SonarQube corren en tu máquina local en Docker. No generan ningún costo — puedes dejarlos corriendo sin problema.

Si quieres liberar memoria y recursos del sistema, los detienes:"

```bash
cd ~/local-ci

# Detener los contenedores — los datos en volúmenes se conservan
docker compose stop
```

"Si quieres eliminar también los volúmenes — toda la configuración de Jenkins y SonarQube:"

```bash
# Solo si quieres empezar desde cero
docker compose down -v
```

"Para el curso recomiendo solo `stop` — si en algún momento quieres retomar la práctica, `docker compose up -d` y todo vuelve exactamente donde lo dejaste."

---

### PASO 2 — Ver el plan de destrucción (2:30 – 5:00)

> *Pantalla: terminal local en el directorio de Terraform de la EC2.*

"Antes de destruir, siempre el plan. El mismo hábito del EP21 — siempre leer antes de confirmar.

Entro al directorio de Terraform de la EC2:"

```bash
cd gitops-infra/infrastructure/terraform/jenkins-ec2
```

"Verifico que el estado apunta a la instancia correcta:"

```bash
terraform state list
# aws_instance.prod_server
# aws_security_group.prod_sg
```

"Dos recursos — la instancia y el Security Group. Eso es exactamente lo que creamos en el EP22.

Veo el plan de destrucción:"

```bash
terraform plan -destroy -var="key_name=aws-key"
```

```
Terraform will perform the following actions:

  # aws_instance.prod_server will be destroyed
  - resource "aws_instance" "prod_server" {
      - ami           = "ami-0c7217cdde317cfec"
      - instance_type = "t2.micro"
      - tags          = { "Name" = "Produccion-K3s" }
    }

  # aws_security_group.prod_sg will be destroyed
  - resource "aws_security_group" "prod_sg" {
      - name = "prod-sg"
    }

Plan: 0 to add, 0 to change, 2 to destroy.
```

"Dos recursos a destruir — exactamente los dos que creamos. Cero sorpresas. Procedo."

---

### PASO 3 — Destruir la EC2 (5:00 – 7:30)

> *Pantalla: terminal ejecutando terraform destroy.*

```bash
terraform destroy -var="key_name=aws-key" -auto-approve
```

"Uso `-auto-approve` porque acabo de revisar el plan y sé exactamente qué va a destruir.

El output muestra la destrucción en orden — Terraform calcula las dependencias y destruye en el orden correcto:

Primero la instancia EC2:"

```
aws_instance.prod_server: Destroying...
aws_instance.prod_server: Still destroying... [30s elapsed]
aws_instance.prod_server: Destruction complete after 32s
```

"Luego el Security Group — que no se puede eliminar mientras la instancia que lo usa existe:"

```
aws_security_group.prod_sg: Destroying...
aws_security_group.prod_sg: Destruction complete after 8s

Destroy complete! Resources: 2 destroyed.
```

"Todo destruido. En AWS, la instancia `Produccion-K3s` está ahora en estado `terminated`. El Security Group `prod-sg` ya no existe."

---

### PASO 4 — Verificar en AWS que no queda nada (7:30 – 9:30)

> *Pantalla: terminal con comandos de verificación.*

"Verifico desde la CLI de AWS que no quedaron recursos activos:

**Instancias EC2 corriendo:**"

```bash
aws ec2 describe-instances \
  --filters "Name=instance-state-name,Values=running" \
  --query "Reservations[].Instances[].InstanceId" \
  --output text
# (silencio — ninguna instancia corriendo)
```

"**Security Groups:**"

```bash
aws ec2 describe-security-groups \
  --filters "Name=group-name,Values=prod-sg" \
  --query "SecurityGroups[].GroupId" \
  --output text
# (silencio — el SG ya no existe)
```

"**Verificación en la consola web de AWS.** Abro la consola en el navegador, voy a EC2 → Instances. La instancia `Produccion-K3s` aparece en estado `terminated` — eso significa que fue eliminada y no genera costos."

---

### PASO 5 — Limpiar el kubectl local (9:30 – 10:30)

> *Pantalla: terminal local.*

"El kubeconfig en mi máquina local todavía tiene la configuración del cluster K3s que ya no existe. Lo limpio para evitar confusión:"

```bash
# Verificar el contexto actual
kubectl config get-contexts
# default   ← el de K3s (ya no existe el cluster)
# minikube  ← sigue disponible si lo usas

# Cambiar al contexto de Minikube si existe, o eliminar el de K3s
kubectl config delete-context default
kubectl config delete-cluster default
kubectl config delete-user default
```

"Si vuelven a crear el cluster K3s en el futuro, el proceso del EP30 reconfigura el kubeconfig desde cero."

---

### PASO 6 — ¿Qué queda en AWS? (10:30 – 12:00)

> *Pantalla: consola de AWS — S3.*

"Hay dos recursos de AWS que no destruimos porque tienen `prevent_destroy = true` y prácticamente no generan costo:

**El bucket S3 del backend de Terraform** — `curso-gitops-terraform-state`. Guarda el state file de Terraform. Si en el futuro quieres retomar el curso o crear recursos nuevos, Terraform encontrará el estado intacto.

**La tabla DynamoDB** — `curso-gitops-terraform-locks`. El mecanismo de locking de Terraform.

El costo mensual de estos dos recursos con el uso del curso es menos de $0.01 USD. Prácticamente cero.

Si quieres eliminarlos también:"

```bash
# Quitar el prevent_destroy del main.tf primero
# Luego:
cd gitops-infra/infrastructure/terraform/backend
terraform destroy -auto-approve

# O directamente con CLI:
aws s3 rm s3://curso-gitops-terraform-state --recursive
aws s3api delete-bucket --bucket curso-gitops-terraform-state
aws dynamodb delete-table --table-name curso-gitops-terraform-locks
```

"La decisión es tuya. Para la gran mayoría de los alumnos, dejar el bucket y la tabla es la opción correcta — su costo es irrisorio y hacen que retomar el curso sea trivial."

---

### CIERRE DEL CURSO (12:00 – 15:00)

> *Pantalla: el diagrama de la arquitectura completa del EP46 — todos los componentes visibles.*

"Eso es el EP49. Y con esto, terminamos el curso de GitOps.

---

Permíteme hacer un recuento de lo que construiste.

Empezaste con una PC y una cuenta de GitHub. Cuarenta y nueve episodios después tienes:

**Un pipeline de CI completo** corriendo en tu máquina local. Jenkins que clona el código, SonarQube que analiza su calidad, Docker que construye la imagen, Trivy que la escanea en busca de vulnerabilidades, y un push automático a Docker Hub con tags versionados y trazables.

**Un cluster de Kubernetes en la nube** — K3s en una EC2 t2.micro de la capa gratuita de AWS. Con Swap configurado para que no colapse bajo la carga, con ArgoCD instalado y expuesto via NodePort, y con kubectl configurado en tu laptop para operarlo remotamente sin SSH.

**Un flujo GitOps funcionando de extremo a extremo.** Un commit en `gitops-app` → Jenkins actualiza `gitops-infra` → ArgoCD detecta el cambio → K3s hace el rolling update → la nueva versión está en producción. Sin intervención manual. Sin tocar el servidor directamente.

**Y dos repositorios** — `gitops-app` con el código y el Jenkinsfile, `gitops-infra` con los manifiestos de Kubernetes y el código de Terraform. La separación de responsabilidades que es el estándar de GitOps en producción.

---

Lo que te llevas de este curso no es configuración específica de K3s o Jenkins local. Lo que te llevas es el patrón. El Jenkinsfile que escribiste funciona en Jenkins en AWS. Los manifiestos YAML que escribiste funcionan en EKS, GKE o AKS. ArgoCD funciona igual en cualquier cluster certificado por la CNCF.

El patrón es el conocimiento. El patrón es lo que vale.

---

Si tienes preguntas, déjalas en los comentarios. Si encontraste algo que no funciona, reporta el issue. Si este curso te ayudó a conseguir trabajo o a mejorar en tu trabajo actual, me encantaría saberlo.

Gracias por llegar hasta aquí.

Nos vemos en el próximo curso."

---

## ✅ Checklist de Verificación Final del Curso

**Módulo 01 — VS Code**
- [ ] VS Code instalado con extensiones de DevOps

**Módulo 02 — Git y GitHub**
- [ ] `gitops-app` y `gitops-infra` en GitHub
- [ ] Gitflow y Conventional Commits como práctica

**Módulo 03 — Docker**
- [ ] App Go containerizada con multi-stage build
- [ ] Imagen en Docker Hub

**Módulo 04 — AWS**
- [ ] Cuenta AWS con MFA y alerta de billing
- [ ] Usuario IAM `admin-curso` con Access Keys

**Módulo 05 — Terraform**
- [ ] Backend S3 + DynamoDB para el state
- [ ] EC2 t2.micro destruida limpiamente

**Módulo 06 — Kubernetes**
- [ ] Conceptos Pod, Deployment, Service dominados
- [ ] `kubectl` como herramienta de operación diaria

**Módulo 07 — K3s**
- [ ] K3s instalado en EC2 Free Tier
- [ ] kubectl configurado para acceso remoto

**Módulo 08 — Jenkins**
- [ ] Pipeline CI completo con 6 stages
- [ ] Patrón de dos repositorios GitOps implementado

**Módulo 09 — ArgoCD**
- [ ] ArgoCD sincronizando `gitops-infra` con K3s
- [ ] Despliegue automático funcionando

**Módulo 10 — Seguridad**
- [ ] SonarQube analizando calidad del código
- [ ] Trivy escaneando vulnerabilidades de la imagen

**Módulo 11 — Proyecto Final**
- [ ] Arquitectura híbrida completa documentada
- [ ] Recursos de AWS limpiados con terraform destroy

---

## 🔧 Si necesitas recrear el entorno

```bash
# 1. Crear la EC2 de nuevo
cd gitops-infra/infrastructure/terraform/jenkins-ec2
terraform apply -var="key_name=aws-key" -auto-approve

# 2. Configurar Swap en la nueva EC2 (EP28)
ssh -i aws-key.pem ubuntu@$(terraform output -raw prod_public_ip)
sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile
sudo mkswap /swapfile && sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

# 3. Instalar K3s (EP29)
curl -sfL https://get.k3s.io | sh -
mkdir -p ~/.kube && sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown ubuntu:ubuntu ~/.kube/config
exit

# 4. Configurar kubectl local (EP30)
scp -i aws-key.pem ubuntu@IP:~/k3s-remote.yaml ~/.kube/k3s-config
export KUBECONFIG=~/.kube/config:~/.kube/k3s-config

# 5. Instalar ArgoCD (EP38-40)
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort", "ports": [{"port": 443, "targetPort": 8080, "nodePort": 30080}]}}'
kubectl apply -f gitops-infra/infrastructure/kubernetes/argocd/application.yaml

# 6. Levantar Jenkins y SonarQube locales
cd ~/local-ci && docker compose up -d
```

"Con práctica, este proceso completo toma menos de 20 minutos."

---

## 🗒️ Notas de Producción
- La apertura mostrando la instancia corriendo y la consola de billing establece el contexto — "esto cuesta, y lo vamos a limpiar".
- El `terraform plan -destroy` antes del `apply -auto-approve` es el hábito correcto aunque ya lo revisaste en el EP21 — mostrarlo en el episodio final lo refuerza.
- El cierre del curso merece tiempo y pausa — no apresurarse. La lista de lo que el alumno construyó, leída en voz alta, es el momento de mayor impacto emocional del curso.
- La sección "si necesitas recrear el entorno" es el regalo práctico del cierre — el alumno sabe que puede volver a practicar en cualquier momento.
- Terminar con el diagrama de arquitectura en pantalla mientras se da el cierre verbal — el último visual que el alumno ve es el sistema completo que construyó.

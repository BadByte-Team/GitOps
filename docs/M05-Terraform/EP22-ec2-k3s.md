# EP 22: EC2 para K3s — Servidor de Producción Gratuito

**Tipo:** PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)
**🔄 MODIFICADO:** En lugar de una t2.medium para Jenkins, se provisiona una t2.micro (capa gratuita) preparada para K3s, abriendo los puertos 30080 (ArgoCD) y 30081 (App).

---

## 🎯 Objetivo
Usar Terraform para crear la EC2 t2.micro definitiva del curso — el servidor que correrá K3s, ArgoCD y la app Go — con su Security Group configurado y disco de 30 GB. Conectar por SSH y verificar los recursos del sistema.

---

## 📋 Prerequisitos
- Backend S3 y DynamoDB configurados (EP17, EP20)
- `terraform init` ejecutado en `jenkins-ec2/` (EP20)
- AWS CLI configurado (EP15)

---

## 💰 ¿Por qué t2.micro y no t2.medium?

| Instancia | RAM | Costo/mes | ¿Gratis? | Rol en el curso |
|---|---|---|---|---|
| t2.medium (original) | 4 GB | ~$33 USD | ❌ | Jenkins en AWS |
| t2.micro (actual) | 1 GB | $0 | ✅ Free Tier | K3s + ArgoCD + App |

La decisión que hace posible el $0: **Jenkins corre en tu máquina local** (EP31). La EC2 solo tiene responsabilidad de producción. Con el Swap del EP28, 1 GB de RAM es suficiente para K3s, ArgoCD, MySQL y la app.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: dos ventanas lado a lado. Izquierda: calculadora de precios de AWS mostrando ~$100/mes para EKS + EC2 t2.medium. Derecha: `variables.tf` con `instance_type = "t2.micro"` y el comentario `(Capa Gratuita)` visible.*

"Bienvenidos al episodio 22.

Llevamos cuatro episodios preparando el terreno: instalamos Terraform, aprendimos los bloques de HCL, conectamos el backend a S3, y dominamos el ciclo de comandos. Todo apuntaba a este momento.

Hoy creamos el servidor de producción del curso. Y quiero que veamos primero esto — el costo de la arquitectura original de este tipo de cursos. EKS para el cluster, una EC2 para Jenkins, un LoadBalancer para exponer la app. Más de cien dólares al mes. Y eso si lo destruyes al terminar cada sesión de práctica.

Ahora miren el `variables.tf` de nuestro proyecto. `instance_type = 't2.micro'`. Capa gratuita de AWS. Costo cero.

¿Cómo es posible con solo 1 GB de RAM en un solo servidor? Porque tomamos decisiones de arquitectura inteligentes. Jenkins no está en AWS — corre en tu máquina local, la vemos en el EP31. ArgoCD se expone con NodePort en lugar de un LoadBalancer. Y K3s en lugar de EKS elimina completamente el costo del control plane.

El resultado es el mismo stack GitOps completo — K3s, ArgoCD, la app, MySQL — a cero dólares. Eso es lo que vamos a crear ahora mismo.

Empecemos."

---

### ¿Qué vamos a crear? (1:30 – 3:30)

> *Pantalla: VS Code con `main.tf` abierto, señalando con el cursor mientras se explica.*

"Dos recursos. Solo dos. Eso es todo lo que hace Terraform en este episodio.

El primero es el **Security Group**. Piénsenlo como el firewall de la instancia — define qué tráfico puede entrar. En nuestro caso, tres puertos:

El **22** es SSH. Para conectarnos al servidor por terminal y administrarlo.

El **30080** es el puerto de ArgoCD. Cuando en el EP39 expongamos ArgoCD usando NodePort, va a escuchar en este puerto. Desde aquí van a abrir el navegador, escribir `http://IP_EC2:30080`, y ver el dashboard de ArgoCD.

El **30081** es el puerto de nuestra app Go. Cuando el EP40 complete el despliegue, van a escribir `http://IP_EC2:30081` y ver la plataforma del curso corriendo en producción.

Lo importante de esto es que estos tres puertos están configurados desde el día uno. No tenemos que volver a Terraform en el EP39 ni en el EP40 a agregar reglas. La infraestructura ya está lista para recibir lo que viene.

El segundo recurso es la **instancia EC2**. Ubuntu 22.04 LTS, tipo `t2.micro`, disco de 30 GB en formato gp3 encriptado. ¿Por qué 30 GB y no los 8 GB del default? Porque K3s descarga imágenes Docker al pull los contenedores — cada imagen puede ser varios cientos de megabytes — y con 8 GB nos quedaríamos sin espacio."

---

### PASO 1 — Crear el Key Pair (3:30 – 5:00)

> *Pantalla: terminal en `jenkins-ec2/`.*

"Antes de aplicar Terraform, necesito el archivo de llave para SSH. Si no lo creaste en el EP16, lo hago ahora con la CLI de AWS:"

```bash
cd gitops-infra/infrastructure/terraform/jenkins-ec2

aws ec2 create-key-pair \
  --key-name aws-key \
  --query 'KeyMaterial' \
  --output text > aws-key.pem

chmod 400 aws-key.pem
```

"El `--query 'KeyMaterial'` extrae solo la clave privada del JSON que devuelve AWS y la redirige directamente al archivo `.pem`. Sin ese flag, el archivo contendría un JSON completo y la clave no serviría para SSH.

El `chmod 400` — ya lo vimos en el EP16 — es obligatorio. SSH rechaza archivos de clave que tengan permisos más permisivos que `400`. Si intentas conectarte con un archivo que tiene `644`, SSH te dice 'Permissions 0644 for aws-key.pem are too open' y se niega a conectar.

Y el `aws-key.pem` está en el `.gitignore` del proyecto. Si hacen `git status` van a ver que el archivo no aparece como untracked — está ignorado explícitamente. Nunca va al repositorio."

---

### PASO 2 — Leer el plan (5:00 – 8:00)

> *Pantalla: terminal.*

"Antes de crear nada, el plan. El hábito del EP21:"

```bash
terraform plan -var="key_name=aws-key"
```

"Leo el output en voz alta.

El Security Group:"

```
# aws_security_group.prod_sg will be created
+ resource "aws_security_group" "prod_sg" {
    + name        = "prod-sg"
    + ingress     = [
        + { from_port = 22,    to_port = 22    }
        + { from_port = 30080, to_port = 30080 }
        + { from_port = 30081, to_port = 30081 }
      ]
  }
```

"Los tres puertos en verde. Todo como esperábamos.

La instancia:"

```
# aws_instance.prod_server will be created
+ resource "aws_instance" "prod_server" {
    + ami           = "ami-0c7217cdde317cfec"
    + instance_type = "t2.micro"
    + tags          = { "Name" = "Produccion-K3s" }
  }
```

"`t2.micro`. El tag `Produccion-K3s` — ese nombre lo usaremos después para filtrar la instancia con la CLI.

Y el resumen:"

```
Plan: 2 to add, 0 to change, 0 to destroy.
```

"Dos recursos a crear. Cero cambios. Cero destrucciones. Exactamente lo que esperábamos. El plan dice lo correcto. Aplico."

---

### PASO 3 — Aplicar (8:00 – 9:30)

> *Pantalla: terminal.*

"Uso `-auto-approve` porque sé exactamente qué está creando — el plan lo acabamos de leer. En un contexto donde no hayas revisado el plan manualmente, no uses esta flag:"

```bash
terraform apply -var="key_name=aws-key" -auto-approve
```

"Terraform crea el Security Group primero. Tiene que existir antes de asignarlo a la instancia — si la instancia se creara antes, no tendría a qué Security Group asignarse. Terraform detecta esa dependencia automáticamente desde el código y lo resuelve en el orden correcto.

Unos segundos después, la instancia.

Al terminar:"

```
Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:
prod_public_ip = "54.x.x.x"
```

"Esa IP — la anoto en algún lugar accesible. La necesito en el EP28 para conectarme por SSH y configurar el Swap, en el EP30 para copiar el kubeconfig, en el EP39 para acceder al dashboard de ArgoCD, y en el EP40 para ver la app corriendo. Es la dirección de producción del curso."

---

### PASO 4 — Verificar en AWS (9:30 – 10:30)

> *Pantalla: terminal.*

"Verifico desde la CLI que todo está como el plan prometió.

El estado de la instancia:"

```bash
aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=Produccion-K3s" \
  --query "Reservations[0].Instances[0].State.Name" \
  --output text
# running
```

"Running. Perfecto.

Los puertos del Security Group:"

```bash
aws ec2 describe-security-groups \
  --filters "Name=group-name,Values=prod-sg" \
  --query "SecurityGroups[0].IpPermissions[*].FromPort" \
  --output table
```

```
-----------
| FromPort |
-----------
| 22       |
| 30080    |
| 30081    |
-----------
```

"Los tres puertos. Exactamente lo que pusimos en el Terraform. El Security Group en AWS refleja exactamente lo que describe el código."

---

### PASO 5 — Conectar por SSH y explorar (10:30 – 13:00)

> *Pantalla: terminal.*

"Espero 30 segundos para que el sistema operativo termine de iniciar — la instancia puede estar en estado `running` pero Ubuntu todavía puede estar completando el arranque:"

```bash
sleep 30

ssh -i aws-key.pem ubuntu@$(terraform output -raw prod_public_ip)
```

"La primera vez aparece el mensaje de verificación del host:"

```
The authenticity of host '54.x.x.x' can't be established.
Are you sure you want to continue connecting (yes/no)? yes
```

"Escribo `yes`. SSH registra la clave del servidor — la próxima vez no vuelve a preguntar.

El prompt cambió. Ya no estoy en mi máquina — estoy en el servidor de AWS. Exploro los tres recursos que más nos importan.

**RAM:**"

```bash
free -h
```

```
               total        used        free
Mem:           981Mi       158Mi       823Mi
Swap:            0B          0B          0B
```

"Ahí está: 981 MB de RAM física, y cero de Swap. Eso es exactamente lo que tenemos que resolver en el EP28. Cuando K3s, ArgoCD y MySQL estén corriendo al mismo tiempo, este 1 GB no va a ser suficiente sin Swap. Pero eso es el EP28 — por ahora la instalación base está como esperábamos.

**Disco:**"

```bash
df -h /
```

```
Filesystem      Size  Used Avail Use% Mounted on
/dev/xvda1       29G  1.8G   28G   7%  /
```

"29 GB disponibles del disco de 30 GB que configuramos en Terraform. El 1 GB del sistema operativo ya está ocupado, lo que es normal. K3s y las imágenes Docker van a ir llenando ese espacio gradualmente, pero tenemos margen de sobra.

**CPU:**"

```bash
nproc
# 1
```

"Un solo procesador. En términos de rendimiento para el curso esto es completamente suficiente — las operaciones de Kubernetes no son CPU-intensivas, son RAM-intensivas. Por eso el Swap importa más que la cantidad de CPUs.

---"

```bash
exit
```

"Regreso a la máquina local."

---

### PASO 6 — Ver el state en S3 (13:00 – 14:00)

> *Pantalla: terminal.*

"El cierre del ciclo que comenzó en el EP17. Verifico que el state se guardó en S3:"

```bash
aws s3 ls s3://curso-gitops-terraform-state/
```

```
                           PRE backend/
                           PRE ec2-prod/
```

```bash
aws s3 ls s3://curso-gitops-terraform-state/ec2-prod/
```

```
2026-03-20 14:30:45       4821 terraform.tfstate
```

"Ahí está. El state de la EC2 ya no está en el disco de mi máquina — está en S3, encriptado y versionado. Si destruyo esta laptop, si alguien más clona `gitops-infra` en otra máquina y hace `terraform init`, el state está ahí y Terraform puede continuar gestionando esta misma instancia.

Eso es exactamente lo que significa tener infraestructura como código."

---

### CIERRE (14:00 – 15:00)

"Eso es el episodio 22.

El servidor de producción del curso está corriendo en AWS. Una sola instancia t2.micro, Free Tier, costo cero. Con su Security Group, sus tres puertos, y su disco de 30 GB. Todo descrito en código, versionado en Git, y el state guardado en S3.

Esta instancia no la vamos a destruir hasta el EP49. A partir de aquí, todos los módulos de Kubernetes y GitOps asumen que este servidor existe.

Lo que viene en los próximos episodios: el EP23 arranca con los conceptos de Kubernetes — Pods, Deployments, Services. El EP27 volvemos a esta EC2 para instalar K3s. Y lo primero que hacemos en el EP28 cuando llegamos a ese servidor es crear los 2 GB de Swap que necesitamos para que K3s no colapse la memoria.

Nos vemos en el EP23."

---

## ✅ Checklist de Verificación
- [ ] `terraform apply` completa con `2 added, 0 changed, 0 destroyed`
- [ ] La instancia aparece como `running` con nombre `Produccion-K3s`
- [ ] Security Group con puertos 22, 30080 y 30081 abiertos
- [ ] `ssh -i aws-key.pem ubuntu@<IP>` conecta sin errores
- [ ] `free -h` muestra ~981 MB de RAM total
- [ ] `df -h` muestra ~29 GB de disco disponible
- [ ] State visible en `s3://curso-gitops-terraform-state/ec2-prod/terraform.tfstate`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Error: No key pair found` | Verificar con `aws ec2 describe-key-pairs --key-names aws-key` — si no existe, ejecutar el Paso 1 |
| `Permission denied (publickey)` | El archivo no tiene los permisos correctos: `chmod 400 aws-key.pem` |
| `UNPROTECTED PRIVATE KEY FILE` | Mismo problema de permisos: `chmod 400 aws-key.pem` |
| `terraform plan` falla con error de backend | `terraform init` no fue ejecutado — ver EP20 |
| La IP pública cambió | La IP de t2.micro cambia al reiniciar la instancia — obtener la nueva con `terraform output prod_public_ip` |
| `Connection timed out` al hacer SSH | La instancia todavía está arrancando — esperar 60 segundos y volver a intentar |

---

## 🗒️ Notas de Producción
- La apertura con las dos ventanas — calculadora de AWS vs el `variables.tf` con `t2.micro` — es el contraste visual más poderoso del módulo. Prepararlo de antemano para que no haya demoras en el momento de la grabación.
- Al leer el `terraform plan` en voz alta, señalar con el cursor en la terminal cada sección mientras la describes — refuerza la conexión entre la voz y el texto en pantalla.
- El `free -h` dentro de la instancia mostrando `Swap: 0B` es el momento exacto para anticipar el EP28 — decirlo en voz alta crea la expectativa sin desviar el foco.
- Al mostrar el state en S3, el tono debe ser de cierre narrativo — "el ciclo que empezó en el EP17 se completa aquí" — no técnico sino emocional. Es el payoff del módulo.
- Decir explícitamente en el cierre que la instancia **no se destruye** hasta el EP49 — el alumno viene de un módulo donde destruía todo al final de cada episodio, y necesita escuchar claramente que este es diferente.

# 🎬 Guión — EP17: S3 y DynamoDB — El Backend de Terraform

**Duración estimada:** 12–15 min
**Tono:** Directo con una pausa conceptual sólida al inicio. Este episodio mezcla teoría importante con práctica de CLI. El alumno tiene que entender el problema antes de ver la solución.

---

## 🎙️ INTRO (0:00 – 0:50)

> *Pantalla: terminal. En una ventana, un archivo `terraform.tfstate` abierto en VS Code — se ve un JSON denso con IPs, ARNs, propiedades de recursos.*

"Bienvenidos al episodio 17.

Este es el archivo más peligroso de un proyecto de infraestructura: el **terraform.tfstate**. Aquí Terraform guarda el registro de todo lo que creó. Qué instancias existen, qué IPs tienen, qué configuración tienen. Sin este archivo, Terraform no sabe qué existe y qué no.

Por defecto, este archivo vive en tu máquina local. Y eso crea tres problemas que vamos a resolver hoy.

El **primero**: si borras el archivo por accidente — cosa que pasa — Terraform pierde el estado de tu infraestructura. No puede actualizar ni destruir lo que creó.

El **segundo**: si trabajas en equipo, cada persona tiene su propia copia del estado. Cuando dos personas hacen `terraform apply` al mismo tiempo, el estado se corrompe.

El **tercero**: el `.tfstate` puede contener contraseñas y claves en texto plano. Guardarlo localmente y commitearlo al repo por accidente es un agujero de seguridad.

La solución es el **backend remoto**: guardar el estado en S3 con encriptación, y usar DynamoDB para asegurarse de que nadie aplique cambios al mismo tiempo.

Eso es lo que creamos hoy. Empecemos."

---

## 🔍 PAUSA CONCEPTUAL — Cómo funciona el backend remoto (0:50 – 3:00)

> *Pantalla: diagrama o slide.*

"El flujo con backend remoto funciona así.

Cuando ejecutas `terraform apply`, Terraform no guarda el estado en tu disco. Lo guarda en **S3**. Específicamente, en un archivo como `ec2-prod/terraform.tfstate` dentro de tu bucket.

Antes de empezar a escribir cualquier cambio, Terraform crea un **lock** en **DynamoDB** — una entrada que dice 'alguien está trabajando en este estado ahora mismo, esperen'. Cuando termina, libera el lock.

```
Tu PC
  │
  │  terraform apply
  ▼
DynamoDB  ←── Terraform pone un lock
  │
S3 Bucket ←── Terraform lee el estado actual
  │
Tu PC     ←── Calcula los cambios
  │
AWS       ←── Aplica los cambios
  │
S3 Bucket ←── Terraform guarda el nuevo estado
  │
DynamoDB  ←── Terraform libera el lock
```

Si alguien más intenta hacer `terraform apply` mientras el lock está activo, Terraform le dice 'espera, hay otro proceso en curso'. Sin DynamoDB, dos personas podrían sobrescribir el estado al mismo tiempo.

---

Dos detalles sobre el costo de estos recursos, porque sé que es lo que todos están pensando.

El bucket S3 para el curso va a tener menos de 1 MB de datos — los archivos `.tfstate` son pequeños. El costo es literalmente cero.

La tabla DynamoDB con `PAY_PER_REQUEST` cobra por solicitud. Terraform hace una o dos operaciones de lock por `apply`. El costo mensual de esto es tan pequeño que ni aparece en el billing.

Son los dos recursos del curso que no tienen costo ni siquiera rozando el Free Tier."

---

## 📌 PASO 1 — Crear el bucket S3 (3:00 – 7:00)

> *Pantalla: terminal.*

"Tres cosas que configurar en el bucket: el versionado, la encriptación y el bloqueo de acceso público. Vamos de a uno.

**Crear el bucket:**"

```bash
aws s3api create-bucket \
  --bucket curso-gitops-terraform-state \
  --region us-east-1

aws s3 ls | grep curso-gitops
# 2026-XX-XX curso-gitops-terraform-state
```

"---

**Activar el versionado:**

Si el `.tfstate` se corrompe o necesitas volver a un estado anterior, el versionado de S3 te permite recuperar versiones previas del archivo."

```bash
aws s3api put-bucket-versioning \
  --bucket curso-gitops-terraform-state \
  --versioning-configuration Status=Enabled

aws s3api get-bucket-versioning \
  --bucket curso-gitops-terraform-state
# "Status": "Enabled"
```

"---

**Activar encriptación en reposo:**

El `.tfstate` puede contener IPs privadas, passwords de bases de datos, claves de acceso. La encriptación lo protege si alguien llegara a tener acceso no autorizado al bucket."

```bash
aws s3api put-bucket-encryption \
  --bucket curso-gitops-terraform-state \
  --server-side-encryption-configuration '{
    "Rules": [{
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }]
  }'
```

"---

**Bloquear el acceso público:**

Un bucket de S3 puede ser accesible públicamente si no lo configuras bien. El estado de Terraform **nunca** debe ser público."

```bash
aws s3api put-public-access-block \
  --bucket curso-gitops-terraform-state \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
```

"Verifico todo junto:"

```bash
echo "=== Versionado ==="
aws s3api get-bucket-versioning \
  --bucket curso-gitops-terraform-state \
  --query "Status" --output text

echo "=== Acceso público bloqueado ==="
aws s3api get-public-access-block \
  --bucket curso-gitops-terraform-state \
  --query "PublicAccessBlockConfiguration.BlockPublicAcls" --output text
```

"Ambos deben responder `Enabled` y `True` respectivamente."

---

## 📌 PASO 2 — Crear la tabla DynamoDB (7:00 – 9:00)

> *Pantalla: terminal.*

"La tabla DynamoDB para locking tiene un único requisito: debe tener una clave primaria llamada exactamente `LockID`. Terraform busca ese nombre específico cuando quiere crear o liberar un lock."

```bash
aws dynamodb create-table \
  --table-name curso-gitops-terraform-locks \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region us-east-1
```

"El `PAY_PER_REQUEST` significa que no se cobra por capacidad reservada — solo por las operaciones que realmente se usan. Para el volumen de este curso, el costo es cero.

Espero a que la tabla esté lista:"

```bash
aws dynamodb wait table-exists \
  --table-name curso-gitops-terraform-locks

aws dynamodb describe-table \
  --table-name curso-gitops-terraform-locks \
  --query "Table.TableStatus" \
  --output text
# ACTIVE
```

---

## 📌 PASO 3 — Ver cómo Terraform usa este backend (9:00 – 11:30)

> *Pantalla: VS Code abierto con el archivo `gitops-infra/infrastructure/terraform/jenkins-ec2/main.tf`.*

"Ahora quiero mostrarte dónde aparecen estos recursos en el proyecto.

Abro el `main.tf` de la EC2 que usaremos en el EP22:"

```hcl
terraform {
  required_providers { ... }

  backend "s3" {
    bucket         = "curso-gitops-terraform-state"   ← el bucket que acabamos de crear
    key            = "ec2-prod/terraform.tfstate"     ← ruta dentro del bucket
    region         = "us-east-1"
    dynamodb_table = "curso-gitops-terraform-locks"   ← la tabla que acabamos de crear
    encrypt        = true
  }
}
```

"Ese bloque `backend "s3"` le dice a Terraform: 'no guardes el estado en disco — guárdalo en S3, y usa DynamoDB para el locking'.

La propiedad `key` es la ruta dentro del bucket. Cada configuración de Terraform del proyecto usa una ruta diferente para no sobreescribirse:

- `backend/terraform.tfstate` → el estado del propio backend
- `ec2-prod/terraform.tfstate` → el estado de la EC2

Así puedo destruir la EC2 con `terraform destroy` sin tocar el backend, y viceversa.

---

También quiero mostrarte el `main.tf` del backend en sí:"

```hcl
resource "aws_s3_bucket" "terraform_state" {
  bucket = var.bucket_name

  lifecycle {
    prevent_destroy = true   ← esto es clave
  }
}
```

"`prevent_destroy = true` protege contra un `terraform destroy` accidental del bucket. Aunque alguien ejecute `terraform destroy` en la carpeta del backend, Terraform se va a negar a borrar el bucket. Tienes que quitar esa línea manualmente y volver a aplicar antes de poder destruirlo. Es una red de seguridad intencional."

---

## 📌 Verificación final (11:30 – 12:30)

> *Pantalla: terminal.*

"Verificación completa de todo lo creado en este episodio:"

```bash
echo "=== Bucket S3 ==="
aws s3 ls | grep curso-gitops-terraform-state

echo "=== Versionado ==="
aws s3api get-bucket-versioning \
  --bucket curso-gitops-terraform-state \
  --query "Status" --output text

echo "=== Acceso público ==="
aws s3api get-public-access-block \
  --bucket curso-gitops-terraform-state \
  --query "PublicAccessBlockConfiguration.BlockPublicAcls" --output text

echo "=== DynamoDB ==="
aws dynamodb describe-table \
  --table-name curso-gitops-terraform-locks \
  --query "Table.TableStatus" --output text
```

"Los cuatro deben responder sin error. Si todos están bien, el backend está listo."

---

## 🎙️ CIERRE (12:30 – 13:30)

"Eso es EP17.

El bucket S3 y la tabla DynamoDB están listos. A partir del EP18 en adelante, cada vez que inicialicemos Terraform en este proyecto, el estado se guardará en la nube y los locks prevendrán conflictos.

Nota algo importante: estos son los **únicos dos recursos de AWS que no vamos a destruir en el EP49**. El bucket tiene `prevent_destroy = true` por una razón. Son la infraestructura de la infraestructura — necesitan existir mientras el proyecto exista.

En el siguiente episodio instalamos Terraform y entendemos su estructura básica. Después de eso, empezamos a usarlo para crear la EC2 de producción.

Nos vemos en el EP18."

---

## 🗒️ Notas de producción

- La apertura mostrando el contenido real de un `.tfstate` es el gancho más efectivo — ayuda a entender por qué necesita estar protegido. Abrir el archivo en VS Code con resaltado de JSON para que se vea denso y complejo.
- El diagrama del flujo con S3 + DynamoDB es el momento más conceptual del episodio. Vale la pena dedicarle unos segundos de pantalla estática — que el espectador lo pueda leer.
- Al mostrar el bloque `backend "s3"` en el `main.tf`, señalar con el cursor cada campo mientras se explica su relación con los recursos recién creados.
- El `prevent_destroy = true` merece un comentario enfático — es la diferencia entre poder recuperarse de un error o no.
- Si el nombre del bucket `curso-gitops-terraform-state` ya existe (nombres de S3 son globales únicos), mencionar que hay que agregar un sufijo único y cambiar el nombre también en el `variables.tf`.

# EP 20: Configuración del Backend Remoto (S3)

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Agregar el bloque `backend "s3"` a la configuración de Terraform para que el state file se guarde en el bucket S3 del EP17, y activar el backend ejecutando `terraform init` sobre el directorio de la EC2.

---

## 📋 Prerequisitos
- Bucket S3 `curso-gitops-terraform-state` y tabla DynamoDB `curso-gitops-terraform-locks` creados (EP17)
- Terraform instalado y bloques de HCL dominados (EP18, EP19)

---

## 🧠 ¿Qué cambia al agregar el backend?

```
Sin backend (EP19):            Con backend (EP20+):
Tu PC                          Tu PC
  │                              │
  │ terraform apply              │ terraform apply
  ▼                              ▼
AWS (crea recursos)            AWS (crea recursos)
  │                              │
  ▼                              ├── S3: guarda el state
terraform.tfstate               │   (encriptado, versionado)
(archivo en tu disco)           │
                                └── DynamoDB: crea lock
                                    (nadie más puede apply al mismo tiempo)
```

El cambio en el código es solo este bloque dentro de `terraform {}`:

```hcl
backend "s3" {
  bucket         = "curso-gitops-terraform-state"
  key            = "ec2-prod/terraform.tfstate"
  region         = "us-east-2"
  dynamodb_table = "curso-gitops-terraform-locks"
  encrypt        = true
}
```

Cada configuración de Terraform usa una `key` diferente para no sobreescribirse entre sí:

| Configuración | key en S3 |
|---|---|
| Backend (S3 + DynamoDB) | `backend/terraform.tfstate` |
| EC2 K3s | `ec2-prod/terraform.tfstate` |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: terminal mostrando el `terraform.tfstate` local que quedó del EP19, con `cat terraform.tfstate | head -20` visible.*

"Bienvenidos al episodio 20.

Este es el `terraform.tfstate` que quedó en el disco local después del EP19. Es el registro de lo que Terraform creó. Si este archivo desaparece — si alguien lo borra accidentalmente, si el disco falla, si trabajas desde otra máquina — Terraform pierde la noción de qué existe y qué no en AWS.

En el EP17 creamos el bucket S3 y la tabla DynamoDB precisamente para resolver ese problema. Hoy los conectamos a Terraform con cinco líneas de configuración.

Es el episodio más corto del módulo, pero uno de los más importantes. Vamos."

---

### 🔍 PAUSA CONCEPTUAL — El bloque `backend "s3"` (1:00 – 3:00)

> *Pantalla: VS Code mostrando el bloque del backend.*

"El cambio que vamos a hacer es agregar este bloque dentro del `terraform {}` de cualquier configuración que quiera usar el backend remoto:

```hcl
backend "s3" {
  bucket         = "curso-gitops-terraform-state"
  key            = "ec2-prod/terraform.tfstate"
  region         = "us-east-2"
  dynamodb_table = "curso-gitops-terraform-locks"
  encrypt        = true
}
```

Voy a explicar cada propiedad porque las van a ver exactamente así en el proyecto.

**`bucket`** — el nombre del bucket S3 del EP17. Ahí es donde Terraform va a guardar el archivo `.tfstate`.

**`key`** — la ruta dentro del bucket donde se guarda el state. No es una clave de seguridad — es literalmente una ruta de archivo, como una carpeta. En nuestro caso, `ec2-prod/terraform.tfstate`. Cada directorio de Terraform del proyecto usa una ruta diferente para no sobreescribirse entre sí. Si el backend guardó su propio state en `backend/terraform.tfstate`, y la EC2 usa `ec2-prod/terraform.tfstate`, son dos archivos completamente independientes. Puedo destruir la EC2 sin tocar el backend, y viceversa.

**`region`** — siempre `us-east-2` en este curso, la misma región donde creamos todo.

**`dynamodb_table`** — la tabla del EP17. Cuando Terraform empieza un `apply`, crea una entrada en esta tabla que dice 'estoy trabajando en este state ahora mismo, esperen'. Cuando termina, la borra. Si dos personas intentan hacer `apply` al mismo tiempo, la segunda ve el lock y espera. Sin esto, las dos podrían sobrescribir el state simultáneamente y corromperlo.

**`encrypt`** — con `true`, S3 encripta el archivo `.tfstate` en reposo con AES-256. El state puede contener IPs privadas, passwords, claves de acceso — necesita estar encriptado."

---

### PASO 1 — Aplicar el backend de S3 y DynamoDB (3:00 – 6:00)

> *Pantalla: terminal en `gitops-infra/infrastructure/terraform/backend/`.*

"Lo primero es crear el propio bucket S3 y la tabla DynamoDB, si no los hiciste ya en el EP17. Esta es la única configuración del proyecto que usa un backend local — porque está creando el backend. No puede guardarse a sí misma en S3 antes de que S3 exista.

Abro el directorio y reviso el archivo antes de aplicar:"

```bash
cd gitops-infra/infrastructure/terraform/backend
cat main.tf
```

"Noto una línea importante en el recurso del bucket S3:"

```hcl
lifecycle {
  prevent_destroy = true
}
```

"Esto es una red de seguridad explícita. Aunque alguien ejecute `terraform destroy` en este directorio — por error, por accidente, por lo que sea — Terraform se va a negar a borrar el bucket. Va a mostrar un error que dice 'prevent_destroy is set to true'. Para eliminarlo habría que quitar esa línea manualmente, hacer un `apply`, y recién entonces el `destroy` procedería. Es una fricción intencional que protege el recurso más crítico del proyecto.

Aplico:"

```bash
terraform init
terraform apply -auto-approve
```

"Verifico que los recursos existen en AWS:"

```bash
aws s3 ls | grep curso-gitops-terraform-state
# Debe aparecer el bucket

aws dynamodb describe-table \
  --table-name curso-gitops-terraform-locks \
  --query "Table.TableStatus" \
  --output text
# ACTIVE
```

"Bucket existente, tabla activa. Todo listo para que Terraform los use."

---

### PASO 2 — Inicializar la EC2 con backend remoto (6:00 – 8:30)

> *Pantalla: terminal en `jenkins-ec2/`.*

"Ahora el paso central del episodio. El `main.tf` de la EC2 ya tiene el bloque `backend "s3"` configurado desde el principio. Solo necesito correr `init` para que Terraform lo active:"

```bash
cd gitops-infra/infrastructure/terraform/jenkins-ec2
terraform init
```

"Detengo el video aquí en el output porque hay una línea específica que quiero que vean:"

```
Initializing the backend...

Successfully configured the backend "s3"! Terraform will automatically
use this backend unless the backend configuration changes.
```

"Esa es la confirmación. `Successfully configured the backend "s3"!`. Terraform encontró el bloque, conectó con el bucket S3, y a partir de ahora cualquier `apply` o `destroy` en este directorio va a leer y escribir el state desde la nube. El archivo local `terraform.tfstate` ya no existe — el estado vive en S3.

Un detalle que vale la pena mencionar: si tuvieras un state local previo del EP19 en este directorio, Terraform te preguntaría aquí si quieres migrarlo a S3. Responderías `yes` y el archivo se copiaría al bucket automáticamente. Como todavía no hemos creado la EC2, no hay state previo, así que no hay migración. Pero es bueno saber que ese flujo existe."

---

### PASO 3 — El lock de DynamoDB en acción (8:30 – 10:00)

> *Pantalla: terminal.*

"Quiero mostrarles cómo funciona el lock en la práctica, aunque todavía no hayamos creado la EC2.

Antes de cualquier `apply`, la tabla de DynamoDB está vacía — no hay ningún lock activo:"

```bash
aws dynamodb scan \
  --table-name curso-gitops-terraform-locks \
  --query "Items" \
  --output json
# []
```

"Cuando en el EP22 ejecutemos `terraform apply`, si abren una segunda terminal y hacen esta misma consulta mientras el apply está corriendo, van a ver una entrada en la tabla con el ID del lock, la ruta del state, y la hora en que se creó. Cuando el apply termina, la entrada desaparece sola.

Si alguna vez el apply se interrumpe a mitad — se corta la conexión, se cierra la terminal, lo que sea — el lock puede quedar activo. En ese caso, Terraform te lo va a decir con un error como este:"

```
Error: Error acquiring the state lock

Lock Info:
  ID:   abc-123-def-456
  Path: curso-gitops-terraform-state/ec2-prod/terraform.tfstate
```

"El ID del lock aparece ahí. Para liberarlo:"

```bash
terraform force-unlock abc-123-def-456
```

"Pero solo haz esto si estás seguro de que no hay otro proceso corriendo. El lock existe por una razón."

---

### CIERRE (10:00 – 11:00)

"Eso es el episodio 20.

El backend remoto está activo. A partir de ahora, cuando ejecutemos `terraform apply` en cualquier directorio del proyecto, el state se guarda en S3 encriptado y DynamoDB previene conflictos. Si destruyes tu máquina local, el estado de la infraestructura está intacto en la nube. Si alguien más clona `gitops-infra` y hace `terraform init`, Terraform encuentra el mismo state y puede continuar desde donde lo dejaste.

En el siguiente episodio dominamos los cinco comandos del ciclo de vida de Terraform — `init`, `validate`, `plan`, `apply` y `destroy`. Los hemos usado ya, pero quiero que sean predecibles antes de llegar al EP22. Que cuando aparezca un output en la terminal, sepan exactamente qué significa.

Nos vemos en el EP21."

---

## ✅ Checklist de Verificación
- [ ] `terraform apply` en `backend/` completa sin errores
- [ ] El bucket `curso-gitops-terraform-state` existe y está activo
- [ ] La tabla `curso-gitops-terraform-locks` está en estado `ACTIVE`
- [ ] `terraform init` en `jenkins-ec2/` muestra `Successfully configured the backend "s3"!`
- [ ] Entiendes para qué sirve la propiedad `key` y por qué es diferente en cada directorio

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Error: Failed to get existing workspaces` | El bucket S3 no existe — ejecutar `terraform apply` en el directorio `backend/` primero |
| `Error: Error acquiring the state lock` | Lock atascado de un apply previo — `terraform force-unlock LOCK_ID` |
| `BucketAlreadyExists` | El nombre del bucket es global en S3 — agregar un sufijo único en `variables.tf` |
| `NoSuchBucket` al hacer init | El bucket fue creado en una región diferente a `us-east-2` — verificar |

---

## 🗒️ Notas de Producción
- La apertura con el `terraform.tfstate` local visible es el gancho — el contraste con "ahora esto va a estar en S3" motiva todo el episodio.
- Al explicar cada propiedad del bloque `backend "s3"`, señalar con el cursor en VS Code mientras describes — el alumno necesita conectar la palabra con la línea de código.
- Hacer zoom en la línea `Successfully configured the backend "s3"!` del output del `init` — es el momento de confirmación más importante del episodio.
- Después del apply, ir a la consola web de AWS y mostrar el bucket S3 creado con sus propiedades (versionado activo, acceso público bloqueado) — la conexión visual entre el código HCL y la realidad de la nube.
- Mencionar explícitamente que en este episodio no se crea la EC2 todavía — eso es el EP22 — para que el alumno no se frustre esperando ver la instancia.

# EP 19: Primeros Pasos con IaC — Provider, Resource y Variables

**Tipo:** PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Entender los cuatro bloques fundamentales de HCL escribiendo una configuración de ejemplo desde cero, para poder leer el `main.tf` real del proyecto sin que ninguna línea sea opaca.

---

## 📋 Prerequisitos
- Terraform instalado (EP18)
- AWS CLI configurado (EP15)

---

## 🧠 Los bloques de HCL

Terraform usa **HCL — HashiCorp Configuration Language**. No es un lenguaje de programación — es declarativo. No le dices *cómo* crear algo, le dices *qué* quieres que exista, y Terraform calcula cómo llegar ahí.

### `terraform {}` — configura el motor
```hcl
terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
}
```
El `~> 5.0` previene que una actualización mayor rompa la configuración.

### `provider {}` — cómo conectarse al servicio
```hcl
provider "aws" { region = "us-east-1" }
```
Terraform usa `~/.aws/credentials` automáticamente. Solo especificamos la región.

### `resource {}` — lo que queremos crear
```hcl
resource "aws_instance" "mi_servidor" {
  ami           = "ami-0c7217cdde317cfec"
  instance_type = "t2.micro"
}
```
Primer argumento: tipo de recurso. Segundo argumento: nombre local dentro de Terraform.

### `variable {}` — parámetros de entrada
```hcl
variable "instance_type" {
  description = "Tipo de instancia EC2"
  type        = string
  default     = "t2.micro"
}
```
Se usan con `var.instance_type`. Si cambias el valor, el cambio se propaga a todos los recursos que lo referencian.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: el archivo `gitops-infra/infrastructure/terraform/jenkins-ec2/main.tf` abierto en VS Code, con todo el contenido visible.*

"Bienvenidos al episodio 19.

Este es el `main.tf` de la EC2 que crearemos en el EP22 — el servidor donde corre K3s, ArgoCD y nuestra app. Tiene un bloque `backend` para guardar el estado en S3, un Security Group con tres reglas de ingress, una instancia EC2 con disco encriptado, y referencias entre recursos.

Si lo miraran ahora mismo sin contexto, probablemente podrían adivinar algunas partes, pero no todas. Y adivinar no es suficiente — necesitamos entender cada línea.

Hoy no vamos a crear eso todavía. Hoy vamos a construir los mismos bloques desde cero, en un ejemplo más simple, hasta que cada pieza sea completamente familiar. La idea es que cuando lleguemos al EP22, abran este archivo y puedan leerlo de principio a fin sin detenerse en nada.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — Los cuatro bloques de HCL (1:00 – 3:30)

> *Pantalla: editor de texto limpio, mostrando cada bloque uno a uno mientras se explica.*

"Terraform usa un lenguaje llamado HCL — HashiCorp Configuration Language. No es programación. Es **declarativo**: describes el estado final que quieres, y Terraform hace el trabajo de llegar ahí.

Hay cuatro bloques que van a aparecer en absolutamente todo lo que hagamos.

---

**El primero es `terraform {}`.**"

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
```

"Este bloque le dice a Terraform qué plugins necesita descargar. En nuestro caso, el provider de AWS de HashiCorp. El `~> 5.0` es una restricción de versión — significa 'acepta cualquier versión 5.x, pero si HashiCorp lanza la versión 6 con cambios incompatibles, no la uses'. Es una protección para que una actualización mayor no rompa la configuración sin aviso.

---

**El segundo es `provider {}`.**"

```hcl
provider "aws" {
  region = "us-east-1"
}
```

"Este bloque configura cómo Terraform se conecta a AWS. Le decimos la región. Las credenciales las toma automáticamente de `~/.aws/credentials` — las mismas que usa la CLI. No hay que repetirlas aquí.

---

**El tercero es `resource {}`, y es el más importante.**"

```hcl
resource "aws_instance" "mi_servidor" {
  ami           = "ami-0c7217cdde317cfec"
  instance_type = "t2.micro"
}
```

"La primera línea tiene dos argumentos separados. El primero — `aws_instance` — es el tipo de recurso. Le dice a Terraform qué clase de objeto crear en AWS. El segundo — `mi_servidor` — es el nombre local, el que nosotros le damos dentro de Terraform para poder referenciar este recurso desde otros lugares.

Dentro van las propiedades del recurso. La AMI, el tipo de instancia. Cada tipo de recurso tiene sus propios atributos — la documentación del provider de AWS en registry.terraform.io lista todos los disponibles.

---

**El cuarto es `variable {}`.**"

```hcl
variable "instance_type" {
  description = "Tipo de instancia EC2"
  type        = string
  default     = "t2.micro"
}
```

"En lugar de escribir `t2.micro` directamente en el resource, lo parametrizamos. El valor por defecto es `t2.micro`, pero se puede sobreescribir con `-var='instance_type=t2.nano'` al ejecutar el plan o apply. Esto hace que la misma configuración funcione para distintos entornos sin duplicar archivos.

Esos son los cuatro bloques. Con esto pueden leer cualquier configuración de Terraform del proyecto."

---

### PASO 1 — Crear el directorio de práctica (3:30 – 4:00)

> *Pantalla: terminal.*

"Creo un directorio de práctica. Todo lo que hagamos aquí es temporal — lo destruimos al final del episodio."

```bash
mkdir ~/terraform-practica && cd ~/terraform-practica
```

---

### PASO 2 — Escribir los tres archivos (4:00 – 8:00)

> *Pantalla: VS Code, creando los archivos uno por uno.*

"Siempre escribo en el mismo orden: variables primero, luego el main, luego los outputs. ¿Por qué ese orden? Porque cuando lees una configuración nueva de alguien más, las variables te dan el contexto — qué parámetros acepta, qué valores usa por defecto — antes de entrar a leer los recursos.

Empiezo por **`variables.tf`**:"

```hcl
variable "region" {
  description = "Region de AWS"
  type        = string
  default     = "us-east-1"
}

variable "ami_id" {
  description = "AMI de Ubuntu 22.04 LTS en us-east-1"
  type        = string
  default     = "ami-0c7217cdde317cfec"
}

variable "instance_type" {
  description = "Tipo de instancia — t2.micro es Free Tier"
  type        = string
  default     = "t2.micro"
}
```

"Tres variables. La región, el ID de la AMI, y el tipo de instancia. Cada una tiene una descripción — eso no es decorativo, es documentación. Cuando alguien vea este archivo en seis meses sin contexto, la descripción le explica qué es cada variable.

Ahora **`main.tf`**:"

```hcl
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

resource "aws_instance" "ejemplo" {
  ami           = var.ami_id
  instance_type = var.instance_type

  tags = {
    Name = "terraform-practica-ep19"
  }
}
```

"Noten cómo el provider y el resource no tienen valores escritos directamente — tienen `var.region`, `var.ami_id`, `var.instance_type`. Son referencias a las variables que acabamos de definir. Si mañana necesito cambiar la región a `us-west-2`, cambio un solo lugar en `variables.tf` y el cambio se propaga automáticamente a todo el archivo.

Y finalmente **`outputs.tf`**:"

```hcl
output "instance_id" {
  description = "ID de la instancia creada por Terraform"
  value       = aws_instance.ejemplo.id
}

output "instance_public_ip" {
  description = "IP publica de la instancia"
  value       = aws_instance.ejemplo.public_ip
}
```

"Los outputs referencian recursos con la sintaxis `tipo_recurso.nombre_local.atributo`. Cuando Terraform termina de crear la instancia, imprime estos dos valores en la terminal. En el EP22, el output más importante va a ser la IP pública — la usaremos para conectarnos por SSH."

---

### PASO 3 — El ciclo init → validate → plan → apply (8:00 – 12:00)

> *Pantalla: terminal, ejecutando cada comando con pausa para ver el output completo.*

"Ahora el ciclo. Siempre en este orden.

**`terraform init`** — siempre primero, siempre en un directorio nuevo. Descarga el provider de AWS:"

```bash
terraform init
```

"Ve el output — `Installing hashicorp/aws v5.x.x`. Terraform descargó el plugin. Solo hay que hacerlo una vez por directorio, o cuando cambias los providers.

---

**`terraform validate`** — verifica la sintaxis antes de hacer cualquier cosa en AWS:"

```bash
terraform validate
# Success! The configuration is valid.
```

"Si hubiera un error de tipeo en el código HCL — un atributo mal escrito, un bloque sin cerrar — este comando lo detecta aquí, antes de hacer ninguna consulta a AWS. Es el compilador de Terraform.

---

**`terraform fmt`** — formatea los archivos al estilo estándar de HCL:"

```bash
terraform fmt
```

"Ajusta la indentación, los espacios entre el signo igual y los valores. No cambia la lógica, solo el estilo. Lo ejecuto siempre antes de hacer commit para que el código se vea consistente.

---

**`terraform plan`** — el más importante. El dry-run:"

```bash
terraform plan
```

"Terraform calcula exactamente qué haría un `apply` y lo muestra sin crear nada. El output tiene dos partes que siempre leo.

La primera es el detalle de cada recurso:"

```
# aws_instance.ejemplo will be created
+ resource "aws_instance" "ejemplo" {
    + ami           = "ami-0c7217cdde317cfec"
    + instance_type = "t2.micro"
    + tags          = { "Name" = "terraform-practica-ep19" }
  }
```

"El signo `+` en verde significa que se va a crear. Más adelante veremos que `-` significa eliminar, y `~` significa modificar.

La segunda parte es el resumen, al final del output:"

```
Plan: 1 to add, 0 to change, 0 to destroy.
```

"Antes de ejecutar cualquier `apply`, siempre leo esta línea primero. Si dice algo diferente a lo que espero — si dice `1 to destroy` cuando no quería destruir nada — es el momento de parar y revisar.

---

**`terraform apply`** — crea la infraestructura:"

```bash
terraform apply
```

"Terraform muestra el plan de nuevo y pide confirmación. Escribo `yes`. Esta pausa es intencional — es tu última oportunidad de leer el plan antes de que ocurra algo en AWS."

"Al terminar, imprime los outputs:"

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:
instance_id        = "i-0abc123..."
instance_public_ip = "54.x.x.x"
```

"Ahí están los dos valores que definimos en `outputs.tf`. El ID de la instancia y la IP pública."

---

### PASO 4 — El state y las variables en la línea de comandos (12:00 – 13:30)

> *Pantalla: terminal.*

"Dos cosas antes de destruir.

**El state local** — voy a mostrarlo para que lo vean por última vez antes de que lo movamos a S3 en el EP20:"

```bash
cat terraform.tfstate | head -40
```

"Ahí está el JSON denso que vimos en el EP17. Terraform guardó el estado de la instancia que acaba de crear. En el EP20 este archivo deja de estar aquí y pasa a vivir en S3.

**Variables en la línea de comandos** — sin editar ningún archivo:"

```bash
terraform plan -var="instance_type=t2.nano"
```

"Terraform recalcula el plan con el nuevo valor. Útil cuando quieres probar distintas configuraciones sin modificar los archivos. En la práctica, esto se usa para diferenciar entornos — dev con `t2.nano`, staging con `t2.micro`, producción con algo más grande — sin duplicar ningún archivo."

---

### PASO 5 — Destruir y conectar con el proyecto real (13:30 – 15:00)

> *Pantalla: terminal, luego VS Code.*

"Destruyo la instancia de práctica:"

```bash
terraform destroy
```

"Muestra el plan inverso — la instancia marcada con `-` en rojo — y pide confirmación. Escribo `yes`.

Limpio el directorio:"

```bash
cd ~ && rm -rf ~/terraform-practica
```

"Ahora la conexión que quería hacer desde el inicio. Abro el `main.tf` real del proyecto:"

```bash
code gitops-infra/infrastructure/terraform/jenkins-ec2/main.tf
```

"Reconocen todo: el bloque `terraform {}` con el `required_providers`, el `provider "aws"`, el `resource "aws_security_group"`, el `resource "aws_instance"`. Los mismos cuatro bloques que acabamos de escribir, con más configuración dentro.

La única línea que todavía no hemos visto es esta:"

```hcl
backend "s3" {
  bucket         = "curso-gitops-terraform-state"
  key            = "ec2-prod/terraform.tfstate"
  region         = "us-east-1"
  dynamodb_table = "curso-gitops-terraform-locks"
  encrypt        = true
}
```

"Ese es el backend remoto — le dice a Terraform que guarde el state en S3 en lugar de en disco local. Eso es exactamente lo que vemos en el EP20."

---

### CIERRE (15:00 – 15:30)

"Eso es el episodio 19. Ya saben leer HCL. Saben qué es un provider, un resource, una variable y un output. Escribieron el ciclo completo de principio a fin.

En el siguiente episodio conectamos el backend remoto. Es el paso que mueve el state de tu disco a S3. Tres líneas de configuración, y el estado de tu infraestructura pasa a estar seguro en la nube.

Nos vemos en el EP20."

---

## ✅ Checklist de Verificación
- [ ] `terraform validate` responde `Success!`
- [ ] `terraform plan` muestra `1 to add, 0 to change, 0 to destroy`
- [ ] `terraform apply` crea la instancia y muestra los outputs
- [ ] `terraform destroy` elimina la instancia sin errores
- [ ] Puedes leer el `main.tf` del proyecto y reconocer todos los bloques

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Error: No valid credential sources found` | Verificar `~/.aws/credentials` con las claves del EP14 |
| `Error: Invalid AMI ID` | La AMI `ami-0c7217cdde317cfec` es de `us-east-1` — verificar que la región en `variables.tf` es la correcta |
| `Error: configuring Terraform AWS Provider` | El usuario IAM `admin-curso` no tiene los permisos necesarios — verificar que tiene `AdministratorAccess` |

---

## 🗒️ Notas de Producción
- La apertura con el `main.tf` real del proyecto como objetivo visual es el gancho — muestra exactamente a dónde llevan los próximos 15 minutos. Mantenerlo en pantalla unos segundos para que el alumno lo procese.
- Escribir los tres archivos en VS Code en lugar de copiar con `cat` — es más natural y muestra el resaltado de sintaxis HCL instalado en el EP02.
- Al ejecutar `terraform plan`, leer en voz alta la línea `Plan: 1 to add, 0 to change, 0 to destroy` y decir explícitamente: "esta es la línea que hay que leer siempre primero".
- El `terraform.tfstate` con `cat | head -40` — no leer el contenido completo, solo mostrarlo brevemente para recordar la motivación del EP17.
- La conexión final con el `main.tf` real del proyecto es el momento más valioso — quedarse 30 segundos señalando con el cursor cada bloque reconocible.
- Destruir y limpiar antes de cerrar el video — mismo hábito del EP16.

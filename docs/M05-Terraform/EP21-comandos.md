# EP 21: Comandos Esenciales — init, validate, plan, apply, destroy

**Tipo:** PRÁCTICA
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Dominar los cinco comandos del ciclo de vida de Terraform entendiendo cuándo y por qué usar cada uno, de modo que en el EP22 el `terraform apply` de la EC2 definitiva sea completamente predecible.

---

## 📋 Prerequisitos
- Terraform instalado (EP18)
- Backend S3 configurado (EP20)
- El directorio `gitops-infra/infrastructure/terraform/jenkins-ec2/` inicializado

---

## 🧠 El ciclo de vida de Terraform

```
init → validate → fmt → plan → apply → (cambios) → plan → apply → destroy
  │                              │
  │                              └── siempre leer antes de confirmar
  └── una vez por directorio, o al cambiar providers
```

Los signos en el output de `plan` y `apply`:

| Símbolo | Color | Significado |
|---|---|---|
| `+` | Verde | Se va a **crear** |
| `-` | Rojo | Se va a **eliminar** |
| `~` | Amarillo | Se va a **modificar** en lugar |
| `-/+` | Rojo/Verde | Se va a **destruir y recrear** — atención especial |

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: terminal con el directorio `jenkins-ec2/` visible. Un `ls` muestra los tres archivos `.tf`.*

"Bienvenidos al episodio 21.

Tenemos la configuración de la EC2 lista. El backend de S3 configurado. Las credenciales de AWS funcionando. Todo está en su lugar.

Antes de ejecutar `terraform apply` en el EP22 y crear el servidor real, quiero que cada uno de los cinco comandos de Terraform sea completamente predecible. Que cuando el output aparezca en la terminal, sepan exactamente qué significa cada línea y qué hacer con ella.

El EP22 es el episodio de payoff del módulo — el momento en que todo lo que aprendimos llega a algo concreto. No quiero que ese momento tenga ninguna sorpresa técnica. Por eso hacemos este episodio primero.

Empecemos."

---

### `terraform init` — inicializar el directorio (1:00 – 3:00)

> *Pantalla: terminal en `jenkins-ec2/`.*

"Primero: `init`.

Este es siempre el primer comando cuando trabajas con un directorio de Terraform nuevo, o cuando alguien más te pasa un proyecto y necesitas prepararlo en tu máquina. Hace dos cosas: descarga los providers declarados en `required_providers` y configura el backend."

```bash
terraform init
```

"Hay dos líneas en el output que siempre verifico.

La primera confirma que el backend S3 está activo:"

```
Successfully configured the backend "s3"!
```

"La segunda confirma que el provider de AWS se descargó:"

```
- Installing hashicorp/aws v5.x.x...
- Installed hashicorp/aws v5.x.x (signed by HashiCorp)
```

"Si ambas aparecen sin errores, el directorio está listo.

**¿Cuándo volver a ejecutar `init`?** Cuando cambias la versión de un provider, cuando agregas un nuevo provider, o cuando cambias la configuración del backend. No es necesario antes de cada `apply` — solo cuando hay cambios estructurales en la configuración."

---

### `terraform validate` — verificar la sintaxis (3:00 – 5:00)

> *Pantalla: terminal.*

"Segundo: `validate`.

Este comando verifica que el código HCL es correcto antes de hacer ninguna consulta a AWS. Verifica la sintaxis, que los atributos requeridos de cada recurso estén presentes, y que las referencias entre recursos sean válidas. Lo que no hace es consultar AWS para verificar que los valores existen — eso ocurre en el `plan`."

```bash
terraform validate
# Success! The configuration is valid.
```

"Para mostrarles qué detecta, voy a introducir un error deliberado. Entro al `main.tf` y cambio el tipo de recurso `aws_instance` por `aws_instancia` — un typo simple, el tipo de error que se comete sin querer:"

```bash
# Editar main.tf con el typo
terraform validate
```

```
│ Error: Invalid resource type
│
│   on main.tf line 18, in resource "aws_instancia" "prod_server":
│   18: resource "aws_instancia" "prod_server" {
│
│ "aws_instancia" is not a valid resource type in provider hashicorp/aws.
```

"Detectado. Me dice exactamente en qué línea está el error y cuál es el problema. Sin haber hecho ninguna llamada a AWS.

Corrijo el typo, vuelvo a validar:"

```bash
terraform validate
# Success!
```

"**¿Cuándo usarlo?** Antes de cada `plan`. En pipelines de CI/CD es el primer paso del proceso — si la sintaxis falla aquí, el pipeline para antes de consultar AWS."

---

### `terraform fmt` — formatear el código (5:00 – 6:00)

> *Pantalla: terminal.*

"Tercero: `fmt`. Corto pero importante.

Aplica el estilo de indentación y espaciado estándar de HCL a todos los archivos `.tf` del directorio. No cambia la lógica — solo el formato."

```bash
terraform fmt -check    # ver qué necesita formateo, sin modificar nada
terraform fmt           # aplicar el formateo
```

"**¿Cuándo usarlo?** Antes de hacer commit a `gitops-infra`. Mantiene el código consistente sin importar quién lo escribió. En un equipo donde varias personas tocan los mismos archivos de Terraform, el `fmt` evita que los commits tengan cambios de formato mezclados con cambios de lógica."

---

### `terraform plan` — el más importante (6:00 – 10:00)

> *Pantalla: terminal.*

"Cuarto — y el más importante de todos los comandos: `plan`.

Voy a decirlo directo: la regla de oro de Terraform, la que más gente rompe y por la que más accidentes ocurren, es esta: **siempre leer el plan antes de confirmar el apply**. Siempre. Sin excepciones. El `plan` es tu único momento de revisión antes de que algo cambie en AWS.

Lo ejecuto:"

```bash
terraform plan
```

"El output tiene tres partes. La primera lista los recursos con lo que va a pasar a cada uno. La segunda son los detalles de los atributos. Y la tercera es el resumen al final, que es lo que hay que leer primero:

```
Plan: 2 to add, 0 to change, 0 to destroy.
```

"En el EP22, cuando lo ejecutemos sobre la EC2 real, debe decir exactamente eso: `2 to add`. El Security Group y la instancia. Si dice cualquier otra cosa — si dice `1 to destroy`, si dice `3 to add` — hay que parar, revisar, y entender por qué antes de continuar.

---

Ahora los signos. Este es el vocabulario del `plan` y vale la pena memorizar lo que significa cada uno.

El `+` en verde es lo más común — significa que Terraform va a crear ese recurso. Sin problema.

El `-` en rojo significa que Terraform va a eliminar ese recurso. Aquí hay que prestar atención y entender por qué.

El `~` en amarillo significa modificación en lugar — el recurso sigue existiendo pero con algún atributo cambiado. Por ejemplo, cambiar las etiquetas de una instancia genera un `~`.

Y el `-/+` es el que más importa entender: Terraform va a **destruir el recurso y recrearlo desde cero**. No puede hacer el cambio en caliente. Aparece cuando cambias algo que AWS no permite modificar en una instancia existente — por ejemplo, el tipo de AMI, o el tipo de cifrado del disco.

¿Por qué importa tanto el `-/+`? Porque significa downtime. Si tienes la app corriendo y aparece un `-/+` en el plan, estás a punto de destruir el servidor y perder la disponibilidad durante el tiempo que tarde en recrearse. Si no lees el plan, no te enteras hasta que ya ocurrió."

---

### `terraform apply` — crear infraestructura (10:00 – 12:00)

> *Pantalla: terminal.*

"Quinto: `apply`."

```bash
terraform apply
```

"Terraform muestra el plan de nuevo — el mismo que viste antes — y pide confirmación explícita:"

```
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value:
```

"Escribo `yes`. Solo `yes`. No Enter solo, no `y`, no `si`. Terraform acepta exactamente la cadena `yes` como confirmación.

Esta pausa existe por una razón. Es el momento de leer el plan una vez más antes de que algo ocurra en AWS. Si algo en el plan te sorprende — si hay un `-/+` que no esperabas, si hay más recursos de los que pensabas — escribe `no`, investiga, y vuelve a correr el plan.

---

Para pipelines automatizados, donde no hay nadie que pueda leer el plan y confirmar manualmente, existe la flag `-auto-approve`:"

```bash
terraform apply -auto-approve
```

"En el EP22 la vamos a usar porque el video sería muy largo parar a escribir `yes` y queremos mostrar el flujo limpio. Pero en cualquier otro contexto, y especialmente en tu trabajo diario, la confirmación manual es el default correcto.

Al terminar el apply, los outputs aparecen al final:"

```
Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:
prod_public_ip = "54.x.x.x"
```

"Para verlos en cualquier momento posterior sin volver a ejecutar apply:"

```bash
terraform output
terraform output prod_public_ip   # solo un output específico
```

---

### `terraform destroy` — eliminar todo (12:00 – 13:30)

> *Pantalla: terminal.*

"Sexto y último: `destroy`.

Este comando elimina todo lo que Terraform creó con esta configuración. Muestra el plan inverso — todos los recursos marcados con `-` en rojo — y pide la misma confirmación explícita."

```bash
terraform destroy
```

"Dos cosas importantes sobre el `destroy` en el contexto del curso.

La primera: en el EP22 crearemos la EC2 definitiva y **no la vamos a destruir hasta el EP49**. Cuando llegue ese momento, el `terraform destroy` en el directorio `jenkins-ec2/` va a eliminar la instancia y el Security Group de forma limpia.

La segunda: el `destroy` no elimina el bucket S3 del backend. Eso es intencional. El `prevent_destroy = true` que vimos en el EP20 lo protege explícitamente. El bucket puede quedarse en tu cuenta de AWS indefinidamente — su costo mensual es prácticamente cero — y te permite retomar el curso en cualquier momento."

---

### Referencia rápida (13:30 – 14:00)

> *Pantalla: slide o terminal con la tabla.*

"Para cerrar, el cheatsheet completo:

| Comando | Cuándo usarlo |
|---|---|
| `terraform init` | Primera vez en un directorio, o al cambiar providers |
| `terraform validate` | Antes de cada plan — detecta errores de sintaxis |
| `terraform fmt` | Antes de hacer commit |
| `terraform plan` | **Siempre** antes del apply — leer el output completo |
| `terraform apply` | Después de revisar el plan |
| `terraform output` | Para ver los valores de salida después del apply |
| `terraform destroy` | Al terminar — elimina todo lo que Terraform creó |"

---

### CIERRE (14:00 – 15:00)

"Eso es el episodio 21.

Ahora el ciclo es predecible. Saben qué hace cada comando, cuándo ejecutarlo, y qué significa cada línea del output. El `+`, el `-`, el `~`, el `-/+` — ya no son signos misteriosos, son información concreta sobre qué va a ocurrir en AWS.

En el siguiente episodio ejecutamos ese ciclo sobre la configuración real: la EC2 t2.micro con el Security Group de tres puertos, el disco de 30 GB, el servidor de producción del curso a costo cero. El episodio de payoff del módulo.

Nos vemos en el EP22."

---

## ✅ Checklist de Verificación
- [ ] Entiendes cuándo ejecutar `terraform init` y por qué no es necesario antes de cada apply
- [ ] `terraform validate` detecta el typo en el tipo de recurso antes del plan
- [ ] Sabes leer el output de `terraform plan` e identificar `+`, `-`, `~` y `-/+`
- [ ] Entiendes la diferencia entre `-` (eliminar) y `-/+` (destruir y recrear)
- [ ] Sabes que `terraform destroy` no borra el bucket S3 del backend gracias al `prevent_destroy`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Error: Backend initialization required` | Ejecutar `terraform init` primero |
| `Error acquiring the state lock` | Lock atascado de un apply previo — `terraform force-unlock LOCK_ID` |
| `Error: Reference to undeclared resource` | Typo en una referencia entre recursos — `terraform validate` lo muestra |
| El plan muestra `-/+` inesperado | Algún atributo cambió que AWS no puede modificar en caliente — revisar qué cambió en los `.tf` |

---

## 🗒️ Notas de Producción
- Para el ejemplo del typo en `validate`, introducirlo en pantalla en vivo — es más didáctico que tenerlo preparado de antemano.
- Al ejecutar `terraform plan`, hacer zoom en el output de la terminal. Los signos `+`, `-`, `~`, `-/+` tienen que ser legibles claramente — si el texto es muy pequeño, el alumno pierde el momento más pedagógico del episodio.
- La tabla de signos merece quedarse en pantalla al menos 5 segundos mientras la comentas — es información de referencia que el alumno va a necesitar recordar.
- La distinción entre `apply` con confirmación manual y `apply -auto-approve` merece énfasis verbal claro: "en el EP22 lo uso con `-auto-approve`, pero en tu trabajo diario siempre con confirmación manual".
- El cheatsheet del cierre puede presentarse como slide fija mientras hablas del EP22 como transición visual.

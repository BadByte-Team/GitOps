# 🎬 Guión — EP15: Instalación y Configuración de AWS CLI

**Duración estimada:** 10–12 min
**Tono:** Práctico y fluido. El alumno ya tiene credenciales — este episodio es técnico pero rápido. Se convierte en una sesión de comandos con explicaciones puntuales.

---

## 🎙️ INTRO (0:00 – 0:40)

> *Pantalla: terminal limpia.*

"Bienvenidos al episodio 15.

En los dos episodios anteriores creamos la cuenta de AWS y el usuario IAM con sus credenciales. Ahora necesitamos la herramienta que nos va a permitir hablar con AWS desde la terminal: **la AWS CLI**.

La CLI es fundamental para este curso. Terraform la usa para autenticarse. Los scripts de instalación la usan para configurar recursos. Y nosotros la vamos a usar para verificar, depurar y limpiar durante todo el camino.

El episodio es rápido: instalar, verificar, y ver los comandos que más vamos a usar. Vamos."

---

## 📌 PASO 1 — Instalar AWS CLI v2 (0:40 – 3:00)

> *Pantalla: terminal en la PC local.*

"La CLI tiene dos versiones — v1 y v2. Siempre la v2. Es más rápida, tiene mejor manejo de perfiles, y es la que AWS sigue actualizando activamente.

Instalación en Linux:"

```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
sudo apt install -y unzip
unzip awscliv2.zip
sudo ./aws/install
rm -rf aws awscliv2.zip
```

"Verifico:"

```bash
aws --version
```

"Debe responder algo como `aws-cli/2.x.x Python/3.x.x Linux/x86_64`. Si ves la v2, perfecto.

---

Dos notas rápidas para otros sistemas.

En **Arch Linux** con pacman: `sudo pacman -S aws-cli`. Si usas yay: `yay -S aws-cli-v2`.

En **macOS** con Homebrew: `brew install awscli`.

En **Windows**: descargar el instalador MSI desde la documentación oficial de AWS y ejecutarlo. Después, verificar en PowerShell con `aws --version`."

---

## 📌 PASO 2 — Verificar las credenciales (3:00 – 4:30)

> *Pantalla: terminal.*

"Si completaste el EP14, las credenciales ya están en `~/.aws/credentials`. La CLI las toma desde ahí automáticamente. No hay que hacer `aws configure` si ya están guardadas.

Verifico que todo está en orden:"

```bash
cat ~/.aws/credentials
cat ~/.aws/config
```

"Debe mostrar el Access Key ID, el Secret Key, la región `us-east-2` y el output `json`.

Si por alguna razón los archivos no existen, lo más rápido es:"

```bash
aws configure
```

"Este comando interactivo pide el Access Key ID, el Secret Access Key, la región y el formato. Llena los mismos valores del EP14."

---

## 📌 PASO 3 — El comando de verificación de oro (4:30 – 5:30)

> *Pantalla: terminal.*

"El primer comando que ejecuto siempre que quiero confirmar que la CLI está funcionando bien con las credenciales correctas es este:"

```bash
aws sts get-caller-identity
```

"STS es el servicio de Security Token Service. Este comando pregunta: '¿quién soy yo desde el punto de vista de AWS?'

La respuesta es un JSON con tres campos:

- **UserId** — el identificador interno del usuario IAM.
- **Account** — el número de tu cuenta de AWS. Lo verás muchas veces en los próximos episodios.
- **Arn** — el ARN completo del usuario. Aquí puedes ver que dice `user/admin-curso` — confirma que estamos autenticados con el usuario correcto, no con el root.

Memoriza este comando. Cada vez que Terraform o la CLI se comporte raro, este es el primer diagnóstico."

---

## 📌 PASO 4 — Comandos que usaremos en el curso (5:30 – 9:30)

> *Pantalla: terminal, ejecutando cada comando con pausa para ver el output.*

"En lugar de explicar la CLI en abstracto, voy a mostrar exactamente los comandos que van a aparecer en los próximos episodios. Así cuando los veas, ya los reconocerás.

---

**EC2 — verificar instancias corriendo**

En el EP22 crearemos la EC2 con Terraform. Cuando quiera verificar que está corriendo y obtener su IP pública, usaré esto:"

```bash
aws ec2 describe-instances \
  --filters "Name=instance-state-name,Values=running" \
  --query "Reservations[].Instances[].[InstanceId,PublicIpAddress,Tags[?Key=='Name'].Value|[0]]" \
  --output table
```

"El flag `--query` es como un filtro JMESPath — te permite sacar solo los campos que necesitas del JSON enorme que devuelve AWS. El `--output table` lo pone en formato legible.

---

**EC2 — obtener la IP de nuestra instancia por nombre**

Esta forma la usaré constantemente para obtener la IP del servidor K3s después de crearlo:"

```bash
aws ec2 describe-instances \
  --filters "Name=tag:Name,Values=Produccion-K3s" \
  --query "Reservations[0].Instances[0].PublicIpAddress" \
  --output text
```

"El tag `Produccion-K3s` es el nombre que le damos a la instancia en el Terraform del EP22. Cuando necesite la IP para hacer SSH o configurar ArgoCD, este es el comando.

---

**S3 — verificar el bucket del backend**

Después de crear el backend en el EP17, para confirmar que existe:"

```bash
aws s3 ls
```

"Y para ver qué hay dentro del bucket de state:"

```bash
aws s3 ls s3://curso-gitops-terraform-state/
```

"---

**Key Pairs — crear el par de claves para la EC2**

En el EP22 crearemos el key pair para SSH con este comando:"

```bash
aws ec2 create-key-pair \
  --key-name aws-key \
  --query 'KeyMaterial' \
  --output text > aws-key.pem

chmod 400 aws-key.pem
```

"Lo muestro ahora para que cuando llegue el momento lo reconozcas. `--query 'KeyMaterial'` extrae solo la clave privada del JSON y la redirige directamente al archivo `.pem`. El `chmod 400` hace que solo tu usuario pueda leerlo — SSH rechaza los archivos con permisos más permisivos.

---

**DynamoDB — verificar la tabla de locking**

Después de crear el backend en el EP17:"

```bash
aws dynamodb list-tables

aws dynamodb describe-table \
  --table-name curso-gitops-terraform-locks \
  --query "Table.TableStatus" \
  --output text
```

"Si responde `ACTIVE`, la tabla está lista para que Terraform la use para bloquear el estado."

---

## 🔍 PAUSA CONCEPTUAL — Sobre el flag --query (9:30 – 10:30)

> *Pantalla: terminal, mostrando la diferencia entre con y sin query.*

"Quiero que noten la diferencia entre ejecutar un comando con y sin el flag `--query`.

Sin query:"

```bash
aws ec2 describe-instances --output json | head -50
```

"Cientos de líneas de JSON. Información sobre cada instancia, cada interfaz de red, cada tag, cada atributo. Útil para debugging, pero abrumador para el uso diario.

Con query:"

```bash
aws ec2 describe-instances \
  --query "Reservations[].Instances[].[InstanceId,PublicIpAddress]" \
  --output table
```

"Solo lo que necesito. Dos columnas, una línea por instancia.

No es necesario que memoricen la sintaxis de JMESPath ahora mismo. Yo ya incluyo el `--query` correcto en cada comando del curso. Pero saber que existe y que es lo que convierte el output de AWS en algo legible — eso vale la pena entenderlo."

---

## 🎙️ CIERRE (10:30 – 11:30)

"Eso es EP15.

AWS CLI instalada, credenciales verificadas, y ya tienes en la cabeza los comandos que más vas a ver en los próximos episodios. Cuando aparezca `aws ec2 describe-instances` o `aws s3 ls`, ya sabes para qué sirven.

En el siguiente episodio vamos a usar la CLI en serio: vamos a crear un Security Group, lanzar una instancia EC2 Ubuntu y conectarnos por SSH por primera vez. No es la instancia definitiva del curso — esa la crearemos con Terraform en el EP22 — pero necesitamos entender los conceptos de EC2 antes de dejarle todo ese trabajo a Terraform.

Nos vemos en el EP16."

---

## 🗒️ Notas de producción

- Este episodio es de ritmo rápido — no hay mucho que pueda salir mal. Mantenerlo ágil.
- La parte del `--query` con la comparación visual del JSON completo vs la tabla limpia es el momento pedagógico más fuerte del episodio. Considera poner ambas terminales en pantalla dividida.
- Al crear el key pair de práctica, asegurarse de que el `chmod 400` se vea ejecutado — ese detalle lo van a repetir en el EP22.
- Si hay tiempo, mostrar `aws configure list` para que vean un resumen de la configuración activa — es un comando de diagnóstico útil que la gente suele no conocer.

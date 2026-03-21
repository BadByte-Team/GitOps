# 🎬 Guión — EP16: EC2 — Conceptos, Lanzar Instancia y Conectar por SSH

**Duración estimada:** 15–18 min
**Tono:** Práctico con pausas conceptuales cortas. El alumno ve por primera vez los componentes de EC2 — AMI, Security Group, Key Pair. La práctica es temporal y explícitamente distinta de la instancia real del EP22.

---

## 🎙️ INTRO (0:00 – 0:50)

> *Pantalla: consola de AWS, servicio EC2.*

"Bienvenidos al episodio 16.

Llevamos varios episodios configurando cosas — la cuenta de AWS, IAM, la CLI. Todo necesario, pero todavía no hemos creado nada en la nube.

Hoy eso cambia. Vamos a lanzar nuestra primera instancia EC2 y a conectarnos a ella por SSH.

Aclaración importante antes de arrancar: la instancia que creamos hoy es de **práctica**. La vamos a lanzar para entender los conceptos, conectarnos, explorar, y luego la terminamos. La instancia **real** del curso — el servidor que va a correr K3s, ArgoCD y nuestra app — la crearemos con Terraform en el EP22. La diferencia entre hacerlo a mano y con Terraform es exactamente la lección de ese episodio.

Pero para entender Terraform, primero hay que entender qué está automatizando. Eso es lo que hacemos aquí.

Empecemos."

---

## 🔍 PAUSA CONCEPTUAL — Los cuatro componentes de una EC2 (0:50 – 3:30)

> *Pantalla: diagrama simple o slide con los cuatro componentes.*

"Antes de ejecutar cualquier comando, necesitas entender los cuatro ingredientes que componen una instancia EC2. Sin esto, los comandos de la CLI son solo magia negra.

---

**AMI — Amazon Machine Image**

La AMI es la fotografía del sistema operativo con la que arranca tu instancia. Es como el ISO de Ubuntu que usarías para instalar Linux — pero ya empaquetado y listo para la nube.

AWS tiene AMIs oficiales de Ubuntu, Amazon Linux, Windows, y muchas más. Nosotros usaremos **Ubuntu 22.04 LTS** en la región us-east-1, cuyo ID es `ami-0c7217cdde317cfec`.

Ese ID lo verán en el Terraform del EP22. Ya saben qué es.

---

**Tipo de instancia**

Define cuánta CPU y RAM tiene tu servidor.

| Tipo | vCPU | RAM | Free Tier |
|---|---|---|---|
| t2.micro | 1 | 1 GB | ✅ 750 hrs/mes |
| t2.medium | 2 | 4 GB | ❌ ~$33/mes |

Usaremos **t2.micro** en todo el curso. Sí, 1 GB de RAM es poco. Por eso en el EP28 vamos a configurar Swap — memoria virtual que le permite a K3s y ArgoCD coexistir sin matar la instancia. Ya llegamos a eso.

---

**Security Group**

Un Security Group es un firewall virtual que controla qué tráfico puede entrar y salir de tu instancia. Es una lista de reglas: `protocolo + puerto + origen`.

Para nuestra instancia definitiva, necesitaremos tres puertos abiertos:

| Puerto | Uso |
|---|---|
| 22 | SSH — para conectarnos |
| 30080 | ArgoCD — interfaz web |
| 30081 | Nuestra app Go |

Hoy, para la práctica, solo abrimos el 22.

---

**Key Pair**

Las instancias Linux en AWS no usan contraseñas para SSH — usan criptografía de clave pública. AWS pone tu clave pública en la instancia al crearla. Tú te conectas con tu clave privada — el archivo `.pem`.

Si pierdes el `.pem`, pierdes el acceso a la instancia. No hay recovery. Por eso el primer comando después de crearlo es siempre `chmod 400`."

---

## 📌 PASO 1 — Crear el Key Pair (3:30 – 5:00)

> *Pantalla: terminal, creando una carpeta de práctica.*

"Voy a trabajar en una carpeta temporal para esta práctica:"

```bash
mkdir ~/practica-ec2 && cd ~/practica-ec2

aws ec2 create-key-pair \
  --key-name practica-ep16 \
  --query 'KeyMaterial' \
  --output text > practica-ep16.pem

chmod 400 practica-ep16.pem

ls -la
# practica-ep16.pem  ← con permisos -r--------
```

"El flag `--query 'KeyMaterial'` extrae solo la clave privada del JSON y la redirige al archivo. Sin ese flag, el JSON completo iría al archivo y el `.pem` no serviría.

Verifico que la key existe en AWS:"

```bash
aws ec2 describe-key-pairs \
  --key-names practica-ep16 \
  --query "KeyPairs[0].KeyName" \
  --output text
# practica-ep16
```

---

## 📌 PASO 2 — Crear el Security Group (5:00 – 7:00)

> *Pantalla: terminal.*

"Necesito el ID de la VPC por defecto — la red virtual donde vive la instancia:"

```bash
VPC_ID=$(aws ec2 describe-vpcs \
  --filters "Name=isDefault,Values=true" \
  --query 'Vpcs[0].VpcId' \
  --output text)

echo "VPC: $VPC_ID"
```

"Ahora creo el Security Group:"

```bash
SG_ID=$(aws ec2 create-security-group \
  --group-name practica-sg-ep16 \
  --description "SG temporal para practica EP16" \
  --vpc-id $VPC_ID \
  --query 'GroupId' \
  --output text)

echo "Security Group: $SG_ID"
```

"Y abro el puerto 22 para SSH:"

```bash
aws ec2 authorize-security-group-ingress \
  --group-id $SG_ID \
  --protocol tcp \
  --port 22 \
  --cidr 0.0.0.0/0
```

"El `0.0.0.0/0` significa 'desde cualquier IP'. Para producción usarías tu IP específica. Para este ejercicio de práctica, lo dejamos abierto."

---

## 📌 PASO 3 — Lanzar la instancia (7:00 – 9:00)

> *Pantalla: terminal.*

"Con el key pair y el Security Group listos, lanzo la instancia:"

```bash
INSTANCE_ID=$(aws ec2 run-instances \
  --image-id ami-0c7217cdde317cfec \
  --instance-type t2.micro \
  --key-name practica-ep16 \
  --security-group-ids $SG_ID \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=practica-ep16}]' \
  --query 'Instances[0].InstanceId' \
  --output text)

echo "Instancia: $INSTANCE_ID"
```

"La instancia tarda entre 30 y 60 segundos en pasar de `pending` a `running`. Esperamos:"

```bash
echo "Esperando que la instancia esté lista..."
aws ec2 wait instance-running --instance-ids $INSTANCE_ID
echo "✅ Lista"

PUBLIC_IP=$(aws ec2 describe-instances \
  --instance-ids $INSTANCE_ID \
  --query 'Reservations[0].Instances[0].PublicIpAddress' \
  --output text)

echo "IP Pública: $PUBLIC_IP"
```

"Guardo esa IP — la necesito para SSH."

---

## 📌 PASO 4 — Conectar por SSH (9:00 – 11:30)

> *Pantalla: terminal.*

"Espero 30 segundos más — la instancia está `running` pero el sistema operativo todavía está terminando de arrancar:"

```bash
sleep 30

ssh -i practica-ep16.pem ubuntu@$PUBLIC_IP
```

"La primera vez pregunta:"

```
The authenticity of host 'X.X.X.X' can't be established.
Are you sure you want to continue connecting (yes/no)? yes
```

"Escribo `yes`. Esto le dice a SSH que confíe en este host la primera vez. A partir de ahora, lo recuerda.

Estoy dentro. El prompt cambió — dice `ubuntu@ip-...`. Estoy en el servidor de AWS.

Exploro rápido:"

```bash
# ¿Qué sistema operativo es?
cat /etc/os-release
# Ubuntu 22.04.x LTS

# ¿Cuánta RAM tiene?
free -h
# Mem: 981Mi — un poco menos de 1GB, como esperábamos

# ¿Cuánto espacio en disco?
df -h
# 8GB por defecto — en el EP22 lo ampliaremos a 30GB con Terraform

# ¿Cuántos CPUs?
nproc
# 1

exit
```

"Esos números los veremos de nuevo en el EP22 y el EP28. Ahora ya saben a qué se refieren cuando mencionamos '1GB de RAM' y por qué el Swap es necesario."

---

## 🔍 PAUSA CONCEPTUAL — EP16 vs EP22 (11:30 – 13:00)

> *Pantalla: tabla comparativa o slide.*

"Quiero ser claro sobre la diferencia entre lo que hicimos hoy y lo que haremos en el EP22.

| Aspecto | Hoy (práctica manual) | EP22 (Terraform real) |
|---|---|---|
| Método | 6 comandos de CLI | 1 comando: `terraform apply` |
| Puertos | Solo SSH (22) | SSH + 30080 (ArgoCD) + 30081 (App) |
| Disco | 8 GB por defecto | 30 GB gp3 encriptado |
| Nombre | `practica-ep16` | `Produccion-K3s` |
| Destino | Se destruye hoy | Vive hasta el EP49 |

Lo que aprendiste hoy — AMI, tipo de instancia, Security Group, Key Pair — es exactamente lo que Terraform va a automatizar en el EP22. Cuando veas el `main.tf` de ese episodio, vas a reconocer cada campo porque lo hiciste a mano primero."

---

## 📌 PASO 5 — Destruir todo al terminar (13:00 – 15:00)

> *Pantalla: terminal.*

"Ahora, lo más importante del episodio: **terminar la instancia**.

Este es el hábito que separa a quien sabe usar AWS de quien recibe facturas sorpresa. Cada instancia EC2 genera costos mientras está `running` — aunque no la estés usando. La única forma de que deje de cobrar es terminarla.

Para esta práctica, con el Free Tier no habría costo. Pero el hábito importa. En el EP22 la instancia definitiva va a ser free tier también — pero cuando lleguemos al EP49, hay que saber destruirla."

```bash
aws ec2 terminate-instances --instance-ids $INSTANCE_ID
echo "Instancia terminada"

# Esperar a que termine completamente
aws ec2 wait instance-terminated --instance-ids $INSTANCE_ID

# Limpiar el Security Group
aws ec2 delete-security-group --group-id $SG_ID

# Limpiar el Key Pair de AWS
aws ec2 delete-key-pair --key-name practica-ep16

# Limpiar el archivo local
rm ~/practica-ec2/practica-ep16.pem
rmdir ~/practica-ec2

echo "✅ Todo limpio"
```

"Verifico que no quedó nada:"

```bash
aws ec2 describe-instances \
  --filters "Name=instance-state-name,Values=running" \
  --query "Reservations[].Instances[].InstanceId" \
  --output text
# (silencio — no hay instancias corriendo)
```

---

## 🎙️ CIERRE (15:00 – 16:00)

"Eso es EP16. Lanzaste tu primera instancia EC2, te conectaste por SSH y la destruiste limpiamente.

Ya conoces los cuatro ingredientes de EC2: AMI, tipo de instancia, Security Group, Key Pair. Y ya tienes el hábito más importante de AWS: terminar lo que no usas.

En el siguiente episodio creamos los dos recursos que necesita Terraform para funcionar de forma remota: el bucket S3 para guardar el estado y la tabla DynamoDB para el locking. Esos dos recursos son permanentes — van a sobrevivir todo el curso.

Nos vemos en el EP17."

---

## 🗒️ Notas de producción

- La pausa conceptual de los cuatro componentes es el núcleo teórico del episodio. Considera una diapositiva con el diagrama en lugar de solo texto en la terminal.
- Al conectar por SSH, mostrar el cambio de prompt en la terminal claramente — es el momento visual más satisfactorio del episodio.
- El `free -h` dentro de la instancia mostrando ~981MB es un buen momento para anticipar el EP28 de Swap: "y aquí es donde en el EP28 vamos a agregar 2GB de memoria virtual para que K3s no colapse esto".
- La tabla comparativa EP16 vs EP22 puede quedar en pantalla unos segundos extra — es el puente pedagógico más importante del episodio.
- Ejecutar la limpieza completa al final del video en lugar de cortarlo antes — refuerza el hábito visualmente.

# EP 18: Instalación de Terraform

**Tipo:** INSTALACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar Terraform desde el repositorio oficial de HashiCorp, verificar la instalación, y entender qué es Terraform y por qué lo usamos en lugar de crear recursos manualmente con la CLI.

---

## 📋 Prerequisitos
- AWS CLI configurado con credenciales de `admin-curso` (EP15)
- Terminal con acceso a internet

---

## 🧠 ¿Qué es Terraform y por qué lo usamos?

En el EP16 creaste una EC2 manualmente con 6 comandos de AWS CLI. Funcionó, pero tiene tres problemas:

| Problema | CLI manual | Terraform |
|---|---|---|
| Reproducibilidad | No — si pierdes los comandos, no sabes cómo recrearlo | Sí — el `.tf` describe exactamente el estado deseado |
| Historial | Ninguno — nadie sabe qué se creó ni cuándo | Git — cada cambio queda registrado en el repositorio |
| Destrucción | Manual — tienes que recordar qué creaste y en qué orden | `terraform destroy` — elimina todo en el orden correcto |

Terraform usa el paradigma **IaC — Infrastructure as Code**. Tu infraestructura se describe en archivos de texto, se versiona en Git, y se aplica de forma reproducible en cualquier máquina con acceso a las credenciales correctas.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: terminal con el historial de comandos del EP16 — el `create-key-pair`, el `create-security-group`, el `authorize-security-group-ingress`, el `run-instances`. Los seis comandos visibles en pantalla.*

"Bienvenidos al episodio 18.

Esto es lo que hicimos en el EP16 para lanzar una instancia EC2. Seis comandos en secuencia. Y eso fue solo la instancia de práctica — sin el Security Group completo con los tres puertos, sin el disco de 30 GB, sin las etiquetas de nombre.

Ahora imagínate esto: terminaste el curso, destruiste todo en el EP49, y tres semanas después decides volver a practicar. Quieres recrear exactamente el mismo servidor. ¿Recuerdas el ID exacto de la AMI de Ubuntu? ¿Qué CIDR tenía el Security Group? ¿En qué orden creaste los recursos para que no hubiera dependencias rotas?

Con comandos de CLI, la respuesta honesta es: probablemente no. Tendrías que buscar todo de nuevo, reconstruirlo de memoria, y rezar para que no te hayas olvidado de algo.

Con Terraform, la respuesta es siempre la misma: está en el repositorio. Abres el `main.tf`, lees exactamente lo que dice, ejecutas `terraform apply`, y en dos minutos tienes exactamente lo mismo. Sin esfuerzo, sin dudas, sin sorpresas.

Eso es lo que vamos a configurar hoy. Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — IaC: el problema que resuelve Terraform (1:30 – 3:30)

> *Pantalla: slide o diagrama con dos columnas — CLI manual a la izquierda, Terraform a la derecha.*

"Antes de instalar nada, quiero que entiendas exactamente qué problema resuelve Terraform, porque si no entiendes el problema, la herramienta parece innecesariamente complicada.

El paradigma se llama **Infrastructure as Code** — Infraestructura como Código. La idea central es esta: tu infraestructura — los servidores, las redes, los discos, los permisos — se describe en archivos de texto igual que el código de tu aplicación. Y esos archivos viven en Git.

Eso tiene tres consecuencias inmediatas que importan para el curso.

La primera es **reproducibilidad**. Si el archivo dice 'quiero una EC2 t3.micro con Ubuntu 22.04 y 30 GB de disco', Terraform va a crear exactamente eso. Hoy, mañana, en otra cuenta de AWS. Siempre lo mismo. El entorno de producción del curso se puede recrear desde cero en dos minutos con un solo comando.

La segunda es **historial**. Cada cambio a los archivos de infraestructura es un commit de Git. Si en algún momento alguien cambia el tipo de instancia de `t3.micro` a algo más caro, ese cambio queda registrado en el repositorio con fecha y mensaje de commit. No hay cambios silenciosos que nadie recuerde haber hecho.

La tercera es **destrucción limpia**. Cuando en el EP49 terminemos el curso, `terraform destroy` va a eliminar todo lo que Terraform creó, en el orden correcto, sin dejar recursos huérfanos que sigan generando costos. Sin tener que recordar manualmente qué existía ni en qué secuencia eliminarlo.

Para este curso, Terraform hace exactamente una cosa: crear la EC2 t3.micro donde vive K3s, con el Security Group correcto. Pero el patrón que aprendes hoy es el mismo que usan equipos de ingeniería para gestionar cientos de recursos en producción."

---

### PASO 1 — Instalar Terraform (3:30 – 7:00)

> *Pantalla: terminal en la PC local.*

"Vamos a instalar desde el repositorio oficial de HashiCorp — la empresa que creó Terraform. Hay paquetes de `terraform` en algunos gestores de paquetes del sistema que pueden estar desactualizados o no ser la versión oficial, así que siempre voy a la fuente.

El proceso tiene tres partes: importar la clave GPG para verificar la autenticidad de los paquetes, agregar el repositorio a la lista de fuentes del sistema, e instalar.

Primero las dependencias necesarias para manejar repositorios externos:"

```bash
sudo apt-get install -y gnupg software-properties-common
```

"Ahora importo la clave criptográfica de HashiCorp. Esto le dice a tu sistema operativo que confíe en los paquetes firmados por HashiCorp. Sin este paso, `apt` rechazaría el repositorio:"

```bash
wget -O- https://apt.releases.hashicorp.com/gpg | \
  gpg --dearmor | \
  sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
```

"Agrego el repositorio oficial a la lista de fuentes. El `$(lsb_release -cs)` detecta automáticamente la versión de Ubuntu que tienes — así el mismo comando funciona en Focal, Jammy, o cualquier versión futura:"

```bash
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
  https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
  sudo tee /etc/apt/sources.list.d/hashicorp.list
```

"Actualizo e instalo:"

```bash
sudo apt-get update && sudo apt-get install -y terraform
```

"Para los que usan **Arch Linux**, es mucho más directo:"

```bash
sudo pacman -S terraform
```

"Y para **macOS** con Homebrew:"

```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

"Sin importar el sistema, verifico de la misma forma:"

```bash
terraform --version
```

"La respuesta debe ser algo como `Terraform v1.7.x`. Si ves la versión, perfecto. Si ves `command not found`, el directorio de instalación no está en el PATH — ejecuta `echo $PATH` y verifica que `/usr/local/bin` aparece en la lista."

---

### PASO 2 — Autocompletado (7:00 – 8:00)

> *Pantalla: terminal.*

"Un paso opcional pero que agradecerás mucho durante los primeros días: el autocompletado.

Terraform tiene bastantes subcomandos — `init`, `plan`, `apply`, `destroy`, `output`, `state`, `fmt`, `validate`... Al principio cuesta recordarlos todos. Con el autocompletado, escribes `terraform` seguido de un espacio, presionas Tab, y la terminal te muestra todas las opciones disponibles. Es como tener el manual integrado en la terminal."

```bash
terraform -install-autocomplete
source ~/.bashrc
```

"Listo. Pruébalo ahora mismo: escribe `terraform` seguido de un espacio y presiona Tab. Deberían aparecer todos los subcomandos disponibles."

---

### PASO 3 — Los tres archivos estándar (8:00 – 9:30)

> *Pantalla: VS Code abriendo `gitops-infra/infrastructure/terraform/jenkins-ec2/`.*

"Antes de escribir cualquier código, quiero mostrarte la estructura que vamos a usar. Es una convención del ecosistema de Terraform — prácticamente todo proyecto en el mundo sigue este mismo patrón de tres archivos.

El `main.tf` es el corazón. Aquí van los recursos que Terraform va a crear en AWS: la EC2, el Security Group, lo que sea.

El `variables.tf` parametriza los valores que pueden cambiar entre entornos. En lugar de escribir `t3.micro` directamente dentro del resource, lo conviertes en una variable y lo usas con `var.instance_type`. Si mañana necesitas cambiar el tipo de instancia, cambias un solo lugar y el cambio se propaga a todo el archivo.

El `outputs.tf` define qué valores te devuelve Terraform cuando termina de crear todo. Por ejemplo, después de crear la EC2, Terraform imprime la IP pública en la terminal. Eso es un output.

Abro el directorio de la EC2 en VS Code. Ahí están los tres archivos. Hoy no los leemos en detalle — eso es el EP19. Solo quiero que vean que el patrón que acabamos de describir ya está implementado. Cuando lleguen a leerlos en el próximo episodio, ya saben exactamente qué rol cumple cada archivo."

---

### PASO 4 — Verificar que Terraform se autentica con AWS (9:30 – 11:00)

> *Pantalla: terminal, creando un directorio de prueba temporal.*

"Por último, un test rápido para confirmar que Terraform puede hablar con AWS usando las credenciales del EP14.

Algo importante que quiero aclarar: Terraform busca las credenciales automáticamente en `~/.aws/credentials`. No hay que configurar nada adicional en los archivos de Terraform, no hay variables de entorno que setear. Si la AWS CLI funciona con tu usuario `admin-curso`, Terraform también funciona con ese mismo usuario — leen del mismo archivo.

Creo un directorio de prueba temporal para el test:"

```bash
mkdir ~/terraform-test && cd ~/terraform-test
```

"Creo un `main.tf` mínimo — solo el bloque del provider de AWS, sin declarar ningún recurso. Solo quiero que Terraform initialize y descargue el plugin:"

```bash
cat > main.tf << 'EOF'
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-2"
}
EOF
```

"Ejecuto `terraform init`. Este comando descarga el provider de AWS — el plugin que contiene toda la lógica para crear, modificar y destruir recursos en la API de Amazon. La primera vez puede tardar 30 segundos porque descarga un binario:"

```bash
terraform init
```

"Verás algo así en el output:"

```
Initializing the backend...
Initializing provider plugins...
- Finding hashicorp/aws versions matching "~> 5.0"...
- Installing hashicorp/aws v5.x.x...
- Installed hashicorp/aws v5.x.x (signed by HashiCorp)

Terraform has been successfully initialized!
```

"Esa última línea, `Terraform has been successfully initialized!`, es la confirmación. El provider de AWS está descargado y las credenciales están accesibles.

Limpio el directorio de prueba — no lo necesito más y los archivos `.terraform/` que genera el init ocupan espacio:"

```bash
cd ~ && rm -rf ~/terraform-test
```

---

### CIERRE (11:00 – 12:00)

"Eso es el episodio 18.

Terraform instalado, autocompletado activado, la estructura de tres archivos clara en la cabeza, y la autenticación con AWS verificada. Todo listo para escribir código real.

En el siguiente episodio vamos a escribir nuestro primer recurso de Terraform desde cero — variables, un provider, una instancia EC2 de ejemplo. El objetivo no es crear infraestructura útil todavía, sino entender cómo se conectan las piezas. Para que cuando lleguemos al EP22 y ejecutemos el `terraform apply` del servidor real, cada línea del `main.tf` sea completamente familiar.

Nos vemos en el EP19."

---

## ✅ Checklist de Verificación
- [ ] `terraform --version` muestra v1.7.x o superior
- [ ] El autocompletado está activado — Tab completa los subcomandos
- [ ] `terraform init` en un directorio con el provider de AWS completa sin errores
- [ ] Entiendes el rol de cada uno de los tres archivos: `main.tf`, `variables.tf`, `outputs.tf`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `terraform: command not found` | Verificar que `/usr/local/bin` está en el PATH: `echo $PATH` |
| `Error: No valid credential sources found` | Verificar que `~/.aws/credentials` existe con las claves del EP14 |
| `wget` falla al importar la clave GPG | Verificar conexión a internet y que `gnupg` está instalado: `sudo apt install gnupg` |
| Versión instalada desactualizada | `sudo apt remove terraform` y reinstalar desde el repositorio de HashiCorp |

---

## 🗒️ Notas de Producción
- La apertura con el historial de comandos del EP16 en pantalla es el gancho narrativo — mantenerla visible mientras hablas durante toda la intro, no cambiar de pantalla.
- En la pausa conceptual de IaC, considera una diapositiva con las tres consecuencias (reproducibilidad, historial, destrucción limpia) — el alumno puede leerlas mientras escucha.
- Al abrir el directorio en VS Code, navegar los tres archivos brevemente señalando con el cursor — solo el nombre y una línea de descripción oral cada uno. No leer el contenido todavía, eso es el EP19.
- El `terraform init` de prueba es el cierre técnico más importante del episodio — muestra que las credenciales del EP14 funcionan con Terraform sin configuración adicional.
- Borrar el directorio de prueba al final del video en vivo — refuerza el hábito de no dejar archivos y carpetas `.terraform/` huérfanos en el sistema.

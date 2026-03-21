# 🎬 Guión — EP14: IAM — Usuarios, Roles y Access Keys

**Duración estimada:** 12–15 min
**Tono:** Directo y con énfasis en seguridad. El alumno ya tiene cuenta AWS — ahora la configura bien antes de tocar cualquier otra cosa.

---

## 🎙️ INTRO (0:00 – 0:50)

> *Pantalla: consola de AWS abierta, logueado como root — se ve el email en la esquina superior derecha.*

"Bienvenidos al episodio 14.

Si miras la esquina superior derecha de la consola, vas a ver tu email de root. Eso es una señal de alerta.

La cuenta root de AWS es como la llave maestra de un edificio: abre absolutamente todo, incluyendo las cosas que nunca querrías abrir por accidente. Puede cancelar la cuenta, vaciar todos los recursos, generar cargos sin límite.

Por eso existe **IAM** — Identity and Access Management. IAM te permite crear usuarios con permisos específicos para el trabajo diario, de forma que el root solo se toca en emergencias reales.

Lo que vamos a hacer en este episodio: crear un usuario llamado `admin-curso`, darle permisos de administrador, y generar las Access Keys que van a usar Terraform y la AWS CLI para autenticarse. Después de este episodio, el root queda guardado.

Empecemos."

---

## 🔍 PAUSA CONCEPTUAL — ¿Cómo funciona IAM? (0:50 – 2:00)

> *Pantalla: diagrama simple o slide.*

"IAM tiene tres conceptos que vale la pena entender antes de crear cualquier cosa.

**Usuarios** — representan a una persona o a un proceso. Tienen credenciales propias: usuario y contraseña para la consola web, o Access Keys para el CLI y herramientas como Terraform.

**Grupos** — conjuntos de usuarios que comparten los mismos permisos. En lugar de asignar permisos uno por uno a cada persona, creas un grupo 'DevOps' con los permisos correctos y agregas usuarios.

**Políticas** — documentos JSON que definen qué acciones están permitidas o negadas sobre qué recursos. AWS tiene políticas predefinidas — como `AdministratorAccess` — y también puedes crear las tuyas.

Para el curso vamos a ser pragmáticos: un solo usuario `admin-curso` con `AdministratorAccess`. En un entorno real de empresa usarías permisos más granulares — solo lo que cada proceso necesita. Pero para aprender GitOps, esto es lo correcto."

---

## 📌 PASO 1 — Ir a IAM y crear el usuario (2:00 – 5:00)

> *Pantalla: consola de AWS. Buscar 'IAM' en la barra de búsqueda.*

"En la barra de búsqueda escribo **'IAM'** y abro el servicio.

En el menú lateral izquierdo, click en **Users** → **Create user**.

Nombre de usuario: `admin-curso`. Sin espacios, todo en minúsculas — es una buena convención para usuarios IAM.

Activo la opción **'Provide user access to the AWS Management Console'**. Esto le da a este usuario una contraseña para entrar a la consola web.

En 'Console password' selecciono **Custom password** e ingreso una contraseña segura. Desactivo la opción que dice 'Users must create a new password at next sign-in' — para el curso no necesitamos ese paso extra.

**Next**."

---

## 📌 PASO 2 — Asignar permisos (5:00 – 7:00)

> *Pantalla: formulario de permisos.*

"En la pantalla de permisos, selecciono **'Attach policies directly'**.

En el buscador escribo `AdministratorAccess` y marco el checkbox.

Esta política da acceso completo a todos los servicios de AWS. Terraform la va a necesitar para crear la EC2, el Security Group, el bucket S3 y la tabla DynamoDB que usaremos en los episodios siguientes.

**Next** → **Create user**.

El usuario `admin-curso` fue creado. AWS me muestra las credenciales de consola:

- **Console sign-in URL:** una URL con el Account ID de tu cuenta. Guárdala — es la URL que usarás para entrar como IAM en lugar del root.
- **User name:** `admin-curso`
- **Password:** la que configuraste.

Descargo el CSV con estas credenciales para tenerlas guardadas."

---

## 📌 PASO 3 — Generar Access Keys (7:00 – 10:00)

> *Pantalla: perfil del usuario admin-curso en IAM.*

"Ahora la parte más importante para el flujo técnico del curso: las **Access Keys**.

Las Access Keys son credenciales programáticas. No son usuario y contraseña para la consola web — son un par de claves que permiten a herramientas como la AWS CLI y Terraform autenticarse con tu cuenta desde la terminal, sin abrir ningún navegador.

Voy a: IAM → Users → `admin-curso` → pestaña **Security credentials** → sección **Access keys** → **Create access key**.

Me pregunta para qué caso de uso. Selecciono **Command Line Interface (CLI)**. Confirmo el aviso y click en **Next**.

En description tag escribo `curso-gitops-local` — esto me ayuda a recordar que estas claves son para mi máquina de desarrollo.

**Create access key**.

---

Y aquí aparecen dos valores que necesito atención total:"

> *Pantalla: se muestran las dos claves. Hacer zoom.*

"**Access Key ID** — empieza con `AKIA`. Es como el nombre de usuario.

**Secret Access Key** — una cadena larga de caracteres. Es la contraseña.

El **Secret Access Key solo se muestra una vez**. Ahora mismo, en esta pantalla. Si cierro sin copiarlo, tendrá que generar un par nuevo — no hay forma de recuperarlo.

Voy a hacer click en **Download .csv file** para tenerlas guardadas localmente. Y también las copio ahora."

---

## 📌 PASO 4 — Guardar las credenciales en la máquina local (10:00 – 12:30)

> *Pantalla: terminal en la PC local.*

"Ahora en la terminal de mi máquina local. La AWS CLI busca las credenciales en `~/.aws/credentials`. Las voy a poner ahí:"

```bash
mkdir -p ~/.aws

cat > ~/.aws/credentials << EOF
[default]
aws_access_key_id = TU_ACCESS_KEY_ID
aws_secret_access_key = TU_SECRET_ACCESS_KEY
EOF

chmod 600 ~/.aws/credentials
```

"El `chmod 600` es importante: solo tu usuario puede leer este archivo. Si otro usuario del sistema pudiera leerlo, tendría acceso a tu cuenta de AWS.

También creo el archivo de configuración regional:"

```bash
cat > ~/.aws/config << EOF
[default]
region = us-east-1
output = json
EOF
```

"Ahora verifico que todo funciona:"

```bash
aws sts get-caller-identity
```

"El resultado debe ser algo así:"

```json
{
    "UserId": "AIDAIOSFODNN7EXAMPLE",
    "Account": "123456789012",
    "Arn": "arn:aws:iam::123456789012:user/admin-curso"
}
```

"Si ves el nombre `admin-curso` en el ARN, las credenciales están funcionando correctamente."

---

## 🔍 PAUSA CONCEPTUAL — Las tres reglas de las Access Keys (12:30 – 13:30)

> *Pantalla: editor de texto o slide con las tres reglas.*

"Antes de cerrar, tres reglas sobre las Access Keys que nunca debes romper.

**Regla uno: nunca las pongas en el código.**"

```python
# ❌ NUNCA hacer esto
AWS_ACCESS_KEY = "AKIAIOSFODNN7EXAMPLE"
```

"Si esto llega a un repositorio público — aunque sea por un segundo — hay bots que monitorean GitHub en tiempo real buscando exactamente eso. En minutos, alguien puede estar usando tu cuenta para minar criptomonedas o lanzar ataques DDoS, y el cargo es tuyo.

**Regla dos: el `.gitignore` ya las protege.** Los archivos `~/.aws/credentials` y `.env` están en el `.gitignore` del proyecto desde el EP07. Nunca los agregues al repo manualmente.

**Regla tres: si sospechas que se filtraron, desactívalas inmediatamente.** IAM → Users → admin-curso → Security credentials → Make inactive. Y crea un par nuevo. No esperes para ver si pasa algo."

---

## 🎙️ CIERRE (13:30 – 14:30)

"Eso es EP14.

Ahora tienes un usuario IAM `admin-curso` con sus credenciales configuradas en la máquina local. Terraform y la AWS CLI van a usar esas credenciales automáticamente cada vez que las invoques — no hay que configurar nada más.

Cierra la sesión de root en el navegador. Desde ahora, cuando uses la consola de AWS, entra con la URL de console sign-in y el usuario `admin-curso`. El root solo cuando sea una emergencia real.

En el siguiente episodio instalamos la AWS CLI y hacemos las primeras pruebas de conectividad — listar regiones, ver el estado de EC2, los primeros comandos que vamos a usar constantemente durante el resto del curso.

Nos vemos en el EP15."

---

## 🗒️ Notas de producción

- Al mostrar las Access Keys, usa claves falsas en el video — borra las reales antes de subir la grabación o tápalas en post-producción. Dejar claves reales en un video de YouTube es uno de los errores más comunes y costosos.
- El `chmod 600` y `aws sts get-caller-identity` son los dos momentos técnicos más importantes del episodio — asegúrate de que la terminal se lea claramente.
- La pausa de las tres reglas puede hacerse más visual con una pantalla dividida mostrando el código incorrecto en rojo a la izquierda y lo correcto a la derecha.
- Al cerrar sesión del root al final del video, hacerlo en vivo y entrar como `admin-curso` — ese detalle refuerza el comportamiento correcto.

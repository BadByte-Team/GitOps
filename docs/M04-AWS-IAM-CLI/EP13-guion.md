# 🎬 Guión — EP13: Crear Cuenta AWS y Activar Free Tier

**Duración estimada:** 15–18 min
**Tono:** Directo, con énfasis especial en seguridad y costos — el alumno viene de un contexto local y esta es su primera interacción con la nube.

---

## 🎙️ INTRO (0:00 – 0:50)

> *Pantalla: calculadora de costos de AWS abierta, mostrando un número grande en rojo — algo así como $120/mes.*

"Bienvenidos al episodio 13. Hasta aquí todo lo que hemos hecho ha sido completamente local: VS Code, Git, GitHub, Docker, Compose. Cero costos. Cero nube.

Hoy entramos a AWS. Y lo primero que quiero mostrarte antes de crear cualquier cuenta es **esto**.

Esto es lo que costaría el stack clásico de un curso de GitOps: EKS para el cluster, una EC2 para Jenkins, un LoadBalancer para exponer la app. Más de cien dólares al mes. Y si olvidas destruirlo un fin de semana, ese cargo llega a tu tarjeta sin aviso.

Este curso no funciona así. En nuestro diseño, la arquitectura completa corre en **una sola instancia t2.micro — la capa gratuita de AWS — y cuesta cero pesos, cero dólares**. Jenkins, SonarQube y Trivy corren en tu máquina local. ArgoCD y tu app corren en K3s sobre esa única EC2 gratuita.

Pero para usar AWS, incluso gratis, necesitas una cuenta. Y esa cuenta necesita configurarse bien desde el día uno para que nunca recibas una sorpresa en el estado de cuenta.

Eso es lo que hacemos hoy. Vamos."

---

## 📌 PASO 1 — Crear la cuenta (0:50 – 3:30)

> *Pantalla: navegador, abriendo aws.amazon.com.*

"Voy a **aws.amazon.com** y hago click en el botón que dice **'Create an AWS Account'** — está en la esquina superior derecha.

El formulario de registro pide tres cosas básicas:

**Email root.** Este email es especial. Es el acceso más privilegiado que existe en AWS — puede cerrar la cuenta, borrar todo, generar cualquier costo. Usa un email al que tengas acceso exclusivo, que no compartas con nadie. Yo recomiendo crear uno específico para esto: algo como `aws-gitops@tudominio.com` o tu correo personal con una contraseña fuerte.

**Nombre de la cuenta.** Ponle algo que te identifique: `tu-nombre-gitops`, o simplemente `curso-gitops`. Este nombre aparece en la consola y en los reportes de billing.

**Contraseña.** Que sea larga y única. No la misma de GitHub ni del correo."

> *Pantalla: formulario llenado, click en continuar.*

"Verifico el email — AWS manda un código de 6 dígitos. Lo ingreso."

---

## 📌 PASO 2 — Información de contacto y método de pago (3:30 – 5:30)

> *Pantalla: formulario de información de contacto.*

"En el tipo de cuenta selecciono **Personal**. Lleno nombre, teléfono y dirección.

Ahora la parte que le genera ansiedad a todo el mundo: **el método de pago**.

AWS requiere una tarjeta de débito o crédito para verificar que eres una persona real. Va a realizar un **cargo de verificación de $1 USD** — ese dólar se devuelve en 3 a 5 días hábiles. No es un cobro real.

Lo importante: si usas Free Tier correctamente — y en este curso te mostraré exactamente cómo hacerlo — **no recibirás ningún cargo**. El único recurso que usaremos que podría generar costo es la EC2 t2.micro, y esa tiene 750 horas mensuales gratis durante el primer año.

Ingreso los datos de la tarjeta y continúo.

---

Una nota si tu banco bloquea el cargo de $1: es una transacción internacional con un comercio llamado 'Amazon Web Services'. Llama a tu banco y pídeles que la autoricen. Es un proceso de 5 minutos."

---

## 📌 PASO 3 — Verificación telefónica y plan de soporte (5:30 – 7:00)

> *Pantalla: formulario de verificación.*

"Verificación por teléfono: AWS llama o manda SMS con un PIN de 4 dígitos. Lo ingreso.

Plan de soporte: selecciono **Basic Support — Free**. Los planes de pago son para empresas con SLAs de producción. Para un curso de aprendizaje, Basic es completamente suficiente: tienes acceso a la documentación, los foros de la comunidad, y el dashboard básico de Trusted Advisor.

Click en 'Complete sign up'. La cuenta puede tardar unos minutos en activarse."

---

## 📌 PASO 4 — Primer acceso a la consola (7:00 – 8:30)

> *Pantalla: console.aws.amazon.com, pantalla de login.*

"Voy a **console.aws.amazon.com** y me logueo con el email root y la contraseña.

Lo primero que hago es seleccionar la región en el menú desplegable de la esquina superior derecha. La región que usaremos en todo el curso es **US East (N. Virginia) — us-east-1**.

¿Por qué us-east-1? Tres razones: es la región más antigua y madura de AWS, tiene la mayor cantidad de servicios disponibles, y tiene los precios más bajos. Todos los recursos del curso — la EC2, el bucket S3, la tabla DynamoDB — los crearemos aquí.

Manten **siempre** esta región seleccionada. Si en algún episodio creas algo y no lo encuentras, lo más probable es que estés mirando la región equivocada."

---

## 📌 PASO 5 — Activar MFA en la cuenta root (8:30 – 11:30)

> *Pantalla: consola de AWS, buscando IAM.*

"Ahora la parte más importante del episodio, y la que más gente saltea: **activar MFA en el root**.

MFA significa Multi-Factor Authentication — autenticación de dos factores. Además de la contraseña, necesitas un código temporal del teléfono para entrar. Si alguien roba tu contraseña, sin el teléfono no puede hacer nada.

La cuenta root de AWS puede hacer **todo**: generar costos masivos, vaciar bases de datos, cerrar la cuenta. Una contraseña filtrada sin MFA es el escenario de pesadilla. No lo dejes para después.

Para activarlo: click en tu nombre de cuenta en la esquina superior derecha → **Security credentials** → sección **Multi-factor authentication** → **Assign MFA device**.

Necesito la app en el teléfono. Si no la tienes, tómate 2 minutos ahora para instalarla:

- **Google Authenticator** — simple y directo.
- **Authy** — más funcional, permite hacer backup del MFA.

Cualquiera de las dos funciona.

---

En AWS: selecciono **Authenticator app** → **Next**. Aparece un código QR.

Abro la app en el teléfono → botón **+** → **Escanear código QR**. Apunto la cámara al QR.

La app genera un código de 6 dígitos que cambia cada 30 segundos. Tengo que ingresar **dos códigos consecutivos** — uno, espero a que cambie, el segundo. Eso confirma que el tiempo del teléfono está sincronizado.

**'Add MFA'**. Listo.

---

Para verificar que funciona: cierro sesión y vuelvo a entrar. Esta vez después de la contraseña, AWS pide el código del autenticador. Lo ingreso. Entro.

A partir de ahora, cada vez que entre como root necesito el teléfono. Eso es exactamente lo que queremos."

---

## 📌 PASO 6 — Configurar la alerta de billing (11:30 – 14:30)

> *Pantalla: consola de AWS, buscando 'Budgets' en la barra de búsqueda.*

"La segunda línea de defensa: una **alerta de billing**. Si por cualquier motivo — un recurso olvidado, un error de configuración — los costos suben, quiero saberlo antes de que llegue el cobro.

En la barra de búsqueda de la consola escribo **'Budgets'** y lo abro.

**'Create budget'** → tipo: **Cost budget** → **Next**.

Configuro:
- **Budget name:** `curso-gitops-alert`
- **Period:** Monthly
- **Budgeted amount:** `$5.00`

¿Por qué $5? Porque en este curso el costo real debería ser $0. Si llega a $5, algo está mal y quiero saberlo. Es un margen de seguridad.

**Next** → **Add an alert threshold**:
- Threshold: **80% of budgeted amount** — eso son $4 USD.
- Email recipients: mi correo.

**Next** → **Create budget**.

---

Una cosa más antes de cerrar. Por defecto, los usuarios IAM — que crearemos en el próximo episodio — no pueden ver la información de facturación. Hay que habilitarlo manualmente.

Click en tu nombre → **Account** → busco la sección **'IAM user and role access to Billing information'** → **Edit** → activo la casilla → **Update**.

Hecho. Ahora el usuario admin que crearemos en el EP14 podrá ver los costos."

---

## 🎙️ CIERRE (14:30 – 15:30)

"Eso es EP13. Cuenta creada, región configurada, MFA activado, alerta de billing en su lugar.

Antes de cerrar, el resumen de lo que cuesta nuestro stack en este curso:

| Recurso | Free Tier | Costo real |
|---|---|---|
| EC2 t2.micro | 750 hrs/mes gratis | $0 |
| S3 backend | 5 GB gratis | $0 |
| DynamoDB | Siempre gratis | $0 |
| Jenkins local | En tu PC | $0 |
| K3s | En la EC2 | $0 |

Cero dólares. Eso es posible porque tomamos decisiones de arquitectura inteligentes: K3s en lugar de EKS, Jenkins local en lugar de una segunda EC2, NodePort en lugar de LoadBalancer. Toda la teoría y práctica de GitOps — exactamente la misma — a costo cero.

En el siguiente episodio vamos a crear el usuario IAM que usaremos para el trabajo diario, y vamos a generar las Access Keys que necesita la AWS CLI. Nunca más usamos el root después de hoy.

Nos vemos en el EP14."

---

## 🗒️ Notas de producción

- La apertura con la calculadora de costos de AWS es el gancho principal — tenerla preparada de antemano con los números reales de EKS + EC2 + ALB.
- Al activar MFA, mostrar el teléfono físico en cámara si es posible, o al menos mostrar la pantalla de la app del autenticador. Es el momento más concreto y visual del episodio.
- Cuando muestras el formulario de la tarjeta, considera tapar los números reales — usa una tarjeta de prueba si la tienes, o cubre el número con una capa en post-producción.
- La tabla de costos al final es candidata a quedarce en pantalla durante el cierre como lower-third o slide.
- Usar la barra de búsqueda de la consola de AWS en lugar de navegar por el menú — es más rápido y enseña un hábito útil.

# EP 28: Configurar Memoria Swap en EC2

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 8–10 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** Episodio completamente nuevo. Enseña a crear memoria virtual (Swap) en Linux — paso vital para que K3s no colapse la instancia t3.micro de 1 GB de RAM.

---

## 🎯 Objetivo
Crear un archivo de Swap de 2 GB en la EC2 y hacerlo persistente entre reinicios, para que K3s y ArgoCD puedan coexistir en una instancia con solo 1 GB de RAM física.

---

## 📋 Prerequisitos
- EC2 t3.micro corriendo (EP22)
- IP pública de la instancia anotada
- Archivo `aws-key.pem` con permisos 400

---

## 🧠 ¿Qué es el Swap y por qué es necesario?

El Swap es espacio en disco que el sistema operativo usa como RAM cuando la memoria física se agota:

```
Sin Swap:                     Con 2 GB Swap:
┌─────────┐                   ┌─────────┐  ┌──────────────┐
│  1 GB   │                   │  1 GB   │  │  2 GB Swap   │
│  RAM    │                   │  RAM    │  │  (en disco)  │
│ [LLENA] │                   │  [uso]  │  │  [desborde]  │
│  ❌ OOM │                   │  ✅ OK  │  │  ✅ OK       │
└─────────┘                   └─────────┘  └──────────────┘
  K3s colapsa                         K3s sobrevive
```

Estimación real de uso de RAM con el stack completo:

| Componente | RAM estimada |
|---|---|
| Sistema operativo Ubuntu | ~150 MB |
| K3s (proceso principal) | ~200 MB |
| ArgoCD (~7 pods) | ~400 MB |
| MySQL | ~200 MB |
| App Go | ~30 MB |
| **Total** | **~980 MB — prácticamente el 100% de la RAM** |

Sin Swap, el kernel activará el **OOM Killer** y terminará procesos sin aviso. Con 2 GB de Swap, el sistema tiene margen suficiente para operar establemente.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: `free -h` ejecutado dentro de la EC2 — mostrando claramente `Swap: 0B`.*

"Bienvenidos al episodio 28.

Estoy conectado al servidor que creamos en el EP22. Y lo primero que hago cada vez que abro el servidor es revisar la memoria disponible:"

```bash
free -h
```

"Miren la línea de Swap: cero. Ni un solo byte de memoria virtual. Solo el 1 GB de RAM física de la instancia t3.micro.

Ahora piensen en lo que vamos a correr en este mismo servidor cuando terminemos el módulo. K3s necesita sus procesos del control plane. ArgoCD tiene aproximadamente siete pods propios. MySQL necesita su memoria para las conexiones y el cache. La app Go. El sistema operativo.

Hice la estimación: la suma total ronda los 980 MB. Prácticamente el 100% de la RAM disponible, sin margen alguno.

¿Qué pasa en Linux cuando el sistema se queda sin memoria? El kernel activa el **OOM Killer** — Out of Memory Killer. Su trabajo es encontrar el proceso que más memoria está consumiendo y terminarlo. Sin aviso, sin mensaje claro en el log. El proceso simplemente desaparece. En nuestro caso, el candidato más probable sería K3s o uno de los pods de ArgoCD.

El Swap resuelve esto. Le reservamos 2 GB del disco SSD de la instancia para usarlos como memoria de desborde. Cuando la RAM se llena, el kernel mueve páginas de memoria poco usadas al disco, liberando espacio para los procesos activos.

¿Es tan rápido como la RAM? No — el disco siempre es más lento. Pero en una instancia con almacenamiento gp3, la diferencia es completamente aceptable para un entorno de aprendizaje. Y la estabilidad que nos da vale completamente la pena.

Este episodio es corto. Cuatro comandos para crear el Swap, uno para hacerlo persistente. Empecemos."

---

### PASO 1 — Verificar espacio disponible (1:30 – 2:30)

> *Pantalla: terminal dentro de la EC2.*

"Antes de crear el archivo de Swap, verifico que hay espacio suficiente en el disco. Necesito al menos 3 GB libres para el archivo de 2 GB más un margen de seguridad:"

```bash
df -h /
```

```
Filesystem      Size  Used Avail Use% Mounted on
/dev/xvda1       29G  1.8G   28G   7%  /
```

"28 GB disponibles — más que suficiente. Ese disco de 30 GB que configuramos en el Terraform del EP22 nos da holgura de sobra.

Y confirmo el estado de la memoria una vez más para tener el punto de partida claro:"

```bash
free -h
```

"Swap en cero. Ese es el antes. Ahora vamos al después."

---

### PASO 2 — Crear el archivo de Swap (2:30 – 5:30)

> *Pantalla: terminal dentro de la EC2. Un comando por bloque, con explicación entre cada uno.*

"Cuatro comandos en secuencia. Los ejecuto uno por uno.

**Primero: crear el archivo de 2 GB.**

`fallocate` es una herramienta de Linux que reserva espacio en disco de forma instantánea — no escribe ceros, simplemente marca ese espacio como reservado. El resultado es inmediato:"

```bash
sudo fallocate -l 2G /swapfile
```

"Verifico que existe:"

```bash
ls -lh /swapfile
# -rw-r--r-- 1 root root 2.0G ... /swapfile
```

"Dos gigabytes. Pero fíjense en los permisos: `-rw-r--r--`. Cualquier usuario puede leer este archivo. Eso es un problema de seguridad — el Swap contiene páginas de la memoria de todos los procesos, incluyendo datos sensibles como passwords o tokens JWT. Solo root debe poder acceder.

---

**Segundo: restringir los permisos.**"

```bash
sudo chmod 600 /swapfile
```

"Ahora `-rw-------`. Solo root puede leer y escribir. Si intentara activar el Swap sin este paso, `swapon` lo rechazaría con un error de permisos inseguros.

---

**Tercero: formatear el archivo como espacio de Swap.**

`mkswap` escribe la cabecera que el kernel de Linux necesita para reconocer el archivo como Swap válido:"

```bash
sudo mkswap /swapfile
```

```
Setting up swapspace version 1, size = 2 GiB (2147479552 bytes)
no label, UUID=abc123...
```

"El UUID que aparece es como el identificador único del espacio de Swap. No hace falta copiarlo ni guardarlo.

---

**Cuarto: activar el Swap.**

`swapon` le dice al kernel que empiece a usar este archivo como memoria virtual:"

```bash
sudo swapon /swapfile
```

"Sin output — eso es correcto. Cuando un comando de Linux no dice nada, generalmente significa que funcionó. Verifico:"

```bash
free -h
```

```
               total        used        free
Mem:           981Mi       165Mi       634Mi
Swap:          2.0Gi         0B       2.0Gi
```

"Ese es el contraste que buscábamos. La línea de Swap ahora muestra `2.0Gi` de total. El sistema tiene memoria virtual disponible.

`USED: 0B` porque todavía no hay presión de memoria — K3s no está instalado. Pero cuando llegue el momento y los pods empiecen a consumir RAM, el sistema tiene el colchón listo."

---

### PASO 3 — Hacer el Swap persistente (5:30 – 7:30)

> *Pantalla: terminal dentro de la EC2.*

"Ahora el paso que más gente olvida — y que hace que todo el trabajo anterior sea inútil si se omite.

El Swap que acabamos de activar está funcionando ahora mismo. Pero si la EC2 se reinicia — por un mantenimiento de AWS, por un apagado accidental, por lo que sea — el Swap desaparece. Al arrancar, Linux no busca archivos de Swap automáticamente a menos que se lo indiques explícitamente.

La forma de hacerlo es a través de `/etc/fstab`. Este archivo controla qué filesystems y espacios de intercambio monta el sistema automáticamente al arrancar. Si agregas una línea aquí, el Swap se activa en cada inicio del servidor.

Agrego la línea:"

```bash
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

"Vale la pena entender el formato de esta línea porque `/etc/fstab` aparece en cualquier trabajo de administración de Linux:

- `/swapfile` — la ruta del archivo o dispositivo
- `none` — el punto de montaje. Para el Swap no aplica, se pone `none`
- `swap` — el tipo: le dice a Linux que este es un espacio de intercambio
- `sw` — opciones de montaje para Swap
- `0 0` — los campos de dump y fsck. Para Swap siempre son cero

Verifico que se escribió:"

```bash
cat /etc/fstab | grep swap
# /swapfile none swap sw 0 0
```

"Ahí está la línea. Ahora el Swap va a sobrevivir cualquier reinicio de la instancia."

---

### PASO 4 — Ajustar swappiness (7:30 – 8:30)

> *Pantalla: terminal dentro de la EC2.*

"Un ajuste más, opcional pero recomendado para un servidor.

El parámetro `vm.swappiness` controla cuán agresivamente el kernel usa el Swap, en una escala del 0 al 100. El valor por defecto en Ubuntu es 60. Eso significa que el kernel empieza a mover páginas al Swap cuando la RAM todavía tiene un 40% disponible — demasiado agresivo para un servidor.

Con un valor de 10, el kernel solo usa el Swap cuando la RAM está prácticamente llena. La RAM siempre es más rápida, así que quiero que el sistema la use al máximo antes de recurrir al disco:"

```bash
# Ver el valor actual
cat /proc/sys/vm/swappiness
# 60

# Cambiar a 10
sudo sysctl vm.swappiness=10

# Hacerlo persistente entre reinicios
echo 'vm.swappiness=10' | sudo tee -a /etc/sysctl.conf
```

---

### PASO 5 — Verificación final y salida (8:30 – 9:30)

> *Pantalla: terminal dentro de la EC2.*

"Antes de salir del servidor, la verificación completa:"

```bash
# Memoria con Swap activo
free -h

# Confirmar que el archivo de Swap está listado
sudo swapon --show

# Confirmar la línea en fstab
grep swap /etc/fstab

# Confirmar swappiness
cat /proc/sys/vm/swappiness
```

"Todo en orden. Salgo:"

```bash
exit
```

---

### CIERRE (9:30 – 10:00)

"Eso es el episodio 28.

La EC2 ahora tiene 1 GB de RAM física más 2 GB de memoria virtual. El OOM Killer no va a matar procesos cuando K3s, ArgoCD y MySQL estén corriendo al mismo tiempo. Esa estabilidad es la que necesitamos para que el resto del módulo funcione sin sorpresas.

En el siguiente episodio instalamos K3s directamente en este servidor. Van a ver el proceso de instalación de Kubernetes más simple que existe — un solo comando `curl` que descarga, instala y configura todo automáticamente.

Nos vemos en el EP29."

---

## ✅ Checklist de Verificación
- [ ] `free -h` dentro de la EC2 muestra `Swap: 2.0Gi`
- [ ] `sudo swapon --show` lista `/swapfile` como activo
- [ ] `grep swap /etc/fstab` muestra la línea de persistencia
- [ ] `cat /proc/sys/vm/swappiness` muestra `10`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `fallocate: Operation not supported` | `sudo dd if=/dev/zero of=/swapfile bs=1M count=2048` como alternativa |
| `swapon: insecure permissions 0644` | `sudo chmod 600 /swapfile` — los permisos deben ser exactamente 600 |
| El Swap desaparece al reiniciar | La línea no está en `/etc/fstab` — verificar con `grep swap /etc/fstab` |
| `No space left on device` | Verificar espacio con `df -h /` — necesitas al menos 3 GB libres |

---

## 🗒️ Notas de Producción
- Abrir el episodio ya conectado a la EC2 con `free -h` en pantalla — el `Swap: 0B` es el gancho visual del episodio.
- Ejecutar cada uno de los cuatro comandos por separado con pausa y explicación — no pegarlos todos de golpe. Cada uno tiene su razón de ser.
- El contraste del `free -h` antes y después del `swapon` es el momento más satisfactorio del episodio — detener el video brevemente para que el alumno lo procese.
- Enfatizar con la voz que el paso del `fstab` es el que más gente olvida y el que hace que todo se pierda al reiniciar.
- Los campos de `/etc/fstab` merecen una explicación breve — ese archivo aparece constantemente en administración de sistemas Linux.

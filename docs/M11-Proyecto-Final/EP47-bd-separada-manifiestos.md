# EP 47: Base de Datos Separada y Manifiestos de Kubernetes

**Tipo:** PRÁCTICA
**Duración estimada:** 18–20 min
**Dificultad:** ⭐⭐⭐ (Avanzado)
**🔄 MODIFICADO:** Explica por qué MySQL no va en el Dockerfile de Go, y recorre en detalle los siete manifiestos que ArgoCD despliega en K3s.

---

## 🎯 Objetivo

Entender por qué la base de datos debe ser un recurso independiente en Kubernetes, y revisar en detalle cada uno de los siete manifiestos del directorio `gitops-infra/infrastructure/kubernetes/app/` — desde el Namespace hasta el NodePort de la app.

---

## 📋 Prerequisitos

- K3s con ArgoCD configurado y sincronizando `gitops-infra` (EP40)
- Los siete archivos YAML presentes en `gitops-infra/infrastructure/kubernetes/app/`

---

## 🧠 La arquitectura de datos en Kubernetes

```
Pod App Go    ──DNS──▶   mysql-svc (ClusterIP)   ──▶   Pod MySQL   ──▶   /var/lib/mysql
(stateless)               (nombre estable)               (stateful)         (datos en disco)
DB_HOST=mysql-svc

                                                    ▲
                                              mysql-configmap
                                              (init.sql al arrancar)

                                              db-credentials Secret
                                              (usuario y contraseña)
```

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:30)

> *Pantalla: el `docker-compose.yml` del proyecto abierto — los servicios `api` y `mysql-db` visibles.*

"Bienvenidos al episodio 47.

Este es el `docker-compose.yml` del proyecto. Dos servicios: la app Go y MySQL. Funcionan juntos porque Docker Compose los conecta en una red compartida y gestiona el ciclo de vida de ambos.

Cuando pasamos al mundo de Kubernetes, la misma lógica aplica — dos servicios, conectados en una red — pero la forma de definirlos cambia completamente. En lugar de un archivo con dos secciones, tenemos siete archivos YAML separados.

¿Por qué siete archivos para dos servicios? Porque en Kubernetes cada responsabilidad tiene su propio objeto: el Pod es una cosa, el Service que lo expone es otra, las credenciales son otra, la configuración inicial de la base de datos es otra. Esa separación parece más compleja al principio, pero hace el sistema mucho más mantenible y seguro.

Hoy recorremos los siete archivos uno por uno. Al terminar, cada línea de cada manifiesto va a tener sentido.

Empecemos."

---

### 🔍 PAUSA CONCEPTUAL — Por qué MySQL NO va en el Dockerfile (1:30 – 4:00)

> *Pantalla: un Dockerfile imaginario con MySQL instalado dentro — marcado con ❌ en rojo.*

"Antes de los manifiestos, el principio más importante de este episodio.

Un error muy común cuando alguien aprende Docker es querer meter todo en una sola imagen. La app y la base de datos juntas en el mismo contenedor. El razonamiento suena lógico: 'si ambos necesitan correrse juntos, ¿por qué no en una sola imagen?'

La respuesta tiene cuatro partes.

**Acoplamiento de ciclos de vida.** Cuando actualizas la app — nuevo código, nueva imagen — tienes que reiniciar el contenedor completo. Y al reiniciar, MySQL también se reinicia. Los datos de las sesiones activas se pierden. Los usuarios experimentan una interrupción. En Kubernetes, el rolling update existe precisamente para evitar eso — reemplaza el pod de la app sin tocar el pod de MySQL.

**Escalabilidad.** Si la app tiene picos de tráfico y quieres escalar a tres réplicas, cada réplica tendría su propio MySQL con sus propios datos. Las tres bases de datos estarían desincronizadas. La app Go no sabría qué datos hay en cuál. Es un caos que no tiene solución simple.

**Persistencia.** Los contenedores son efímeros por diseño. Cuando un pod muere y se recrea, pierde todo lo que tenía en su sistema de archivos. Si MySQL estuviera en el mismo pod que la app, cada redeploy borraría todos los datos. En Kubernetes, la persistencia se gestiona con volúmenes que sobreviven al pod — eso solo funciona si MySQL tiene su propio pod.

**El principio de responsabilidad única.** Una imagen, un proceso, una responsabilidad. La imagen de la app Go sabe servir HTTP. La imagen de MySQL sabe gestionar una base de datos. Cada una puede actualizarse de forma independiente sin afectar a la otra.

La arquitectura correcta es esta: un pod para la app, un pod para MySQL, un Service que actúa como dirección estable entre los dos. La app Go no sabe ni le importa la IP del pod de MySQL — solo conoce el nombre `mysql-svc`, y ese nombre siempre resuelve a MySQL sin importar cuántas veces se haya recreado el pod."

---

### Recorrido por los siete manifiestos (4:00 – 16:00)

> *Pantalla: VS Code con el directorio `gitops-infra/infrastructure/kubernetes/app/` abierto.*

"Abro el directorio. Siete archivos YAML. Los voy a recorrer en el orden en que Kubernetes los necesita aplicar."

---

#### Archivo 1 — `namespace.yaml` (4:00 – 5:00)

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: curso-gitops
```

"Simple. Crea el namespace `curso-gitops` que agrupa todos los recursos del proyecto. Sin el namespace, los demás archivos no saben en qué espacio de nombres crear sus recursos.

ArgoCD crea el namespace automáticamente gracias a `CreateNamespace=true` en la `syncPolicy` de la Application."

---

#### Archivo 2 — `secrets.yaml` (5:00 – 7:00)

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
  namespace: curso-gitops
type: Opaque
data:
  username: Y3Vyc29fYXBw          # curso_app
  password: QzRyczBfUzNjdXIzX1BAc3Mh  # C4rs0_S3cur3_P@ss!
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
  namespace: curso-gitops
type: Opaque
data:
  jwt-secret: Z2s4c19wcjBkX3MzY3IzdF9jaDRuZzNfbTMh
```

"Los Secrets guardan información sensible en Base64 — no es cifrado, es solo codificación. La ventaja real de los Secrets sobre los ConfigMaps es el control de acceso: Kubernetes puede restringir qué pods pueden leer qué Secrets.

Para generar los valores en Base64:"

```bash
echo -n "curso_app" | base64
# Y3Vyc29fYXBw

echo -n "C4rs0_S3cur3_P@ss!" | base64
# QzRyczBfUzNjdXIzX1BAc3Mh
```

"El flag `-n` es importante — sin él, `echo` agrega un salto de línea al final que se incluye en el Base64 y la contraseña resultante tiene un carácter extra invisible."

---

> ⚠️ **ADVERTENCIA DE SEGURIDAD — Base64 NO es cifrado**
>
> Los valores en `secrets.yaml` están codificados en Base64, no cifrados. Cualquier persona con acceso al repositorio `gitops-infra` puede decodificarlos con `echo "Y3Vyc29fYXBw" | base64 -d` y obtener las credenciales en texto plano.
>
> **En producción, nunca se commitean Secrets en Base64 a Git.** Alternativas reales:
>
> - **Sealed Secrets (Bitnami):** Cifra los Secrets con una llave pública. Solo el cluster puede descifrarlos. El archivo cifrado sí se puede commitear a Git de forma segura.
> - **External Secrets Operator:** Los Secrets viven en un gestor externo (AWS Secrets Manager, HashiCorp Vault) y Kubernetes los lee dinámicamente.
> - **SOPS (Mozilla):** Cifra los valores del YAML con llaves GPG o KMS antes de commitear.
>
> Para el curso usamos Base64 directo por simplicidad, pero es importante entender que esta práctica es inaceptable en cualquier entorno de producción.

---

#### Archivo 3 — `mysql-configmap.yaml` (7:00 – 8:30)

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-init-config
  namespace: curso-gitops
data:
  init.sql: |
    CREATE DATABASE IF NOT EXISTS curso_db;
    USE curso_db;
    CREATE TABLE users (
      id INT AUTO_INCREMENT PRIMARY KEY,
      username VARCHAR(50) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      role ENUM('admin', 'student') NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE modules (...);
    CREATE TABLE episodes (...);
```

"El ConfigMap guarda configuración no sensible. En este caso, el script SQL que inicializa la base de datos la primera vez que arranca el pod de MySQL.

MySQL tiene un mecanismo especial: cualquier archivo `.sql` que esté en `/docker-entrypoint-initdb.d/` se ejecuta automáticamente al arrancar por primera vez. El Deployment de MySQL monta este ConfigMap en esa ruta."

---

#### Archivo 4 — `mysql-deployment.yaml` (8:30 – 11:00)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
  namespace: curso-gitops
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "300m"
        env:
        - name: MYSQL_ROOT_PASSWORD
          value: "r00t_S3cur3_P@ss!"
        - name: MYSQL_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: MYSQL_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        volumeMounts:
        - name: init-script
          mountPath: /docker-entrypoint-initdb.d
      volumes:
      - name: init-script
        configMap:
          name: mysql-init-config
```

"Tres cosas importantes aquí.

**Los `resources.limits`** — críticos en una t3.micro. Sin límites, MySQL podría consumir toda la RAM disponible y dejar sin memoria al pod de la app o a los componentes del sistema. Con `limits.memory: 256Mi`, Kubernetes termina el pod si intenta usar más — es el comportamiento preferible a un colapso total del nodo.

**El `secretKeyRef`** — en lugar de poner el usuario y la contraseña directamente en el Deployment, los lee del Secret `db-credentials`. Si alguien lee el Deployment, no ve las credenciales. Solo ve una referencia al Secret.

**El volumen del ConfigMap** — monta el `init.sql` del ConfigMap en `/docker-entrypoint-initdb.d/`. MySQL ejecutará ese script al arrancar la primera vez, creando las tablas que necesita la app."

> ⚠️ **ADVERTENCIA DE SEGURIDAD — Contraseña de root hardcodeada**
>
> A diferencia de `MYSQL_USER` y `MYSQL_PASSWORD` (que usan `secretKeyRef`), la contraseña de root está en texto plano directamente en el Deployment: `value: "r00t_S3cur3_P@ss!"`. Cualquier persona que lea este archivo tiene acceso root al MySQL.
>
> En producción, `MYSQL_ROOT_PASSWORD` también debe venir de un Secret:
>
> ```yaml
> - name: MYSQL_ROOT_PASSWORD
>   valueFrom:
>     secretKeyRef:
>       name: db-credentials
>       key: root-password
> ```

---

#### Archivo 5 — `mysql-service.yaml` (11:00 – 12:00)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-svc
  namespace: curso-gitops
spec:
  ports:
  - port: 3306
  selector:
    app: mysql
```

"El Service más simple del proyecto. Tipo `ClusterIP` por defecto — solo accesible desde dentro del cluster. El nombre `mysql-svc` es el DNS interno que la app Go usa para conectarse: `DB_HOST=mysql-svc`. Ese nombre siempre resuelve al pod de MySQL, sin importar cuántas veces se haya recreado."

---

#### Archivo 6 — `deployment.yaml` (12:00 – 14:30)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: curso-gitops
  namespace: curso-gitops
spec:
  replicas: 1
  selector:
    matchLabels:
      app: curso-gitops
  template:
    spec:
      containers:
      - name: curso-gitops
        image: TU_USUARIO/curso-gitops:latest  # Jenkins actualiza esta línea
        resources:
          limits:
            memory: "128Mi"
            cpu: "200m"
        env:
        - name: DB_HOST
          value: "mysql-svc"
        - name: DB_USER
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: password
        - name: DB_NAME
          value: "curso_db"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
```

"Este es el archivo que Jenkins actualiza con `sed` en cada pipeline. La línea `image:` cambia de tag en cada build.

`DB_HOST=mysql-svc` — la app Go se conecta a MySQL usando el nombre del Service, no una IP. Ese nombre es estable a través de reinicios y recreaciones del pod.

Las credenciales vienen de los Secrets — `db-credentials` y `app-secrets`. La app nunca tiene las contraseñas hardcodeadas."

---

#### Archivo 7 — `service.yaml` (14:30 – 16:00)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: curso-gitops-svc
  namespace: curso-gitops
spec:
  type: NodePort
  selector:
    app: curso-gitops
  ports:
  - name: http
    protocol: TCP
    port: 80
    targetPort: 8080
    nodePort: 30081
```

"El único Service de tipo `NodePort` del stack. El tráfico que llega a `<IP_EC2>:30081` lo redirige al puerto 8080 del contenedor de la app Go. El puerto 30081 fue abierto en el Security Group de Terraform en el EP22."

---

### La relación entre los siete archivos (16:00 – 17:30)

> *Pantalla: diagrama de dependencias.*

"Para cerrar, la relación entre los siete archivos:"

```
namespace.yaml
    │ todos los recursos viven aquí
    ▼
secrets.yaml ──────────────────────────────────┐
    │ referenciado por                          │
    ▼                                           ▼
mysql-deployment.yaml          deployment.yaml (app Go)
    │ usa:                          │ usa:
    ├── secretKeyRef → db-credentials   ├── secretKeyRef → db-credentials
    └── configMap → mysql-init-config   ├── secretKeyRef → app-secrets
                                        └── DB_HOST=mysql-svc ──▶ mysql-svc

mysql-configmap.yaml ──▶ mysql-deployment.yaml
mysql-service.yaml   ──▶ accesible como mysql-svc en la red del cluster
service.yaml (NodePort 30081) ──▶ expone la app al exterior
```

---

### CIERRE (17:30 – 18:30)

"Eso es el EP47.

Los siete manifiestos recorridos en detalle. Por qué existen, qué crean, cómo se referencian entre sí. El Namespace que agrupa todo, los Secrets que guardan credenciales, el ConfigMap con el SQL de inicialización, los Deployments que gestionan los pods, los Services que conectan todo.

En el siguiente episodio vemos el pipeline completo en acción una vez más — con todos los stages de seguridad activos, con ArgoCD sincronizando, con el rolling update en tiempo real. El run final del curso.

Nos vemos en el EP48."

---

## ✅ Checklist de Verificación

- [ ] Entiendes por qué MySQL no puede estar en el mismo Dockerfile que la app Go
- [ ] Sabes por qué se usa `secretKeyRef` en lugar de poner credenciales directamente
- [ ] Entiendes la diferencia entre `ClusterIP` (mysql-svc) y `NodePort` (curso-gitops-svc)
- [ ] Puedes explicar el rol del ConfigMap con el `init.sql`
- [ ] Los siete archivos existen en `gitops-infra/infrastructure/kubernetes/app/`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| Pod MySQL en `CrashLoopBackOff` | `kubectl logs mysql-POD -n curso-gitops` — probablemente el Secret tiene valores Base64 incorrectos |
| App en `CrashLoopBackOff` con error de DB | `kubectl logs curso-gitops-POD -n curso-gitops` — verificar que `DB_HOST=mysql-svc` resuelve |
| `init.sql` no se ejecutó | El ConfigMap no está montado correctamente — verificar `volumeMounts` en mysql-deployment.yaml |
| `base64 -d` devuelve caracteres raros | El valor fue codificado con `echo` sin `-n` — regenerar con `echo -n "valor" \| base64` |

---

## 🗒️ Notas de Producción

- Abrir con el `docker-compose.yml` y conectarlo explícitamente con lo que vienen los manifiestos — el alumno conoce Compose desde el EP12.
- La pausa conceptual de "por qué no juntar MySQL y la app" es el núcleo pedagógico del episodio — tomarse el tiempo necesario para cada uno de los cuatro puntos.
- Al mostrar el `secretKeyRef`, hacer zoom en las líneas relevantes y comparar visualmente con poner el valor directo — el contraste muestra la ventaja de seguridad.
- El diagrama de relaciones entre los siete archivos puede presentarse como slide al final — da una vista de pájaro antes de cerrar.
- Verificar en vivo que ArgoCD muestra todos los recursos como `Synced` después del recorrido — confirma que los manifiestos son válidos.

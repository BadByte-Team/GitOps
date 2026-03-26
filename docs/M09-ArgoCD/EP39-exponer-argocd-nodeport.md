# EP 39: Exponer ArgoCD usando NodePort

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 8–10 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** En lugar de un LoadBalancer de AWS (~$20/mes), ArgoCD se expone de forma gratuita usando NodePort en el puerto 30080.

---

## 🎯 Objetivo

Cambiar el Service de ArgoCD de `ClusterIP` a `NodePort` en el puerto 30080, obtener y cambiar la contraseña inicial, y acceder al dashboard de ArgoCD desde el navegador.

---

## 📋 Prerequisitos

- ArgoCD instalado y todos los pods en `Running` (EP38)
- Puerto 30080 abierto en el Security Group de la EC2 (configurado en el EP22 con Terraform)

---

## 🧠 LoadBalancer vs NodePort

| Tipo | Cómo funciona | Costo |
|---|---|---|
| `LoadBalancer` | Crea automáticamente un AWS ALB/ELB | ~$20 USD/mes |
| `NodePort` | Expone el servicio en un puerto del nodo directamente | **$0** |

El puerto 30080 ya está abierto en el Security Group desde el EP22. El Terraform que usamos abrió ese puerto específicamente anticipando este momento. No hay que tocar nada de infraestructura.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: `kubectl get svc argocd-server -n argocd` mostrando el tipo `ClusterIP`.*

"Bienvenidos al episodio 39.

ArgoCD está corriendo en el cluster, pero su interfaz web no es accesible desde fuera. El Service `argocd-server` es de tipo `ClusterIP` — eso significa que solo puede recibir tráfico desde dentro del cluster.

Para acceder desde el navegador necesitamos exponerlo hacia afuera. En EKS la forma estándar sería cambiar el tipo a `LoadBalancer`, lo que crearía automáticamente un Application Load Balancer en AWS. Costo: alrededor de $20 al mes, solo por tener ese balanceador existiendo.

Nuestra solución es más simple y más barata: `NodePort`. En lugar de un balanceador de carga externo, el tráfico entra directamente por un puerto de la instancia EC2. Y ese puerto — el 30080 — ya lo abrimos en el Terraform del EP22.

Tres comandos y ArgoCD está accesible. Vamos."

---

### PASO 1 — Cambiar el Service a NodePort (1:00 – 3:30)

> *Pantalla: terminal local.*

"El comando `kubectl patch` modifica un recurso existente en el cluster sin tener que editar y volver a aplicar un YAML completo:"

```bash
kubectl patch svc argocd-server -n argocd \
  -p '{"spec": {"type": "NodePort", "ports": [{"port": 443, "targetPort": 8080, "nodePort": 30080}]}}'
```

"Desglose del comando:

- `patch svc argocd-server -n argocd` — modifica el Service llamado `argocd-server` en el namespace `argocd`
- `-p '...'` — el parche en formato JSON. Cambia el tipo a `NodePort` y mapea el puerto 443 del service al puerto 8080 del contenedor, exponiéndolo externamente en el puerto 30080

Verifico el cambio:"

```bash
kubectl get svc argocd-server -n argocd
```

```
NAME            TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
argocd-server   NodePort   10.43.x.x      <none>        80:XXXXX/TCP,443:30080/TCP   5m
```

"El tipo cambió de `ClusterIP` a `NodePort`. El `443:30080/TCP` confirma que el puerto 443 interno está mapeado al 30080 externo.

`EXTERNAL-IP` muestra `<none>` — eso es correcto con NodePort. La dirección externa no es una IP asignada al Service sino la IP de la instancia EC2 más el puerto."

---

### PASO 2 — Obtener la contraseña inicial (3:30 – 5:30)

> *Pantalla: terminal.*

"ArgoCD genera una contraseña aleatoria durante la instalación y la guarda como un Secret en Kubernetes. La obtengo con este comando:"

```bash
kubectl -n argocd get secret argocd-initial-admin-secret \
  -o jsonpath="{.data.password}" | base64 -d
echo
```

"El pipeline de comandos hace tres cosas:

1. `get secret argocd-initial-admin-secret` — obtiene el Secret que contiene la contraseña
2. `-o jsonpath="{.data.password}"` — extrae solo el campo `password` del JSON
3. `| base64 -d` — decodifica de Base64 (Kubernetes guarda los secrets en Base64)
4. `echo` — agrega una nueva línea para que la contraseña no quede pegada al prompt

El resultado es algo como `xK9mP2rQvN3hL8wJ`. Copio esa cadena — la necesito en el siguiente paso."

---

### PASO 3 — Acceder a la interfaz web (5:30 – 7:30)

> *Pantalla: navegador.*

"Obtengo la IP pública de la EC2 — si no la recuerdas:"

```bash
cd gitops-infra/infrastructure/terraform/jenkins-ec2
terraform output prod_public_ip
```

"Abro en el navegador: **`http://<IP_PUBLICA_EC2>:30080`**

El navegador puede mostrar una advertencia de certificado SSL — ArgoCD usa un certificado autofirmado por defecto. Click en 'Avanzado' → 'Continuar de todas formas'. En un entorno de producción real usarías un certificado válido, pero para el curso el autofirmado es suficiente.

La pantalla de login de ArgoCD aparece.

- **Usuario:** `admin`
- **Contraseña:** la que copié en el paso anterior

Login."

---

### PASO 4 — Cambiar la contraseña (7:30 – 8:30)

> *Pantalla: dashboard de ArgoCD.*

"Buena práctica: cambio la contraseña inmediatamente. La inicial es generada aleatoriamente y no es memorable.

En ArgoCD: click en el usuario `admin` en la esquina superior izquierda → **'User Info'** → **'Update Password'**.

Ingreso la contraseña actual, la nueva (que voy a recordar), y confirmo.

---

También elimino el Secret de la contraseña inicial — ya no lo necesito:"

```bash
kubectl delete secret argocd-initial-admin-secret -n argocd
```

"Esto es una buena práctica de seguridad: el Secret inicial debería existir solo hasta que se cambie la contraseña."

> ⚠️ **ADVERTENCIA DE SEGURIDAD — ArgoCD expuesto sin cifrado**
>
> El dashboard de ArgoCD está accesible en `http://<IP>:30080` — sin TLS/HTTPS. Esto significa que el usuario, la contraseña, y toda la comunicación con ArgoCD **viajan en texto plano por internet**. Un atacante en la red puede capturar las credenciales con un simple sniff de tráfico.
>
> **En producción, ArgoCD SIEMPRE debe usarse con HTTPS.** Alternativas:
>
> - **SSH Tunnel (sin costo):** Acceder vía `ssh -i aws-key.pem -L 8443:localhost:30080 ubuntu@<IP>` y abrir `http://localhost:8443`. Esto cifra todo el tráfico y elimina la necesidad de exponer el puerto 30080 a internet.
> - **Ingress con TLS:** Usar el Traefik incluido en K3s con cert-manager y Let's Encrypt para obtener un certificado SSL válido (requiere un dominio).
>
> Para el curso usamos HTTP directo por simplicidad, pero es importante entender que esta práctica es inaceptable en cualquier entorno con datos reales.

---

### PASO 5 — Explorar el dashboard (8:30 – 9:30)

> *Pantalla: dashboard de ArgoCD.*

"El dashboard está vacío — no hay Applications todavía. Eso es lo que configuramos en el EP40.

Las secciones principales que vamos a usar:

**Applications** — el listado de las aplicaciones que ArgoCD gestiona. En nuestro caso habrá una: `curso-gitops`.

**Settings → Repositories** — donde conectamos `gitops-infra`. Lo hacemos en el EP40.

**Settings → Clusters** — el cluster K3s aparece aquí automáticamente como `in-cluster` porque ArgoCD corre dentro del mismo cluster que gestiona."

---

### CIERRE (9:30 – 10:00)

"Eso es el EP39.

ArgoCD accesible en `http://<IP_EC2>:30080`. Sin LoadBalancer, sin costo adicional. El puerto 30080 que abrimos en Terraform en el EP22 cumplió exactamente el propósito para el que lo configuramos.

En el siguiente episodio conectamos ArgoCD al repositorio privado `gitops-infra` usando el Personal Access Token de GitHub. Y creamos la Application que va a observar el directorio `infrastructure/kubernetes/app/` y sincronizarlo automáticamente con el cluster K3s.

Nos vemos en el EP40."

---

## ✅ Checklist de Verificación

- [ ] `kubectl get svc argocd-server -n argocd` muestra tipo `NodePort`
- [ ] El puerto externo es `30080` en la columna `PORT(S)`
- [ ] La interfaz web carga en `http://<IP_EC2>:30080`
- [ ] Login con `admin` funciona
- [ ] La contraseña fue cambiada desde la inicial
- [ ] El Secret `argocd-initial-admin-secret` fue eliminado

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| El navegador no carga la URL | Verificar que el puerto 30080 está abierto: `aws ec2 describe-security-groups --filters "Name=group-name,Values=prod-sg" --query "SecurityGroups[0].IpPermissions[*].FromPort"` |
| `kubectl patch` falla | ArgoCD no está instalado — verificar EP38 |
| Contraseña incorrecta | Asegurarse de copiar la contraseña completa del `base64 -d` sin espacios ni saltos de línea extra |
| El `patch` no cambia el tipo del Service | Verificar que el JSON del patch es válido — las comillas pueden tener problemas de escape en algunas terminales |

---

## 🗒️ Notas de Producción

- La intro mostrando el Service en ClusterIP y explicando el costo del LoadBalancer es el gancho narrativo del episodio — la solución NodePort contrasta perfectamente.
- El `kubectl patch` es un comando que mucha gente no conoce — explicar su función general brevemente antes de mostrar el comando específico.
- Al mostrar la advertencia de certificado SSL en el navegador, explicar que es normal y seguro para el entorno del curso — muchos alumnos se asustan y detienen aquí.
- Eliminar el Secret inicial después de cambiar la contraseña — hacerlo en vivo da el hábito correcto al alumno.
- Mostrar brevemente el dashboard vacío al final y nombrar las secciones que se van a usar en el EP40 — crea la anticipación del siguiente episodio.

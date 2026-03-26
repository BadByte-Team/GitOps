# EP 43: SonarQube Local — Proyecto y Token de Análisis

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)
**🔄 MODIFICADO:** SonarQube corre en la misma máquina local que Jenkins, levantado desde el `docker-compose.yml` del EP31. No requiere EC2 adicional.

---

## 🎯 Objetivo
Verificar que SonarQube está operativo, crear el proyecto `curso-gitops`, generar el token de análisis que Jenkins usará para autenticarse, y entender qué mide un Quality Gate.

---

## 📋 Prerequisitos
- `docker compose up -d` ejecutado en `~/local-ci/` (EP31)
- `vm.max_map_count=262144` configurado en el sistema (EP31)

---

## 🧠 ¿Qué analiza SonarQube?

SonarQube analiza el código fuente — no la imagen Docker. Busca:

| Categoría | Qué detecta |
|---|---|
| **Bugs** | Código que casi seguramente produce comportamiento incorrecto |
| **Vulnerabilities** | Patrones de código con riesgos de seguridad |
| **Code Smells** | Código que funciona pero es difícil de mantener |
| **Coverage** | Porcentaje del código cubierto por tests |
| **Duplications** | Código duplicado que debería refactorizarse |

El **Quality Gate** es un conjunto de umbrales configurables. Si el código no los supera, el pipeline falla. Por ejemplo: "no pasar si hay más de 0 vulnerabilidades nuevas" o "no pasar si la cobertura de tests baja del 80%".

En el curso usamos el Quality Gate por defecto de SonarQube — es suficientemente estricto para nuestro código Go sin ser bloqueante en el aprendizaje.

---

## 🎬 Guión del Video

### INTRO (0:00 – 1:00)

> *Pantalla: `docker compose ps` en `~/local-ci/` mostrando los tres contenedores — jenkins, sonarqube, sonar-db — todos en `Up`.*

"Bienvenidos al episodio 43.

SonarQube ya está corriendo desde el EP31. Lo levantamos junto con Jenkins en el mismo `docker-compose.yml`. En este episodio lo configuramos para el curso: creamos el proyecto, generamos el token de análisis, y entendemos qué es un Quality Gate.

El EP44 conecta este token con Jenkins. El EP45 agrega los stages de seguridad al Jenkinsfile. Este episodio es el puente entre los dos.

Empecemos."

---

### PASO 1 — Verificar que SonarQube está operativo (1:00 – 2:30)

> *Pantalla: terminal y navegador.*

"Primero verifico que los tres contenedores están corriendo:"

```bash
cd ~/local-ci
docker compose ps
```

```
NAME         IMAGE                    STATUS
jenkins      jenkins/jenkins:lts      Up
sonarqube    sonarqube:lts-community  Up
sonar-db     postgres:15-alpine       Up
```

"Si algún contenedor no está `Up`, lo levanto:"

```bash
docker compose up -d
```

"Si SonarQube se reinicia constantemente, el `vm.max_map_count` no está configurado:"

```bash
sudo sysctl -w vm.max_map_count=262144
echo 'vm.max_map_count=262144' | sudo tee -a /etc/sysctl.conf
docker compose restart sonarqube
```

"Abro SonarQube en el navegador: **`http://localhost:9000`**

Si la pantalla dice 'SonarQube is starting...' — esperar 1-2 minutos. El proceso de arranque interno de SonarQube tarda más que el arranque del contenedor.

Login con `admin` / `admin`. Si ya cambié la contraseña en el EP31, uso la nueva."

---

### 🔍 PAUSA CONCEPTUAL — El Quality Gate (2:30 – 4:30)

> *Pantalla: la pantalla de Quality Gates en SonarQube.*

"Antes de crear el proyecto, quiero que entiendan qué es el Quality Gate porque va a aparecer en el pipeline del EP45.

Un Quality Gate es un conjunto de condiciones que el código debe cumplir para que el análisis sea considerado exitoso. Si alguna condición falla, el Quality Gate reporta 'Failed' y el pipeline puede detenerse.

SonarQube incluye un Quality Gate predeterminado llamado 'Sonar way' que define condiciones sobre el **código nuevo** — el que se agrega en cada commit. Las condiciones típicas son:

- 0 bugs nuevos de severidad blocker o critical
- 0 vulnerabilidades nuevas
- Cobertura de tests del nuevo código ≥ 80%
- Duplicación del nuevo código ≤ 3%

Para el curso de GitOps, el análisis de cobertura requeriría tests unitarios en el código Go — algo que está fuera del scope del curso. Por eso vamos a configurar el pipeline para que el stage de SonarQube **reporte** el análisis pero no bloquee el pipeline por cobertura insuficiente.

Lo veremos en detalle en el EP45."

---

### PASO 2 — Crear el proyecto en SonarQube (4:30 – 7:00)

> *Pantalla: navegador en SonarQube.*

"En SonarQube: **Projects** → **Create Project** → **Manually**.

Completo el formulario:
- **Project display name:** `curso-gitops` — el nombre visible en la UI
- **Project key:** `curso-gitops` — el identificador técnico. Debe coincidir exactamente con el `-Dsonar.projectKey=curso-gitops` del Jenkinsfile
- **Main branch name:** `main`

Click en **Next**.

SonarQube pregunta cómo configurar el análisis. Selecciono **'Use the global setting'** — usamos el Quality Gate 'Sonar way' que viene por defecto.

Click en **Create project**.

El proyecto aparece vacío — sin análisis todavía. El primer análisis llegará cuando el pipeline de Jenkins corra con el stage de SonarQube activo."

---

### PASO 3 — Generar el token de análisis (7:00 – 9:30)

> *Pantalla: navegador en SonarQube — pantalla de tokens.*

"Ahora el token que Jenkins va a usar para autenticarse con SonarQube.

	En SonarQube: click en tu usuario en la esquina superior derecha → **My Account** → pestaña **Security** → sección **Generate Tokens**.

Completo:
- **Name:** `jenkins-token` — descripción para recordar para qué es
- **Type:** `Global Analysis Token` — permite analizar cualquier proyecto
- **Expires in:** `No expiration` — para el curso no necesitamos que expire

Click en **Generate**.

El token aparece **una sola vez**. Empieza con `squ_`. Lo copio inmediatamente en un lugar seguro — si cierro esta pantalla sin copiarlo, tendré que generar uno nuevo.

El token se ve así: `squ_abc123def456ghi789jkl012mno345pqr678`

No lo pego en el Jenkinsfile ni en ningún archivo que vaya al repositorio. En el EP44 lo vamos a guardar como credencial en Jenkins — cifrado, no visible."

---

### PASO 4 — Explorar el dashboard de calidad (9:30 – 11:00)

> *Pantalla: navegador en SonarQube — pantalla del proyecto vacío.*

"El proyecto existe pero sin datos todavía. Cuando el pipeline corra por primera vez con el análisis activo, aquí vamos a ver:

**Overall Code** — el estado general del repositorio: cuántos bugs, vulnerabilidades, code smells y duplicaciones tiene el código base completo.

**New Code** — el estado del código agregado en el último período. Este es el que el Quality Gate evalúa — solo miramos hacia adelante, no juzgamos el código legado.

**Activity** — el historial de análisis. Con el tiempo, verás si la calidad mejora o empeora commit a commit.

Para el código Go del curso, que tiene autenticación JWT, handlers HTTP y queries a MySQL, SonarQube debería encontrar code smells menores pero probablemente ninguna vulnerabilidad crítica — el código fue escrito teniendo en cuenta las mejores prácticas desde el inicio."

---

### CIERRE (11:00 – 12:00)

"Eso es el EP43.

SonarQube operativo, el proyecto `curso-gitops` creado, y el token de análisis generado. Todo listo para la integración con Jenkins.

En el siguiente episodio configuramos Jenkins para que sepa dónde está SonarQube — la URL del servidor y el token de autenticación. Con eso, el stage de análisis del Jenkinsfile va a poder enviar los resultados automáticamente.

Nos vemos en el EP44."

---

## ✅ Checklist de Verificación
- [ ] `docker compose ps` muestra jenkins, sonarqube y sonar-db en `Up`
- [ ] SonarQube responde en `http://localhost:9000`
- [ ] La contraseña de `admin` fue cambiada desde la inicial
- [ ] El proyecto `curso-gitops` existe con project key `curso-gitops`
- [ ] El token `jenkins-token` fue generado y copiado
- [ ] Entiendes qué es un Quality Gate y qué condiciones evalúa

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| SonarQube se reinicia constantemente | `sudo sysctl -w vm.max_map_count=262144` y `docker compose restart sonarqube` |
| La pantalla muestra 'SonarQube is starting...' por más de 5 min | Revisar logs: `docker compose logs sonarqube --tail=30` |
| `http://localhost:9000` no carga | Verificar que el contenedor está en `Up`: `docker compose ps` |
| Se perdió la configuración del proyecto | Los volúmenes no están definidos en el compose — verificar `docker-compose.yml` del EP31 |
| El token no aparece después de generarlo | Solo aparece una vez — si se perdió, hacer click en 'Revoke' y generar uno nuevo |

---

## 🗒️ Notas de Producción
- La pausa conceptual del Quality Gate es importante para que el alumno entienda por qué en el EP45 el pipeline puede fallar por razones de calidad — no es un error técnico.
- Al generar el token, hacer una pausa dramática antes de navegar fuera de la página — "este token solo aparece una vez, lo copio ahora".
- Mostrar el proyecto vacío al final y anticipar cómo va a verse después del primer análisis — crea expectativa del EP44 y EP45.
- La explicación de que el Quality Gate evalúa código *nuevo* vs código *total* es el concepto que más confunde a los alumnos — explicarlo con calma.

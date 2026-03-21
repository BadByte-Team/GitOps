# EP 32: Instalación de Plugins Indispensables

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 10–12 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar los plugins que el pipeline de CI necesita para funcionar: Docker Pipeline, SonarQube Scanner, NodeJS Plugin y Eclipse Temurin (JDK). Sin estos plugins, los stages del Jenkinsfile no van a encontrar las herramientas que buscan.

---

## 📋 Prerequisitos
- Jenkins corriendo en `http://localhost:8080` (EP31)
- Usuario `admin` creado

---

## 🧠 ¿Por qué plugins?

Jenkins por sí solo es un motor de automatización genérico. Los plugins son los que le dan capacidades específicas. Sin el plugin de Docker, Jenkins no puede ejecutar `docker build`. Sin el de SonarQube, no puede enviar el análisis. Sin el de JDK, no puede compilar Java ni usar herramientas que lo requieran.

| Plugin | Para qué lo usamos |
|---|---|
| Docker Pipeline | Permite usar comandos Docker en el Jenkinsfile |
| SonarQube Scanner | Integra el análisis de código con SonarQube local |
| NodeJS Plugin | Herramientas de Node.js disponibles en el pipeline |
| Eclipse Temurin installer | Gestiona la instalación de JDK dentro de Jenkins |

---

## 🎬 Guión del Video

### INTRO (0:00 – 0:45)

> *Pantalla: el dashboard de Jenkins recién configurado del EP31.*

"Bienvenidos al episodio 32.

Jenkins está corriendo pero está vacío. Si ejecutáramos el pipeline ahora mismo, fallaría en el primer stage que necesite Docker, o en el que necesite el scanner de SonarQube. Los plugins son lo que convierte este Jenkins básico en la herramienta que necesitamos.

Son cuatro plugins. El proceso es el mismo para todos. Vamos directo."

---

### PASO 1 — Ir a la gestión de plugins (0:45 – 1:30)

> *Pantalla: navegador en Jenkins.*

"Desde el dashboard de Jenkins:

**Manage Jenkins** → **Plugins** → pestaña **Available plugins**

Aquí está el listado completo de plugins disponibles para instalar. La barra de búsqueda en la parte superior filtra en tiempo real."

---

### PASO 2 — Instalar los plugins (1:30 – 7:00)

> *Pantalla: navegador buscando e instalando cada plugin.*

"Busco e instalo los cuatro. Los marco todos antes de instalar para hacerlo en un solo paso.

---

**Docker Pipeline**

Busco `Docker Pipeline` en la barra de búsqueda. Aparece como 'Docker Pipeline' del editor CloudBees. Marco el checkbox.

Este plugin es el que permite escribir comandos Docker directamente en el Jenkinsfile. Sin él, la línea `docker build -t ${DOCKER_IMAGE}:${BUILD_TAG} .` simplemente no existe para Jenkins.

---

**SonarQube Scanner**

Busco `SonarQube Scanner`. Marco el checkbox.

Este plugin hace dos cosas: primero, le dice a Jenkins cómo hablar con un servidor de SonarQube — en nuestro caso el que corre localmente en `http://localhost:9000`. Segundo, hace disponible la herramienta `sonar-scanner` en el pipeline para que el stage de análisis pueda ejecutarla.

---

**NodeJS Plugin**

Busco `NodeJS`. Marco el checkbox.

El scanner de SonarQube para algunos análisis necesita Node.js. Este plugin gestiona la instalación de Node.js dentro del entorno de Jenkins y lo hace disponible en el pipeline a través del bloque `tools { nodejs 'node18' }`.

---

**Eclipse Temurin installer**

Busco `Eclipse Temurin`. Marco el checkbox.

Este plugin es el que le permite a Jenkins descargar y gestionar versiones de JDK automáticamente. En el Jenkinsfile tenemos `tools { jdk 'jdk17' }` — ese nombre `jdk17` va a corresponder a una configuración que crearemos en el EP33 usando este plugin.

---

Con los cuatro marcados, click en **'Install'** en la parte superior o inferior de la lista.

Jenkins muestra el progreso de instalación de cada plugin. Tarda entre 2 y 4 minutos. Al final aparece la opción de reiniciar — marco la casilla **'Restart Jenkins when installation is complete and no jobs are running'**.

Jenkins se reinicia automáticamente cuando termina."

---

### PASO 3 — Verificar la instalación (7:00 – 8:30)

> *Pantalla: navegador en Jenkins → Manage Jenkins → Plugins → Installed plugins.*

"Después del reinicio, vuelvo a **Manage Jenkins → Plugins → Installed plugins**.

En la barra de búsqueda verifico que los cuatro plugins aparecen como instalados:"

```
Docker Pipeline        ✅
SonarQube Scanner      ✅
NodeJS                 ✅
Eclipse Temurin        ✅
```

"Si alguno no aparece, volver a la pestaña 'Available' y buscarlo de nuevo — puede que no se haya instalado correctamente en el primer intento."

---

### PASO 4 — Un plugin más: Pipeline Utility Steps (8:30 – 9:30)

> *Pantalla: navegador buscando el plugin.*

"Mientras estamos en la gestión de plugins, instalo uno adicional que es útil para trabajar con archivos en el pipeline:

**Pipeline Utility Steps** — proporciona funciones helper como leer archivos, trabajar con JSON, etc. No lo usamos directamente en el Jenkinsfile del curso, pero está en la mayoría de los pipelines de producción y vale la pena tenerlo.

Misma operación: buscar, marcar, instalar."

---

### CIERRE (9:30 – 10:00)

"Eso es el EP32.

Los cuatro plugins esenciales instalados y verificados. Jenkins ya tiene las herramientas que el pipeline va a necesitar.

En el siguiente episodio configuramos las herramientas — Tools en Manage Jenkins. Ahí le decimos a Jenkins exactamente qué versión de JDK usar y dónde encontrar el scanner de SonarQube. Es el paso que une los plugins que acabamos de instalar con los nombres que usa el Jenkinsfile.

Nos vemos en el EP33."

---

## ✅ Checklist de Verificación
- [ ] Docker Pipeline aparece en 'Installed plugins'
- [ ] SonarQube Scanner aparece en 'Installed plugins'
- [ ] NodeJS aparece en 'Installed plugins'
- [ ] Eclipse Temurin aparece en 'Installed plugins'
- [ ] Jenkins se reinició correctamente después de la instalación

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| El plugin no aparece en 'Available' | Actualizar el catálogo: Manage Jenkins → Plugins → pestaña 'Advanced' → botón 'Check now' |
| La instalación falla con error de red | El contenedor de Jenkins no tiene acceso a internet — verificar la red de Docker |
| Jenkins no reinicia después de la instalación | Reiniciar manualmente: `docker compose restart jenkins` |
| Un plugin instalado no aparece en 'Installed' | Buscar por el nombre exacto — el nombre visible puede diferir del nombre interno |

---

## 🗒️ Notas de Producción
- Buscar e instalar los cuatro plugins de forma consecutiva antes de hacer click en 'Install' — es más eficiente y hace el video más fluido que instalarlos uno a uno.
- Mientras Jenkins instala, explicar brevemente para qué sirve cada plugin en lugar de simplemente esperar en silencio.
- El reinicio automático al final es el momento correcto para hacer una pausa narrativa y anticipar el EP33.
- Mostrar la verificación en 'Installed plugins' — da certeza visual al alumno de que todo está bien.

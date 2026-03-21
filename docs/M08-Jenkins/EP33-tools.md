# EP 33: Configuración de Tools — JDK, Node y SonarQube Scanner

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 8–10 min
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Configurar las herramientas en **Manage Jenkins → Tools** para que el Jenkinsfile pueda referenciarlas por nombre. El pipeline usa `tools { jdk 'jdk17' }` y `tool('sonar-scanner')` — esos nombres tienen que existir aquí primero.

---

## 📋 Prerequisitos
- Plugins instalados: Eclipse Temurin, NodeJS, SonarQube Scanner (EP32)

---

## 🧠 ¿Qué es Tools en Jenkins?

El bloque `tools` del Jenkinsfile permite declarar herramientas que Jenkins gestiona automáticamente. En lugar de que el pipeline asuma que JDK está instalado en el servidor, Jenkins descarga y configura la versión exacta que le indiques.

El flujo es:
1. Aquí en Tools: "cuando alguien pida `jdk17`, usa Eclipse Temurin versión 17"
2. En el Jenkinsfile: `tools { jdk 'jdk17' }` → Jenkins busca la configuración por ese nombre

Si el nombre en Tools no coincide exactamente con el nombre en el Jenkinsfile, el pipeline falla con `Tool type "jdk" does not have an install of "jdk17"`.

---

## 🎬 Guión del Video

### INTRO (0:00 – 0:45)

> *Pantalla: el bloque `tools` del Jenkinsfile del proyecto abierto en VS Code.*

"Bienvenidos al episodio 33.

Miren este bloque del Jenkinsfile. `jdk 'jdk17'`, `nodejs 'node18'`, `tool('sonar-scanner')`. Tres referencias a herramientas por nombre.

Esos nombres tienen que estar configurados en Jenkins antes de que el pipeline los pueda usar. Si el pipeline intenta resolver `jdk17` y Jenkins no sabe qué es eso, el build falla inmediatamente.

Eso es lo que configuramos en este episodio. Es rápido. Vamos directo."

---

### PASO 1 — Ir a Tools (0:45 – 1:15)

> *Pantalla: navegador en Jenkins.*

"Desde el dashboard:

**Manage Jenkins** → **Tools**

Esta sección centraliza todas las herramientas que Jenkins puede gestionar automáticamente."

---

### PASO 2 — Configurar JDK 17 (1:15 – 3:30)

> *Pantalla: navegador en la sección JDK de Tools.*

"Busco la sección **JDK installations** y click en **'Add JDK'**.

Completo el formulario:
- **Name:** `jdk17` — exactamente así, en minúsculas, sin espacios. Este es el nombre que el Jenkinsfile va a buscar.
- **Install automatically:** ✅ activado
- En el menú de instaladores aparece **'Install from adoptium.net'** — lo selecciono
- **Version:** `jdk-17.x.x+x` — la versión 17 más reciente disponible

Click en **'Apply'** — no en 'Save' todavía, porque voy a configurar las otras herramientas en el mismo paso.

El nombre `jdk17` importa mucho. Si el Jenkinsfile dice `jdk 'jdk17'` y aquí escriben `JDK17` con mayúscula o `jdk-17` con guion, el pipeline va a fallar. Tiene que ser letra por letra idéntico."

---

### PASO 3 — Configurar Node.js 18 (3:30 – 5:30)

> *Pantalla: navegador en la sección NodeJS de Tools.*

"Busco la sección **NodeJS installations** y click en **'Add NodeJS'**.

- **Name:** `node18` — de nuevo, exactamente como aparece en el Jenkinsfile
- **Install automatically:** ✅ activado
- **Version:** `NodeJS 18.x.x` — la versión 18 más reciente

Click en **'Apply'**."

---

### PASO 4 — Configurar SonarQube Scanner (5:30 – 7:30)

> *Pantalla: navegador en la sección SonarQube Scanner de Tools.*

"Busco la sección **SonarQube Scanner installations** y click en **'Add SonarQube Scanner'**.

- **Name:** `sonar-scanner` — tal cual aparece en el Jenkinsfile con el método `tool('sonar-scanner')`
- **Install automatically:** ✅ activado
- La versión más reciente disponible

Click en **'Save'** — ahora sí guardo todo."

---

### PASO 5 — Verificar que los nombres coinciden con el Jenkinsfile (7:30 – 9:00)

> *Pantalla: VS Code con el Jenkinsfile del proyecto abierto lado a lado con Jenkins Tools.*

"Antes de cerrar el episodio, hago la verificación más importante: comparar los nombres configurados aquí con los nombres usados en el Jenkinsfile.

El Jenkinsfile del proyecto tiene:"

```groovy
tools {
    jdk 'jdk17'       ← debe coincidir con el Name en JDK installations
    nodejs 'node18'   ← debe coincidir con el Name en NodeJS installations
}

environment {
    SCANNER_HOME = tool('sonar-scanner')  ← debe coincidir con SonarQube Scanner
}
```

"Jenkins Tools tiene:
- JDK: `jdk17` ✅
- NodeJS: `node18` ✅
- SonarQube Scanner: `sonar-scanner` ✅

Los tres coinciden. Si en algún episodio futuro el pipeline falla con un error como `Tool type "jdk" does not have an install of "..."`, este es el primer lugar donde venir a verificar."

---

### CIERRE (9:00 – 9:30)

"Eso es el EP33.

JDK 17, Node.js 18 y SonarQube Scanner configurados con los nombres exactos que el Jenkinsfile va a buscar.

En el siguiente episodio configuramos las credenciales — Docker Hub y GitHub. Son los dos recursos externos a los que el pipeline necesita autenticarse para hacer push de imágenes y push de commits.

Nos vemos en el EP34."

---

## ✅ Checklist de Verificación
- [ ] JDK installations tiene una entrada llamada exactamente `jdk17`
- [ ] NodeJS installations tiene una entrada llamada exactamente `node18`
- [ ] SonarQube Scanner installations tiene una entrada llamada exactamente `sonar-scanner`
- [ ] Los nombres coinciden carácter por carácter con el Jenkinsfile del proyecto

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| No aparece la sección JDK en Tools | El plugin Eclipse Temurin no está instalado — volver al EP32 |
| No aparece la sección NodeJS | El plugin NodeJS no está instalado — volver al EP32 |
| El pipeline falla con `Tool type "jdk" does not have an install of "jdk17"` | El nombre en Tools no coincide con el del Jenkinsfile — verificar mayúsculas y guiones |
| La descarga automática del JDK falla | El contenedor de Jenkins no tiene acceso a internet — verificar `docker compose ps` |

---

## 🗒️ Notas de Producción
- Abrir el Jenkinsfile del proyecto en VS Code antes de empezar y mantenerlo visible — los nombres en el Jenkinsfile son la referencia que guía cada configuración.
- Al escribir el nombre de cada herramienta, decirlo en voz alta letra por letra — `j-d-k-1-7` — para enfatizar que la coincidencia exacta es crítica.
- La comparación final side-by-side entre Tools y el Jenkinsfile es el momento de mayor valor pedagógico del episodio — tomarse tiempo en ella.

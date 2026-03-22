# EP 44: Integrar SonarQube con Jenkins — Token y Webhook

**Tipo:** CONFIGURACIÓN
**Duración estimada:** 12–15 min
**Dificultad:** ⭐⭐ (Intermedio)

---

## 🎯 Objetivo
Conectar Jenkins con SonarQube de forma bidireccional: Jenkins envía el análisis a SonarQube usando el token del EP43, y SonarQube notifica a Jenkins el resultado del Quality Gate mediante un webhook. Sin esta integración, el stage `SonarQube Analysis` del Jenkinsfile no puede comunicarse con SonarQube.

---

## 📋 Prerequisitos
- Jenkins corriendo en `http://localhost:8080` (EP31)
- Plugin SonarQube Scanner instalado (EP32)
- SonarQube corriendo en `http://localhost:9000` con el token generado (EP43)

---

## 🧠 La comunicación bidireccional

La integración tiene dos sentidos y ambos son necesarios:

```
Jenkins ──────── token ──────────▶ SonarQube
         (envía código a analizar)  (analiza y guarda resultados)

SonarQube ────── webhook ──────▶ Jenkins
            (notifica: Quality Gate pasó o falló)
```

**Sin el token:** Jenkins no puede autenticarse con SonarQube. El análisis falla en el primer intento.

**Sin el webhook:** Jenkins envía el análisis pero nunca sabe si el Quality Gate pasó o falló. El stage `waitForQualityGate` quedaría esperando indefinidamente.

---

## 🎬 Guión del Video

### INTRO (0:00 – 0:50)

> *Pantalla: el stage `SonarQube Analysis` del Jenkinsfile abierto en VS Code con `withSonarQubeEnv('sonarqube-server')` visible.*

"Bienvenidos al episodio 44.

Esta línea del Jenkinsfile — `withSonarQubeEnv('sonarqube-server')` — está intentando usar una configuración de SonarQube llamada `sonarqube-server`. Si esa configuración no existe en Jenkins, el stage falla inmediatamente.

Y `waitForQualityGate` necesita que SonarQube le notifique a Jenkins el resultado del análisis. Sin un webhook configurado, espera para siempre.

Hoy configuramos las dos partes. Son 10 minutos de configuración que hacen que el análisis de calidad funcione completamente.

Empecemos."

---

### PASO 1 — Agregar el token de SonarQube como credencial en Jenkins (0:50 – 3:30)

> *Pantalla: navegador en Jenkins.*

"Primero guardo el token del EP43 como credencial en Jenkins, de modo que Jenkins pueda usarlo para autenticarse con SonarQube sin exponerlo en texto plano.

**Manage Jenkins** → **Credentials** → **System** → **Global credentials** → **Add Credentials**:

- **Kind:** `Secret text`
- **Scope:** Global
- **Secret:** pego el token `squ_abc123...` del EP43
- **ID:** `sonarqube-token`
- **Description:** `SonarQube Analysis Token`

Click en **Create**."

---

### PASO 2 — Configurar el servidor de SonarQube en Jenkins (3:30 – 7:00)

> *Pantalla: navegador en Jenkins → Manage Jenkins → System.*

"Ahora le digo a Jenkins dónde está el servidor de SonarQube y cómo autenticarse.

**Manage Jenkins** → **System** → busco la sección **SonarQube servers** → **Add SonarQube**:

- **Name:** `sonarqube-server` — este nombre tiene que coincidir exactamente con el que usa el Jenkinsfile en `withSonarQubeEnv('sonarqube-server')`. Si ponen `SonarQube` con mayúscula o `sonarqube_server` con guion bajo, el pipeline no lo encontrará.
- **Server URL:** `http://localhost:9000` — la URL donde corre SonarQube localmente
- **Server authentication token:** selecciono `sonarqube-token` — la credencial que acabo de crear

Click en **Save**.

---

Un momento para aclarar algo que confunde a mucha gente aquí. La URL es `http://localhost:9000`. Pero Jenkins corre dentro de un contenedor Docker, y `localhost` dentro de ese contenedor es el propio contenedor Jenkins, no la máquina host.

¿Por qué funciona entonces? Porque el contenedor de Jenkins y el de SonarQube comparten la misma red Docker — la red `local-ci_default` que Docker Compose creó. Dentro de esa red, `localhost` sí resuelve al host. También se puede usar el nombre del contenedor `http://sonarqube:9000` — ambas opciones funcionan porque comparten la red de Compose."

---

### PASO 3 — Configurar el webhook en SonarQube (7:00 – 10:30)

> *Pantalla: navegador en SonarQube.*

"La segunda parte de la integración: el webhook que SonarQube usa para notificar a Jenkins cuando el análisis termina.

En SonarQube: **Administration** → **Configuration** → **Webhooks** → **Create**:

- **Name:** `jenkins`
- **URL:** `http://localhost:8080/sonarqube-webhook/`
- **Secret:** (dejar vacío para el curso — en producción se usaría un secret para verificar la autenticidad de la notificación)

Click en **Create**.

El webhook aparece en la lista. Hay un botón de prueba — **'Test'** — que SonarQube puede usar para verificar que la URL de Jenkins es alcanzable.

---

La misma aclaración sobre `localhost` aplica aquí: SonarQube también corre en un contenedor Docker dentro de la misma red. `localhost:8080` resuelve al servicio Jenkins de esa red. Si esto no funciona en tu configuración específica, reemplaza `localhost` por la IP de tu máquina en la red local — `192.168.x.x` — o por el nombre del servicio `http://jenkins:8080/sonarqube-webhook/`."

---

### PASO 4 — Verificar la integración con un análisis de prueba (10:30 – 13:00)

> *Pantalla: Jenkins ejecutando el pipeline y SonarQube mostrando resultados.*

"La única forma real de verificar que todo está bien conectado es ejecutar el pipeline.

En Jenkins: abro el pipeline `curso-gitops-ci` → **Build Now**.

Mientras el stage `SonarQube Analysis` corre, abro SonarQube en el navegador. En unos momentos el proyecto `curso-gitops` debería mostrar que hay un análisis en progreso — un indicador de 'Computing'.

Cuando el análisis termina:
1. SonarQube procesa los resultados
2. Evalúa el Quality Gate
3. Envía el webhook a `localhost:8080/sonarqube-webhook/`
4. Jenkins recibe la notificación
5. El stage `waitForQualityGate` recibe el resultado y continúa o falla

En el Console Output de Jenkins busco esta línea — confirma que la comunicación funcionó:"

```
ANALYSIS SUCCESSFUL, you can find the results at:
http://localhost:9000/dashboard?id=curso-gitops
```

"Y en SonarQube aparece el primer análisis del proyecto con los resultados del Quality Gate. Verde significa que el código pasó todas las condiciones. Naranja o rojo significa que algo falló y habría que revisarlo antes de que el pipeline llegue a producción."

---

### CIERRE (13:00 – 14:00)

"Eso es el EP44.

La integración bidireccional está completa: Jenkins puede enviar análisis a SonarQube y SonarQube puede notificar a Jenkins el resultado del Quality Gate.

En el siguiente episodio revisamos el Jenkinsfile completo con todos los stages de seguridad — SonarQube, Quality Gate y Trivy — para entender el flujo completo y ejecutamos el pipeline viendo cómo cada stage de seguridad cumple su rol.

Nos vemos en el EP45."

---

## ✅ Checklist de Verificación
- [ ] Credencial `sonarqube-token` existe en Jenkins como Secret text
- [ ] Servidor `sonarqube-server` está configurado en Manage Jenkins → System
- [ ] La URL del servidor es `http://localhost:9000`
- [ ] El webhook `jenkins` está configurado en SonarQube apuntando a `http://localhost:8080/sonarqube-webhook/`
- [ ] El pipeline ejecuta el stage `SonarQube Analysis` sin errores
- [ ] El proyecto `curso-gitops` en SonarQube muestra el primer análisis con su Quality Gate

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `Could not connect to SonarQube server` en Jenkins | La URL del servidor es incorrecta o SonarQube no está corriendo — `docker compose ps` |
| `Invalid credentials` al analizar | El token de SonarQube es incorrecto — verificar que `sonarqube-token` tiene el valor correcto |
| `waitForQualityGate` espera indefinidamente | El webhook no está configurado en SonarQube o la URL es incorrecta |
| El Quality Gate falla con `No conditions` | El proyecto no tiene condiciones definidas — asignarle el Quality Gate 'Sonar way' en la configuración del proyecto |
| SonarQube no puede alcanzar `localhost:8080` | Usar `http://jenkins:8080/sonarqube-webhook/` — el nombre del servicio de Compose |

---

## 🗒️ Notas de Producción
- La explicación de la comunicación bidireccional con el diagrama es el corazón teórico del episodio — el alumno necesita entender por qué se necesitan las dos configuraciones (token Y webhook).
- La aclaración sobre `localhost` dentro de contenedores Docker es un punto de confusión frecuente — explicarla con calma antes de seguir.
- Ejecutar el pipeline al final y ver el primer análisis en SonarQube es el cierre más satisfactorio del episodio — muestra que la integración funciona.
- Al mostrar los resultados en SonarQube, señalar el Quality Gate y su estado — es el concepto del EP43 materializándose visualmente.

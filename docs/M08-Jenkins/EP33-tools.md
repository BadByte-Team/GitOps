# EP 33: Configurar Herramientas en Manage Jenkins > Tools

**Tipo:** CONFIGURACION
## Objetivo
Configurar JDK, NodeJS y Docker en Jenkins para usarlos en pipelines.

## Pasos Detallados

Ir a: **Manage Jenkins → Tools**

### 1. JDK
- Click "Add JDK"
- Name: `jdk17`
- Install automatically: ✅
- Installer: "Install from adoptium.net" → Version: `17`

### 2. NodeJS
- Click "Add NodeJS"
- Name: `node18`
- Install automatically: ✅
- Version: `18.x`

### 3. SonarQube Scanner
- Click "Add SonarQube Scanner"
- Name: `sonar-scanner`
- Install automatically: ✅

### 4. Docker (si no está en PATH)
Docker ya está instalado en la EC2 por el script. Solo verificar:
```bash
docker --version  # En la EC2
```

## Verificación
- [ ] Las 3 herramientas están configuradas en Tools
- [ ] Los nombres coinciden con los del Jenkinsfile (jdk17, node18, sonar-scanner)

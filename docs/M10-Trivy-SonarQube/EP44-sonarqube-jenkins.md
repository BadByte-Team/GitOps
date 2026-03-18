# EP 44: Integrar SonarQube con Jenkins — Token y Webhook

**Tipo:** CONFIGURACION
## Objetivo
Conectar SonarQube con Jenkins para análisis de código en el pipeline.

## Pasos Detallados

### 1. Generar Token en SonarQube
1. SonarQube → My Account → Security
2. Generate Token → nombre: `jenkins` → Generate
3. **Copiar el token** (solo se muestra una vez)

### 2. Agregar Token en Jenkins
1. Manage Jenkins → Credentials → Global
2. Add Credentials:
   - Kind: **Secret text**
   - ID: `sonarqube-token`
   - Secret: (pegar el token)

### 3. Configurar SonarQube Server en Jenkins
1. Manage Jenkins → System → SonarQube servers
2. Add SonarQube:
   - Name: `sonarqube-server`
   - Server URL: `http://localhost:9000`
   - Server authentication token: seleccionar `sonarqube-token`

### 4. Configurar Webhook en SonarQube
1. SonarQube → Administration → Configuration → Webhooks
2. Create:
   - Name: `jenkins`
   - URL: `http://localhost:8080/sonarqube-webhook/`

## Verificación
- [ ] Token generado en SonarQube
- [ ] Credencial `sonarqube-token` en Jenkins
- [ ] Server `sonarqube-server` configurado en Jenkins
- [ ] Webhook de SonarQube apuntando a Jenkins

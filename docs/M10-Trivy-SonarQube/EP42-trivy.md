# EP 42: Instalación de Trivy en EC2 Jenkins

**Tipo:** INSTALACION / PRACTICA
## Objetivo
Instalar Trivy en la EC2 Jenkins y escanear una imagen Docker.

## Pasos Detallados

### 1. Instalar Trivy
```bash
# Usar el script del proyecto
chmod +x infrastructure/scripts/install-trivy.sh
./infrastructure/scripts/install-trivy.sh

# Verificar
trivy --version
```

### 2. Escanear una Imagen
```bash
# Escanear imagen local
trivy image TU_USUARIO/curso-gitops:latest

# Solo vulnerabilidades HIGH y CRITICAL
trivy image --severity HIGH,CRITICAL TU_USUARIO/curso-gitops:latest

# Fallar si hay vulnerabilidades CRITICAL
trivy image --exit-code 1 --severity CRITICAL TU_USUARIO/curso-gitops:latest
```

### 3. Entender el Reporte
```
┌───────────────┬──────────┬──────────┬──────────────────┐
│    Library    │ Severity │  Status  │     Fixed In     │
├───────────────┼──────────┼──────────┼──────────────────┤
│ libcrypto3    │ HIGH     │ fixed    │ 3.1.4-r5         │
│ libssl3       │ CRITICAL │ fixed    │ 3.1.4-r5         │
└───────────────┴──────────┴──────────┴──────────────────┘
```

## Archivos Involucrados
- `infrastructure/scripts/install-trivy.sh`
- `infrastructure/jenkins/Jenkinsfile` (stage "Trivy Scan")

## Verificación
- [ ] Trivy está instalado y funciona
- [ ] Puedes escanear imágenes y leer el reporte

# EP 15: Instalación y Configuración de AWS CLI

**Tipo:** INSTALACION / CONFIGURACION
## Objetivo
Instalar AWS CLI v2 y configurar las credenciales IAM.

## Pasos Detallados

### 1. Instalar AWS CLI v2
```bash
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
sudo apt install -y unzip
unzip awscliv2.zip
sudo ./aws/install
rm -rf aws awscliv2.zip

# Verificar
aws --version
```

### 2. Configurar Credenciales
```bash
aws configure
# AWS Access Key ID: (pegar tu Access Key)
# AWS Secret Access Key: (pegar tu Secret Key)
# Default region name: us-east-1
# Default output format: json
```

### 3. Verificar Conexión
```bash
aws sts get-caller-identity
# Debe mostrar tu Account ID, ARN y User
```

## Archivos Involucrados
- `infrastructure/scripts/install-jenkins.sh` — ya incluye instalación de AWS CLI

## Verificación
- [ ] `aws --version` muestra v2
- [ ] `aws sts get-caller-identity` muestra tu cuenta

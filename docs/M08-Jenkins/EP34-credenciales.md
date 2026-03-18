# EP 34: Credenciales — Docker Hub y GitHub Token

**Tipo:** CONFIGURACION
## Objetivo
Agregar credenciales de Docker Hub y GitHub en Jenkins para usarlas en los pipelines.

## Pasos Detallados

Ir a: **Manage Jenkins → Credentials → System → Global credentials → Add Credentials**

### 1. Docker Hub
- Kind: **Username with password**
- ID: `docker-hub-credentials`
- Username: tu usuario de Docker Hub
- Password: tu contraseña o Access Token
- Description: "Docker Hub"

### 2. GitHub Token
- En GitHub: Settings → Developer settings → Personal access tokens → Generate new token
  - Scopes: `repo` (full control)
- En Jenkins:
  - Kind: **Secret text**
  - ID: `github-token`
  - Secret: (pegar el token)
  - Description: "GitHub Token"

### 3. AWS Credentials (para Terraform destroy)
- Kind: **AWS Credentials**
- ID: `aws-credentials`
- Access Key ID: (tu access key)
- Secret Access Key: (tu secret key)

## Verificación
- [ ] 3 credenciales aparecen en la lista global
- [ ] Los IDs coinciden con los del Jenkinsfile

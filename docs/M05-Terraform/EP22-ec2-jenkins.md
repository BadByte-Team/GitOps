# EP 22: Crear la EC2 Jenkins con Terraform

**Tipo:** PRACTICA
## Objetivo
Usar el código Terraform del proyecto para provisionar la instancia Jenkins con user-data.

## Prerequisitos
- Backend S3 configurado (EP17, EP20)
- AWS CLI configurado (EP15)

## Pasos Detallados

### 1. Crear Key Pair en AWS (si no existe)
```bash
aws ec2 create-key-pair --key-name jenkins-key --query 'KeyMaterial' --output text > jenkins-key.pem
chmod 400 jenkins-key.pem
```

### 2. Aplicar Terraform
```bash
cd infrastructure/terraform/jenkins-ec2

terraform init
terraform plan -var="key_name=jenkins-key"
terraform apply -var="key_name=jenkins-key"
```

### 3. Esperar la Instalación (~5 min)
El user-data script instala: Java, Jenkins, Docker, kubectl, Terraform, AWS CLI.

### 4. Acceder a Jenkins
```bash
# Obtener IP
terraform output jenkins_url

# Conectar por SSH para ver la contraseña
ssh -i jenkins-key.pem ubuntu@$(terraform output -raw jenkins_public_ip)
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
```

### 5. Abrir en el Navegador
- Ir a `http://IP:8080`
- Pegar la contraseña inicial
- Instalar plugins sugeridos
- Crear usuario admin

## Archivos Involucrados
- `infrastructure/terraform/jenkins-ec2/main.tf`
- `infrastructure/terraform/jenkins-ec2/variables.tf`
- `infrastructure/terraform/jenkins-ec2/outputs.tf`
- `infrastructure/scripts/install-jenkins.sh`

## Verificación
- [ ] La EC2 está running
- [ ] Jenkins responde en `http://IP:8080`
- [ ] Docker funciona en la EC2: `docker ps`

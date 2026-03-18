# EP 16: EC2 — Lanzar Instancia y Conectar por SSH

**Tipo:** PRACTICA
## Objetivo
Crear un Security Group, lanzar una instancia EC2 Ubuntu y conectar por SSH.

## Pasos Detallados

### 1. Crear Key Pair
```bash
aws ec2 create-key-pair --key-name curso-gitops-key --query 'KeyMaterial' --output text > curso-gitops-key.pem
chmod 400 curso-gitops-key.pem
```

### 2. Crear Security Group
```bash
# Obtener VPC ID por defecto
VPC_ID=$(aws ec2 describe-vpcs --filters "Name=isDefault,Values=true" --query 'Vpcs[0].VpcId' --output text)

# Crear SG
SG_ID=$(aws ec2 create-security-group --group-name curso-gitops-sg --description "SG para curso GitOps" --vpc-id $VPC_ID --query 'GroupId' --output text)

# Abrir SSH
aws ec2 authorize-security-group-ingress --group-id $SG_ID --protocol tcp --port 22 --cidr 0.0.0.0/0

echo "Security Group: $SG_ID"
```

### 3. Lanzar EC2
```bash
INSTANCE_ID=$(aws ec2 run-instances   --image-id ami-0c7217cdde317cfec   --instance-type t2.micro   --key-name curso-gitops-key   --security-group-ids $SG_ID   --query 'Instances[0].InstanceId' --output text)

echo "Instance ID: $INSTANCE_ID"

# Esperar a que esté running
aws ec2 wait instance-running --instance-ids $INSTANCE_ID

# Obtener IP pública
PUBLIC_IP=$(aws ec2 describe-instances --instance-ids $INSTANCE_ID --query 'Reservations[0].Instances[0].PublicIpAddress' --output text)
echo "IP Publica: $PUBLIC_IP"
```

### 4. Conectar por SSH
```bash
ssh -i curso-gitops-key.pem ubuntu@$PUBLIC_IP
```

### 5. Terminar la Instancia (al finalizar)
```bash
aws ec2 terminate-instances --instance-ids $INSTANCE_ID
```

## Verificación
- [ ] La instancia está en estado "running"
- [ ] Puedes conectar por SSH
- [ ] Terminas la instancia al finalizar para evitar costos

# EP 43: Instalar SonarQube en EC2 con Docker

**Tipo:** INSTALACION / CONFIGURACION
## Objetivo
Levantar SonarQube en la EC2 Jenkins usando Docker Compose.

## Pasos Detallados

### 1. Copiar el Compose file
```bash
# En la EC2 Jenkins
mkdir -p ~/sonarqube
cp infrastructure/docker/sonarqube/docker-compose.yml ~/sonarqube/
```

### 2. Configurar requisitos del sistema
```bash
# SonarQube requiere más memoria virtual
sudo sysctl -w vm.max_map_count=262144
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### 3. Levantar SonarQube
```bash
cd ~/sonarqube
docker compose up -d

# Verificar
docker compose ps
# Esperar ~2 minutos a que arranque
```

### 4. Acceder
- URL: `http://IP-DE-EC2:9000`
- Usuario: `admin`
- Contraseña: `admin` → te pide cambiarla

### 5. Cambiar Contraseña
- Login → cambiar a una contraseña segura

## Archivos Involucrados
- `infrastructure/docker/sonarqube/docker-compose.yml`

## Verificación
- [ ] SonarQube responde en `http://IP:9000`
- [ ] Puedes loguearte como admin
- [ ] La contraseña fue cambiada

# EP 31: Instalación de Jenkins en EC2 Ubuntu

**Tipo:** INSTALACION
## Objetivo
Instalar Jenkins en una instancia EC2, abrir el puerto 8080 y obtener la contraseña inicial.

## Prerequisitos
- EC2 Jenkins creada con Terraform (EP22) o manualmente

## Pasos Detallados

### Si usaste Terraform (EP22)
Jenkins ya está instalado via user-data. Solo necesitas:
```bash
ssh -i jenkins-key.pem ubuntu@$(terraform output -raw jenkins_public_ip)
sudo cat /var/lib/jenkins/secrets/initialAdminPassword
```

### Si instalas manualmente
```bash
# Usar el script del proyecto
chmod +x infrastructure/scripts/install-jenkins.sh
sudo ./infrastructure/scripts/install-jenkins.sh
```

### Configuración Inicial
1. Abrir `http://IP:8080` en el navegador
2. Pegar la contraseña inicial
3. "Install suggested plugins" → esperar
4. Crear usuario administrador
5. Confirmar la URL de Jenkins

## Archivos Involucrados
- `infrastructure/scripts/install-jenkins.sh`

## Verificación
- [ ] Jenkins responde en el puerto 8080
- [ ] Puedes loguearte con el usuario admin
- [ ] Docker funciona en la máquina: `docker ps`

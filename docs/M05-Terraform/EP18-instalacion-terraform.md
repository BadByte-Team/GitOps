# EP 18: Instalación de Terraform en Local y en EC2

**Tipo:** INSTALACION
## Objetivo
Instalar Terraform desde el repositorio de HashiCorp y verificar la versión.

## Pasos Detallados

### Instalar en Ubuntu (local o EC2)
```bash
# Agregar repositorio
sudo apt-get install -y gnupg software-properties-common
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list

# Instalar
sudo apt-get update && sudo apt-get install -y terraform

# Verificar
terraform --version
```

### Agregar al PATH (si es necesario)
```bash
echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bashrc
source ~/.bashrc
```

### Autocompletado
```bash
terraform -install-autocomplete
source ~/.bashrc
```

## Verificación
- [ ] `terraform --version` muestra la versión
- [ ] El autocompletado funciona (escribir `terraform` + TAB)

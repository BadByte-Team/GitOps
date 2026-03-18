#!/bin/bash
# ══════════════════════════════════════
#  Instalar Trivy en Ubuntu
# ══════════════════════════════════════
set -euo pipefail

echo ">>> Instalando Trivy..."

sudo apt-get install -y wget apt-transport-https gnupg lsb-release

wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor | sudo tee /usr/share/keyrings/trivy.gpg > /dev/null

echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/trivy.list

sudo apt-get update -y
sudo apt-get install -y trivy

echo ">>> Trivy instalado: $(trivy --version)"
echo ""
echo "Uso basico:"
echo "  trivy image nginx:latest"
echo "  trivy image --severity HIGH,CRITICAL tu-imagen:tag"

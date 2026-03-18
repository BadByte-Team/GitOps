#!/bin/bash
# ══════════════════════════════════════════════
#  User-data: Jenkins + Docker + kubectl + Terraform
#  Para EC2 Ubuntu 22.04+
# ══════════════════════════════════════════════
set -euo pipefail

export DEBIAN_FRONTEND=noninteractive

echo ">>> Actualizando sistema..."
apt-get update -y && apt-get upgrade -y

# ── Java 17 (requerido por Jenkins) ──
echo ">>> Instalando Java 17..."
apt-get install -y fontconfig openjdk-17-jre

# ── Jenkins ──
echo ">>> Instalando Jenkins..."
curl -fsSL https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key | tee /usr/share/keyrings/jenkins-keyring.asc > /dev/null
echo "deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] https://pkg.jenkins.io/debian-stable binary/" | tee /etc/apt/sources.list.d/jenkins.list > /dev/null
apt-get update -y
apt-get install -y jenkins
systemctl enable jenkins
systemctl start jenkins

# ── Docker ──
echo ">>> Instalando Docker..."
apt-get install -y ca-certificates curl gnupg lsb-release
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Agregar jenkins y ubuntu al grupo docker
usermod -aG docker jenkins
usermod -aG docker ubuntu

# ── kubectl ──
echo ">>> Instalando kubectl..."
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl

# ── Terraform ──
echo ">>> Instalando Terraform..."
apt-get install -y gnupg software-properties-common
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list
apt-get update -y
apt-get install -y terraform

# ── AWS CLI v2 ──
echo ">>> Instalando AWS CLI v2..."
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
apt-get install -y unzip
unzip -q awscliv2.zip
./aws/install
rm -rf aws awscliv2.zip

# ── jq (util para scripts) ──
apt-get install -y jq

echo "══════════════════════════════════════════════"
echo "  ✅ Instalacion completa"
echo "  Jenkins:   http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):8080"
echo "  Password:  sudo cat /var/lib/jenkins/secrets/initialAdminPassword"
echo "══════════════════════════════════════════════"

#!/bin/bash
# ══════════════════════════════════════════════
#  Limpieza completa de recursos AWS
#  Ejecutar desde tu maquina local
# ══════════════════════════════════════════════
set -euo pipefail

echo "⚠️  Este script eliminará TODOS los recursos del curso GitOps en AWS"
echo "    - ArgoCD Application"
echo "    - Cluster EKS (via Terraform)"
echo "    - EC2 Jenkins (via Terraform)"
echo ""
read -p "¿Continuar? (yes/no): " CONFIRM
if [ "$CONFIRM" != "yes" ]; then
    echo "Cancelado."
    exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INFRA_DIR="$(dirname "$SCRIPT_DIR")"

# 1. Eliminar app en ArgoCD
echo ""
echo ">>> [1/4] Eliminando aplicacion en ArgoCD..."
kubectl delete application curso-gitops -n argocd 2>/dev/null || echo "   (no encontrada o ya eliminada)"
kubectl delete namespace curso-gitops 2>/dev/null || echo "   (namespace ya eliminado)"

# 2. Desinstalar ArgoCD
echo ""
echo ">>> [2/4] Desinstalando ArgoCD..."
kubectl delete namespace argocd 2>/dev/null || echo "   (ArgoCD ya eliminado)"

# 3. Destruir EKS
echo ""
echo ">>> [3/4] Destruyendo cluster EKS..."
cd "$INFRA_DIR/terraform/eks-cluster"
terraform init -input=false
terraform destroy -auto-approve

# 4. Destruir Jenkins EC2
echo ""
echo ">>> [4/4] Destruyendo Jenkins EC2..."
cd "$INFRA_DIR/terraform/jenkins-ec2"
terraform init -input=false
terraform destroy -auto-approve

echo ""
echo "══════════════════════════════════════════════"
echo "  ✅ Todos los recursos han sido eliminados"
echo "  💡 Revisa la consola de AWS para confirmar"
echo "  💡 El bucket S3 del state NO se elimina (prevent_destroy)"
echo "══════════════════════════════════════════════"

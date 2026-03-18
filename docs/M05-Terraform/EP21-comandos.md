# EP 21: Comandos — init, validate, plan, apply, destroy

**Tipo:** PRACTICA
## Objetivo
Dominar el ciclo completo de Terraform y entender cada comando.

## Ciclo de Vida de Terraform

```
init → validate → plan → apply → (cambios) → plan → apply → destroy
```

### Comandos
| Comando | Qué hace |
|---|---|
| `terraform init` | Descarga providers y configura backend |
| `terraform validate` | Verifica que la sintaxis es correcta |
| `terraform fmt` | Formatea archivos `.tf` |
| `terraform plan` | Muestra qué cambios va a hacer (dry-run) |
| `terraform apply` | Ejecuta los cambios |
| `terraform destroy` | Elimina todos los recursos |
| `terraform state list` | Lista recursos en el state |
| `terraform output` | Muestra los outputs |

### Buenas Prácticas
1. **Siempre** ejecutar `plan` antes de `apply`
2. **Leer el plan completamente** antes de confirmar
3. Usar `-auto-approve` solo en CI/CD, nunca manualmente
4. Formatear con `terraform fmt` antes de commitear

### Ejemplo Completo
```bash
cd infrastructure/terraform/jenkins-ec2

terraform init
terraform validate
terraform fmt -check
terraform plan -out=tfplan
terraform apply tfplan

# Al finalizar
terraform destroy
```

## Verificación
- [ ] Puedes ejecutar el ciclo completo sin errores
- [ ] Entiendes la diferencia entre plan y apply

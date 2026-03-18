# EP 13: Crear Cuenta AWS y Activar Free Tier

**Tipo:** CONFIGURACION
## Objetivo
Crear una cuenta de AWS, activar el Free Tier y proteger la cuenta root con MFA.

## Pasos Detallados

### 1. Registro en AWS
1. Ir a https://aws.amazon.com → "Create an AWS Account"
2. Ingresar email y nombre de cuenta
3. Agregar método de pago (tarjeta — no cobra si usas Free Tier)
4. Verificar identidad por teléfono
5. Seleccionar plan **Basic Support (Free)**

### 2. Activar MFA en Root
1. Ir a IAM → Dashboard → "Activate MFA"
2. Elegir "Virtual MFA device"
3. Escanear QR con Google Authenticator o Authy
4. Ingresar dos códigos consecutivos

### 3. Entender el Panel de Costos
- **AWS Cost Explorer**: ver gastos en tiempo real
- **Budgets**: crear alerta cuando passes $5 USD
  - Billing → Budgets → Create Budget → Cost Budget → $5

## Verificación
- [ ] Puedes acceder a la consola de AWS
- [ ] MFA está activado en la cuenta root
- [ ] Tienes un budget de alerta configurado

## Notas
> ⚠️ **NUNCA** uses la cuenta root para el trabajo diario — crear un usuario IAM en el EP14
> ⚠️ Configura una alerta de billing para evitar cargos inesperados

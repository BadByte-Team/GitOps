# EP 14: IAM — Usuarios, Roles y Access Keys

**Tipo:** CONFIGURACION
## Objetivo
Crear un usuario IAM con permisos de administrador y generar Access Keys para usar con CLI.

## Pasos Detallados

### 1. Crear Usuario IAM
1. IAM → Users → Create User
2. Nombre: `admin-curso`
3. Enable console access: ✅
4. Auto-generated password → copiarla

### 2. Agregar Permisos
1. Attach policies directly
2. Buscar y seleccionar `AdministratorAccess`
3. Next → Create User

### 3. Generar Access Key
1. IAM → Users → `admin-curso` → Security Credentials
2. Create Access Key
3. Seleccionar "Command Line Interface (CLI)"
4. **Copiar Access Key ID y Secret Access Key** (solo se muestra una vez)

### 4. Guardar Credenciales de forma Segura
```bash
# Guardar en un archivo local seguro
mkdir -p ~/.aws
cat > ~/.aws/credentials << EOF
[default]
aws_access_key_id = TU_ACCESS_KEY_ID
aws_secret_access_key = TU_SECRET_ACCESS_KEY
EOF
chmod 600 ~/.aws/credentials
```

## Verificación
- [ ] Puedes hacer login en la consola con el usuario IAM
- [ ] Tienes el Access Key ID y Secret guardados
- [ ] Nunca commiteas credenciales en Git

## Notas
> ⚠️ NUNCA commitear Access Keys en un repositorio
> El `.gitignore` del proyecto ya excluye `*.key` y `.env`

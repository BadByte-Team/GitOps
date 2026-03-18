# EP 08: Instalación de Docker en Ubuntu (EC2)

**Tipo:** INSTALACIÓN  
**Duración estimada:** 10–12 min  
**Dificultad:** ⭐ (Básico)

---

## 🎯 Objetivo
Instalar Docker Engine en Ubuntu usando el repositorio oficial de Docker (no el paquete de Ubuntu que está desactualizado) y configurar permisos para ejecutar sin `sudo`.

---

## 📋 Prerequisitos
- Ubuntu 20.04+ (local o EC2)
- Acceso con permisos de superusuario

---

## 📺 Paso a Paso para el Video

### Paso 1 — Desinstalar Versiones Anteriores

**Ejecutar en:** Terminal de la máquina donde instalarás Docker (tu local o EC2 por SSH)

```bash
# Remover instalaciones previas
sudo apt-get remove -y docker docker-engine docker.io containerd runc 2>/dev/null || true
```

---

### Paso 2 — Instalar Docker desde el Repositorio Oficial

> 📌 **Documentación oficial:** https://docs.docker.com/engine/install/ubuntu/

```bash
# 1. Actualizar paquetes e instalar dependencias
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# 2. Agregar la clave GPG oficial de Docker
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# 3. Agregar el repositorio oficial
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 4. Instalar Docker Engine + Compose plugin
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

---

### Paso 3 — Configurar Docker sin `sudo` (Post-Install)

```bash
# Agregar tu usuario al grupo docker
sudo usermod -aG docker $USER

# Activar el grupo para la sesión actual (evita cerrar sesión)
newgrp docker
```

> ⚠️ **Si estás en EC2:** Después de `usermod` necesitas cerrar la sesión SSH y reconectar para que tome efecto.

---

### Paso 4 — Verificar

```bash
# Sin sudo — debe funcionar
docker run hello-world

# Versiones
docker --version
# Docker version 27.x.x

docker compose version
# Docker Compose version v2.x.x
```

**Resultado esperado de `docker run hello-world`:**
```
Hello from Docker!
This message shows that your installation appears to be working correctly.
```

---

### Paso 5 — Comandos de Verificación Adicionales

```bash
# Estado del servicio
sudo systemctl status docker
# Debe decir: active (running)

# Información completa
docker info
```

---

## 🗂️ Archivos del Proyecto Relacionados
| Archivo | Relevancia |
|---|---|
| `infrastructure/scripts/install-jenkins.sh` | Ya incluye la instalación completa de Docker (líneas 36-49) |
| `curso-gitops/Dockerfile` | Lo usaremos en el EP10 para construir nuestra imagen |

---

## ✅ Checklist de Verificación
- [ ] `docker run hello-world` funciona **sin sudo**
- [ ] `docker --version` muestra versión 27+
- [ ] `docker compose version` muestra v2+
- [ ] El servicio está activo: `sudo systemctl status docker`

---

## 🔧 Troubleshooting

| Problema | Solución |
|---|---|
| `permission denied` al ejecutar sin sudo | Ejecutar `newgrp docker` o cerrar/abrir sesión |
| `Cannot connect to Docker daemon` | `sudo systemctl start docker` |
| El paquete `docker-ce` no se encuentra | Verificar que el repositorio de Docker se agregó correctamente |
| Conflicto con `docker.io` de Ubuntu | Desinstalar primero: `sudo apt remove docker.io` |

---

## 📌 Notas para el Video
- Mostrar la diferencia entre `docker.io` (paquete viejo de Ubuntu) y `docker-ce` (oficial de Docker)
- Ejecutar cada comando paso a paso
- Mostrar que `docker run hello-world` falla con permisos antes del `usermod`
- Después del `newgrp`, mostrar que funciona sin sudo
- Terminar con `docker --version` y `docker compose version`

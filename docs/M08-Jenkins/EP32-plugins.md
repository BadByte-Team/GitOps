# EP 32: Instalación de Plugins Indispensables

**Tipo:** CONFIGURACION
## Objetivo
Instalar los plugins necesarios para el pipeline CI/CD.

## Plugins a Instalar

Ir a: **Manage Jenkins → Plugins → Available plugins**

| Plugin | Para qué |
|---|---|
| Docker Pipeline | Build y push de imágenes Docker |
| NodeJS Plugin | Ejecutar herramientas Node.js |
| Eclipse Temurin installer | JDK administrado por Jenkins |
| Pipeline: AWS Steps | Interacción con AWS desde pipelines |
| SonarQube Scanner | Análisis de código estático |
| OWASP Dependency-Check | Escaneo de dependencias |

### Instalación
1. Manage Jenkins → Plugins → Available plugins
2. Buscar cada plugin por nombre
3. Marcar checkbox → "Install without restart"
4. Al final: reiniciar Jenkins

## Verificación
- [ ] Todos los plugins aparecen en "Installed plugins"
- [ ] Jenkins se reinició correctamente

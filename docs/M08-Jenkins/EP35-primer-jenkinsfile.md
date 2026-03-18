# EP 35: Primer Jenkinsfile — Pipeline Declarativo

**Tipo:** PRACTICA
## Objetivo
Entender la estructura de un pipeline declarativo y crear tu primer Jenkinsfile.

## Estructura de un Jenkinsfile

```groovy
pipeline {
    agent any                    // Dónde ejecutar

    tools {                      // Herramientas
        jdk 'jdk17'
    }

    environment {                // Variables de entorno
        MY_VAR = "valor"
    }

    stages {                     // Etapas del pipeline
        stage('Build') {
            steps {
                sh 'echo "Building..."'
            }
        }

        stage('Test') {
            steps {
                sh 'echo "Testing..."'
            }
        }
    }

    post {                       // Acciones post-ejecución
        success { echo '✅ OK' }
        failure { echo '❌ FAIL' }
        always  { cleanWs() }
    }
}
```

### Crear Pipeline en Jenkins
1. New Item → "Pipeline" → nombre: `test-pipeline`
2. Pipeline → Definition: "Pipeline script from SCM"
3. SCM: Git → URL del repo
4. Script Path: `Jenkinsfile`
5. Build Now

## Archivos Involucrados
- `infrastructure/jenkins/Jenkinsfile`

## Verificación
- [ ] El pipeline ejecuta y muestra los stages en verde
- [ ] Entiendes la estructura: agent, stages, steps, post

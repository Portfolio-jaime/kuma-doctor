# Kuma Doctor 🩺

`kuma-doctor` es una herramienta de línea de comandos (CLI) escrita en Go para realizar un diagnóstico completo y robusto de un service mesh de [Kuma](https://kuma.io/) corriendo sobre Kubernetes (EKS).

Permite a los operadores de plataforma y desarrolladores validar rápidamente la salud, seguridad, resiliencia y configuración general del mesh de forma interactiva o a través de comandos para automatización y CI/CD.

## Diagrama de Arquitectura

El flujo de la herramienta está diseñado para ser modular y extensible, separando la interfaz de usuario, la lógica de análisis y la generación de reportes.

```
+------------------+
|      Usuario     |
+------------------+
         |
         v
+------------------+      +-----------------------------------------+
|   kuma-doctor    |----->|         Interfaz (TUI / Comandos)         |
|     (Binario)    |      |         (Cobra & Survey)                  |
+------------------+      +-----------------------------------------+
                                     |
                                     v
                  +-----------------------------------------+
                  |       Lógica de Análisis (Paquetes)       |
                  |  - Dataplanes     - Resiliencia         |
                  |  - Políticas      - Observabilidad      |
                  |  - mTLS           - ...                 |
                  +-----------------------------------------+
                                     |
                                     v
                  +-----------------------------------------+
                  |     Cliente de Kubernetes (client-go)     |
                  +-----------------------------------------+
                                     |
                                     v
+------------------+      +-----------------------------------------+
| Cluster EKS      |<---->|          Servidor API de K8s          |
+------------------+      +-----------------------------------------+
         ^                                   | (Resultados)
         |                                   |
         |                            +------------------+
         +----------------------------| Report Generator |
                                      | (txt, md, json)  |
                                      +------------------+
```

## Características Principales

- **Menú Interactivo:** Una interfaz amigable para guiar al usuario a través de los diferentes análisis.
- **Modo No Interactivo:** Subcomandos para cada análisis, perfectos para scripts y pipelines de CI/CD.
- **Reportes Multi-formato:** Genera reportes en formato de texto plano (para la consola), Markdown (para documentación) y JSON (para integración con otras herramientas).
- **Análisis Comprensivo:**
  - ✅ **Resumen General:** Vista de pájaro del estado del mesh.
  - ✅ **Estado de Dataplanes:** Verifica la conectividad de cada proxy del mesh.
  - ✅ **Políticas de Tráfico:** Detecta servicios sin protección y reglas demasiado permisivas.
  - ✅ **Seguridad mTLS:** Valida que el cifrado de tráfico esté activado y correctamente configurado.
  - ✅ **Políticas de Resiliencia:** Busca brechas en la configuración de reintentos, timeouts y circuit breakers.
  - ✅ **Políticas de Observabilidad:** Comprueba si las políticas de logging, métricas y tracing están en su lugar.

## Instalación

Para usar `kuma-doctor` como un comando global en tu sistema, asegúrate de tener Go instalado y configurado correctamente.

1.  **Clona el repositorio (si aplica):**
    ```bash
    git clone <tu-url-de-repositorio>
    cd kuma-doctor
    ```

2.  **Instala la herramienta:**
    El comando `go install` compilará el código y moverá el ejecutable a tu `$GOPATH/bin`, que debería estar en tu `PATH`.

    ```bash
    go install .
    ```

## Uso

Una vez instalado, `kuma-doctor` estará disponible desde cualquier lugar en tu terminal.

#### Modo Interactivo
Simplemente ejecuta el comando sin argumentos para lanzar el menú.
```bash
kuma-doctor
```
Navega con las flechas y presiona Enter para seleccionar un análisis.

#### Modo No Interactivo (Comandos)
Usa el subcomando `check` para ejecutar análisis específicos.

```bash
# Revisar el estado de todos los dataplanes
kuma-doctor check dataplanes

# Revisar la configuración de mTLS y guardar un reporte en Markdown
kuma-doctor check mtls --output md --file reporte-mtls.md

# Revisar las políticas de resiliencia y generar un reporte en JSON
kuma-doctor check resilience -o json -f reporte-resiliencia.json
```

#### Opciones Globales
- `-o, --output`: Especifica el formato de salida (`txt`, `md`, `json`). Por defecto es `txt`.
- `-f, --file`: Especifica un archivo para guardar el reporte. Si se omite, se imprime en la consola.

## Estructura del Proyecto

El código está organizado para facilitar su mantenimiento y extensión:

- `/cmd`: Contiene la definición de los comandos de la CLI (usando Cobra).
- `/internal`: Código interno de la aplicación no destinado a ser importado por otros proyectos.
  - `/kubernetes`: Lógica para la conexión con el clúster.
  - `/report`: Generadores para los diferentes formatos de reporte.
  - `/tui`: Lógica para el menú interactivo (usando Survey).
- `/pkg`: Paquetes con la lógica de negocio principal que podrían ser reutilizados.
  - `/analysis`: Contiene toda la lógica de validación para cada aspecto del mesh.

## Cómo Contribuir

¡Las contribuciones son bienvenidas! Por favor, abre un *issue* para discutir tu idea o envía un *pull request* con tu mejora.

## Licencia

Este proyecto está bajo la Licencia MIT.
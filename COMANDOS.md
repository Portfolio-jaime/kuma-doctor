# Referencia de Comandos de Kuma-Doctor

Este documento proporciona una guía detallada de todos los comandos disponibles en la CLI `kuma-doctor`, sus objetivos y ejemplos de uso.

`kuma-doctor` puede ser operado de dos maneras principales:
1.  **Modo Interactivo:** Una interfaz amigable que te guía a través de los análisis disponibles. Se activa ejecutando `kuma-doctor` sin argumentos.
2.  **Modo No Interactivo:** Comandos y subcomandos específicos que pueden ser usados en scripts para automatización, con salidas en diferentes formatos.

## Opciones Globales

Estas opciones (o *flags*) están disponibles en casi todos los comandos:

- `-o, --output <formato>`: Especifica el formato de salida.
  - `txt`: Texto plano con colores, optimizado para la consola (por defecto).
  - `md`: Markdown, ideal para generar documentación.
  - `json`: Formato estructurado, perfecto para integración con otras herramientas.
- `-f, --file <ruta>`: Guarda el reporte en el archivo especificado en lugar de mostrarlo en la consola.
- `-h, --help`: Muestra un mensaje de ayuda para cualquier comando o subcomando.

---

## Comandos Principales

### `kuma-doctor`

Este es el comando raíz. Si se ejecuta sin subcomandos, lanza el menú interactivo.

- **Objetivo:** Proveer una experiencia de usuario sencilla y guiada para realizar diagnósticos sin necesidad de memorizar subcomandos.
- **Funcionalidad:** Inicia una sesión con un menú que lista todas las categorías de análisis disponibles.
- **Ejemplo de Uso:**
  ```bash
  # Iniciar el menú interactivo
  kuma-doctor
  ```

### `kuma-doctor report`

Ejecuta todos los análisis disponibles de forma secuencial y los consolida en un único reporte.

- **Objetivo:** Obtener un diagnóstico completo y exhaustivo del estado del mesh con un solo comando. Ideal para revisiones periódicas o para obtener una "fotografía" completa de la salud del sistema.
- **Funcionalidad:** Llama internamente a cada una de las funciones de análisis (`Summary`, `Dataplanes`, `MTP`, `mTLS`, `Resilience`, `Observability`) y une sus resultados en un solo documento.
- **Ejemplos de Uso:**
  ```bash
  # Generar el reporte completo en la consola
  kuma-doctor report

  # Generar un reporte completo en formato Markdown y guardarlo en un archivo
  kuma-doctor report --output md --file REPORTE_KUMA_COMPLETO.md
  ```

### `kuma-doctor check`

Este es un comando "padre" que agrupa todos los análisis individuales para su ejecución no interactiva. No hace nada por sí solo, pero contiene los siguientes subcomandos.

- **Ejemplo de Uso:**
  ```bash
  # Ver todos los chequeos disponibles
  kuma-doctor check --help
  ```

---

## Subcomandos de `check`

### `check dataplanes`

- **Objetivo:** Verificar la salud y conectividad fundamental de cada proxy de servicio (`kuma-dp`) en el mesh. Es el chequeo más básico e importante.
- **Funcionalidades Clave:**
    - Itera sobre todos los recursos `Dataplane`.
    - Analiza el campo `health: { ready: true }` dentro de cada `inbound` en la especificación del networking.
    - Clasifica cada Dataplane como `Online`, `Offline`, `Degraded` o `Info` (si no tiene inbounds).
- **Ejemplos de Uso:**
  ```bash
  # Ejecutar el análisis y mostrar en consola
  kuma-doctor check dataplanes

  # Guardar el resultado en un archivo JSON
  kuma-doctor check dataplanes -o json -f dataplanes.json
  ```

### `check traffic-permissions`

- **Alias:** `mtp`
- **Objetivo:** Auditar la configuración de seguridad del tráfico, encontrando posibles servicios aislados o reglas demasiado permisivas.
- **Funcionalidades Clave:**
    - Obtiene una lista de todos los servicios (`kuma.io/service`) del clúster.
    - Revisa todas las políticas `MeshTrafficPermission`.
    - **Alerta (🚨)** si encuentra servicios que no son destino de ninguna política, lo que podría dejarlos sin tráfico entrante.
    - **Informa (✅)** si una política permite tráfico desde cualquier origen (`from: '*'`), para revisión manual.
- **Ejemplos de Uso:**
  ```bash
  # Usar el alias 'mtp' para un análisis rápido
  kuma-doctor check mtp
  ```

### `check mtls`

- **Objetivo:** Verificar que la encriptación de tráfico mTLS, una de las principales características de seguridad de un service mesh, esté correctamente activada y forzada.
- **Funcionalidades Clave:**
    - Comprueba si mTLS está activado en el recurso `Mesh` (`spec.mtls.enabledBackend`).
    - Valida que el backend de mTLS activado esté correctamente definido en la lista de `backends`.
    - **Advierte (⚠️)** si mTLS está activo pero existen políticas `MeshTrafficPermission` que usan la acción `Allow` en lugar de `AllowWithMTLS`, creando potenciales brechas de seguridad.
- **Ejemplos de Uso:**
  ```bash
  # Ejecutar la auditoría de mTLS
  kuma-doctor check mtls
  ```

### `check resilience`

- **Objetivo:** Asegurar que las aplicaciones dentro del mesh sean robustas y puedan soportar fallos de red o sobrecargas temporales.
- **Funcionalidades Clave:**
    - Revisa la cobertura de las políticas `MeshRetry`, `MeshTimeout` y `MeshCircuitBreaker`.
    - **Advierte (⚠️)** sobre cada servicio que no esté cubierto por alguno de estos tres tipos de políticas, ya que podría no recuperarse de errores transitorios o ser vulnerable a fallas en cascada.
- **Ejemplos de Uso:**
  ```bash
  # Revisar qué servicios carecen de políticas de resiliencia
  kuma-doctor check resilience
  ```

### `check observability`

- **Alias:** `obs`
- **Objetivo:** Confirmar que el mesh está configurado para ser observable, lo cual es vital para la monitorización, la depuración y el entendimiento del comportamiento del sistema.
- **Funcionalidades Clave:**
    - Verifica la existencia de políticas a nivel de mesh para `MeshLog` (logs de acceso), `MeshMetric` (métricas para Prometheus) y `MeshTrace` (tracing distribuido).
    - **Advierte (⚠️)** si falta alguna de estas políticas globales, ya que implicaría una pérdida de visibilidad en esa área.
- **Ejemplos de Uso:**
  ```bash
  # Usar el alias 'obs' para revisar la configuración de telemetría
  kuma-doctor check obs
  ```
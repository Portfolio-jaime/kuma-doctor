# Kuma Doctor 🩺

`kuma-doctor` es una herramienta de línea de comandos (CLI) escrita en Go para realizar un diagnóstico completo y robusto de un service mesh de [Kuma](https://kuma.io/) corriendo sobre Kubernetes (EKS).

Permite a los operadores de plataforma y desarrolladores validar rápidamente la salud, seguridad, resiliencia y configuración general del mesh de forma interactiva o a través de comandos para automatización y CI/CD.

## Tecnologías Utilizadas

Este proyecto fue construido utilizando las siguientes librerías y frameworks principales:

- **[Go](https://golang.org/):** El lenguaje de programación base.
- **[Cobra](https://github.com/spf13/cobra):** Framework líder para la creación de aplicaciones CLI robustas.
- **[Survey](https://github.com/AlecAivazis/survey):** Para la creación de los menús interactivos y amigables.
- **[Client-Go](https://github.com/kubernetes/client-go):** La librería oficial de Kubernetes para interactuar con la API del clúster.
- **[Fatih/Color](https://github.com/fatih/color):** Para añadir colores distintivos a la salida en la terminal, mejorando la legibilidad.

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
  - ✅ **Reporte Completo:** Ejecuta todos los análisis anteriores de una sola vez y genera un informe consolidado.

## Instalación

Para usar `kuma-doctor` como un comando global en tu sistema, sigue estos pasos.

### Requisitos Previos
1.  **Go (versión 1.18+):** Necesitas tener Go instalado.
2.  **Kubeconfig:** Tu `kubeconfig` debe estar configurado para apuntar al clúster de EKS que deseas analizar (ej. vía `aws eks update-kubeconfig ...`).

### 1. Instrucciones de Instalación

La forma recomendada es usar `go install`, que compila el binario y lo coloca en la ruta correcta automáticamente.

```bash
# Opcional: Clona el repositorio si no lo tienes
# git clone <tu-url-de-repositorio>
# cd kuma-doctor

# Instala la herramienta con go install
go install .
```
Este comando compilará el proyecto e instalará el ejecutable `kuma-doctor` en tu directorio de binarios de Go (usualmente `$HOME/go/bin`).

### 2. Solución de Problemas: `command not found`

Después de ejecutar `go install`, si abres una nueva terminal y recibes el error `zsh: command not found: kuma-doctor` (o similar en bash), significa que el directorio de binarios de Go no está en tu `PATH` (la lista de directorios donde tu terminal busca programas).

Sigue estos pasos para solucionarlo permanentemente:

**Paso 1: Verifica dónde instaló Go el programa**

Confirma que el ejecutable existe en la carpeta de binarios de Go.

```bash
ls -l $(go env GOPATH)/bin
```
Deberías ver `kuma-doctor` en la lista.

**Paso 2: Revisa tu `PATH` actual**

Echa un vistazo a las rutas que tu terminal ya conoce.
```bash
echo $PATH
```
Lo más seguro es que la ruta `.../go/bin` no aparezca aquí.

**Paso 3: Añade la carpeta de Go a tu `PATH`**

Ejecuta el comando correspondiente a tu terminal para añadir la ruta de forma permanente a tu archivo de configuración.

* **Para Zsh (usada por defecto en macOS moderno):**
    ```bash
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
    ```
* **Para Bash (usada por defecto en la mayoría de sistemas Linux):**
    ```bash
    echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
    ```

**Paso 4: Aplica los Cambios**

La configuración se carga al iniciar una nueva terminal. Tienes dos opciones:
- **(Recomendado)** Cierra por completo tu ventana de terminal y abre una nueva.
- **(Alternativa)** En la misma terminal, ejecuta `source ~/.zshrc` o `source ~/.bashrc` para recargar la configuración.

**Paso 5: ¡Verificación Final!**

En la nueva terminal, comprueba que el sistema ahora sí encuentra el comando:
```bash
which kuma-doctor
```
Debería devolverte la ruta completa (ej. `/home/usuario/go/bin/kuma-doctor`). ¡Listo!

## Uso

Una vez instalado, `kuma-doctor` estará disponible desde cualquier lugar en tu terminal.

#### Modo Interactivo
Simplemente ejecuta el comando sin argumentos para lanzar el menú.
```bash
kuma-doctor
```

#### Modo No Interactivo (Comandos)
Usa los subcomandos `check` y `report` para análisis específicos.

```bash
# Generar el reporte completo y guardarlo en un archivo Markdown
kuma-doctor report --output md --file REPORTE_COMPLETO.md

# Revisar solo la configuración de mTLS
kuma-doctor check mtls

# Ver todos los comandos disponibles
kuma-doctor --help
kuma-doctor check --help
```

---
**Para una guía completa y detallada de todos los comandos y sus funcionalidades, por favor consulta la [Referencia de Comandos](COMMANDS.md).**
---

## Estructura del Proyecto

El código está organizado para facilitar su mantenimiento y extensión:

- `/cmd`: Contiene la definición de los comandos de la CLI (usando Cobra).
- `/internal`: Código interno de la aplicación no destinado a ser importado por otros proyectos.
- `/pkg`: Paquetes con la lógica de negocio principal que podrían ser reutilizados.

## Cómo Contribuir

¡Las contribuciones son bienvenidas! Por favor, abre un *issue* para discutir tu idea o envía un *pull request* con tu mejora.

## Autor

**Jaime A. Henao**
*Cloud Engineer/Devops*
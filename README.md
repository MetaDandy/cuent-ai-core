# Cuent AI Core

**Backend para generación de audio y video con IA y API REST**

## Descripción

Cuent AI Core es un servicio backend escrito en Go cuyo foco principal es la **generación de audio y video mediante IA**. Este proyecto integra:

* **Generación de texto**: procesamiento de contenido mediante Google Gemini para adaptar el texto. 
* **Síntesis de voz**: conversión de texto a audio usando ElevenLabs para voces naturales y configurables, junto a procesamiento de efectos sonoros. 🔗[https://elevenlabs.io](https://elevenlabs.io)
* **Generación de video**: creación de clips a partir de guiones y assets, utilizando herramientas como FFmpeg 🔗[https://ffmpeg.org](https://ffmpeg.org)
* **API REST**: rutas organizadas bajo `/api/v1` para solicitar generación de texto, audio o video y consultar el estado de los trabajos

## Características principales

* **Procesamiento de texto**: endpoint `/api/v1/text/generate` que recibe un prompt y devuelve un texto refinado por Gemini.
* **Síntesis de voz**: endpoint `/api/v1/audio/synthesize` que acepta texto y opciones de voz, y retorna URL o buffer de audio generado con ElevenLabs.
* **Generación de video**: endpoint `/api/v1/video/create` que combina texto, audio e imágenes para producir un video final procesado con FFmpeg.
* **Colaborativo** (WebSocket): canal `/ws` para colaboración al crear proyectos manualmente.
* **Contenedor de dependencias**: diseño modular con inyección de handlers desde `src/Container`.
* **Docker y Docker Compose**: orquestación completa para desarrollo y despliegue.

## Tecnologías

* **Go 1.21** (última versión estable) 🔗[https://go.dev/doc/go1.21](https://go.dev/doc/go1.21)
* **Fiber v2** (API REST y WebSocket) 🔗[https://docs.gofiber.io/api/v2/introduction](https://docs.gofiber.io/api/v2/introduction)
* **ElevenLabs API** para síntesis de voz 🔗[https://elevenlabs.io](https://elevenlabs.io)
* **Google Gemini** para procesamiento y generación de texto (v1.0)
* **FFmpeg** para renderizado y edición de video 🔗[https://ffmpeg.org](https://ffmpeg.org)
* **Docker** y **Docker Compose** para contenerización

## Requisitos previos

* Go 1.21 o superior
* Cuenta y API key de ElevenLabs
* Acceso a Google Gemini API
* Docker 20.10+ y Docker Compose 2.x
* FFmpeg instalado localmente o en contenedor

## Instalación y ejecución

1. Clonar el repositorio:

   ```bash
   git clone https://github.com/MetaDandy/cuent-ai-core.git
   cd cuent-ai-core
   ```
2. Obtener variables de entorno en un archivo `.env`:

   ```dotenv
   GEMINI_API_KEY=tu_api_key_gemini
   ELEVENLABS_API_KEY=tu_api_key_elevenlabs
   ```
3. Instalar dependencias Go:

   ```bash
   go mod download
   ```
4. Compilar la aplicación:

   ```bash
   go build -o cuent-ai-core ./cmd
   ```
5. Ejecutar localmente:

   ```bash
   ./cuent-ai-core
   ```
6. Con Docker Compose:

   ```bash
   docker-compose up --build
   ```

## Estructura del proyecto

```
cmd/                  # Punto de entrada
src/                  # Lógica principal, handlers y servicios de IA
config/               # Configuración y definiciones de entornos
helper/               # Wrappers de clientes de Gemini, ElevenLabs y FFmpeg
middleware/           # Middleware de Fiber para logs, CORS
Dockerfile            # Imagen de producción
docker-compose.yml    # Orquestación de contenedores
todo.md               # Tareas pendientes
```

## Colaboradores

* **[MetaDandy](https://github.com/MetaDandy)** (mantenedor principal)
* **[MikelKen](https://github.com/MikelKen)** (colaborador)

## Contribución

1. Crear un issue describiendo tu propuesta.
2. Realizar fork y una rama con nombre descriptivo.
3. Abrir un pull request con tus cambios.

## Licencia

Este proyecto no cuenta con licencia específica. Consulta con los mantenedores antes de reutilizar su código.

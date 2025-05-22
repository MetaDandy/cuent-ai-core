# Documentación de pruebas API – CUENT-AI Postman Collection

Este documento describe las pruebas de API definidas en la colección Postman **CUENT-AI**. Incluye detalles de cada endpoint, método HTTP, headers, body de ejemplo y descripción.

---

## Variables de entorno

| Variable       | Descripción             | Ejemplo                        |
| -------------- | ----------------------- | ------------------------------ |
| `{{cuent-ai}}` | URL base de la API      | `http://localhost:8000/api/v1` |
| `{{token}}`    | Token Bearer para auth. | `eyJhbGciOiJIUzI1NiI…`         |

---

## Autenticación

Todas las rutas protegidas requieren el header HTTP:

```
Authorization: Bearer {{token}}
```

---

## Recursos y Endpoints

### 1. ElevenLabs

| Nombre           | Método | Ruta                                  | Descripción                           |
| ---------------- | ------ | ------------------------------------- | ------------------------------------- |
| **tts**          | POST   | `/elevenlabs/SFX`                     | Síntesis de voz simple (TTS)          |
| **SFX avanzado** | POST   | `/elevenlabs/SFX`                     | Generación de efectos de sonido (SFX) |
| **Listar voces** | GET    | `https://api.elevenlabs.io/v1/voices` | Obtener lista de voces ElevenLabs     |

#### 1.1. tts

* **URL**: `{{cuent-ai}}/elevenlabs/SFX`
* **Método**: `POST`
* **Headers**:

  * `Content-Type: application/json`
* **Body**:

  ```json
  {
    "text": "hello "
  }
  ```
* **Descripción**: Genera un archivo de audio con la voz TTS que recita el texto.

#### 1.2. SFX avanzado

* **URL**: `{{cuent-ai}}/elevenlabs/SFX`
* **Método**: `POST`
* **Headers**:

  * `Content-Type: application/json`
* **Body**:

  ```json
  {
    "description": "Heavy rain, realistic ambient noise, no voices, no music",
    "duration_seconds": 2.0,
    "prompt_influence": 1.0,
    "output_format": "mp3_44100_128"
  }
  ```
* **Descripción**: Genera un efecto de sonido a partir de parámetros avanzados.

#### 1.3. Listar voces

* **URL**: `https://api.elevenlabs.io/v1/voices`
* **Método**: `GET`
* **Headers**:

  * `xi-api-key: <tu_api_key>`
* **Descripción**: Recupera todas las voces disponibles en la API de ElevenLabs.

---

### 2. FLOW

| Nombre   | Método | Ruta            | Descripción                     |
| -------- | ------ | --------------- | ------------------------------- |
| **flow** | POST   | `/cuentai/flow` | Procesamiento de flujo de texto |

#### 2.1. flow

* **URL**: `{{cuent-ai}}/cuentai/flow`
* **Método**: `POST`
* **Headers**:

  * `Content-Type: application/json`
* **Body**:

  ```json
  {
    "text": "hello"
  }
  ```
* **Descripción**: Envía texto para procesamiento mediante el endpoint “flow” de Cuent-AI.

---

### 3. SupabaseTest

| Nombre     | Método | Ruta        | Descripción                 |
| ---------- | ------ | ----------- | --------------------------- |
| **upload** | POST   | `/supabase` | Subir un archivo a Supabase |

#### 3.1. upload

* **URL**: `{{cuent-ai}}/supabase`

* **Método**: `POST`

* **Body (form-data)**:

  | Key    | Valor                 | Tipo |
  | ------ | --------------------- | ---- |
  | bucket | `audio`               | text |
  | path   | `test/prueba`         | text |
  | mime   | `image/jpeg`          | text |
  | file   | (seleccionar archivo) | file |

* **Descripción**: Prueba de subida de archivos a un bucket de Supabase.

---

### 4. Usuarios (`users`)

| Nombre                     | Método | Ruta                          | Descripción                            |
| -------------------------- | ------ | ----------------------------- | -------------------------------------- |
| **Listar usuarios**        | GET    | `/users`                      | Obtener todos los usuarios             |
| **Perfil**                 | GET    | `/users/profile`              | Obtener perfil del usuario autenticado |
| **Suscripción**            | GET    | `/users/subscription`         | Obtener datos de suscripción           |
| **Obtener usuario por ID** | GET    | `/users/:id`                  | Obtener un usuario específico          |
| **Iniciar sesión**         | POST   | `/users/sign-in`              | Autenticar usuario                     |
| **Registrar usuario**      | POST   | `/users/sign-up`              | Crear nuevo usuario                    |
| **Agregar suscripción**    | POST   | `/users/:id/add-subscription` | Añadir suscripción a un usuario        |
| **Cambiar contraseña**     | PATCH  | `/users/change-password`      | Actualizar contraseña de usuario       |

#### 4.1. Listar usuarios

* **URL**: `{{cuent-ai}}/users`
* **Método**: `GET`
* **Descripción**: Recupera la lista completa de usuarios.

#### 4.2. Perfil

* **URL**: `{{cuent-ai}}/users/profile`
* **Método**: `GET`
* **Descripción**: Obtiene los datos del perfil del usuario autenticado.

#### 4.3. Suscripción

* **URL**: `{{cuent-ai}}/users/subscription`
* **Método**: `GET`
* **Descripción**: Muestra el estado de suscripción del usuario.

#### 4.4. Obtener usuario por ID

* **URL**: `{{cuent-ai}}/users/{{userId}}`
* **Método**: `GET`
* **Descripción**: Recupera los datos de un usuario por su UUID.

#### 4.5. Iniciar sesión

* **URL**: `{{cuent-ai}}/users/sign-in`
* **Método**: `POST`
* **Body**:

  ```json
  {
    "email": "admin@gmail.com",
    "password": "changeme1234"
  }
  ```
* **Descripción**: Devuelve token JWT y datos de usuario.

#### 4.6. Registrar usuario

* **URL**: `{{cuent-ai}}/users/sign-up`
* **Método**: `POST`
* **Body**:

  ```json
  {
    "name": "user 1",
    "email": "user@gmail.com",
    "password": "changeme123"
  }
  ```
* **Descripción**: Crea un nuevo usuario en el sistema.

#### 4.7. Agregar suscripción

* **URL**: `{{cuent-ai}}/users/{{userId}}/add-subscription`
* **Método**: `POST`
* **Descripción**: Asocia una suscripción a un usuario existente.

#### 4.8. Cambiar contraseña

* **URL**: `{{cuent-ai}}/users/change-password`
* **Método**: `PATCH`
* **Body**:

  ```json
  {
    "old_password": "changeme123",
    "new_password": "changeme1234",
    "confirm_password": "changeme1234"
  }
  ```
* **Descripción**: Actualiza la contraseña del usuario autenticado.

---

### 5. Proyectos (`projects`)

| Nombre                  | Método | Ruta            | Descripción                     |
| ----------------------- | ------ | --------------- | ------------------------------- |
| **Crear proyecto**      | POST   | `/projects`     | Crear un nuevo proyecto         |
| **Listar proyectos**    | GET    | `/projects`     | Obtener todos los proyectos     |
| **Ver proyecto**        | GET    | `/projects/:id` | Obtener detalles de un proyecto |
| **Eliminar proyecto**   | DELETE | `/projects/:id` | Eliminar un proyecto            |
| **Actualizar proyecto** | PATCH  | `/projects/:id` | Modificar datos de un proyecto  |

#### 5.1. Crear proyecto

* **URL**: `{{cuent-ai}}/projects`
* **Método**: `POST`
* **Body**:

  ```json
  {
    "name": "Project 1",
    "description": "description the project 1",
    "user_id": "c372a76d-fb0f-461a-b558-818ab6e426d0"
  }
  ```
* **Descripción**: Registra un nuevo proyecto asociado a un usuario.

#### 5.2. Listar proyectos

* **URL**: `{{cuent-ai}}/projects`
* **Método**: `GET`
* **Descripción**: Recupera todos los proyectos existentes.

#### 5.3. Ver proyecto

* **URL**: `{{cuent-ai}}/projects/{{projectId}}`
* **Método**: `GET`
* **Descripción**: Muestra los detalles de un proyecto por su UUID.

#### 5.4. Eliminar proyecto

* **URL**: `{{cuent-ai}}/projects/{{projectId}}`
* **Método**: `DELETE`
* **Descripción**: Borra el proyecto especificado.

#### 5.5. Actualizar proyecto

* **URL**: `{{cuent-ai}}/projects/{{projectId}}`
* **Método**: `PATCH`
* **Body**:

  ```json
  {
    "name": "Editando el proyecto"
  }
  ```
* **Descripción**: Modifica campos de un proyecto existente.

---

### 6. Scripts (`scripts`)

| Nombre                         | Método | Ruta                      | Descripción                     |
| ------------------------------ | ------ | ------------------------- | ------------------------------- |
| **Crear script**               | POST   | `/scripts`                | Añadir nuevo guion              |
| **Listar scripts**             | GET    | `/scripts`                | Obtener todos los guiones       |
| **Ver script**                 | GET    | `/scripts/:id`            | Detalles de un guion            |
| **Regenerar script**           | PATCH  | `/scripts/:id/regenerate` | Regenera contenido del guion    |
| **Mezclar asset en script**    | POST   | `/scripts/:id/mixed`      | Añade mezcla de assets al guion |
| **Eliminar carpeta de script** | DELETE | `/scripts/:id/folder`     | Elimina carpeta asociada        |

#### 6.1. Crear script

* **URL**: `{{cuent-ai}}/scripts`
* **Método**: `POST`
* **Body**:

  ```json
  {
    "text_entry": "EDIPO…",
    "project_id": "fd6c480e-bfa0-4042-9161-1de0881398ed"
  }
  ```
* **Descripción**: Inserta un nuevo guion de texto para un proyecto.

#### 6.2. Listar scripts

* **URL**: `{{cuent-ai}}/scripts`
* **Método**: `GET`
* **Descripción**: Devuelve la lista de todos los guiones.

#### 6.3. Ver script

* **URL**: `{{cuent-ai}}/scripts/{{scriptId}}`
* **Método**: `GET`
* **Descripción**: Muestra el contenido de un guion por su ID.

#### 6.4. Regenerar script

* **URL**: `{{cuent-ai}}/scripts/{{scriptId}}/regenerate`
* **Método**: `PATCH`
* **Descripción**: Vuelve a generar el texto del guion con IA.

#### 6.5. Mezclar asset en script

* **URL**: `{{cuent-ai}}/scripts/{{scriptId}}/mixed`
* **Método**: `POST`
* **Descripción**: Añade una mezcla de assets al guion especificado.

#### 6.6. Eliminar carpeta de script

* **URL**: `{{cuent-ai}}/scripts/{{scriptId}}/folder`
* **Método**: `DELETE`
* **Descripción**: Elimina la carpeta generada para el script.

---

### 7. Assets (`assets`)

| Nombre                      | Método | Ruta                         | Descripción                         |
| --------------------------- | ------ | ---------------------------- | ----------------------------------- |
| **Listar assets**           | GET    | `/assets`                    | Obtener todos los assets            |
| **Ver asset**               | GET    | `/assets/:id`                | Detalles de un asset                |
| **Subir asset**             | POST   | `/assets/:id`                | Subir o actualizar un asset         |
| **Generar todo**            | POST   | `/assets/:id/generate_all`   | Genera todos los assets para script |
| **Regenerar todo**          | POST   | `/assets/:id/regenerate_all` | Vuelve a generar todos los assets   |
| **Generar video**           | POST   | `/assets/:id/generate_video` | Crea video a partir de assets       |
| **Obtener script de asset** | GET    | `/assets/:id/script`         | Recupera script asociado a asset    |

---

### 8. Suscripciones (`subscription`)

| Nombre                     | Método | Ruta                | Descripción                     |
| -------------------------- | ------ | ------------------- | ------------------------------- |
| **Listar suscripciones**   | GET    | `/subscription`     | Obtener todas las suscripciones |
| **Ver suscripción por ID** | GET    | `/subscription/:id` | Detalles de una suscripción     |

---

### 9. Gemini Formatter

| Nombre        | Método | Ruta                 | Descripción           |
| ------------- | ------ | -------------------- | --------------------- |
| **Formatter** | POST   | `/cuentai/formatter` | Formatea texto con IA |

#### 9.1. Formatter

* **URL**: `{{cuent-ai}}/cuentai/formatter`
* **Método**: `POST`
* **Body**:

  ```json
  {
    "text": "EDIPO…"
  }
  ```
* **Descripción**: Ajusta y formatea el texto entregado.

---

### 10. Aloha (raíz)

* **URL**: `http://localhost:8000/`
* **Método**: `GET`
* **Descripción**: Endpoint de comprobación (ping básico al servidor).

---

> **Nota**: Sustituye `{{userId}}`, `{{projectId}}`, `{{scriptId}}` y `{{token}}` por los valores reales según tu entorno de pruebas.

# TODO DEL BACKEND

- [x] Solo generar los audios en el generate all cuando esta en status pending y error.
- [x] Mixear los audios de un script en uno solo.
- [x] Crear las suscripciones y definir los tokens.
- [x] Vincular las suscripciones a la tabla usario y hacer los handlers correspondiendtes.
- [x] Autenticar el usuario, y que por defecto tenga la suscripcion free.
- [x] Protejer las rutas por medio de un jwt.
- [x] Que todo el flujo de generacion use los tokens de la suscripcion.
- [x] Usar un modelo para la generación de los scripts.
- [x] Cuando un Script se regenera, borrar logicamente los assets vinculados a el, junto con el borrado de las url de la tabla assets.
- [x] Implementar monetizacion y meter pasarela de pago.
- [x] Dejar que el usuario cree de forma manual un asset sin cobrar porque no se usa la ia.
- [x] Implementar generacion de video.
- [x] Poner el precio en susbcription, tambien en los dtos, resetear la db.
- [] La creación manual del script debe tener websocket.
- [] Ver si hacer la edicion del script con sus assets, para editar un asset, su line, su tipo y cambiarlos de posición, combinado con websockets.
- [] Poner las claves de stripe para el front.
- [] Subir a aws.

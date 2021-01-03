# Tarea3Dist

Conexión entre admin y broker, luego admin y DNSs LOCAL

# Admin

- Al ejecutar el admin pide un comando. Si este es `create`, `update` o `delete` lo manda al broker. Este le responde con las IPs en las que estarán los DNS
- Luego se conecta al DNS con esta IP y le manda el mismo comando.
- Los 3 comandos (`create`, `delete` y `update`) están funcionando.
- El `create` debe ser de la siguiente forma: `create nombre.dominio IP`. Ejemplo: `create google.cl 123.0.0.1`. Los demás comandos son tal cual en el PDF de la tarea.
- Consistencia "Read your Writes" hecha.
- Puede hacer updates y deletes incluso si no ha hecho creates.

# DNS

- Crean todo lo que tengan que crear (registros ZF, logs y relojes) y hacen un merge cada 5 minutos (consistencia eventual hecha).

# Cliente

- El cliente se conecta al Broker si hace una petición `get nombre.dominio` y este le responde con el ip del DNS + Reloj + Ip solicitada
- Un ejemplo seria `get test.cl` recibiendo como respuesta ["test.cl 9001 250,0,0 20.0.0.1"] y mostrandolo por consola.
- Si se recibe `get test.cl` nuevamente se hara la peticion y actualizara la informacion correspondiente llevada en memoria.
- Si se realiza otra solicitud `get google.cl` por ejemplo se realizara la peticion y si esta no existe ya en momeria simplemente se agregara a esta.
- Consistencia "Monotonic Reads" realizada.
- Si se ingresa `exit` mostrara por pantalla toda la informacion que de los dominios que a solicitado.

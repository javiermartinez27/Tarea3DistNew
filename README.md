# Admin

- Al ejecutar el admin pide un comando. Si este es `create`, `update` o `delete` lo manda al broker. Este le responde con las IPs en las que estarán los DNS
- Luego se conecta al DNS con esta IP y le manda el mismo comando.
- El `create` debe ser de la siguiente forma: `create nombre.dominio IP`. Ejemplo: `create google.cl 123.0.0.1`. Los demás comandos son tal cual en el PDF de la tarea.

# DNS

- Crean todo lo que tengan que crear (registros ZF, logs y relojes) y hacen un merge cada 5 minutos.

# Cliente

- El cliente se conecta al Broker si hace una petición `get nombre.dominio` y este le responde con el ip del DNS + Reloj + Ip solicitada
- Un ejemplo seria `get test.cl` recibiendo como respuesta ["test.cl 9001 250,0,0 20.0.0.1"] y mostrandolo por consola.
- Si se recibe `get test.cl` nuevamente se hara la peticion y actualizara la informacion correspondiente llevada en memoria.
- Si se realiza otra solicitud `get google.cl` por ejemplo se realizara la peticion y si esta no existe ya en momeria simplemente se agregara a esta.
- Si se ingresa `exit` mostrara por pantalla toda la informacion que de los dominios que a solicitado.

# IPs

- Broker: 10.10.28.154:9000
- DNS1: 10.10.28.155:9001
- DNS2: 10.10.28.156:9002
- DNS3: 10.10.28.157:9003
- Administradores y Clientes están en todas las máquinas.

# PARA CORRER
- Ir a la carpeta que se desea y ejecutar el comando `make`.
- De preferencia, ejecutar primero los DNS2 y 3, y finalmente el 1. En otro caso, si se ejecuta primero el DNS1, los otros DNS se deben conectar en menos de 5 minutos, sino el DNS1 fallará al tratar de hacer la consistencia.

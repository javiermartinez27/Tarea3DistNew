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

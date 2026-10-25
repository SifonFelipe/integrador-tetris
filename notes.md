# Indicaciones / notas

## Personalizaciones

1. Agregamos que con `KeyDown` la pieza pueda caer más rápido.

2. Agregamos la siguiente pieza en el "cuadrado" existente para mostrarla.

3. Agregamos un programa que actua como servidor (en la misma red), el cuál maneja puntajes, partidas y estados. Cada sesión
   de tetris (localhost corriendo) se puede conectar una vez, con el ID asignado por el servidor.
   Tiene bastantes problemas de """seguridad""", podés modificar tu puntaje con un `curl`/`request`, sin
   verificar. Lo ideal hubiera sido usar sockets para manejar eventos en tiempo real pero decidimos
   usar un servidor http. En el `README` están los detalles de como jugar "online".

4. Al ir a `gameover.html`, si los 2 jugadores ya perdieron, va a mostrar quién ganó. Sino va a
   haber un mensaje de `Esperando rival...`.

5. Arreglamos el servidor dado previamente. Anteriormente, lo que pasaba era que, al refrescar la
   página o crear una nueva sesión, esta se sumaba a las conexiones a las que el juego tenia que
mandar los frames del tablero. Aunque una sesión no exista más (cerraste la pestaña, refrescaste la
página), este seguía mandando frames, y los repartía entre las sesiones existentes, no mandando a
todos lo mismo, provocando el efecto de *lag* o *trabado*. Esto lo reemplazamos por un **HUB** que
permite tener varias sesiones al mismo tiempo, pudiendo reconocer cuando una sesión deja de existir,
sacándola de las sesiones a las que tiene que mandar mensajes. Además, incluso si hay varias
sesiones conectadas al juego, esta no tendría ese efecto de *lag* porque **broadcastea** los frames a
todos por igual.

## Usos de IA

Usamos IA para ver sobre los errores de `CORS` (que nos saltaban al fetchear desde el frontend del
tetris al multiplayer server), qué eran y cómo solucionarlo respectivamente.

# Tetris — trabajo por etapas

El objetivo es completar las funciones marcadas con `//PROGRAMAR` en `logica.go`. La interfaz, el bucle y el temporizador ya están resueltos.

## Reglas de esta versión

- La pieza desciende automáticamente cada 800 ms.
- Izquierda y derecha mueven; arriba intenta rotar. Si el giro choca, se cancela.
- Cada pieza tiene cuatro bloques. Sus coordenadas se guardan como `[fila, columna]`.
- Una pieza que no puede bajar se fija. Las filas completas se eliminan y las superiores bajan.
- Se suman 1, 3, 5 u 8 puntos al eliminar 1, 2, 3 o 4 filas en una misma fijación.
- El nivel aumenta cada diez líneas, como indicador; la caída mantiene el mismo intervalo.
- La partida termina si la nueva pieza no puede aparecer en sus cuatro posiciones iniciales.

## Funciones provistas

`obtenerFormaRotacion` entrega las coordenadas locales de una forma. `colisionaConTablero` comprueba límites, bordes y bloques fijos. `verificarLineasCompletas` coordina la revisión de filas y calcula puntos y nivel usando las dos funciones que completarás en la etapa 5.

## Orden sugerido

1. **Tablero:** `generarTablero`. Verificá los cuatro bordes y las celdas interiores.
2. **Aparición:** `generarNuevaPieza` y `actualizarTablero`. Verificá cuatro bloques y la conservación de los bordes y bloques fijos.
3. **Movimiento:** `calcularNuevaPosicionPieza`. Si un bloque de la propuesta choca, ninguno se mueve.
4. **Fijación y fin:** `fijarPieza` y `verificarFinDeJuego`. Probá mover la pieza justo antes de tocar el piso y bloquear una posición de aparición en fila 2 o 3.
5. **Líneas:** `filaCompleta` y `eliminarFila`. Primero una fila completa; después dos consecutivas. Al borrar, copiá desde abajo hacia arriba y vaciá la fila 1 sin tocar las paredes.
6. **Rotación:** `rotarPieza`. Usá el origen de la forma para calcular la propuesta. Probá O y cuatro giros consecutivos de cada pieza. Si hay colisión, conservá las coordenadas y la rotación anteriores.

La fila del array de pieza identifica un bloque; no es una fila del tablero. Por ejemplo, `piezaActiva[2][0]` contiene la fila del tablero donde está el tercer bloque.

## Ejecutar

Desde esta carpeta:

```text
go run main.go logica.go
```

Abrí `http://localhost:8080/`. El esqueleto compila, pero necesita las implementaciones para jugar. Ejecutá una sola versión del proyecto a la vez. Para empezar otra partida después de Game Over, detené el servidor con Ctrl+C y volvé a ejecutarlo.


## Online

Para poder jugar online entre 2 computadoras en la misma red, vas a tener que ejecutar en una sesión de terminal aparte el servidor. El código está en `multiplayer-server/`:

```bash
go run main.go
```

Ahí ya va a estar corriendo el servidor y te va a decir en qué IP. Por ejemplo:

```bash
Servidor escuchando en http://192.168.68.103:9000
```

Luego, cuando vayas a ejecutar el script de tetris, vas a agregarle una `flag` al comando para agregar a qué servidor se va a conectar.

```bash
# dentro de tetris/

go run main.go logica.go --server=http://192.168.68.103:9000
```

Y esto desde las 2 computadoras y listo!. Podes fijarte el estado del servidor corriendo:

```bash
# si es desde la misma PC que corres el servidor
curl http://localhost:9000/match

# si es desde otra PC en la red (o la misma)
curl http://192.168.68.103:9000
```

## Problemas a Arreglar:
1. Cuando carga la página busca conectarse al servidor. Si tenés el localhost abierto de antes y corrés después, no va a intentar conectarse.

2. Si hacés un refresh de la sesión, se olvida tu sesión anterior en el servidor y te intentás conectar como un nuevo jugador.

3. Fácilmente exploiteable, el puntaje se manda por requests HTTP, sin ningún token.

## TODOs:
1. Pasar a websockets para ver actualizaciones del rival, manejar sesiones.

2. Guardar en `sessionStorage` el ID dado para volver a conectarte como ese mismo player

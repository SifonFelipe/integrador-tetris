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


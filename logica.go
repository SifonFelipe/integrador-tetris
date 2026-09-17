package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ═══════════════════════════════════════════════════════════════════════════════
// CONSTANTES
// ═══════════════════════════════════════════════════════════════════════════════

const (
	// Dimensiones del tablero (incluyen bordes)
	constCantFilasTablero    = 22
	constCantColumnasTablero = 14

	// Índices del array piezaActiva [fila, columna]
	constCantColumnasPieza = 2
	constPiezaFila         = 0
	constPiezaCol          = 1

	// Tipos de piezas tetrominó
	constPiezaI = 0
	constPiezaO = 1
	constPiezaT = 2
	constPiezaS = 3
	constPiezaZ = 4
	constPiezaL = 5
	constPiezaJ = 6

	// Símbolos del tablero
	constSimboloVacio      = ""
	constSimboloBorde      = "X"
	constSimboloBloqueFijo = "B"
	constSimboloPieza      = "P"

	// Puntuación
	constPuntos1Linea  = 1
	constPuntos2Lineas = 3
	constPuntos3Lineas = 5
	constPuntos4Lineas = 8

	// Intervalo fijo de caída (ms); el nivel es informativo
	constIntervaloDescenso = 800
	constLineasPorNivel    = 10
)

// ═══════════════════════════════════════════════════════════════════════════════
// VARIABLES GLOBALES DE CONTROL (leídas desde main.go por keyPressHandler)
// ═══════════════════════════════════════════════════════════════════════════════

var (
	direccionCol   int  // -1 = izquierda, 0 = quieto, 1 = derecha
	rotarPiezaFlag bool // true cuando el usuario presiona ↑
)

// ═══════════════════════════════════════════════════════════════════════════════
// LOOP PRINCIPAL DEL JUEGO
// ═══════════════════════════════════════════════════════════════════════════════

// generarEventos contiene el loop principal del juego.
// Es llamado como goroutine desde main.go.
func generarEventos() {
	var (
		tablero          [constCantFilasTablero][constCantColumnasTablero]string
		piezaActiva      [4][constCantColumnasPieza]int
		tipoPiezaActiva  int
		rotacionActual   int
		puntos           int
		nivel            int
		lineasEliminadas int
	)

	rand.Seed(time.Now().Unix())

	// ── Inicialización ────────────────────────────────────────────────────────
	nivel = 1
	puntos = 0
	lineasEliminadas = 0
	rotacionActual = 0

	// Inicializar variables de control de teclado
	direccionCol = 0
	rotarPiezaFlag = false

	// Generar tablero con bordes
	tablero = generarTablero()

	// Generar primera pieza
	piezaActiva, tipoPiezaActiva = generarNuevaPieza(constCantColumnasTablero)

	// Comprobar la aparición antes de dibujar.
	if verificarFinDeJuego(tablero, piezaActiva) {
		enviarActualizacionTablero(tablero)
		enviarGameOver(puntos)
		return
	}

	// Dibujar la primera pieza en el tablero
	actualizarTablero(&tablero, piezaActiva)

	// Enviar estado inicial al cliente
	enviarActualizacionTablero(tablero)
	enviarActualizacionTexto(fmt.Sprint("Puntos: ", puntos, " | Nivel: ", nivel))

	// Temporizador de descenso automático
	ultimoDescenso := time.Now()

	// ── Loop principal ────────────────────────────────────────────────────────
	for {
		// ── Procesar entrada del usuario ──────────────────────────────────────

		// Movimiento lateral
		if direccionCol != 0 {
			calcularNuevaPosicionPieza(tablero, &piezaActiva, 0, direccionCol)
			direccionCol = 0
		}

		// Rotación
		if rotarPiezaFlag {
			rotarPieza(tablero, &piezaActiva, tipoPiezaActiva, &rotacionActual)
			rotarPiezaFlag = false
		}

		// ── Descenso automático con intervalo fijo ─────────────────────
		if time.Since(ultimoDescenso) >= time.Duration(constIntervaloDescenso)*time.Millisecond {
			ultimoDescenso = time.Now()

			if !calcularNuevaPosicionPieza(tablero, &piezaActiva, 1, 0) {
				// La pieza tocó el fondo o una pieza fija → fijar
				fijarPieza(&tablero, piezaActiva)
				verificarLineasCompletas(&tablero, &puntos, &lineasEliminadas, &nivel)

				// Generar nueva pieza
				piezaActiva, tipoPiezaActiva = generarNuevaPieza(constCantColumnasTablero)
				rotacionActual = 0
				// La pieza nueva no debe sobrescribir bloques ya fijados.
				if verificarFinDeJuego(tablero, piezaActiva) {
					enviarActualizacionTablero(tablero)
					enviarGameOver(puntos)
					return
				}
			}
		}

		// ── Actualizar tablero y enviar al cliente ────────────────────────────
		actualizarTablero(&tablero, piezaActiva)
		enviarActualizacionTablero(tablero)
		enviarActualizacionTexto(fmt.Sprint(
			"Puntos: ", puntos,
			" | Nivel: ", nivel,
			" | Líneas: ", lineasEliminadas,
		))

		// Esperar un ciclo antes del próximo frame (50ms ≈ 20 FPS)
		time.Sleep(50 * time.Millisecond)
	}
}

// ═══════════════════════════════════════════════════════════════════════════════
// FUNCIONES A IMPLEMENTAR
// ═══════════════════════════════════════════════════════════════════════════════

// generarTablero genera la matriz tablero inicializando todas las celdas como
// vacías y colocando los bordes (piso, paredes laterales y tope).
//
// Retorna: la matriz tablero generada.
func generarTablero() [constCantFilasTablero][constCantColumnasTablero]string {
	var tablero [constCantFilasTablero][constCantColumnasTablero]string

	//PROGRAMAR

	return tablero
}

// obtenerFormaRotacion devuelve los offsets (desplazamientos [fila, col] relativos
// al origen de la forma) de los 4 bloques para el tipo de pieza y estado de rotación dados.
//
// Parámetros:
//
//	tipoPieza: tipo de pieza (0=I, 1=O, 2=T, 3=S, 4=Z, 5=L, 6=J)
//	rotacion:  estado de rotación (0 a 3)
//
// Retorna: array de 4 pares [fila, col] con los offsets de cada bloque.
func obtenerFormaRotacion(tipoPieza int, rotacion int) [4][constCantColumnasPieza]int {
	// [7 piezas][4 rotaciones][4 bloques][fila, col]
	// Las coordenadas son relativas al origen de la forma, no al bloque [1].
	formas := [7][4][4][constCantColumnasPieza]int{

		// Pieza I (tipo 0). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   ■■■■    □□■□    □□□□    □■□□
		// f=1   □□□□    □□■□    □□□□    □■□□
		// f=2   □□□□    □□■□    ■■■■    □■□□
		// f=3   □□□□    □□■□    □□□□    □■□□
		{
			{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
			{{0, 2}, {1, 2}, {2, 2}, {3, 2}},
			{{2, 0}, {2, 1}, {2, 2}, {2, 3}},
			{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
		},

		// Pieza O (tipo 1). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   ■■□□    ■■□□    ■■□□    ■■□□
		// f=1   ■■□□    ■■□□    ■■□□    ■■□□
		// f=2   □□□□    □□□□    □□□□    □□□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
			{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
		},

		// Pieza T (tipo 2). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   □■□□    ■□□□    □□□□    □■□□
		// f=1   ■■■□    ■■□□    ■■■□    ■■□□
		// f=2   □□□□    ■□□□    □■□□    □■□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 1}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 0}, {1, 0}, {1, 1}, {2, 0}},
			{{1, 0}, {1, 1}, {1, 2}, {2, 1}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 1}},
		},

		// Pieza S (tipo 3). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   □■■□    ■□□□    □■■□    ■□□□
		// f=1   ■■□□    ■■□□    ■■□□    ■■□□
		// f=2   □□□□    □■□□    □□□□    □■□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
			{{0, 1}, {0, 2}, {1, 0}, {1, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
		},

		// Pieza Z (tipo 4). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   ■■□□    □■□□    ■■□□    □■□□
		// f=1   □■■□    ■■□□    □■■□    ■■□□
		// f=2   □□□□    ■□□□    □□□□    ■□□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
			{{0, 0}, {0, 1}, {1, 1}, {1, 2}},
			{{0, 1}, {1, 0}, {1, 1}, {2, 0}},
		},

		// Pieza L (tipo 5). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   ■□□□    ■■■□    ■■□□    □□■□
		// f=1   ■□□□    ■□□□    □■□□    ■■■□
		// f=2   ■■□□    □□□□    □■□□    □□□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 0}, {1, 0}, {2, 0}, {2, 1}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 0}},
			{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
			{{0, 2}, {1, 0}, {1, 1}, {1, 2}},
		},

		// Pieza J (tipo 6). ■ = bloque; □ = celda vacía.
		// Cada dibujo conserva las coordenadas locales de la tabla (4 x 4).
		//       R0      R1      R2      R3
		// f=0   □■□□    ■□□□    ■■□□    ■■■□
		// f=1   □■□□    ■■■□    ■□□□    □□■□
		// f=2   ■■□□    □□□□    ■□□□    □□□□
		// f=3   □□□□    □□□□    □□□□    □□□□
		{
			{{0, 1}, {1, 1}, {2, 0}, {2, 1}},
			{{0, 0}, {1, 0}, {1, 1}, {1, 2}},
			{{0, 0}, {0, 1}, {1, 0}, {2, 0}},
			{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
		},
	}

	return formas[tipoPieza][rotacion]
}

// generarNuevaPieza genera una pieza aleatoria (tipo 0 a 6) centrada en la
// parte superior del tablero (fila 1).
//
// Parámetros:
//
//	cantColumnasTablero: ancho total del tablero (para calcular el centrado)
//
// Retorna:
//   - Array de 4 bloques con las posiciones [fila, col] de la nueva pieza.
//   - Entero con el tipo de pieza generada (0 a 6).
func generarNuevaPieza(cantColumnasTablero int) ([4][constCantColumnasPieza]int, int) {
	var pieza [4][constCantColumnasPieza]int

	//PROGRAMAR

	return pieza, 0
}

// actualizarTablero limpia los símbolos de la pieza activa anterior del tablero
// y vuelve a dibujarla en su posición actual.
// Escribir P solo en celdas interiores vacías; conservar bordes y bloques fijos.
//
// Parámetros:
//
//	tablero:     puntero a la matriz tablero (por referencia)
//	piezaActiva: posiciones actuales de los 4 bloques (por valor)
func actualizarTablero(
	tablero *[constCantFilasTablero][constCantColumnasTablero]string,
	piezaActiva [4][constCantColumnasPieza]int,
) {
	//PROGRAMAR
}

// FUNCIÓN PROVISTA: usarla para comprobar movimientos, giros y aparición.
// colisionaConTablero es una función auxiliar que verifica si alguno de los
// bloques propuestos choca con un borde o un bloque fijo.
// Las celdas que contienen constSimboloPieza (la propia pieza activa) se ignoran.
func colisionaConTablero(
	tablero [constCantFilasTablero][constCantColumnasTablero]string,
	bloques [4][constCantColumnasPieza]int,
) bool {
	for i := 0; i < 4; i++ {
		f := bloques[i][constPiezaFila]
		c := bloques[i][constPiezaCol]

		// Verificar límites del array
		if f < 0 || f >= constCantFilasTablero || c < 0 || c >= constCantColumnasTablero {
			return true
		}

		celda := tablero[f][c]
		if celda == constSimboloBorde || celda == constSimboloBloqueFijo {
			return true
		}
	}
	return false
}

// calcularNuevaPosicionPieza intenta mover la pieza activa según la dirección dada.
// Solo actualiza la pieza si el movimiento no genera colisiones con bordes o bloques fijos.
//
// Parámetros:
//
//	tablero:      estado actual del tablero (por valor)
//	piezaActiva:  puntero a la pieza activa (por referencia)
//	direccionFila: 0 = sin movimiento vertical, 1 = bajar una fila
//	direccionCol:  -1 = izquierda, 0 = sin movimiento, 1 = derecha
//
// Retorna: true si el movimiento fue válido, false si hubo colisión.
func calcularNuevaPosicionPieza(
	tablero [constCantFilasTablero][constCantColumnasTablero]string,
	piezaActiva *[4][constCantColumnasPieza]int,
	direccionFila int,
	direccionCol int,
) bool {
	//PROGRAMAR

	return false
}

// rotarPieza intenta rotar la pieza activa 90° en sentido horario.
// Recuperar el origen restando al bloque [1] sus coordenadas en la forma actual.
// Sumar a ese origen las coordenadas de la nueva forma; el bloque [1] no es el origen.
// Si la propuesta colisiona, conservar la pieza y su estado de rotación.
//
// Parámetros:
//
//	tablero:        estado actual del tablero (por valor)
//	piezaActiva:    puntero a la pieza activa (por referencia)
//	tipoPieza:      tipo de la pieza activa
//	rotacionActual: puntero al estado de rotación actual (por referencia)
func rotarPieza(
	tablero [constCantFilasTablero][constCantColumnasTablero]string,
	piezaActiva *[4][constCantColumnasPieza]int,
	tipoPieza int,
	rotacionActual *int,
) {
	//PROGRAMAR
}

// fijarPieza convierte los bloques de la pieza activa en bloques fijos en el tablero.
// Se llama con una pieza válida cuando no puede seguir bajando.
// Primero limpiar todos los símbolos P anteriores. Después escribir B en las
// cuatro coordenadas actuales, aunque esas celdas todavía no contengan P.
//
// Parámetros:
//
//	tablero:     puntero a la matriz tablero (por referencia)
//	piezaActiva: posiciones actuales de los 4 bloques (por valor)
func fijarPieza(
	tablero *[constCantFilasTablero][constCantColumnasTablero]string,
	piezaActiva [4][constCantColumnasPieza]int,
) {
	//PROGRAMAR
}

// filaCompleta indica si todas las celdas interiores de una fila contienen B.
// La fila recibida está entre 1 y constCantFilasTablero-2; no revisar paredes.
func filaCompleta(tablero [constCantFilasTablero][constCantColumnasTablero]string, fila int) bool {
	//PROGRAMAR
	return false
}

// eliminarFila copia las filas superiores una posición hacia abajo y vacía la fila 1.
// Recorrer desde la fila eliminada hacia arriba. Conservar las paredes y el piso.
func eliminarFila(tablero *[constCantFilasTablero][constCantColumnasTablero]string, fila int) {
	//PROGRAMAR
}

// FUNCIÓN PROVISTA: coordina la revisión de filas y calcula puntos y nivel.
// verificarLineasCompletas detecta y elimina filas llenas. Actualiza puntos y nivel.
func verificarLineasCompletas(
	tablero *[constCantFilasTablero][constCantColumnasTablero]string,
	puntos *int,
	lineasEliminadas *int,
	nivel *int,
) int {
	eliminadas := 0

	// Recorrer desde el piso hacia arriba (saltando los bordes)
	fila := constCantFilasTablero - 2
	for fila >= 1 {
		if filaCompleta(*tablero, fila) {
			eliminarFila(tablero, fila)

			eliminadas++
			// No decrementar fila: la nueva fila[fila] vino de arriba y hay que revisarla
		} else {
			fila--
		}
	}

	// Calcular puntos según cantidad de líneas eliminadas en esta llamada
	switch eliminadas {
	case 1:
		*puntos += constPuntos1Linea
	case 2:
		*puntos += constPuntos2Lineas
	case 3:
		*puntos += constPuntos3Lineas
	case 4:
		*puntos += constPuntos4Lineas
	}

	// Actualizar contadores globales
	*lineasEliminadas += eliminadas
	*nivel = *lineasEliminadas/constLineasPorNivel + 1

	return eliminadas
}

// verificarFinDeJuego indica si la nueva pieza no puede aparecer.
// Se llama ANTES de dibujarla. Revisar sus cuatro coordenadas con colisionaConTablero.
// No alcanza con revisar la fila 1: una pieza puede ocupar también las filas 2 y 3.
func verificarFinDeJuego(
	tablero [constCantFilasTablero][constCantColumnasTablero]string,
	nuevaPieza [4][constCantColumnasPieza]int,
) bool {
	//PROGRAMAR
	return false
}

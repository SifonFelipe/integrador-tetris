package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

// Canal para enviar actualizaciones al cliente vía SSE
var updates = make(chan string)

func main() {
	// Página principal del juego
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Ruta para manejar la solicitud POST del evento de teclado
	http.HandleFunc("/keypress", keyPressHandler)

	// Ruta SSE: envía actualizaciones del tablero al cliente en tiempo real
	http.HandleFunc("/updates", updatesHandler)

	// Pantalla de Game Over
	http.HandleFunc("/gameover", gameoverHandler)

	// Pantalla de Win (tablero completado)
	http.HandleFunc("/win", winHandler)

	// Inicia la goroutine con el loop principal del juego
	go generarEventos()

	fmt.Println("Abrir navegador e ingresar a http://localhost:8080/")

	// Inicia el servidor en el puerto 8080
	http.ListenAndServe(":8080", nil)
}

// ── Handlers de rutas ────────────────────────────────────────────────────────

// gameoverHandler sirve la pantalla de game over con el puntaje final.
func gameoverHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Points string
	}
	data := PageData{
		Points: r.URL.Query().Get("points"),
	}
	tmpl, err := template.ParseFiles("gameover.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// winHandler sirve la pantalla de victoria (tablero lleno de líneas completadas — no aplica en Tetris,
// pero se deja disponible para la personalización).
func winHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Points string
	}
	data := PageData{
		Points: r.URL.Query().Get("points"),
	}
	tmpl, err := template.ParseFiles("win.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// updatesHandler mantiene la conexión SSE abierta y envía actualizaciones al cliente.
func updatesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		update := <-updates
		fmt.Fprintf(w, "data: %s\n\n", update)
		w.(http.Flusher).Flush()
	}
}

// keyPressHandler recibe las teclas presionadas por el usuario y actualiza
// las variables globales de control del juego.
func keyPressHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error al leer los datos JSON", http.StatusBadRequest)
		return
	}

	keyPressed, ok := data["key"]
	if !ok {
		http.Error(w, "Tecla no proporcionada", http.StatusBadRequest)
		return
	}

	// Actualizar variables de control según la tecla presionada
	switch keyPressed {
	case "ArrowLeft":
		direccionCol = -1 // mover pieza a la izquierda
	case "ArrowRight":
		direccionCol = 1 // mover pieza a la derecha
	case "ArrowUp":
		rotarPiezaFlag = true // rotar pieza
	}

	w.WriteHeader(http.StatusOK)
}

// ── Funciones de comunicación con el cliente ─────────────────────────────────

// enviarActualizacionTablero convierte el tablero en JSON y lo envía al cliente.
func enviarActualizacionTablero(tablero [constCantFilasTablero][constCantColumnasTablero]string) {
	update, err := json.Marshal(tablero)
	if err != nil {
		fmt.Println("Error al convertir la matriz en JSON:", err)
		return
	}
	updates <- string(update)
}

// enviarActualizacionTexto envía un mensaje de texto al cliente (puntos, nivel, etc.).
func enviarActualizacionTexto(text string) {
	updates <- "{\"is_text\": true, \"text\": \"" + text + "\"}"
}

// enviarGameOver envía la señal de fin de juego con el puntaje final.
func enviarGameOver(points int) {
	texto := fmt.Sprint("{\"game_over\": true, \"points\": \"", points, "\"}")
	updates <- texto
	fmt.Println("Game Over. Points:", points)
}

// enviarWin envía la señal de victoria con el puntaje (para personalizaciones).
func enviarWin(points int) {
	texto := fmt.Sprint("{\"win\": true, \"points\": \"", points, "\"}")
	updates <- texto
	fmt.Println("Win. Points:", points)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

// Canal para enviar actualizaciones al cliente vía SSE
var updates = make(chan string)

const (
	constMatchWaiting = "waiting"
	constMatchPlaying = "playing"
	constMatchFinished = "finished"
	constMatchDraw = "draw"
)

type Player struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Points int `json:"points"`
	GameOver bool `json:"game_over"`
}

type Match struct {
	Player1 *Player `json:"player1"`
	Player2 *Player `json:"player2"`
}

type MatchResponse struct {
	Player1 *Player `json:"player1"`
	Player2 *Player `json:"player2"`
	Status string `json:"status"`
	Result string `json:"result"`
}

var currentMatch = Match{}


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

	// route to connect to server
	http.HandleFunc("/join", joinHandler)

	// route to send score
	http.HandleFunc("/score", scoreHandler)

	// route to get the match status
	http.HandleFunc("/match", matchHandler)

	// route to say gameover to the server
	http.HandleFunc("/gameover-server", multiplayerGameOverHandler)

	// Inicia la goroutine con el loop principal del juego
	go generarEventos()

	fmt.Println("Abrir navegador e ingresar a http://localhost:8080/")

	// Inicia el servidor en el puerto 8080
	http.ListenAndServe(":8080", nil)
}

// -- utils para match --
func getMatchResult() string {
	if currentMatch.Player1 == nil || currentMatch.Player2 == nil {
		return constMatchWaiting
	}

	if !currentMatch.Player1.GameOver || !currentMatch.Player2.GameOver {
		return constMatchPlaying
	}

	if currentMatch.Player1.Points > currentMatch.Player2.Points {
		return currentMatch.Player1.ID
	}

	if currentMatch.Player2.Points > currentMatch.Player1.Points {
		return currentMatch.Player2.ID
	}

	return constMatchDraw
}


// ── Handlers de rutas ────────────────────────────────────────────────────────

func joinHandler(w http.ResponseWriter, r *http.Request) {
	/* joinHandler handles the requests to join to the game.
	It accepts POST requests with a JSON body containing the player's name. */

	if r.Method != http.MethodPost {  // only accept POST requests
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Parse the JSON body to get the player's name
	var data struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	var player Player

	// ask for a empty slot in the current match
	if currentMatch.Player1 == nil {
		player = Player{ID: "1", Name: data.Name + " (Jugador 1)"}
		currentMatch.Player1 = &player

	} else if currentMatch.Player2 == nil {
		player = Player{ID: "2", Name: data.Name + " (Jugador 2)"}
		currentMatch.Player2 = &player

	} else {
		http.Error(w, "El juego ya tiene dos jugadores", http.StatusForbidden)
		return
	}

	// Return the player information as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}


func scoreHandler(w http.ResponseWriter, r *http.Request) {
	/* scoreHandler handles the requests to update the score of a player. */

	// NOTE: this is easily exploitable, curl'ing POST requests to this endpoint
	// can change the score of any player.

	if r.Method != http.MethodPost {  // only accept POST requests
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		PlayerID string `json:"player_id"`
		Points int `json:"points"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// set points for the player in the current match
	fmt.Println(data)
	fmt.Println(currentMatch)
	if currentMatch.Player1 != nil && currentMatch.Player1.ID == data.PlayerID {
		currentMatch.Player1.Points = data.Points

	} else if currentMatch.Player2 != nil && currentMatch.Player2.ID == data.PlayerID {
		currentMatch.Player2.Points = data.Points
	
	} else {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}


func matchHandler(w http.ResponseWriter, r *http.Request) {
	/* matchHandler handles the requests to get the current match status. */
	
	if r.Method != http.MethodGet {  // only accept GET requests
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	result := getMatchResult()
	status := constMatchPlaying

	if result == constMatchWaiting {
		status = constMatchWaiting
	} else if result != constMatchPlaying {
		status = constMatchFinished
	}

	response := MatchResponse{
		Player1: currentMatch.Player1,
		Player2: currentMatch.Player2,
		Status: status,
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func multiplayerGameOverHandler(w http.ResponseWriter, r *http.Request) {
	/* multiplayerGameOverHandler handles the requests to set the game over status for a player. */
	
	if r.Method != http.MethodPost {  // only accept POST requests
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		PlayerID string `json:"player_id"`
		Points int `json:"points"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	var player *Player

	if currentMatch.Player1 != nil && currentMatch.Player1.ID == data.PlayerID {
		player = currentMatch.Player1

	} else if currentMatch.Player2 != nil && currentMatch.Player2.ID == data.PlayerID {
		player = currentMatch.Player2
	
	} else {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	player.Points = data.Points
	player.GameOver = true

	w.WriteHeader(http.StatusOK)
}


// gameoverHandler sirve la pantalla de game over con el puntaje final.
func gameoverHandler(w http.ResponseWriter, r *http.Request) {
	type PageData struct {
		Points string
		PlayerID string
	}

	data := PageData{
		Points: r.URL.Query().Get("points"),
		PlayerID: r.URL.Query().Get("player_id"),
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
	case "ArrowDown":
		acelerarCaidaFlag = true // acelerar caída de la pieza
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

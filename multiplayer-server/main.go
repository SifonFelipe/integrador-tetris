package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

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


type RivalResponse struct {
	Rival *Player `json:"rival"`
}


var currentMatch = Match{}


func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	/* middleware function to enable CORS for the given handler. */
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}


func getLocalIp() (string, error) {
	// Get the local IP address of the machine
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, addr := range addrs {

		// verify if the address is an IP address and not a loopback address (127.0.0.1)
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no local IP address found")
}


func main() {
	http.HandleFunc("/join", enableCORS(joinHandler))  // join the game
	http.HandleFunc("/score", enableCORS(scoreHandler))  // update the score of a player
	http.HandleFunc("/gameover", enableCORS(gameOverHandler))  // set the game over status for a player
	http.HandleFunc("/match", enableCORS(matchHandler))  // get the current match status
	http.HandleFunc("/rival", enableCORS(rivalHandler))  // get the rival's information

	port := ":9000"
	ipLocal, err := getLocalIp()

	if err != nil {
		fmt.Println("Error al obtener la IP local:", err)
		ipLocal = "localhost"
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Servidor escuchando en http://%s%s\n", ipLocal, port)

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}

}


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


// === Route Handlers ===

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


func gameOverHandler(w http.ResponseWriter, r *http.Request) {
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


func rivalHandler(w http.ResponseWriter, r *http.Request) {
	/* rivalHandler handles the requests to get the rival's information. */
	if r.Method != http.MethodGet {  // only accept GET requests
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		PlayerID string `json:"player_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	var rival *Player

	if currentMatch.Player1 != nil && currentMatch.Player1.ID == data.PlayerID {
		rival = currentMatch.Player2
	} else if currentMatch.Player2 != nil && currentMatch.Player2.ID == data.PlayerID {
		rival = currentMatch.Player1
	} else {
		http.Error(w, "Jugador no encontrado", http.StatusNotFound)
		return
	}

	response := RivalResponse{Rival: rival}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

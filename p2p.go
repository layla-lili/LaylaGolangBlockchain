package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var peers = make(map[*websocket.Conn]bool) // Connected peers
var mu sync.Mutex                          // Mutex for thread safety

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Handle incoming WebSocket connections
func HandleP2PConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading WebSocket:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	peers[conn] = true
	mu.Unlock()

	fmt.Println("New peer connected!")
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(peers, conn)
			mu.Unlock()
			fmt.Println("Peer disconnected")
			return
		}
		fmt.Println("Received message:", string(msg))
	}
}

// Broadcast blockchain updates to all peers
// func BroadcastBlockchain() {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	for conn := range peers {
// 		err := conn.WriteJSON(Blockchain)
// 		if err != nil {
// 			conn.Close()
// 			delete(peers, conn)
// 		}
// 	}
// }

// Start WebSocket server
func StartP2PServer(port string) {
	http.HandleFunc("/ws", HandleP2PConnection)
	go http.ListenAndServe(":"+port, nil)
	fmt.Println("P2P Server running on port", port)
}

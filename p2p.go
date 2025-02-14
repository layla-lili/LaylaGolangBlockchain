package main

import (
	"fmt"
	"net/http"
)


// Start WebSocket server
func StartP2PServer(port string) {
	go http.ListenAndServe(":"+port, nil)
	fmt.Println("P2P Server running on port", port)
}

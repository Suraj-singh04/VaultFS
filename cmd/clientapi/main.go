package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Suraj-singh04/vaultfs/internal/clientapi"
)

func main() {
	nameNodeAddr := os.Getenv("NAMENODE_ADDR")
	if nameNodeAddr == "" {
		nameNodeAddr = "localhost:50051"
	}

	coordinator, err := clientapi.NewCoordinator(nameNodeAddr)
	if err != nil {
		log.Fatalf("Failed to connect to NameNode: %v", err)
	}

	handler := clientapi.NewHandler(coordinator)

	http.HandleFunc("/upload", handler.Upload)
	http.HandleFunc("/fetch", handler.Fetch)

	log.Println("Client API server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	router "backend/src/infra"
	"log"
)

func main() {
	r := router.NewRouter()
	
	log.Println("Server running on localhost:80")
	if err := r.Run(":80"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

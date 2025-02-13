package main

import (
	"fmt"
	"log"
)
func main(){
	store, err := NewPostgresStorage()
	if err != nil {
		log.Fatalf("Error creating storage: %v", err)
	}
	if err := store.Init(); err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}
	fmt.Printf("store: %v\n", store)
	server := NewAPIServer(":8080", store)
	server.Start()

}
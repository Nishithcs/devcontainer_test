package main

import (
	"clusterix-code/commands"
	"log"
)

func main() {
	if err := commands.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}

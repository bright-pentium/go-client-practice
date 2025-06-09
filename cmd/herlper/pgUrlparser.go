package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bright-pentium/go-client-practice/internal/configs"
)

func main() {
	// Check for the .env file path as a command-line argument.
	if len(os.Args) != 2 {
		log.Printf("Usage: %s path/to/.env\n", os.Args[0])
		os.Exit(1)
	}
	envFilePath := os.Args[1]

	// Load the PostgreSQL configuration from the .env file.
	cfg, err := configs.LoadPostgresConfig(envFilePath)
	if err != nil {
		log.Fatalf("Error: %v\n", err) // Use log.Fatalf for non-recoverable errors
	}

	// Construct the database URL.
	dbURL := cfg.DBURL()
	fmt.Println(dbURL) // Print to standard output, consistent with the shell script
}

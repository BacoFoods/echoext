package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting example server...")

	// Start blocks until an interrupt/terminate signal is received and then
	// gracefully shuts down both the main and metrics servers.
	if err := StartServer(); err != nil {
		fmt.Printf("server error: %v\n", err)
	}

	fmt.Println("Server stopped.")
}

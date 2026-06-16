package main

import (
	"fmt"
	"log"

	"github.com/riazahmedshah/go-booking/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("✅ Config loaded successfully! Environment: %s\n", cfg.Database.User)
}

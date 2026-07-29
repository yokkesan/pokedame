package main

import (
	"log"

	"api-generated/database"
	_ "api-generated/routers"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	if err := database.Initialize(); err != nil {
		log.Fatalf("application startup failed: %v", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("database shutdown failed: %v", err)
		}
	}()

	log.Println("database connection established")

	beego.Run()
}

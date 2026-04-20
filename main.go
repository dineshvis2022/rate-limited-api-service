package main

import (
	"api-service/routes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Welcome to Rate-Limiter Api Service!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("[ERROR] loading env file: ", err)
	}

	r := routes.SetupRouter()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("[ERROR] Missing port!: ")
	}
	r.Run(":" + port)

}

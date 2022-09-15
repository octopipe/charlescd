package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/auth"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()
	e = auth.NewController(e)
	e.Logger.Fatal(e.Start(":8080"))
}

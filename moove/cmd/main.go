package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Nerzal/gocloak/v11"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/modules"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	keycloak := gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_PATH"), gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms"))
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorization := c.Request().Header.Get("Authorization")
			if authorization == "" {
				log.Println("Invalid header!")
				return c.JSON(401, "Unauthorized")
			}
			res, err := keycloak.RetrospectToken(c.Request().Context(),
				strings.TrimPrefix(authorization, "Bearer "),
				os.Getenv("KEYCLOAK_CLIENT_ID"),
				os.Getenv("KEYCLOAK_CLIENT_SECRET"),
				os.Getenv("KEYCLOAK_REALM"))

			if err != nil {
				return c.JSON(400, "Invalid token")
			}

			fmt.Println(res)

			if !*res.Active {
				return c.JSON(401, "Unauthorized")
			}

			return nil
		}
	})
	e = modules.NewController(e)
	e.Logger.Fatal(e.Start(":8080"))
}

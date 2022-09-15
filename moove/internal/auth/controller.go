package auth

import (
	"context"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v11"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	keycloak     gocloak.GoCloak
	clientId     string
	clientSecret string
	realm        string
}

type SigninRequest struct {
	Username string
	Password string
}

func NewController(e *echo.Echo) *echo.Echo {
	c := Controller{
		keycloak:     gocloak.NewClient(os.Getenv("KEYCLOAK_BASE_PATH"), gocloak.SetAuthAdminRealms("admin/realms"), gocloak.SetAuthRealms("realms")),
		clientId:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		clientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		realm:        os.Getenv("KEYCLOAK_REALM"),
	}
	g := e.Group("/auth")
	g.POST("/signin", c.signin)
	return e
}

// http://localhost:58236/realms/ciam/protocol/openid-connect/token

func (controller Controller) signin(c echo.Context) error {
	r := new(SigninRequest)
	if err := c.Bind(r); err != nil {
		log.Println(err)
		return err
	}

	jwt, err := controller.keycloak.Login(context.Background(),
		controller.clientId,
		controller.clientSecret,
		controller.realm,
		r.Username,
		r.Password)

	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(200, jwt)
}

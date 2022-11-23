package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type AuthMiddleware struct {
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewAuthMiddleware(provider *oidc.Provider, oauth2Config oauth2.Config, verifier *oidc.IDTokenVerifier) AuthMiddleware {
	return AuthMiddleware{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}
}

func (a AuthMiddleware) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := a.getToken(c)
		if token == "" {
			zap.S().Error(errors.New("invalid token"))
			return echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid token").Error())
		}

		parser := jwt.NewParser(jwt.WithoutClaimsValidation())
		var unverifiedClaims jwt.RegisteredClaims
		_, _, err := parser.ParseUnverified(token, &unverifiedClaims)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		var idToken *oidc.IDToken
		for _, aud := range unverifiedClaims.Audience {
			verifier := a.provider.Verifier(&oidc.Config{ClientID: aud})
			idToken, err = verifier.Verify(c.Request().Context(), token)
			if err == nil {
				break
			}
		}

		if idToken == nil {
			zap.S().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			zap.S().Error(err)
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return next(c)
	}
}

func (a AuthMiddleware) getToken(c echo.Context) string {
	authorization := c.Request().Header.Get("authorization")

	token := strings.TrimPrefix(authorization, "Bearer ")

	if token != "" && tokenIsValid(token) {
		return token
	}

	return ""
}

func tokenIsValid(token string) bool {
	return len(strings.SplitN(token, ".", 3)) == 3
}

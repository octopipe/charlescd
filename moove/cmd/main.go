package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/octopipe/charlescd/moove/internal/circle"
	circleHandler "github.com/octopipe/charlescd/moove/internal/circle/handler"
	"github.com/octopipe/charlescd/moove/internal/core/grpcclient"
	"github.com/octopipe/charlescd/moove/internal/workspace"
	workspaceHandler "github.com/octopipe/charlescd/moove/internal/workspace/handler"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction(zap.AddCaller())
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	// provider, err := oidc.NewProvider(context.Background(), "http://127.0.0.1:53192/realms/master")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// clientID := "charles-cd"
	// clientSecret := "jwR8cBaDLlGlY8cEzYsq05MEgLtRy4mv"

	// oauth2Config := oauth2.Config{
	// 	ClientID:     clientID,
	// 	ClientSecret: clientSecret,
	// 	RedirectURL:  "",
	// 	Endpoint:     provider.Endpoint(),
	// 	Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	// }

	// verifier := provider.Verifier(&oidc.Config{ClientID: clientID})

	db, err := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Table("workspaces").AutoMigrate(&workspace.WorkspaceModel{})

	grpcClient, err := grpcclient.NewGrpcClient()
	if err != nil {
		log.Fatalln(err)
	}

	workspaceRepository := workspace.NewRepository(db)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	circleRepository := circle.NewRepository(grpcClient)
	circleUseCase := circle.NewUseCase(workspaceRepository, circleRepository)

	e := echo.New()
	e.Use(middleware.CORS())
	// authMiddleware := auth.NewAuthMiddleware(provider, oauth2Config, verifier)
	// e.Use(authMiddleware.Handle)
	workspaceHandler.NewEchohandler(e, workspaceUseCase)
	circleHandler.NewEchohandler(e, circleUseCase)
	e.Logger.Fatal(e.Start(":8080"))
}

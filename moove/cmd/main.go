package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/octopipe/charlescd/moove/internal/circle"
	circleHandler "github.com/octopipe/charlescd/moove/internal/circle/handler"
	"github.com/octopipe/charlescd/moove/internal/core/httpclient"
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

	db, err := gorm.Open(sqlite.Open("sample.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.Table("workspaces").AutoMigrate(&workspace.WorkspaceModel{})

	httpClient := httpclient.NewHttpClient()

	workspaceRepository := workspace.NewRepository(db)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	circleRepository := circle.NewRepository(httpClient)
	circleUseCase := circle.NewUseCase(circleRepository)

	e := echo.New()
	workspaceHandler.NewEchohandler(e, workspaceUseCase)
	circleHandler.NewEchohandler(e, circleUseCase)
	e.Logger.Fatal(e.Start(":3000"))
}

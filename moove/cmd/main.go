package main

import (
	"log"

	"github.com/labstack/echo/v4"
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

	workspaceRepository := workspace.NewRepository(db)
	workspaceUseCase := workspace.NewUseCase(workspaceRepository)

	e := echo.New()
	workspaceHandler.NewEchohandler(e, workspaceUseCase)
	e.Logger.Fatal(e.Start(":8080"))
}

package appserver

import (
	"github.com/Mrityunjoy99/sample-go/src/application"
	"github.com/Mrityunjoy99/sample-go/src/common/config"
	"github.com/Mrityunjoy99/sample-go/src/domain/service"
	"github.com/Mrityunjoy99/sample-go/src/infrastructure/database"
	"github.com/Mrityunjoy99/sample-go/src/repository"
	"github.com/Mrityunjoy99/sample-go/src/tools/logger"
	"github.com/gin-gonic/gin"
)

func Start() {
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(c)
	if err != nil {
		panic(err)
	}

	r := repository.NewRepository(db)
	domainService, gerr := service.NewServiceRegistry(c)
	if gerr != nil {
		panic(gerr.Error())
	}

	logger, err := logger.NewZapLogger(c.App.Name)
	if err != nil {
		panic(err)
	}

	appService, err := application.NewService(c, r, logger, domainService)
	if err != nil {
		panic(err.Error())
	}

	g := gin.Default()
	RegisterRoutes(g, logger, *appService, *domainService)

	err = g.Run()
	if err != nil {
		panic(err)
	}
}

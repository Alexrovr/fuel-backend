package api

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"heat-backend/internal/app/handler"
	"heat-backend/internal/app/repository"
	"heat-backend/internal/app/storage"
)

// StartServer собирает приложение и регистрирует маршруты.
func StartServer(logger *logrus.Logger) error {
	fuelRepository := repository.NewFuelRepository()
	minioMedia := storage.NewMinioMedia()
	fuelHandler := handler.NewFuelHandler(fuelRepository, minioMedia, logger)

	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"heatForVolume": func(heatOfCombustionKJ int, volumeM3 float64) float64 {
			return float64(heatOfCombustionKJ) * volumeM3
		},
	})
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/resources", "./resources")

	router.GET("/fuel_feed", fuelHandler.GetFuelFeed)
	router.GET("/fuel_feed/:fuel_id", fuelHandler.GetFuelFeed)
	router.GET("/fuel_draft", fuelHandler.GetFuelDraft)
	router.GET("/fuel_grid", fuelHandler.GetFuelGrid)

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/fuel_grid")
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3030"
	}

	logger.Infof("сервер расчёта теплоты сгорания запущен на http://localhost:%s", port)
	return router.Run(":" + port)
}

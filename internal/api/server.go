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

// StartServer собирает приложение: репозиторий -> обработчики -> роутинг.
func StartServer(logger *logrus.Logger) error {
	fuelRepository := repository.NewFuelRepository()
	minioMedia := storage.NewMinioMedia()
	fuelHandler := handler.NewFuelHandler(fuelRepository, minioMedia, logger)

	router := gin.Default()

	// Шаблонизатор. Функция heatForVolume считает количество теплоты,
	// выделившейся при полном сгорании заданного объёма топлива при н.у.
	router.SetFuncMap(template.FuncMap{
		"heatForVolume": func(heatOfCombustionKJ int, volumeM3 float64) float64 {
			return float64(heatOfCombustionKJ) * volumeM3
		},
	})
	router.LoadHTMLGlob("templates/*.html")

	// Стили и вспомогательная статика. Изображения и видео карточек
	// отдаются не отсюда, а из Minio.
	router.Static("/resources", "./resources")

	// Три GET-метода приложения
	router.GET("/fuel_feed", fuelHandler.GetFuelFeed)          // лента, открытая из панели вкладок
	router.GET("/fuel_feed/:fuel_id", fuelHandler.GetFuelFeed) // лента по ID (+ ?next=true)
	router.GET("/fuel_draft", fuelHandler.GetFuelDraft)        // получение черновика
	router.GET("/fuel_grid", fuelHandler.GetFuelGrid)          // список всех карточек (+ ?min_heat=)

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

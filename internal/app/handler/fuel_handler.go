package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"heat-backend/internal/app/models"
	"heat-backend/internal/app/repository"
	"heat-backend/internal/app/storage"
)

// FuelHandler — слой контроллеров. Обработчики обращаются к репозиторию,
// репозиторий к обработчикам — никогда (однонаправленная зависимость).
type FuelHandler struct {
	repository *repository.FuelRepository
	media      *storage.MinioMedia
	logger     *logrus.Logger
}

func NewFuelHandler(
	fuelRepository *repository.FuelRepository,
	media *storage.MinioMedia,
	logger *logrus.Logger,
) *FuelHandler {
	return &FuelHandler{
		repository: fuelRepository,
		media:      media,
		logger:     logger,
	}
}

// toCardView переносит карточку из коллекции в модель представления:
// здесь, в контроллере-обработчике, вычисляется количество лайков
// и по ключам Minio собираются url изображения и видео.
func (h *FuelHandler) toCardView(fuel models.Fuel) models.HeatCardView {
	return models.HeatCardView{
		FuelID:             fuel.FuelID,
		FuelName:           fuel.FuelName,
		ChemicalFormula:    fuel.ChemicalFormula,
		CombustionNote:     fuel.CombustionNote,
		HeatOfCombustionKJ: fuel.HeatOfCombustionKJ,
		IgnitionTempC:      fuel.IgnitionTempC,
		FlameTempC:         fuel.FlameTempC,
		AirDemandM3:        fuel.AirDemandM3,
		ImageKey:           fuel.ImageKey,
		VideoKey:           fuel.VideoKey,
		ImageURL:           h.media.ObjectURL(fuel.ImageKey),
		VideoURL:           h.media.ObjectURL(fuel.VideoKey),
		FuelStatus:         fuel.FuelStatus,
		LikesCount:         len(fuel.LikedByUserIDs),
	}
}

// GetFuelFeed — GET /fuel_feed и GET /fuel_feed/:fuel_id
//
// Лента горения в вертикальной развёртке. Идентификатор карточки в адресе
// не обязателен: из панели вкладок лента открывается без него и показывает
// первое опубликованное топливо. Параметр ?next=true открывает следующую
// карточку после указанного идентификатора.
func (h *FuelHandler) GetFuelFeed(ctx *gin.Context) {
	fuelIDParam := ctx.Param("fuel_id")
	nextParam := ctx.Query("next")

	var (
		fuel models.Fuel
		err  error
	)

	switch {
	case fuelIDParam == "":
		fuel, err = h.repository.GetFirstPublishedFuel()
	default:
		fuelID, convErr := strconv.Atoi(fuelIDParam)
		if convErr != nil {
			h.logger.Warnf("некорректный идентификатор топлива в url: %q", fuelIDParam)
			h.renderFeedError(ctx, http.StatusBadRequest, "Некорректный идентификатор вида топлива")
			return
		}
		if nextParam == "true" {
			fuel, err = h.repository.GetNextPublishedFuel(fuelID)
		} else {
			fuel, err = h.repository.GetPublishedFuelByID(fuelID)
		}
	}

	if err != nil {
		h.logger.Warnf("лента горения: %v (id=%q, next=%q)", err, fuelIDParam, nextParam)
		h.renderFeedError(ctx, http.StatusNotFound, "Такого вида топлива нет в справочнике")
		return
	}

	card := h.toCardView(fuel)
	ctx.HTML(http.StatusOK, "fuel_feed.html", gin.H{
		"PageTitle":  "Лента горения",
		"ActiveTab":  "feed",
		"HeatCard":   card,
		"HasCard":    true,
		"NextFuelID": card.FuelID,
	})
}

func (h *FuelHandler) renderFeedError(ctx *gin.Context, status int, message string) {
	ctx.HTML(status, "fuel_feed.html", gin.H{
		"PageTitle":    "Лента горения",
		"ActiveTab":    "feed",
		"HasCard":      false,
		"FeedErrorMsg": message,
	})
}

// GetFuelDraft — GET /fuel_draft
//
// Страница добавления карточки топлива. В первой лабораторной сохранение
// не реализуется: форма только показывает поля единственного черновика
// из коллекции, разнесённые по отдельным блокам.
func (h *FuelHandler) GetFuelDraft(ctx *gin.Context) {
	draft, err := h.repository.GetDraftFuel()
	if err != nil {
		h.logger.Warnf("страница добавления: %v", err)
		ctx.HTML(http.StatusOK, "fuel_draft.html", gin.H{
			"PageTitle": "Добавление топлива",
			"ActiveTab": "draft",
			"HasDraft":  false,
		})
		return
	}

	ctx.HTML(http.StatusOK, "fuel_draft.html", gin.H{
		"PageTitle": "Добавление топлива",
		"ActiveTab": "draft",
		"HasDraft":  true,
		"HeatCard":  h.toCardView(draft),
	})
}

// GetFuelGrid — GET /fuel_grid?min_heat=...
//
// Плитка карточек топлива в два столбца. Фильтрация по теме выполняется
// на сервере по числовому полю «теплота сгорания, кДж/м³ при н.у.»,
// введённое значение возвращается в шаблон и остаётся в поле ввода.
func (h *FuelHandler) GetFuelGrid(ctx *gin.Context) {
	minHeatQuery := ctx.Query("min_heat")

	minHeatKJ := 0
	filterErrorMsg := ""
	if minHeatQuery != "" {
		parsed, err := strconv.Atoi(minHeatQuery)
		switch {
		case err != nil:
			filterErrorMsg = "Теплота сгорания задаётся целым числом в кДж/м³"
		case parsed < 0:
			filterErrorMsg = "Теплота сгорания не может быть отрицательной"
		default:
			minHeatKJ = parsed
		}
	}

	fuels := h.repository.GetPublishedFuels(minHeatKJ)
	cards := make([]models.HeatCardView, 0, len(fuels))
	for _, fuel := range fuels {
		cards = append(cards, h.toCardView(fuel))
	}

	h.logger.Infof("плитка топлив: min_heat=%q, найдено карточек: %d", minHeatQuery, len(cards))

	ctx.HTML(http.StatusOK, "fuel_grid.html", gin.H{
		"PageTitle":      "Виды топлива",
		"ActiveTab":      "grid",
		"HeatCards":      cards,
		"MinHeatFilter":  minHeatQuery,
		"FilterErrorMsg": filterErrorMsg,
	})
}

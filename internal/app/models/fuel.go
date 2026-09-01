package models

// FuelStatus — статус карточки топлива в коллекции.
// По ТЗ лабораторной у каждой услуги ровно один из трёх статусов.
type FuelStatus string

const (
	FuelStatusDraft     FuelStatus = "черновик"
	FuelStatusPublished FuelStatus = "опубликован"
	FuelStatusDeleted   FuelStatus = "удален"
)

// Fuel — услуга предметной области «Расчёт количества теплоты в реакции горения».
// Одна карточка = один вид топлива, для которого известна теплота сгорания
// при нормальных условиях (н.у.: 0 °C, 101 325 Па).
//
// Все поля по предметной области атомарны (1НФ): только числа, даты и текст.
// Единственный вложенный массив — идентификаторы пользователей, поставивших лайк
// (связь м-м «пользователь — топливо», которая в ЛР №2 станет отдельной таблицей).
type Fuel struct {
	FuelID int

	FuelName        string // Метан
	ChemicalFormula string // CH4
	CombustionNote  string // краткое описание реакции горения

	// Поля по предметной области
	HeatOfCombustionKJ int     // низшая теплота сгорания, кДж/м³ при н.у.
	IgnitionTempC      int     // температура воспламенения, °C
	FlameTempC         int     // температура пламени в воздухе, °C
	AirDemandM3        float64 // теоретический расход воздуха, м³ на 1 м³ топлива

	// Ключи объектов в Minio (латиница), хранятся двумя отдельными полями
	ImageKey string
	VideoKey string

	FuelStatus     FuelStatus
	LikedByUserIDs []int // м-м: кто лайкнул карточку
}

// HeatCardView — модель представления одной карточки для шаблонов.
// Собирается в контроллере-обработчике: там же вычисляется количество лайков
// и абсолютные url изображения и видео в Minio.
type HeatCardView struct {
	FuelID int

	FuelName        string
	ChemicalFormula string
	CombustionNote  string

	HeatOfCombustionKJ int
	IgnitionTempC      int
	FlameTempC         int
	AirDemandM3        float64

	ImageKey string // ключ объекта в Minio, как он лежит в коллекции
	VideoKey string
	ImageURL string // абсолютный адрес объекта, собранный обработчиком
	VideoURL string

	FuelStatus FuelStatus
	LikesCount int // вычисляется по коллекции в обработчике
}

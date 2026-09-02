package models

type FuelStatus string

const (
	FuelStatusDraft     FuelStatus = "черновик"
	FuelStatusPublished FuelStatus = "опубликован"
	FuelStatusDeleted   FuelStatus = "удален"
)

type Fuel struct {
	FuelID int

	FuelName        string // Метан
	ChemicalFormula string // CH4
	CombustionNote  string // краткое описание реакции горения

	HeatOfCombustionKJ int
	IgnitionTempC      int
	FlameTempC         int
	AirDemandM3        float64

	ImageKey string
	VideoKey string

	FuelStatus     FuelStatus
	LikedByUserIDs []int
}

type HeatCardView struct {
	FuelID int

	FuelName        string
	ChemicalFormula string
	CombustionNote  string

	HeatOfCombustionKJ int
	IgnitionTempC      int
	FlameTempC         int
	AirDemandM3        float64

	ImageKey string
	VideoKey string
	ImageURL string
	VideoURL string

	FuelStatus FuelStatus
	LikesCount int
}

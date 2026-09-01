package repository

import (
	"errors"

	"heat-backend/internal/app/models"
)

// ErrFuelNotFound возвращается, когда карточка топлива отсутствует в коллекции
// или скрыта от интерфейса статусом «удален».
var ErrFuelNotFound = errors.New("вид топлива не найден")

// ErrDraftNotFound возвращается, когда у пользователя нет карточки-черновика.
var ErrDraftNotFound = errors.New("черновик карточки топлива не найден")

// FuelRepository — слой доступа к данным. В первой лабораторной база данных
// не используется: единственная модель-коллекция хранится в памяти процесса.
type FuelRepository struct {
	fuels []models.Fuel
}

// NewFuelRepository наполняет коллекцию реальными справочными данными
// по теплоте сгорания газообразных топлив при нормальных условиях.
func NewFuelRepository() *FuelRepository {
	return &FuelRepository{
		fuels: []models.Fuel{
			{
				FuelID:          1,
				FuelName:        "Метан",
				ChemicalFormula: "CH4",
				CombustionNote: "Основной компонент природного газа. Полное сгорание идёт по уравнению " +
					"CH4 + 2O2 -> CO2 + 2H2O. При нормальных условиях один кубометр метана отдаёт " +
					"35 800 кДж теплоты, поэтому метан используют как эталон при расчёте тепловой мощности " +
					"бытовых котлов и промышленных горелок. Горит спокойным голубым пламенем без копоти, " +
					"что говорит о полном окислении углерода до диоксида.",
				HeatOfCombustionKJ: 35800,
				IgnitionTempC:      537,
				FlameTempC:         1957,
				AirDemandM3:        9.52,
				ImageKey:           "methane.jpg",
				VideoKey:           "methane.mp4",
				FuelStatus:         models.FuelStatusPublished,
				LikedByUserIDs:     []int{1, 2, 3, 5, 8, 13, 21},
			},
			{
				FuelID:          2,
				FuelName:        "Пропан-бутан",
				ChemicalFormula: "C3H8 + C4H10",
				CombustionNote: "Сжиженный углеводородный газ, смесь пропана и бутана в соотношении 50/50. " +
					"Реакция горения пропана: C3H8 + 5O2 -> 3CO2 + 4H2O. Кубометр смеси при н.у. выделяет " +
					"около 108 000 кДж — втрое больше метана, поэтому баллонный газ применяют там, где " +
					"нужна высокая тепловая мощность при малом объёме хранения. Требует втрое большего " +
					"притока воздуха, иначе горение становится неполным и появляется сажа.",
				HeatOfCombustionKJ: 108000,
				IgnitionTempC:      470,
				FlameTempC:         1970,
				AirDemandM3:        27.37,
				ImageKey:           "propane_butane.jpg",
				VideoKey:           "propane_butane.mp4",
				FuelStatus:         models.FuelStatusPublished,
				LikedByUserIDs:     []int{2, 4, 7, 9, 11, 14, 16, 19, 23},
			},
			{
				FuelID:          3,
				FuelName:        "Ацетилен",
				ChemicalFormula: "C2H2",
				CombustionNote: "Топливо газовой сварки и резки металлов. Полное сгорание: " +
					"2C2H2 + 5O2 -> 4CO2 + 2H2O. Кубометр ацетилена при н.у. даёт 56 000 кДж, но " +
					"главное его достоинство не в теплоте, а в температуре пламени: в смеси с кислородом " +
					"оно достигает 3150 °C — выше, чем у любого другого промышленного газа. " +
					"Самая низкая температура воспламенения в подборке, всего 335 °C.",
				HeatOfCombustionKJ: 56000,
				IgnitionTempC:      335,
				FlameTempC:         3150,
				AirDemandM3:        11.91,
				ImageKey:           "acetylene.jpg",
				VideoKey:           "acetylene.mp4",
				FuelStatus:         models.FuelStatusPublished,
				LikedByUserIDs:     []int{1, 6, 10, 12, 18},
			},
			{
				FuelID:          4,
				FuelName:        "Водород",
				ChemicalFormula: "H2",
				CombustionNote: "Единственное топливо подборки, при сгорании которого не образуется " +
					"диоксид углерода: 2H2 + O2 -> 2H2O. Объёмная теплота сгорания самая низкая — " +
					"10 800 кДж/м³ при н.у., потому что молекула водорода очень лёгкая. " +
					"Зато на единицу массы водород выделяет 120 000 кДж/кг, вне конкуренции среди " +
					"всех химических топлив. Пламя почти бесцветное и требует всего 2,38 м³ воздуха.",
				HeatOfCombustionKJ: 10800,
				IgnitionTempC:      510,
				FlameTempC:         2130,
				AirDemandM3:        2.38,
				ImageKey:           "hydrogen.jpg",
				VideoKey:           "hydrogen.mp4",
				FuelStatus:         models.FuelStatusPublished,
				LikedByUserIDs:     []int{3, 5, 15},
			},
			{
				FuelID:          5,
				FuelName:        "Бутан",
				ChemicalFormula: "C4H10",
				CombustionNote: "Чистый бутан выведен из подборки: при температуре ниже нуля он " +
					"перестаёт испаряться из баллона, поэтому в справочнике оставлена только " +
					"пропан-бутановая смесь. Карточка помечена статусом «удален» и в интерфейсе " +
					"не отображается.",
				HeatOfCombustionKJ: 123770,
				IgnitionTempC:      405,
				FlameTempC:         1970,
				AirDemandM3:        30.94,
				ImageKey:           "butane.jpg",
				VideoKey:           "butane.mp4",
				FuelStatus:         models.FuelStatusDeleted,
				LikedByUserIDs:     []int{4},
			},
			{
				FuelID:          6,
				FuelName:        "Метано-водородная смесь",
				ChemicalFormula: "CH4 + H2",
				CombustionNote: "Черновик карточки: смесь природного газа с 20 % водорода, которую " +
					"испытывают как переходное топливо для действующих газовых сетей. Справочные " +
					"значения ещё уточняются, поэтому карточка не опубликована и доступна только " +
					"на странице добавления.",
				HeatOfCombustionKJ: 30800,
				IgnitionTempC:      520,
				FlameTempC:         2000,
				AirDemandM3:        8.09,
				ImageKey:           "methane_hydrogen.jpg",
				VideoKey:           "methane_hydrogen.mp4",
				FuelStatus:         models.FuelStatusDraft,
				LikedByUserIDs:     []int{},
			},
		},
	}
}

// GetPublishedFuels возвращает опубликованные карточки топлива.
// Фильтрация по теме выполняется на сервере: minHeatKJ — нижняя граница
// теплоты сгорания в кДж/м³ при н.у. Значение 0 означает «без фильтра».
func (r *FuelRepository) GetPublishedFuels(minHeatKJ int) []models.Fuel {
	result := make([]models.Fuel, 0, len(r.fuels))
	for _, fuel := range r.fuels {
		if fuel.FuelStatus != models.FuelStatusPublished {
			continue
		}
		if fuel.HeatOfCombustionKJ < minHeatKJ {
			continue
		}
		result = append(result, fuel)
	}
	return result
}

// GetPublishedFuelByID отдаёт одну опубликованную карточку по идентификатору.
func (r *FuelRepository) GetPublishedFuelByID(fuelID int) (models.Fuel, error) {
	for _, fuel := range r.fuels {
		if fuel.FuelID == fuelID && fuel.FuelStatus == models.FuelStatusPublished {
			return fuel, nil
		}
	}
	return models.Fuel{}, ErrFuelNotFound
}

// GetNextPublishedFuel возвращает следующую опубликованную карточку после
// указанной. Список закольцован: после последней карточки идёт первая.
func (r *FuelRepository) GetNextPublishedFuel(fuelID int) (models.Fuel, error) {
	published := r.GetPublishedFuels(0)
	if len(published) == 0 {
		return models.Fuel{}, ErrFuelNotFound
	}
	for i, fuel := range published {
		if fuel.FuelID == fuelID {
			return published[(i+1)%len(published)], nil
		}
	}
	return models.Fuel{}, ErrFuelNotFound
}

// GetFirstPublishedFuel нужен ленте, открытой из панели вкладок,
// когда идентификатор карточки в url не указан.
func (r *FuelRepository) GetFirstPublishedFuel() (models.Fuel, error) {
	published := r.GetPublishedFuels(0)
	if len(published) == 0 {
		return models.Fuel{}, ErrFuelNotFound
	}
	return published[0], nil
}

// GetDraftFuel отдаёт единственную карточку в статусе «черновик»,
// которая отображается на странице добавления.
func (r *FuelRepository) GetDraftFuel() (models.Fuel, error) {
	for _, fuel := range r.fuels {
		if fuel.FuelStatus == models.FuelStatusDraft {
			return fuel, nil
		}
	}
	return models.Fuel{}, ErrDraftNotFound
}

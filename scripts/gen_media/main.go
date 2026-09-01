// Генератор медиафайлов карточек топлива.
//
// Для каждого вида топлива рисуется процедурное пламя своего цвета:
//   - постер   resources/media/<key>.jpg      — используется в плитке,
//     в форме добавления и как
//     poster видео в ленте;
//   - кадры    build/frames/<key>/f_%03d.png  — исходники короткого ролика,
//     из них ffmpeg собирает mp4.
//
// Запуск:  go run ./scripts/gen_media
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
)

const (
	frameWidth  = 540
	frameHeight = 960
	frameCount  = 72 // 3 секунды при 24 кадрах в секунду
)

// flameProfile описывает внешний вид пламени одного вида топлива.
type flameProfile struct {
	MediaKey string // methane, propane_butane, ...
	Core     [3]float64
	Mid      [3]float64
	Edge     [3]float64
	Width    float64 // относительная ширина факела
	Height   float64 // относительная высота факела
	Turb     float64 // интенсивность турбулентности
	Sooty    bool    // коптящее пламя: тёмный дымный шлейф сверху
}

// Цвета подобраны по реальному виду пламени каждого топлива.
var flameProfiles = []flameProfile{
	{
		// Метан горит спокойным голубым пламенем без копоти
		MediaKey: "methane",
		Core:     [3]float64{235, 250, 255},
		Mid:      [3]float64{70, 165, 255},
		Edge:     [3]float64{20, 55, 170},
		Width:    0.30, Height: 0.62, Turb: 0.55,
	},
	{
		// Пропан-бутан даёт мощный оранжевый факел
		MediaKey: "propane_butane",
		Core:     [3]float64{255, 246, 210},
		Mid:      [3]float64{253, 43, 6},
		Edge:     [3]float64{130, 25, 5},
		Width:    0.36, Height: 0.74, Turb: 0.85,
		Sooty: true,
	},
	{
		// Ацетилен: узкое ослепительно яркое ядро сварочного пламени
		MediaKey: "acetylene",
		Core:     [3]float64{255, 255, 255},
		Mid:      [3]float64{255, 200, 120},
		Edge:     [3]float64{190, 70, 15},
		Width:    0.22, Height: 0.80, Turb: 0.40,
		Sooty: true,
	},
	{
		// Водород горит почти бесцветным бледно-фиолетовым пламенем
		MediaKey: "hydrogen",
		Core:     [3]float64{250, 245, 255},
		Mid:      [3]float64{170, 160, 240},
		Edge:     [3]float64{60, 45, 130},
		Width:    0.26, Height: 0.55, Turb: 0.65,
	},
	{
		// Бутан (карточка со статусом «удален»)
		MediaKey: "butane",
		Core:     [3]float64{255, 240, 195},
		Mid:      [3]float64{247, 153, 0},
		Edge:     [3]float64{140, 45, 5},
		Width:    0.34, Height: 0.70, Turb: 0.80,
		Sooty: true,
	},
	{
		// Метано-водородная смесь (карточка-черновик)
		MediaKey: "methane_hydrogen",
		Core:     [3]float64{240, 253, 255},
		Mid:      [3]float64{80, 210, 235},
		Edge:     [3]float64{20, 80, 150},
		Width:    0.28, Height: 0.66, Turb: 0.60,
	},
}

func main() {
	mediaDir := filepath.Join("resources", "media")
	framesRoot := filepath.Join("build", "frames")

	if err := os.MkdirAll(mediaDir, 0o755); err != nil {
		panic(err)
	}

	for _, profile := range flameProfiles {
		// Постер: кадр из середины анимации.
		poster := renderFrame(profile, 0.5)
		posterPath := filepath.Join(mediaDir, profile.MediaKey+".jpg")
		if err := writeJPEG(posterPath, poster); err != nil {
			panic(err)
		}

		// Кадры ролика.
		frameDir := filepath.Join(framesRoot, profile.MediaKey)
		if err := os.MkdirAll(frameDir, 0o755); err != nil {
			panic(err)
		}
		for i := 0; i < frameCount; i++ {
			phase := float64(i) / float64(frameCount)
			framePath := filepath.Join(frameDir, fmt.Sprintf("f_%03d.png", i))
			if err := writePNG(framePath, renderFrame(profile, phase)); err != nil {
				panic(err)
			}
		}

		fmt.Printf("%-18s постер %s и %d кадров в %s\n",
			profile.MediaKey, posterPath, frameCount, frameDir)
	}
}

// renderFrame рисует один кадр пламени. phase меняется от 0 до 1 и
// подобран так, чтобы кадр phase=1 совпадал с кадром phase=0 —
// тогда ролик зацикливается без рывка.
func renderFrame(profile flameProfile, phase float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, frameWidth, frameHeight))
	t := phase * 2 * math.Pi

	burnerTop := float64(frameHeight) * 0.86
	flameBase := burnerTop
	flameTop := burnerTop - float64(frameHeight)*profile.Height

	for y := 0; y < frameHeight; y++ {
		fy := float64(y)

		// Фон: тёмная лабораторная камера с лёгким подсветом снизу.
		bgGlow := math.Max(0, 1-math.Abs(fy-flameBase)/(float64(frameHeight)*0.55))
		bgGlow = math.Pow(bgGlow, 3) * 0.30

		for x := 0; x < frameWidth; x++ {
			fx := float64(x)

			r := 12 + profile.Mid[0]*bgGlow*0.35
			g := 12 + profile.Mid[1]*bgGlow*0.35
			b := 16 + profile.Mid[2]*bgGlow*0.35

			// Горелка в нижней части кадра.
			if fy >= burnerTop {
				nozzleHalf := float64(frameWidth) * 0.055
				depth := (fy - burnerTop) / (float64(frameHeight) - burnerTop)
				if math.Abs(fx-float64(frameWidth)/2) < nozzleHalf+depth*float64(frameWidth)*0.10 {
					shade := 58.0 - depth*26.0
					r, g, b = shade, shade*0.97, shade*0.95
				}
			}

			// Интенсивность пламени в точке.
			intensity := flameIntensity(profile, fx, fy, flameBase, flameTop, t)
			if intensity > 0 {
				fr, fg, fb := flameColor(profile, intensity)
				r = r*(1-math.Min(1, intensity)) + fr*math.Min(1, intensity)
				g = g*(1-math.Min(1, intensity)) + fg*math.Min(1, intensity)
				b = b*(1-math.Min(1, intensity)) + fb*math.Min(1, intensity)
			}

			// Коптящий шлейф над факелом.
			if profile.Sooty && fy < flameTop {
				smoke := smokeDensity(fx, fy, flameTop, t)
				r = r*(1-smoke) + 46*smoke
				g = g*(1-smoke) + 42*smoke
				b = b*(1-smoke) + 40*smoke
			}

			img.Set(x, y, color.RGBA{
				R: clampByte(r),
				G: clampByte(g),
				B: clampByte(b),
				A: 255,
			})
		}
	}

	return img
}

// flameIntensity возвращает яркость факела в точке: 0 — вне пламени,
// 1 и выше — ядро. Форма факела задаётся конусом, который искажается
// суммой синусоид (дешёвая замена шуму Перлина).
func flameIntensity(profile flameProfile, fx, fy, base, top, t float64) float64 {
	if fy > base {
		return 0
	}

	height := base - top
	rise := (base - fy) / height // 0 у горелки, 1 у вершины факела
	if rise < 0 || rise > 1.15 {
		return 0
	}

	// Полуширина факела: расширяется у основания, сходится к вершине.
	halfWidth := float64(frameWidth) * profile.Width *
		math.Sin(math.Pi*math.Pow(math.Min(rise, 1), 0.75)) * 0.9
	halfWidth += float64(frameWidth) * 0.035

	// Колебание оси факела и «языки» пламени.
	sway := profile.Turb * float64(frameWidth) * 0.055 *
		(math.Sin(t+rise*5.1) + 0.45*math.Sin(2*t+rise*9.4) + 0.25*math.Sin(3*t+rise*15.7))
	axis := float64(frameWidth)/2 + sway*rise

	tongue := profile.Turb * 0.22 *
		(math.Sin(2*t+rise*13.0) + 0.6*math.Sin(3*t+rise*21.0))
	halfWidth *= 1 + tongue*rise

	if halfWidth <= 1 {
		return 0
	}

	// Поперечный профиль: ядро в центре, мягкий спад к краям.
	across := math.Abs(fx-axis) / halfWidth
	if across >= 1 {
		return 0
	}
	radial := math.Pow(1-across*across, 1.7)

	// Продольный профиль: максимум чуть выше горелки, затухание к вершине.
	vertical := math.Sin(math.Pi*math.Pow(math.Min(rise, 1), 0.85)) * 1.15
	if rise > 1 {
		vertical *= math.Max(0, 1-(rise-1)/0.15)
	}

	// Мерцание.
	flicker := 1 + 0.10*math.Sin(4*t+rise*7.0)

	return radial * vertical * flicker
}

// flameColor переводит яркость в цвет: край -> середина -> ядро.
func flameColor(profile flameProfile, intensity float64) (float64, float64, float64) {
	i := math.Min(intensity, 1.35)
	switch {
	case i < 0.55:
		k := i / 0.55
		return mix(profile.Edge, profile.Mid, k)
	default:
		k := math.Min(1, (i-0.55)/0.70)
		return mix(profile.Mid, profile.Core, k)
	}
}

func mix(from, to [3]float64, k float64) (float64, float64, float64) {
	return from[0] + (to[0]-from[0])*k,
		from[1] + (to[1]-from[1])*k,
		from[2] + (to[2]-from[2])*k
}

// smokeDensity рисует дымный шлейф над коптящим факелом.
func smokeDensity(fx, fy, flameTop, t float64) float64 {
	above := (flameTop - fy) / (flameTop + 1)
	if above <= 0 {
		return 0
	}
	drift := 34*math.Sin(t+above*4.0) + 16*math.Sin(2*t+above*7.5)
	axis := float64(frameWidth)/2 + drift
	spread := float64(frameWidth)*0.10 + above*float64(frameWidth)*0.30
	across := math.Abs(fx-axis) / spread
	if across >= 1 {
		return 0
	}
	return math.Pow(1-across, 2) * math.Max(0, 0.45-above*0.55)
}

func clampByte(v float64) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return uint8(v)
	}
}

func writeJPEG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, img, &jpeg.Options{Quality: 88})
}

func writePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

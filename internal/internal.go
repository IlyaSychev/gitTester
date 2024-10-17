package internal

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// 6 hours -> seconds
const Cycle = 60 * 60 * 6

type Point struct {
	x                float64
	y                float64
	h                float64
	Path             [][]float64
	EnterAndExitTime map[int]float64
	PercentTimeEnter float64
}

func NewPoint(startX, startY, startH float64) *Point {
	tmp := make([][]float64, 0)
	tmp = append(tmp, []float64{startX, startY, startH})
	return &Point{startX, startY, startH, tmp, map[int]float64{}, 0}
}

func (aerostat *Point) Polinom(VelocitiesBefore, VelocitiesAfter map[[3]float64][3]float64, dt float64, i float64) {

	arrH := []float64{200, 150, 100, 70, 50, 40, 30, 20, 10, 7, 5, 3, 2}

	//fmt.Println("coordinates x and y = ", aerostat.x, aerostat.y)

	x1, x2 := round(aerostat.x)
	y1, y2 := round(aerostat.y)
	h1, h2 := roundH(aerostat.h, arrH)

	//fmt.Println("local sc x1 x2 = ", x1, x2)
	//fmt.Println("local sc y1 y2 = ", y1, y2)
	//fmt.Println("local sc h1 h2 = ", h1, h2)

	xCenter := (x2 + x1) / 2
	yCenter := (y2 + y1) / 2
	hCenter := (h1 + h2) / 2

	eta := aerostat.x - xCenter
	eta = eta / (x2 - x1) * 2
	ksi := aerostat.y - yCenter
	ksi = ksi / (y2 - y1) * 2
	nu := aerostat.h - hCenter
	nu = nu / (h2 - h1) * 2

	//fmt.Println("local eta1 = ", eta)
	//fmt.Println("local ksi = ", ksi)
	//fmt.Println("local nu = ", nu)
	//fmt.Println()

	fnForm1 := (ksi + 1) * (eta - 1) * (nu - 1) / 8
	fnForm2 := (-1) * 1 * (ksi + 1) * (eta + 1) * (nu - 1) / 8
	fnForm3 := (-1) * 1 * (ksi - 1) * (eta - 1) * (nu - 1) / 8
	fnForm4 := (ksi - 1) * (eta + 1) * (nu - 1) / 8
	fnForm5 := (-1) * 1 * (ksi + 1) * (eta - 1) * (nu + 1) / 8
	fnForm6 := (ksi + 1) * (eta + 1) * (nu + 1) / 8
	fnForm7 := (ksi - 1) * (eta - 1) * (nu + 1) / 8
	fnForm8 := (-1) * (ksi - 1) * (eta + 1) * (nu + 1) / 8

	if x2 == 360 {
		x2 = 0
	}

	//fmt.Println("fnForms: ", fnForm1, fnForm2, fnForm3, fnForm4, fnForm5, fnForm6, fnForm7, fnForm8)

	V1u, V2u, V3u, V4u, V5u, V6u, V7u, V8u := FindApproxTimeVelocity(VelocitiesBefore, VelocitiesAfter, dt, i, x1, x2, y1, y2, h1, h2, 0)
	V1v, V2v, V3v, V4v, V5v, V6v, V7v, V8v := FindApproxTimeVelocity(VelocitiesBefore, VelocitiesAfter, dt, i, x1, x2, y1, y2, h1, h2, 1)

	u := fnForm1*V1u + fnForm2*V2u + fnForm3*V3u + fnForm4*V4u + fnForm5*V5u + fnForm6*V6u + fnForm7*V7u + fnForm8*V8u
	v := fnForm1*V1v + fnForm2*V2v + fnForm3*V3v + fnForm4*V4v + fnForm5*V5v + fnForm6*V6v + fnForm7*V7v + fnForm8*V8v
	//hvel := 0
	//hvel := fnForm1*Velocities[[3]float64{x1, y2, h1}][2] + fnForm2*Velocities[[3]float64{x2, y2, h1}][2] + fnForm3*Velocities[[3]float64{x1, y1, h1}][2] + fnForm4*Velocities[[3]float64{x2, y1, h1}][2] + fnForm5*Velocities[[3]float64{x1, y2, h2}][2] + fnForm6*Velocities[[3]float64{x2, y2, h2}][2] + fnForm7*Velocities[[3]float64{x1, y1, h2}][2] + fnForm8*Velocities[[3]float64{x2, y1, h2}][2]
	//hvel /= 100 // pa/s -> mb/s

	//fmt.Println("u, v , hvel = ", u, v, hvel)

	aerostat.x = aerostat.x + MetresToDegrees(aerostat.y, u*dt)
	aerostat.y = aerostat.y + MetresToDegrees(0, v*dt)
	//aerostat.h = aerostat.h + hvel*dt

	if aerostat.y > 90 {
		aerostat.y = 90 - math.Mod(aerostat.y, 90)
		aerostat.x = aerostat.x + 180
	}
	if aerostat.y < -90 {
		aerostat.y = -90 - math.Mod(aerostat.y, 90)
		aerostat.x = aerostat.x + 180
	}
	if aerostat.x >= 360 {
		aerostat.x = math.Mod(aerostat.x, 360)
	}
	if aerostat.x < 0 {
		aerostat.x = 360 + aerostat.x
	}

	aerostat.Path = append(aerostat.Path, []float64{aerostat.x, aerostat.y, aerostat.h})
	//fmt.Println(aerostat.x, aerostat.y)
}

// функция округления
func round(coordinate float64) (float64, float64) {
	coordinateMin := float64(math.Floor(float64(coordinate)))
	coordinateMax := coordinateMin + 1
	return coordinateMin, coordinateMax
}

func round05(coordinate float64) (float64, float64) {
	coordinateMin := math.Floor(coordinate*2) / 2
	coordinateMax := coordinateMin + 0.5
	return coordinateMin, coordinateMax
}

func MetresToDegrees(y, S float64) float64 {
	return S / (math.Pi * 6378000 / 180 * math.Cos(y/180*math.Pi))
}

func roundH(h float64, arrH []float64) (float64, float64) {
	for i := 0; i < len(arrH)-1; i++ {
		if arrH[i] >= h && arrH[i+1] <= h {
			return arrH[i+1], arrH[i]
		}
	}
	return 0, 0
}

func FindApproxTimeVelocity(VelocitiesBefore, VelocitiesAfter map[[3]float64][3]float64, dt float64, i float64, x1, x2, y1, y2, h1, h2 float64, component int) (float64, float64, float64, float64, float64, float64, float64, float64) {
	V1 := VelocitiesBefore[[3]float64{x1, y2, h1}][component] + (VelocitiesAfter[[3]float64{x1, y2, h1}][component]-VelocitiesBefore[[3]float64{x1, y2, h1}][component])/Cycle*dt*i
	V2 := VelocitiesBefore[[3]float64{x2, y2, h1}][component] + (VelocitiesAfter[[3]float64{x2, y2, h1}][component]-VelocitiesBefore[[3]float64{x2, y2, h1}][component])/Cycle*dt*i
	V3 := VelocitiesBefore[[3]float64{x1, y1, h1}][component] + (VelocitiesAfter[[3]float64{x1, y1, h1}][component]-VelocitiesBefore[[3]float64{x1, y1, h1}][component])/Cycle*dt*i
	V4 := VelocitiesBefore[[3]float64{x2, y1, h1}][component] + (VelocitiesAfter[[3]float64{x2, y1, h1}][component]-VelocitiesBefore[[3]float64{x2, y1, h1}][component])/Cycle*dt*i
	V5 := VelocitiesBefore[[3]float64{x1, y2, h2}][component] + (VelocitiesAfter[[3]float64{x1, y2, h2}][component]-VelocitiesBefore[[3]float64{x1, y2, h2}][component])/Cycle*dt*i
	V6 := VelocitiesBefore[[3]float64{x2, y2, h2}][component] + (VelocitiesAfter[[3]float64{x2, y2, h2}][component]-VelocitiesBefore[[3]float64{x2, y2, h2}][component])/Cycle*dt*i
	V7 := VelocitiesBefore[[3]float64{x1, y1, h2}][component] + (VelocitiesAfter[[3]float64{x1, y1, h2}][component]-VelocitiesBefore[[3]float64{x1, y1, h2}][component])/Cycle*dt*i
	V8 := VelocitiesBefore[[3]float64{x2, y1, h2}][component] + (VelocitiesAfter[[3]float64{x2, y1, h2}][component]-VelocitiesBefore[[3]float64{x2, y1, h2}][component])/Cycle*dt*i
	//fmt.Println("VxBefore VxAfter i", VelocitiesBefore[[3]float64{x1, y2, h1}][component], VelocitiesAfter[[3]float64{x1, y2, h1}][component], i, V1)
	return V1, V2, V3, V4, V5, V6, V7, V8
}

func CreateImage(arr [][]float64) error {
	// Открыть исходный файл
	f, err := os.Open("grid_map2.png")
	if err != nil {
		return err
	}
	defer f.Close()

	// Декодировать исходное изображение
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	// Создать копию изображения
	newImg := image.NewRGBA(img.Bounds())
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			newImg.Set(x, y, img.At(x, y))
		}
	}

	// Нарисовать точки на копии изображения
	for i := range arr {
		x1 := int(arr[i][0] * 3)
		y1 := int(math.Abs(arr[i][1]*3 - 270))
		newImg.Set(x1, y1, color.RGBA{255, 0, 0, 255}) // Рисуем красную точку
	}

	// Сохранить копию изображения в новый файл
	out, err := os.Create("results.png")
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, newImg)
}

func (p *Point) GetCoordinates() (float64, float64, float64) {
	return p.x, p.y, p.h
}

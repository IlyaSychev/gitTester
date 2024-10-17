package main

import (
	"fmt"
	"gitTester/internal"
	"gitTester/output"
	"gitTester/parsejson"
	"gitTester/smp"
	"strconv"
)

const (
	t0   = 0.0
	dt   = 1
	tMax = 31 * 24 * 60 * 60

	sep = 100
	x1  = 33 + 6.87*0
	x2  = 33 + 6.87*1
	x3  = 33 + 6.87*2
	x4  = 33 + 6.87*3
	x5  = 33 + 6.87*4
	x6  = 33 + 6.87*5
	x7  = 33 + 6.87*6
	x8  = 33 + 6.87*7
	x9  = 33 + 6.87*8
	x10 = 33 + 6.87*9
	x11 = 33 + 6.87*10
	x12 = 33 + 6.87*11
	x13 = 33 + 6.87*12
	x14 = 33 + 6.87*13
	x15 = 33 + 6.87*14
	x16 = 33 + 6.87*15
	x17 = 33 + 6.87*16
	x18 = 33 + 6.87*17
	x19 = 33 + 6.87*18
	x20 = 33 + 6.87*19
	x21 = 33 + 6.87*20
	x22 = 33 + 6.87*21
	x23 = 33 + 6.87*22
	x24 = 33 + 6.87*23

	//y0 = 65
	y1 = 70
	//y1 = 72.5
	y2 = 80
)

//var arrCoordX = [24]float64{x1, x2, x3, x4, x5, x6, x7, x8, x9, x10, x11, x12, x13, x14, x15, x16, x17, x18, x19, x20, x21, x22, x23, x24} //для 48 аэростатов

var arrCoordX = [12]float64{x1, x3, x5, x7, x9, x11, x13, x15, x17, x19, x21, x23} // для 24 аэростатов

func main() {
	namesArr := parsejson.ReadDir()
	fmt.Println(namesArr)
	//aerostat1 := internal.NewPoint(30, 70, 21.88)
	smpPoints := smp.InitPoints()

	aerostats := InitAerostats()

	var i float64 = 0
	var index int = 0

	VelocitiesBefore, VelocitiesAfter := parsejson.CallParse(namesArr, index, nil)

	for t := t0; t <= tMax; t += dt {
		if i > internal.Cycle/dt {
			i = 1
			index++
			VelocitiesBefore, VelocitiesAfter = parsejson.CallParse(namesArr, index, VelocitiesAfter)
			fmt.Printf("%f of %d\n", t, tMax)
		}

		for _, aerostat := range aerostats {
			aerostat.Polinom(VelocitiesBefore, VelocitiesAfter, dt, i)
		}

		for i, point := range smpPoints {
			for j, aerostat := range aerostats {
				x, y, _ := aerostat.GetCoordinates()
				if isInCircle(x, y, point) && !point.EnterStatus {
					aerostat.EnterAndExitTime[i+1] = t
					point.EnterAndExitTime[j+1] = append(point.EnterAndExitTime[j+1], t)
					point.EnterStatus = true
					point.EnteredAerostat = j + 1
					if point.MaxDisEnteredTime < point.TempDisEnteredTime {
						point.MaxDisEnteredTime = point.TempDisEnteredTime
					}
					point.TempDisEnteredTime = 0
				} else if !isInCircle(x, y, point) && point.EnterStatus && point.EnteredAerostat == j+1 {
					point.EnterStatus = false
					aerostat.EnterAndExitTime[i+1] = t
					point.EnterAndExitTime[j+1] = append(point.EnterAndExitTime[j+1], t)
					point.Count = point.Count + (point.EnterAndExitTime[j+1][len(point.EnterAndExitTime[j+1])-1] - point.EnterAndExitTime[j+1][len(point.EnterAndExitTime[j+1])-2])
					point.TempDisEnteredTime = 0
				}
			}
			if !point.EnterStatus {
				point.TempDisEnteredTime += dt
				if point.MaxDisEnteredTime < point.TempDisEnteredTime {
					point.MaxDisEnteredTime = point.TempDisEnteredTime
				}
			}
		}
		i++
	}
	//internal.CreateImage(aerostat.Path)
	for _, point := range smpPoints {
		for _, enterExitTimes := range point.EnterAndExitTime {
			if len(enterExitTimes)%2 != 0 {
				point.Count = point.Count + (tMax - enterExitTimes[len(enterExitTimes)-1])
			}
		}
	}

	allCoordinates := make([][][]float64, 0)

	for i, aerostat := range aerostats {
		lightPath, err := output.SaveCoordinatesToGeoJSON(CleanZpath(aerostat.Path), "aerostat"+strconv.Itoa(i+1)+".geojson", sep)
		if err != nil {
			fmt.Println(err)
			return
		}
		allCoordinates = append(allCoordinates, lightPath)
	}

	InitPercentes(smpPoints)
	smp.WriteJson(smpPoints)
	output.SaveAllCoordinatesToGeoJSON(allCoordinates, "alltrajectories.geojson")
}

func CleanZpath(arr [][]float64) [][]float64 {
	for i := range arr {
		if arr[i][0] > 180 {
			arr[i][0] -= 360
		}
		arr[i] = arr[i][:2]
	}
	return arr
}

func isInCircle(aerostatX, aerostatY float64, point *smp.SmpPoint) bool {
	a := internal.MetresToDegrees(point.Latitude, 400000)
	b := internal.MetresToDegrees(0, 400000)
	return (aerostatX-point.Longtitude)*(aerostatX-point.Longtitude)/a/a+(aerostatY-point.Latitude)*(aerostatY-point.Latitude)/b/b <= 1
}

func InitPercentes(arr []*smp.SmpPoint) {
	for _, point := range arr {
		point.PersentTimeAccess = float64(point.Count) * 100 / float64(tMax)
	}
}

func InitAerostats() []*internal.Point {
	arr := make([]*internal.Point, 0)
	for _, x := range arrCoordX {
		arr = append(arr, internal.NewPoint(x, y2, 21.88))
		//arr = append(arr, internal.NewPoint(x, y1, 21.88))
	}
	for _, x := range arrCoordX {
		//arr = append(arr, internal.NewPoint(x, y2, 21.88))
		arr = append(arr, internal.NewPoint(x, y1, 21.88))
	}
	// for _, x := range arrCoordX {
	// 	//arr = append(arr, internal.NewPoint(x, y2, 21.88))
	// 	arr = append(arr, internal.NewPoint(x, y0, 21.88))
	// }
	return arr
}

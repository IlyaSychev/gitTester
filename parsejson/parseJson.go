package parsejson

import (
	"encoding/json"
	"fmt"
	"os"
)

type Parametres struct {
	Date      string
	Parameter string
	Level     float64
	Value     float64
	Lat       float64
	Lon       float64
	Unit      string
}

func Parse(fileName string) map[[3]float64][3]float64 {
	var Inpt []Parametres
	data, _ := os.ReadFile("C:/Users/ikqw/go projects/gitTester/grbFiles/" + fileName)
	err := json.Unmarshal(data, &Inpt)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
		//return
	}

	Velocities := make(map[[3]float64][3]float64, 0)

	for i, data := range Inpt {
		if data.Parameter == "Vertical velocity" {
			tmp := Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}]
			tmp[2] = data.Value
			Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}] = tmp
		} else if data.Parameter == "U component of wind" {
			tmp := Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}]
			tmp[0] = data.Value
			Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}] = tmp
		} else if data.Parameter == "V component of wind" {
			tmp := Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}]
			tmp[1] = data.Value
			Velocities[[3]float64{Inpt[i].Lon, Inpt[i].Lat, Inpt[i].Level}] = tmp
		}
	}

	return Velocities
}

// func ChangeVelocitiesData(VelocitiesBefore, VelociitiesAfter map[[3]float64][3]float64) (map[[3]float64][3]float64, map[[3]float64][3]float64) {
// 	return VelociitiesAfter, Parse()
// }

func ReadDir() []string {
	arr := make([]string, 0)
	var root string = "C:/Users/ikqw/go projects/gitTester/grbFiles"
	fileInfo, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range fileInfo {
		arr = append(arr, file.Name())
	}
	return arr
}

func CallParse(namesArr []string, index int, VelocitiesAfter map[[3]float64][3]float64) (map[[3]float64][3]float64, map[[3]float64][3]float64) {
	if index > 0 {
		return VelocitiesAfter, Parse(namesArr[index+1])
	}
	return Parse(namesArr[index]), Parse(namesArr[index+1])
}

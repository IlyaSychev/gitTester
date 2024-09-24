package main

import (
	"fmt"
	"gitTester/internal"
	"gitTester/output"
	"gitTester/parsejson"
)

const dt = 600

func main() {
	namesArr := parsejson.ReadDir()
	fmt.Println(namesArr)
	aerostat1 := internal.NewPoint(33, 69, 10)
	aerostat2 := internal.NewPoint(33, 69, 25)
	aerostat3 := internal.NewPoint(33, 69, 55)
	aerostat4 := internal.NewPoint(33, 69, 121)

	var i float64 = 0
	var index int = 0

	VelocitiesBefore, VelocitiesAfter := parsejson.CallParse(namesArr, index, nil)

	for t := 0; t <= 30*24*60*60; t += dt {
		if i > internal.Cycle/dt {
			fmt.Println("index, i = ", index, i, t)
			i = 0
			index++
			VelocitiesBefore, VelocitiesAfter = parsejson.CallParse(namesArr, index, VelocitiesAfter)
		}
		//fmt.Println("index, i = ", index, i, t)
		aerostat1.Polinom(VelocitiesBefore, VelocitiesAfter, dt, i)
		aerostat2.Polinom(VelocitiesBefore, VelocitiesAfter, dt, i)
		aerostat3.Polinom(VelocitiesBefore, VelocitiesAfter, dt, i)
		aerostat4.Polinom(VelocitiesBefore, VelocitiesAfter, dt, i)
		i++
	}
	//internal.CreateImage(aerostat.Path)
	output.SaveCoordinatesToGeoJSON(CleanZpath(aerostat1.Path), "aerostat1.geojson")
	output.SaveCoordinatesToGeoJSON(CleanZpath(aerostat2.Path), "aerostat2.geojson")
	output.SaveCoordinatesToGeoJSON(CleanZpath(aerostat3.Path), "aerostat3.geojson")
	output.SaveCoordinatesToGeoJSON(CleanZpath(aerostat4.Path), "aerostat4.geojson")
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

package output

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Coordinate struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type LineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Feature struct {
	Type       string     `json:"type"`
	ID         int        `json:"id"`
	Geometry   LineString `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	Description   string  `json:"description"`
	Stroke        string  `json:"stroke"`
	StrokeWidth   string  `json:"stroke-width"`
	StrokeOpacity float64 `json:"stroke-opacity"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Metadata Metadata  `json:"metadata"`
	Features []Feature `json:"features"`
}

type Metadata struct {
	Name    string `json:"name"`
	Creator string `json:"creator"`
}

func SaveCoordinatesToGeoJSON(coordinates [][]float64, filename string, separate int) ([][]float64, error) {
	// Создаем структуру GeoJSON
	smallArrOfCoordinates := make([][]float64, 0)
	for i := range coordinates {
		if i == 0 || i%separate == 0 || i == len(coordinates)-1 {
			smallArrOfCoordinates = append(smallArrOfCoordinates, coordinates[i])
		}
	}
	geoJSON := FeatureCollection{
		Type: "FeatureCollection",
		Metadata: Metadata{
			Name:    "path",
			Creator: "Yandex Map Constructor",
		},
		Features: []Feature{
			{
				Type: "Feature",
				ID:   0,
				Geometry: LineString{
					Type:        "LineString",
					Coordinates: smallArrOfCoordinates,
				},
				Properties: Properties{
					Description:   "etomesto.ru track 266470",
					Stroke:        "#ed1d24",
					StrokeWidth:   "2",
					StrokeOpacity: 1,
				},
			},
		},
	}

	// Открываем файл для записи
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Кодируем структуру в JSON и записываем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(geoJSON); err != nil {
		return nil, err
	}

	return smallArrOfCoordinates, nil
}

func SaveAllCoordinatesToGeoJSON(coordinates [][][]float64, filename string) error {
	// Создаем структуру GeoJSON
	arr := make([]Feature, len(coordinates))
	for i := range coordinates {
		arr[i] = Feature{
			Type: "Feature",
			ID:   0,
			Geometry: LineString{
				Type:        "LineString",
				Coordinates: coordinates[i],
			},
			Properties: Properties{
				Description:   "etomesto.ru track 266470",
				Stroke:        randomHexColor(),
				StrokeWidth:   "3",
				StrokeOpacity: 1,
			},
		}
	}

	geoJSON := FeatureCollection{
		Type: "FeatureCollection",
		Metadata: Metadata{
			Name:    "path",
			Creator: "Yandex Map Constructor",
		},
		Features: arr,
	}

	// Открываем файл для записи
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Кодируем структуру в JSON и записываем в файл
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(geoJSON); err != nil {
		return err
	}

	return nil
}

func randomHexColor() string {

	r := rand.Intn(256)
	g := rand.Intn(256)
	b := rand.Intn(256)

	hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return hex
}

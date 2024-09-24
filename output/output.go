package output

import (
	"encoding/json"
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

func SaveCoordinatesToGeoJSON(coordinates [][]float64, filename string) error {
	// Создаем структуру GeoJSON
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
					Coordinates: coordinates,
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

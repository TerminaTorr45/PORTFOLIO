package groupietracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Artist est la structure pour stocker les informations sur les artistes
type Artist struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationmembers"`
	FirstAlbum   string   `json:"firstAlbum"`
	Preview      string
	DeezerName   string
}
type Location struct {
	ID    int      `json:"id"`
	Locat []string `json:"locations"`
	Data  string   `json:"dates"`
	Lat   float64  `json:"lat"`
	Lng   float64  `json:"lng"`
}

type GeoJSON struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

type Feature struct {
	Type       string    `json:"type"`
	Geometry   *Geometry `json:"geometry"`
	Properties *Property `json:"properties"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Property struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetArtists() []Artist {
	var artists []Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &artists)
	return artists
}
func GetArtist(id string) Artist {
	var artists Artist
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + id)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &artists)
	return artists
}
func GetLocations() []Location {
	var locations []Location
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(locations)
	return locations
}

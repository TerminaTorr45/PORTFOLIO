package main

import (
	"encoding/json"
	"fmt"
	"groupietracker/groupietracker"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// var loca = regexp.MustCompile(`/\d`)
var artists []groupietracker.Artist

type DeezerData struct {
	Data []Data `json:"data"`
}

type Data struct {
	Preview string `json:"preview"`
}

func SearchLocations(locations []groupietracker.Location, search string) []groupietracker.Location {
	var results []groupietracker.Location
	for _, loc := range locations {
		if contains(loc.Locat, search) {
			results = append(results, loc)
		}
	}
	return results
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func replaceSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "%20")
}

func SetDeezerName(artists []groupietracker.Artist) []groupietracker.Artist {
	for i := 0; i < len(artists); i++ {
		artists[i].DeezerName = replaceSpaces(artists[i].Name)
	}
	return artists
}

func GetDeezerPreviews(artists []groupietracker.Artist) []groupietracker.Artist {
	for i := 0; i < len(artists); i++ {
		nom := artists[i].DeezerName
		response, err := http.Get("https://api.deezer.com/search?q=" + nom)
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var responseObject DeezerData

		json.Unmarshal(responseData, &responseObject)

		if len(responseObject.Data) != 0 {
			artists[i].Preview = responseObject.Data[0].Preview
		}
	}
	return artists
}
func main() {
	// Code existant pour récupérer les artistes, définir les noms Deezer et récupérer les aperçus
	artists := groupietracker.GetArtists()
	SetDeezerName(artists)
	GetDeezerPreviews(artists)

	// Création d'une instance de FileServer pour servir les fichiers statiques
	fs := http.FileServer(http.Dir("groupietracker/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Création de la route pour la page d'accueil
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, artists)
	})
	http.HandleFunc("/id", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("./groupietracker/static/template/informationartiste.html", "./groupietracker/static/template/song.html", "groupietracker/static/template/barrederecherche.html", "./groupietracker/static/localisations/localisations.html")
		value := r.URL.Query().Get("id")

		fmt.Println(groupietracker.GetArtist(value))
		fmt.Println(value)

		tmpl.Execute(w, groupietracker.GetArtist(value))
	})
	http.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("./groupietracker/static/home/home.html", "./groupietracker/static/template/song.html", "groupietracker/static/template/barrederecherche.html", "./groupietracker/static/localisations/localisations.html", "./groupietracker/static/template/filtre.html")
		tmpl.Execute(w, artists)
	})
	// Création de la route pour la page de détail d'un artiste
	http.HandleFunc("./groupietracker/static/template/infosartistes.json", func(w http.ResponseWriter, r *http.Request) {
		// Analyse de l'ID de l'artiste dans l'URL
		detId := strings.TrimPrefix(r.URL.Path, "/artist/")
		id, err := strconv.Atoi(detId)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		// Vérification que l'ID est valide
		if id <= 0 || id > len(artists) {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		// Récupération de l'artiste correspondant à l'ID et affichage de la page de détail
		artist := artists[id-1]
		Detail(w, r, artist)

	})

	// Création de la route pour la recherche de localisations
	http.HandleFunc("/locations/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			// Afficher la page d'accueil avec les locations
			tmpl, err := template.ParseFiles("./groupietracker/static/localisations/localisations.html", "./groupietracker/static/template/song.html", "./groupietracker/static/template/barrederecherche.html", "./groupietracker/static/artistes/artistes.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			// Récupérer la recherche de l'utilisateur et afficher les résultats
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			//search := r.PostForm.Get("search")
			//results := SearchLocations(groupietracker.GetLocations(), search)
			tmpl, err := template.ParseFiles("./groupietracker/static/localisations/localisations.html", "./groupietracker/static/template/song.html", "./groupietracker/static/template/barrederecherche.html", "./groupietracker/static/template/song.html", "./groupietracker/static/template/barrederecherche.html", "./groupietracker/static/home/home.html", "./groupietracker/static/artistes/artistes.html", "./groupietracker/static/template/filtre.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// Démarrer le serveur web
	fmt.Println("et c'est partiiiii")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func Home(w http.ResponseWriter, r *http.Request, artists []groupietracker.Artist) {
	tmpl, _ := template.ParseFiles("./groupietracker/static/home/home.html", "./groupietracker/static/template/song.html", "groupietracker/static/template/barrederecherche.html", "./groupietracker/static/template/filtre.html")
	tmpl.Execute(w, artists)
}

func Detail(w http.ResponseWriter, r *http.Request, artist groupietracker.Artist) {
	tmpl, _ := template.ParseFiles("./groupietracker/static/home/home.html")
	tmpl.Execute(w, artist)
}

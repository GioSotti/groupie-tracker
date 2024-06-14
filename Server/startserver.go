package Server

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type InfoBloc struct {
	Name       string
	Image      string
	FirstAlbum string
}
type ArtistAllData struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []string            `json:"locations"`
	ConcertDates []string            `json:"concertDates"`
	Relations    map[string][]string `json:"relations"`
}
type Artist struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    string              `json:"locations"`
	ConcertDates string              `json:"concertDates"`
	Relations    map[string][]string `json:"relations"`
}
type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type LocationData struct {
	Index []Location `json:"index"`
}
type ConcertDate struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type DateconcertData struct {
	Index []ConcertDate `json:"index"`
}
type RelationDate struct {
	ID            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}
type RelationData struct {
	Index []RelationDate `json:"index"`
}

const Url = "https://groupietrackers.herokuapp.com/api"

var Artists []Artist
var DataArtist []ArtistAllData
var LocationsData LocationData
var DatesconcertData DateconcertData
var DataRelation RelationData

func getDataArtist() error {
	resp, err := http.Get(Url + "/artists")
	if err != nil {
		return errors.New("Erreur dans le get Location")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Erreur dans le read Location")
	}
	json.Unmarshal(bytes, &DataArtist)
	return nil
}
func getDataLocation() error {
	resp, err := http.Get(Url + "/locations")
	if err != nil {
		return errors.New("Erreur dans le get Location")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Erreur dans le read Location")
	}
	json.Unmarshal(bytes, &LocationsData)
	return nil
}
func getDataDates() error {
	resp, err := http.Get(Url + "/dates")
	if err != nil {
		return errors.New("Erreur dans le get Dates")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Erreur dans le read Dates")
	}
	json.Unmarshal(bytes, &DatesconcertData)
	return nil
}
func getDataRelation() error {
	resp, err := http.Get(Url + "/relation")
	if err != nil {
		return errors.New("Erreur dans le get Relation")
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Erreur dans le read relation")
	}
	json.Unmarshal(bytes, &DataArtist)
	return nil
}
func getFullData() {
	getDataRelation()
	getDataDates()
	getDataArtist()
	getDataLocation()
	var calque ArtistAllData

	for i := range Artists {
		calque.Name = Artists[i].Name
		calque.Image = Artists[i].Image
		calque.Members = Artists[i].Members
		calque.FirstAlbum = Artists[i].FirstAlbum
		calque.Locations = LocationsData.Index[i].Locations
		calque.CreationDate = Artists[i].CreationDate
		calque.ConcertDates = DatesconcertData.Index[i].Dates
		calque.Relations = DataRelation.Index[i].DatesLocation
		calque.ID = i + 1

		DataArtist = append(DataArtist, calque)

	}
}
func filterArtists(Group bool) []ArtistAllData {
	var data []ArtistAllData
	var temp ArtistAllData

	for i := range DataArtist {
		temp.Members = DataArtist[i].Members
		if len(temp.Members) == 1 && Group {
			data = append(data, DataArtist[i])
		} else if !(Group) && len(temp.Members) > 1 {
			data = append(data, DataArtist[i])
		}
	}
	return data
}
func filterSize(groupNbr int) []ArtistAllData {

	var data []ArtistAllData
	var temp ArtistAllData
	for i := range DataArtist {
		temp.Members = DataArtist[i].Members
		if len(temp.Members) == groupNbr {
			data = append(data, DataArtist[i])
		}
	}
	return data
}
func FilterDate(DateMin int, DateMax int) []ArtistAllData {
	var data []ArtistAllData
	var temp ArtistAllData

	for i := range DataArtist {
		temp.CreationDate = DataArtist[i].CreationDate
		if temp.CreationDate >= DateMin && temp.CreationDate <= DateMax {

			data = append(data, DataArtist[i])
		}
	}
	return data
}
func FilterADate(DateMin int, DateMax int) []ArtistAllData {
	var data []ArtistAllData

	for i := range DataArtist {
		parts := strings.Split(DataArtist[i].FirstAlbum, "-")
		yearStr := parts[len(parts)-1] // on récupère le dernier élément de la liste, qui doit être l'année
		yearInt, err := strconv.Atoi(yearStr)
		if err != nil {
			// gestion d'erreur si la conversion échoue
			continue
		}
		if yearInt >= DateMin && yearInt <= DateMax {
			data = append(data, DataArtist[i])
		}
	}
	return data
}

func FilterFunc(filter string, DateMin int, DateMax int, GroupSize int, DateAMin int, DateAMax int) []ArtistAllData {
	var data []ArtistAllData

	if filter == "artist" {

		data = filterArtists(true)
	} else if filter == "group" {
		data = filterArtists(false)
	}
	if GroupSize != 0 {
		data = filterSize(GroupSize)
	}
	if DateMax != 0 && DateMin != 0 {
		data = FilterDate(DateMin, DateMax)
	}
	if DateAMax != 0 && DateAMin != 0 {
		data = FilterADate(DateAMin, DateAMax)

	}
	return data

}
func Search(search string) []ArtistAllData {
	if search == "" {
		return DataArtist
	}

	var data []ArtistAllData

	for i, _ := range DataArtist {
		if strings.HasPrefix(strings.ToLower(DataArtist[i].Name), strings.ToLower(search)) {
			data = append(data, DataArtist[i])
		}
	}
	return data
}

func DisplayAccueil(w http.ResponseWriter, r *http.Request) {
	custTemplate, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Filter := r.FormValue("Filter")
	Tab := r.FormValue("Tab")
	DateMin, _ := strconv.Atoi(r.FormValue("DateMin"))
	DateMax, _ := strconv.Atoi(r.FormValue("DateMax"))
	GroupSize, _ := strconv.Atoi(r.FormValue("size"))
	DateAMin, _ := strconv.Atoi(r.FormValue("DateAMin"))
	DateAMax, _ := strconv.Atoi(r.FormValue("DateAMax"))
	if err != nil {
	}

	Data := DataArtist

	if Tab != "" && len(Data) != 0 && Filter == "" {
		Data = Search(Tab)
		if Data == nil {
		}

	} else if Filter != "" || (DateMax != 0 && DateMin != 0) || GroupSize != 0 && len(Data) != 0 || (DateAMax != 0 && DateAMin != 0) {
		Data = FilterFunc(Filter, DateMin, DateMax, GroupSize, DateAMin, DateAMax)
	}

	err = custTemplate.Execute(w, Data)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func DisplayArtist(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		http.Error(w, "Missing artist ID in request URL", http.StatusBadRequest)
		return
	}
	idStr := segments[len(segments)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID in request URL", http.StatusBadRequest)
		return
	}

	if id < 1 || id > len(DataArtist) {
		http.Error(w, "Invalid artist ID in request URL", http.StatusBadRequest)
		return
	}
	artistData := DataArtist[id-1]

	custTemplate, err := template.ParseFiles("./templates/Artist.html")
	if err != nil {
		http.Error(w, "Failed to load template file", http.StatusInternalServerError)
		return
	}
	err = custTemplate.Execute(w, artistData)
	if err != nil {
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}
func Handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	custTemplate, err := template.ParseFiles("./templates/404.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = custTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func StartServer() {
	println("server started on http://localhost:8080/accueil")
	getFullData()
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./templates/css"))))
	http.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("./templates/JS"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./assets/fonts"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./templates/assets"))))
	http.HandleFunc("/", Handle404)
	http.HandleFunc("/accueil", DisplayAccueil)
	http.HandleFunc("/Artist/", DisplayArtist)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

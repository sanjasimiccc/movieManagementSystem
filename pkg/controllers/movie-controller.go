package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/models"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/utils"
)

var NewMovie models.Movie //struct Movie u okviru models paketa

func GetMovies(w http.ResponseWriter, r *http.Request) {

	newMovies := models.GetAllMovies()
	res, _ := json.Marshal(newMovies) //to sto dobijemo od baze zelimo da konvertujemo u json i to je nas response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res) //ovo je glavno, omogucava da posaljemo na frontend tj postman
}

func GetMovieById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) //da bismo pristupili movieId-ju koji je i okviru request-a
	movieId := vars["movieId"]
	id, err := strconv.ParseInt(movieId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
		return
	}

	movieDetails, _ := models.GetMovieById(id) //koristim _, tj blank jer vraca db, a ja sad ne zelim da je koristim
	res, _ := json.Marshal(movieDetails)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	CreateMovie := &models.Movie{}
	//od usera ce da dobijemo json, a mi hocemo da parsiramo u nesto sto ce db da razume
	utils.ParseBody(r, CreateMovie)
	m := CreateMovie.CreateMovie() //vratice ono sto je i kreirano u bazi, znaci isti movie record poslat od usera
	res, _ := json.Marshal(m)      //pa taj record konvertujem u json

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieId := vars["movieId"]
	id, err := strconv.ParseInt(movieId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}
	movie := models.DeleteMovie(id)
	res, _ := json.Marshal(movie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// ovo bezveze, ne mora ovako nego eto
func UpdateWholeMovie(w http.ResponseWriter, r *http.Request) {
	var updateMovie = &models.Movie{}
	utils.ParseBody(r, updateMovie) //radi unmarshal json-a u strukturu koju golang razume

	vars := mux.Vars(r)
	movieId := vars["movieId"]
	id, err := strconv.ParseInt(movieId, 0, 0)
	if err != nil {
		fmt.Println("error while parsing")
	}

	movieDetails, db := models.GetMovieById(id)
	if updateMovie.Title != "" {
		movieDetails.Title = updateMovie.Title
	}
	if updateMovie.Year != 0 {
		movieDetails.Year = updateMovie.Year
	}
	if updateMovie.Genre != "" {
		movieDetails.Genre = updateMovie.Genre
	}
	if updateMovie.Director != "" {
		movieDetails.Director = updateMovie.Director
	}
	if updateMovie.Plot != "" {
		movieDetails.Plot = updateMovie.Plot
	}

	db.Save(&movieDetails)
	res, _ := json.Marshal(movieDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

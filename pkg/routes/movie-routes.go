// rute na koje ce korisnik da "udari" sa frontenda/iz postmana,...
package routes

import (
	//GoLang koristi apsolutne putanje za import!
	"github.com/gorilla/mux"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/controllers" //ovo ce pomoci da importujem moj controller i da pristupim fajlu
)

// ova funkcija ce imati sve moje rute
var RegisterMovieStoreRoutes = func(router *mux.Router) {
	//ovde kreiram svoje handlere
	router.HandleFunc("/movies/", controllers.CreateMovie).Methods("POST") //znaci ta funkcija mi je u kontroleru
	router.HandleFunc("/movies/", controllers.GetMovies).Methods("GET")
	router.HandleFunc("/movies/{movieId}", controllers.GetMovieById).Methods("GET")
	router.HandleFunc("/movies/{movieId}", controllers.UpdateWholeMovie).Methods("PUT")
	//router.HandleFunc("/movies/{movieId}", controllers.UpdateMoviePartially).Methods("PATCH")
	router.HandleFunc("movies/{movieId}", controllers.DeleteMovieById).Methods("DELETE")
}

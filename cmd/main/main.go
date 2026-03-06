package main

//1. ovde cemo kreirati server - definisati localhost
//2. on ce reci golangu da su rute u routes...

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	//mora i da importujem rute!
	"github.com/sanjasimiccc/movieManagementSystem/pkg/routes"
)

func main() {
	r := mux.NewRouter() //prom r koja ce inicijalizovati ruter
	routes.RegisterMovieStoreRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9010", r)) //ListenAndServe funkcija nam pomaze da kreiramo server, a prosledjujemo joj adresu i port na kome hocemo da pokrenemo server

}

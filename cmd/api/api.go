package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sanjasimiccc/movieManagementSystem/middleware"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/controllers"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/external"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/repositories"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/services"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer { //kreirace novu instancu API servera
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter() //pa ako se api menja i imamo novu verziju samo ovde promenimo u npr. v2
	//ovde sada mozemo da kreiramo nase handlere, ali to ne bi bilo uredno ako ovde krenemo da kreiramo endpointe
	//zato cemo sve da podelimo u servise

	//svaki put kad hocemo da registrujemo rute, dodjemo ovde i REGISTRUJEMO NOVI SERVIS!!
	movieStore := repositories.NewStore(s.db) //pozivam konstruktor za store i prosledjujem pravi pointer ka bazi
	movieProviderFactory := external.NewProviderFactory()
	movieService := services.NewService(movieStore, movieProviderFactory)
	movieHandler := controllers.NewHandler(movieService) //pozivam handle konstruktor gde sad prosledjujem ovaj konkretni store koji bi trebao da implementira sve metode interfejsa
	movieHandler.RegisterRoutes(subrouter)

	return http.ListenAndServe(s.addr, middleware.APIKeyMiddleware(router))
}

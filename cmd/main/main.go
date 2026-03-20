package main

//1. ovde cemo kreirati server - definisati localhost
//2. on ce reci golangu da su rute u routes...

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	//mora i da importujem rute!
	"github.com/sanjasimiccc/movieManagementSystem/cmd/api"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/config"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/controllers"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/external"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/repositories"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/services"
)

func main() {
	config.Connect()
	db := config.GetDB()
	config.Migrate()

	movieStore := repositories.NewStore(db)
	movieProviderFactory := external.NewProviderFactory()
	movieService := services.NewService(movieStore, movieProviderFactory)
	movieHandler := controllers.NewHandler(movieService)

	//kreiranje servera
	server := api.NewAPIServer(":9010", movieHandler)

	//pokretanje servera
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

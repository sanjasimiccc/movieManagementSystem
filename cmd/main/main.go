package main

//1. ovde cemo kreirati server - definisati localhost
//2. on ce reci golangu da su rute u routes...

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	//mora i da importujem rute!
	"github.com/sanjasimiccc/movieManagementSystem/cmd/api"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/config"
)

func main() {
	//inicijalizacija DB konekcije
	config.Connect()
	db := config.GetDB()

	//migracija
	config.Migrate()

	//kreiranje servera
	server := api.NewAPIServer(":9010", db)

	//pokretanje servera
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

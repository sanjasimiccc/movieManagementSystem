package config

//cela poenta ovog fajla je da vrati promenljivu 'db' koja ce da pomogne ostalim fajlovima da interaguju sa db
import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //u dokumentaciji je _ ispred -> to se koristi da se pokrene init() funkcija iz paketa i registruje MySQL dialect u GORM-u, iako direktno ne koristim nista iz tog paketa
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

var db *gorm.DB //globalna promenljiva koja cuva konekciju ka bazi, sve ostale fje mogu da je koriste preko GetDB(

func Connect() { //pomaze da otvorimo konekciju ka bazi - konkretno ka mysql bazi
	//gorm mi pomaze da pricam sa sql-lite
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		Envs.DBUser,
		Envs.DBPassword,
		Envs.DBAddress,
		Envs.DBName,
	)
	d, err := gorm.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	db = d

	if err := db.DB().Ping(); err != nil { //zapocinje/proverava konekciju ka bazi
		panic(err)
	}
}

// vraca konekciju za koriscenje u repo
func GetDB() *gorm.DB {
	return db
}

func Migrate() {
	if err := db.AutoMigrate(&types.Movie{}).Error; err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

package config

//cela poenta ovog fajla je da vrati promenljivu 'db' koja ce da pomogne ostalim fajlovima da interaguju sa db
import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //u dokumentaciji je _ ispred -> to se koristi da se pokrene init() funkcija iz paketa i registruje MySQL dialect u GORM-u, iako direktno ne koristim nista iz tog paketa
)

var (
	db *gorm.DB //globalna promenljiva koja cuva konekciju ka bazi, sve ostale fje mogu da je koriste preko GetDB()
)

func Connect() { //pomaze da otvorimo konekciju ka bazi - konkretno ka mysql bazi
	//gorm mi pomaze da pricam sa sql-lite
	d, err := gorm.Open("mysql", "root:asdf1234S!@tcp(localhost:3307)/movies?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
}

// funkcija koju cemo koristiti u drugim fajlovima
func GetDB() *gorm.DB {
	return db
}

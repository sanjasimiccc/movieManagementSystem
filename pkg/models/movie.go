package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/config"
)

var db *gorm.DB

// model nam daje strukturu kako cemo nesto skladistiti u bazi
type Movie struct {
	gorm.Model
	ID         int       `gorm:"" json:"id"`
	Title      string    `gorm:"" json:"title"`
	Year       int       `json:"year"`
	Genre      string    `json:"genre"`
	Director   string    `json:"director"`
	Plot       string    `json:"plot"`
	IMDBRating *float64  `json:"imdb_rating"`
	ExternalID *string   `json:"external_id"`
	Source     *string   `json:"source"` //odakle su podaci il mi to ne treba?
	CreatedAt  time.Time `json:"created_at"`
}

func init() { //inicijalizacija db
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Movie{}) //!!!
}

// funkcije za komunikaciju sa bazom!
func (m *Movie) CreateMovie() *Movie { //we receive smth of type Movie, and also return that same book that we created
	//preko db komuniciramo sa bazom
	db.NewRecord(m) //funkcija koja postoji u okviru GORM-a, pa ne moramo da pisemo upite
	db.Create(&m)
	return m
}

func GetAllMovies() []Movie { //vracam slice/listu
	var Movies []Movie
	db.Find(&Movies)
	return Movies
}

func GetMovieById(id int64) (*Movie, *gorm.DB) {
	var getMovie Movie
	db := db.Where("ID=?", id).Find(&getMovie)
	return &getMovie, db
}

func DeleteMovie(id int64) Movie {
	var movie Movie
	db.Where("ID=?", id).Delete(movie)
	return movie
}

//Update?
//pronaci cemo knjigu koju zelimo po id, a onda cemo je obrisati i nove podatke koristiti za kreiranje nove
//get + delete + create ?????

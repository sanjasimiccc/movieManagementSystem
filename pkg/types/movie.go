package types

import "github.com/jinzhu/gorm"

type MovieService interface {
	GetAllMovies() ([]Movie, error)
	CreateMovie(payload CreateMoviePayload) (*Movie, error)
	GetMovieById(id int64) (*Movie, error)
	DeleteMovieById(id int64) error
	UpdateMovie(id int64, payload UpdateMoviePayloadPUT) (*Movie, error)
	PatchMovie(id int64, payload PatchMoviePayload) (*Movie, error)
	SearchMovies(params MovieSearchAndFilterParams) ([]Movie, error)
}

type MovieStore interface {
	GetMovieById(id int64) (*Movie, error)
	GetAllMovies() ([]Movie, error)

	DeleteMovie(id int64) error
	CreateMovie(movie Movie) (*Movie, error)
	UpdateMovie(id int64, movie Movie) (*Movie, error)
	UpdateMovieFields(id int64, updates map[string]any) (*Movie, error)
	SearchMovies(params MovieSearchAndFilterParams) ([]Movie, error)
	ExistsMovie(title string, director string) (bool, error)
}

type Movie struct {
	gorm.Model
	//ID         int       `gorm:"" json:"id"`
	Title      string   `gorm:"" json:"title"`
	Year       int      `json:"year"`
	Genre      string   `json:"genre"`
	Director   string   `json:"director"`
	Plot       string   `json:"plot"`
	IMDBRating *float64 `json:"imdb_rating"`
	ExternalID *string  `json:"external_id"`
	Source     *string  `json:"source"` //odakle su podaci il mi to ne treba?
	//CreatedAt  time.Time `json:"created_at"`
}

type CreateMoviePayload struct {
	Title    string `json:"title" validate:"required"`
	Year     int    `json:"year" validate:"required"`
	Genre    string `json:"genre" validate:"required"`
	Director string `json:"director" validate:"required"`
	Plot     string `json:"plot" validate:"required"`
}

type UpdateMoviePayloadPUT struct {
	Title    string `json:"title" validate:"required"`
	Year     int    `json:"year" validate:"required"`
	Genre    string `json:"genre" validate:"required"`
	Director string `json:"director" validate:"required"`
	Plot     string `json:"plot" validate:"required"`
}

type PatchMoviePayload struct {
	Title    *string `json:"title"`
	Year     *int    `json:"year"`
	Genre    *string `json:"genre"`
	Director *string `json:"director"`
	Plot     *string `json:"plot"`
}

type MovieSearchAndFilterParams struct { //ne treba mi json jer su parametri iz query stringa
	Title    *string
	Genre    *string
	Director *string
	Year     *int

	Page  int
	Limit int
}

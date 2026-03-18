package types

import "github.com/jinzhu/gorm"

type MovieService interface {
	GetAllMovies(page int, limit int) ([]Movie, error)
	CreateMovie(payload CreateMoviePayload) (*Movie, error)
	GetMovieById(id int64) (*Movie, error)
	DeleteMovieById(id int64) error
	UpdateMovie(id int64, payload UpdateMoviePayloadPUT) (*Movie, error)
	PatchMovie(id int64, payload PatchMoviePayload) (*Movie, error)
	SearchMovies(params MovieSearchAndFilterParams) ([]Movie, error)
	FetchAndStoreMovies() ([]Movie, error)
}

type MovieStore interface {
	GetMovieById(id int64) (*Movie, error)
	GetAllMovies(page int, limit int) ([]Movie, error)

	DeleteMovie(id int64) error
	CreateMovie(movie Movie) (*Movie, error)
	CreateMovies(movies []Movie) error
	UpdateMovie(id int64, movie Movie) (*Movie, error)
	UpdateMovieFields(id int64, updates map[string]any) (*Movie, error)
	SearchMovies(params MovieSearchAndFilterParams) ([]Movie, error)
	ExistsMovie(title string, director string) (bool, error)
	ExistsMovieExcludingID(title string, director string, excludeID int64) (bool, error)
}

type Movie struct {
	gorm.Model
	//ID         int       `gorm:"" json:"id"`
	Title      string  `gorm:"" json:"title"`
	Year       int     `json:"year"`
	Genre      string  `json:"genre"`
	Director   string  `json:"director"`
	Plot       string  `json:"plot"`
	IMDBRating float64 `json:"imdb_rating"`
	ExternalID *string `gorm:"uniqueIndex:idx_external_source" json:"external_id"`
	Source     *string `gorm:"uniqueIndex:idx_external_source" json:"source"`
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

type PaginationParams struct {
	Page  int
	Limit int
}

type FetchResult struct {
	Source string
	ID     string
	Data   string
}

type MovieProvider interface {
	FetchMovie(id string) (Movie, error)
}

type OMDbFetchData struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Plot       string `json:"Plot"`
	IMDBRating string `json:"imdbRating"`
}

type TMDbFetchData struct {
	Title  string `json:"title"`
	Year   string `json:"release_date"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
	Director string  `json:"director"` // TMDb možda nema direktno, ostaviti prazno ili kasnije fetch-ovati crew
	Plot     string  `json:"overview"`
	Rating   float64 `json:"vote_average"`
}

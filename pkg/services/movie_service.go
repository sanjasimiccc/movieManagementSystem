package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/external"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
	utils "github.com/sanjasimiccc/movieManagementSystem/pkg/utils/json"
)

type MovieService struct {
	repo                 types.MovieStore
	movieProviderFactory *external.ProviderFactory
}

func NewService(repo types.MovieStore, movieProviderFactory *external.ProviderFactory) *MovieService {
	return &MovieService{
		repo:                 repo,
		movieProviderFactory: movieProviderFactory,
	}
}

func (s *MovieService) CreateMovie(payload types.CreateMoviePayload) (*types.Movie, error) {

	if err := utils.Validate.Struct(payload); err != nil {
		validationErr := types.ValidationError{Err: err}
		return nil, validationErr
	}

	exists, err := s.repo.ExistsMovie(payload.Title, payload.Director)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, types.ErrMovieAlreadyExists
	}

	movie := types.Movie{
		Title:    payload.Title,
		Year:     payload.Year,
		Genre:    payload.Genre,
		Director: payload.Director,
		Plot:     payload.Plot,
	}

	return s.repo.CreateMovie(movie)
}

func (s *MovieService) GetMovieById(id int64) (*types.Movie, error) {
	movie, err := s.repo.GetMovieById(id)

	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) GetAllMovies() ([]types.Movie, error) {
	movies, err := s.repo.GetAllMovies()

	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) DeleteMovieById(id int64) error {
	err := s.repo.DeleteMovie(id)

	if err != nil {
		return err
	}
	return nil
}

func (s *MovieService) UpdateMovie(id int64, payload types.UpdateMoviePayloadPUT) (*types.Movie, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		validationErr := types.ValidationError{Err: err}
		return nil, validationErr
	}

	// exists, err := s.repo.ExistsMovie(payload.Title, payload.Director)
	// if err != nil {
	// 	return nil, err
	// }
	// if exists {
	// 	return nil, types.ErrMovieAlreadyExists
	// }

	movie := types.Movie{
		Title:    payload.Title,
		Year:     payload.Year,
		Genre:    payload.Genre,
		Director: payload.Director,
		Plot:     payload.Plot,
	}

	return s.repo.UpdateMovie(id, movie)
}

func (s *MovieService) PatchMovie(id int64, payload types.PatchMoviePayload) (*types.Movie, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		validationErr := types.ValidationError{Err: err}
		return nil, validationErr
	}

	updates := make(map[string]any)

	if payload.Title != nil {
		updates["title"] = *payload.Title
	}
	if payload.Year != nil {
		updates["year"] = *payload.Year
	}
	if payload.Genre != nil {
		updates["genre"] = *payload.Genre
	}
	if payload.Director != nil {
		updates["director"] = *payload.Director
	}
	if payload.Plot != nil {
		updates["plot"] = *payload.Plot
	}

	return s.repo.UpdateMovieFields(id, updates)
}

func (s *MovieService) SearchMovies(params types.MovieSearchAndFilterParams) ([]types.Movie, error) {
	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Limit <= 0 {
		params.Limit = 10
	}

	//radi sigurnosti
	if params.Limit > 100 {
		params.Limit = 100
	}

	return s.repo.SearchMovies(params)
}

func (s *MovieService) FetchAndStoreMovies() ([]types.Movie, error) {
	file, err := os.Open("../../movie_sources.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var movieSources map[string]string
	if err := json.NewDecoder(file).Decode(&movieSources); err != nil {
		return nil, err
	}

	results := make(chan types.Movie, len(movieSources)) //buffered channel
	var wg sync.WaitGroup                                //kaze Main gorutini da saceka da sve gorutine zavrse posao pre nego sto nastavi dalje

	//worker gorutine za svaki film
	for id, source := range movieSources {
		wg.Add(1) //postavlja broj gorutina i thread-ova koje main thread treba da ceka. Uvecava brojac za prosledjenu vrednost

		go func(movieID string, apiSource string) {
			defer wg.Done() //kada gorutina zavrsi posao, mora reci da je zavrsila sa wg.Done(), to smanjuje brojac u WaitGroup za 1
			//defer znaci odlozi izvrsenje funkcije do kraja trenutne funkcije. Kada gorutina zavrsi posao ili se desi greska, wg.Done() ce biti automatski pozvan

			movieProvider, err := s.movieProviderFactory.Get(source)
			if err != nil {
				fmt.Println(err)
				return
			}

			movie, err := movieProvider.FetchMovie(movieID)
			if err != nil {
				fmt.Println("error during fetch operation:", err)
				return
			}

			results <- movie
			fmt.Printf("%s done %s\n", source, movieID)

		}(id, source)

	}

	//gorutina koja zatvara kanal kada su sve gorutine zavrsile
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All goroutines finished, closing channel")
	}()

	//mapiranje u Movie struct
	var movies []types.Movie
	for movie := range results {
		movies = append(movies, movie)

		if _, err := s.repo.CreateMovie(movie); err != nil {
			fmt.Println("Error inserting movie: ", err)
		} else {
			fmt.Println("movie inserted")
		}
	}

	// //batch insert u bazu
	// if err := s.repo.CreateMovies(movies); err != nil {
	// 	return nil, err
	// }
	return movies, nil
}

// func fetchOMDb(id string) string {
// 	apiKey := "5257cb6d" //skloni ga u .env ili posto ti sale reko da koriste .json, njega koristi
// 	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", id, apiKey)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return fmt.Sprintf("error: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Sprintf("OMDb read error: %v", err)
// 	}

// 	return string(body)
// }

// func fetchTMDb(id string) string {
// 	apiKey := "081ceea64560bb90d333e9fb727f1927"
// 	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", id, apiKey)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return fmt.Sprintf("error: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Sprintf("TMDb read error: %v", err)
// 	}

// 	return string(body)
// }

// // parsira raw JSON iz Result.Data u tvoj Movie struct
// func mapToMovieStruct(res types.FetchResult) (types.Movie, error) {
// 	var movie types.Movie

// 	switch res.Source {
// 	case "OMDb":
// 		var omdbData struct {
// 			Title      string `json:"Title"`
// 			Year       string `json:"Year"`
// 			Genre      string `json:"Genre"`
// 			Director   string `json:"Director"`
// 			Plot       string `json:"Plot"`
// 			IMDBRating string `json:"imdbRating"`
// 		}
// 		if err := json.Unmarshal([]byte(res.Data), &omdbData); err != nil {
// 			fmt.Printf("OMDb JSON unmarshal failed for ID %s: %v\n", res.ID, err)
// 			return movie, err
// 		}

// 		// Parse Year i IMDBRating u odgovarajuće tipove
// 		year, _ := strconv.Atoi(omdbData.Year)
// 		var imdbRating float64
// 		if r, err := strconv.ParseFloat(omdbData.IMDBRating, 64); err == nil {
// 			imdbRating = r
// 		}

// 		movie = types.Movie{
// 			Title:      omdbData.Title,
// 			Year:       year,
// 			Genre:      omdbData.Genre,
// 			Director:   omdbData.Director,
// 			Plot:       omdbData.Plot,
// 			IMDBRating: imdbRating,
// 			Source:     res.Source,
// 			ExternalID: res.ID,
// 		}

// 	case "TMDb":
// 		var tmdbData struct {
// 			Title  string `json:"title"`
// 			Year   string `json:"release_date"`
// 			Genres []struct {
// 				Name string `json:"name"`
// 			} `json:"genres"`
// 			Director string  `json:"director"` // TMDb možda nema direktno, ostaviti prazno ili kasnije fetch-ovati crew
// 			Plot     string  `json:"overview"`
// 			Rating   float64 `json:"vote_average"`
// 		}
// 		if err := json.Unmarshal([]byte(res.Data), &tmdbData); err != nil {
// 			fmt.Printf("TMDb JSON unmarshal failed for ID %s: %v\n", res.ID, err)
// 			return movie, err
// 		}

// 		year := 0
// 		if len(tmdbData.Year) >= 4 {
// 			year, _ = strconv.Atoi(tmdbData.Year[:4])
// 		}

// 		var genreNames []string
// 		for _, g := range tmdbData.Genres {
// 			genreNames = append(genreNames, g.Name)
// 		}

// 		if tmdbData.Director == "" {
// 			tmdbData.Director = "Unknown" // ili npr. "N/A"
// 		}

// 		movie = types.Movie{
// 			Title:      tmdbData.Title,
// 			Year:       year,
// 			Genre:      strings.Join(genreNames, ", "),
// 			Director:   tmdbData.Director,
// 			Plot:       tmdbData.Plot,
// 			IMDBRating: tmdbData.Rating,
// 			Source:     res.Source,
// 			ExternalID: res.ID,
// 		}
// 	}

// 	// Debug štampanje za proveru šta se kreiralo
// 	fmt.Printf("Mapped Movie: %+v\n\n", movie)
// 	if movie.Title == "" || movie.Year == 0 {
// 		fmt.Printf("⚠️ Warning: Movie seems invalid! ID=%s, Source=%s\n", res.ID, res.Source)
// 	}

// 	return movie, nil
// }

package services

import (
	"fmt"
	"sync"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/config"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/external"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/utils"
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

func (s *MovieService) GetAllMovies(page int, limit int) ([]types.Movie, error) {
	movies, err := s.repo.GetAllMovies(page, limit)

	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) DeleteMovieById(id int64) error {
	return s.repo.DeleteMovie(id)
}

func (s *MovieService) UpdateMovie(id int64, payload types.UpdateMoviePayloadPUT) (*types.Movie, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		validationErr := types.ValidationError{Err: err}
		return nil, validationErr
	}

	exists, err := s.repo.ExistsMovieExcludingID(payload.Title, payload.Director, id)
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

	return s.repo.UpdateMovie(id, movie)
}

func (s *MovieService) PatchMovie(id int64, payload types.PatchMoviePayload) (*types.Movie, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		validationErr := types.ValidationError{Err: err}
		return nil, validationErr
	}

	//zbog provere duplikata
	existingMovie, err := s.repo.GetMovieById(id)
	if err != nil {
		return nil, err
	}
	title := existingMovie.Title
	director := existingMovie.Director

	if payload.Title != nil {
		title = *payload.Title
	}
	if payload.Director != nil {
		director = *payload.Director
	}

	if payload.Title != nil || payload.Director != nil {
		exists, err := s.repo.ExistsMovieExcludingID(title, director, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, types.ErrMovieAlreadyExists
		}
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

	if params.Limit > 100 {
		params.Limit = 100
	}

	return s.repo.SearchMovies(params)
}

func (s *MovieService) FetchAndStoreMovies() ([]types.Movie, error) {
	movieSources, err := config.LoadMovieSources()
	if err != nil {
		return nil, err
	}

	results := make(chan types.Movie, len(movieSources))
	var wg sync.WaitGroup

	for id, source := range movieSources {
		wg.Add(1)

		go func(movieID string, apiSource string) {
			defer wg.Done()

			movieProvider, err := s.movieProviderFactory.Get(apiSource)
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
			fmt.Printf("%s done %s\n", apiSource, movieID)

		}(id, source)

	}

	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All goroutines finished, closing channel")
	}()

	var movies []types.Movie
	for movie := range results {
		movies = append(movies, movie)

		if _, err := s.repo.CreateMovie(movie); err != nil {
			fmt.Println("Error inserting movie: ", err)
		} else {
			fmt.Println("movie inserted")
		}
	}

	return movies, nil
}

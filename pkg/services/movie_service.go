package services

import (
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/utils"
)

type MovieService struct {
	repo types.MovieStore
}

func NewService(repo types.MovieStore) *MovieService {
	return &MovieService{repo: repo}
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

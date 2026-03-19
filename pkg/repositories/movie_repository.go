package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateMovie(movie types.Movie) (*types.Movie, error) { //we receive smth of type Movie, and also return that same book that we created
	result := s.db.Create(&movie)

	if err := result.Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Store) CreateMovies(movies []types.Movie) error {
	if len(movies) == 0 {
		return nil
	}

	result := s.db.Create(&movies) // batch insert
	if result.Error != nil {
		return result.Error
	}

	fmt.Println("Rows affected:", result.RowsAffected)
	return nil
}

func (s *Store) GetAllMovies(page int, limit int) ([]types.Movie, error) {
	var movies []types.Movie

	offset := (page - 1) * limit
	result := s.db.
		Limit(limit).
		Offset(offset).
		Find(&movies) //ignorise soft-deleted zapise

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", result.Error)
	}
	return movies, nil
}

func (s *Store) GetMovieById(id int64) (*types.Movie, error) {
	var movie types.Movie
	//result := s.db.Where("ID=?", id).Find(&movie)
	result := s.db.First(&movie, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, types.ErrMovieNotFound
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &movie, nil
}

func (s *Store) DeleteMovie(id int64) error {
	//s.db.Unscoped().Delete(...) --> hard delete
	result := s.db.Delete(&types.Movie{}, id) //DELETE FROM movies WHERE ID=id;
	if err := result.Error; err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return types.ErrMovieNotFound
	}

	return nil
}

func (s *Store) UpdateMovie(id int64, movie types.Movie) (*types.Movie, error) {

	var oldMovie types.Movie

	if err := s.db.First(&oldMovie, id).Error; err != nil {
		return nil, types.ErrMovieNotFound
	}

	if err := s.db.Model(&oldMovie).Updates(movie).Error; err != nil {
		//kad prosledjujes struct u Updates, azuriraju se samo non-zero polja po defaultu
		return nil, err
	}

	return &oldMovie, nil
}

func (s *Store) UpdateMovieFields(id int64, updates map[string]any) (*types.Movie, error) {

	var movie types.Movie

	if err := s.db.First(&movie, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrMovieNotFound
		}
		return nil, err
	}

	//using map to update atributes or Select to specify fields do update
	//with map Updates function updates only specified fields

	if len(updates) > 0 {
		if err := s.db.Model(&movie).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &movie, nil
}

func (s *Store) SearchMovies(params types.MovieSearchAndFilterParams) ([]types.Movie, error) {
	var movies []types.Movie
	db := s.db.Model(&types.Movie{})

	if params.Title != nil && *params.Title != "" {
		//case insensitive search omogucen
		db = db.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(*params.Title)+"%")
	}
	//za filtere zelim exact match, jer bi na frontendu to bili ponudjene opcije za filtriranje
	if params.Genre != nil && *params.Genre != "" {
		db = db.Where("LOWER(genre) LIKE ?", "%"+strings.ToLower(*params.Genre)+"%") //ovde nije striktno, jer mi je genre string koji moze da sadrzi vise, ali ako se projekat zakomplikuje to bih izmenila
	}
	if params.Director != nil && *params.Director != "" {
		db = db.Where("director = ?", *params.Director)
	}

	if params.Year != nil && *params.Year != 0 {
		db = db.Where("year = ?", *params.Year)
	}

	offset := (params.Page - 1) * params.Limit
	result := db.
		Offset(offset).
		Limit(params.Limit).
		Find(&movies)

	if result.Error != nil {
		return nil, result.Error
	}

	return movies, nil
}

func (s *Store) ExistsMovie(title string, director string) (bool, error) {
	var count int64
	if err := s.db.Model(&types.Movie{}).
		Where("title = ? AND director = ?", title, director).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) ExistsMovieExcludingID(title string, director string, excludeID int64) (bool, error) {
	var count int64
	if err := s.db.Model(&types.Movie{}).
		Where("title = ? AND director = ? AND id <> ?", title, director, excludeID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

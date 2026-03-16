package types

import (
	"errors"
	"fmt"
)

var ErrMovieNotFound = errors.New("movie not found")
var ErrMovieAlreadyExists = errors.New("movie with this title and director already exists")

type ValidationError struct {
	Err error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Err)
}

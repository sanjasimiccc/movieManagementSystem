package utils

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

func ParsePaginationParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	page := 1
	limit := 10
	var err error

	if pageStr := query.Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return 0, 0, errors.New("invalid page")
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, errors.New("invalid limit")
		}

		maxLimit := 100
		if limit > maxLimit {
			limit = maxLimit
		}
	}

	return page, limit, nil
}

func ParseMovieSearchParams(r *http.Request) (types.MovieSearchAndFilterParams, error) {
	params := types.MovieSearchAndFilterParams{}

	page, limit, err := ParsePaginationParams(r)
	if err != nil {
		return params, err
	}
	params.Page = page
	params.Limit = limit

	query := r.URL.Query()

	if title := query.Get("title"); title != "" {
		params.Title = &title
	}
	if genre := query.Get("genre"); genre != "" {
		params.Genre = &genre
	}
	if director := query.Get("director"); director != "" {
		params.Director = &director
	}
	if yearStr := query.Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return params, errors.New("invalid year")
		}
		params.Year = &year
	}
	return params, nil
}

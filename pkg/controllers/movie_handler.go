package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
	utils "github.com/sanjasimiccc/movieManagementSystem/pkg/utils/json"
)

type MovieHandler struct {
	store   types.MovieStore
	service types.MovieService
}

func NewHandler(store types.MovieStore, service types.MovieService) *MovieHandler {
	return &MovieHandler{store: store, service: service}
}

func (h *MovieHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/movies/", h.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/", h.GetMovies).Methods("GET")
	router.HandleFunc("/movies/{movieId}", h.GetMovieById).Methods("GET")
	router.HandleFunc("/movies/{movieId}", h.UpdateWholeMovie).Methods("PUT")
	router.HandleFunc("/movies/{movieId}", h.UpdateMoviePartially).Methods("PATCH")
	router.HandleFunc("/movies/{movieId}", h.DeleteMovieById).Methods("DELETE")
	router.HandleFunc("/movies/search/", h.SearchMovies).Methods("GET")
	router.HandleFunc("/movies/fetchData/", h.FetchData).Methods("POST")
}

func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := h.service.GetAllMovies()

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //da bismo pristupili movieId-ju koji je i okviru request-a; za trenutni request vraca promenljive iz rute
	movieId := vars["movieId"]

	id, err := strconv.ParseInt(movieId, 10, 64) //ili strconv.Atoin(movieId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid movie id"))
		return
	}

	movieDetails, err := h.service.GetMovieById(id)

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, movieDetails)
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {

	var payload types.CreateMoviePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	movie, err := h.service.CreateMovie(payload)

	if err != nil {
		if errors.Is(err, types.ErrMovieAlreadyExists) {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		var validationErr types.ValidationError
		if errors.As(err, &validationErr) {
			utils.WriteError(w, http.StatusBadRequest, validationErr)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieId := vars["movieId"]

	id, err := strconv.ParseInt(movieId, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid movie id"))
		return
	}

	err = h.service.DeleteMovieById(id)

	if err != nil {
		if errors.Is(err, types.ErrMovieNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "movie deleted successfully",
	})
}

func (h *MovieHandler) UpdateWholeMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	movieIDStr := vars["movieId"]
	movieIDInt, err := strconv.ParseInt(movieIDStr, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid movie id"))
		return
	}

	var moviePayload types.UpdateMoviePayloadPUT
	if err := utils.ParseJSON(r, &moviePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	movie, err := h.service.UpdateMovie(movieIDInt, moviePayload)

	if err != nil {
		if errors.Is(err, types.ErrMovieNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		// if errors.Is(err, types.ErrMovieAlreadyExists) {
		// 	utils.WriteError(w, http.StatusBadRequest, err)
		// 	return
		// }
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movie)
}

func (h *MovieHandler) UpdateMoviePartially(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["movieId"]

	idInt, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid movie id"))
		return
	}

	var moviePayload types.PatchMoviePayload
	if err := utils.ParseJSON(r, &moviePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	movie, err := h.service.PatchMovie(idInt, moviePayload)

	if err != nil {
		var validationErr types.ValidationError

		if errors.As(err, &validationErr) {
			utils.WriteError(w, http.StatusBadRequest, validationErr)
			return
		}
		if errors.Is(err, types.ErrMovieNotFound) {
			utils.WriteError(w, http.StatusNotFound, err)
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movie)

}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	// citanje query parametara
	query := r.URL.Query()

	params := types.MovieSearchAndFilterParams{
		//postavi default vrednosti
		Page:  1,
		Limit: 2,
	}
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
		if err == nil {
			params.Year = &year
		}
	}

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			params.Page = page
			fmt.Println(page)
		}
	}

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			params.Limit = limit
			fmt.Println(limit)
		}
	}

	movies, err := h.service.SearchMovies(params)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movies)
}

func (h *MovieHandler) FetchData(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.FetchAndStoreMovies()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, movies)
}

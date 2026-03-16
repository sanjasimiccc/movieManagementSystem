package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/utils"
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
	router.HandleFunc("/movies/fetchData/", h.FetchMovieData).Methods("GET") //POST
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

func (h *MovieHandler) FetchMovieData(w http.ResponseWriter, r *http.Request) {
	ids := []string{"tt0111161", "tt0137523"}

	type Result struct {
		Source string
		ID     string
		Data   string
	}

	results := make(chan Result)
	var wg sync.WaitGroup //kaze Main gorutini da saceka da sve gorutine zavrse posao pre nego sto nastavi dalje

	for _, id := range ids {
		wg.Add(2) //postavlja broj gorutina i thread-ova koje main thread treba da ceka. Uvecava brojac za prosledjenu vrednost

		//OMDb gorutina
		go func(movieID string) {
			fmt.Println("OMDb start", movieID)
			defer wg.Done() //kada gorutina zavrsi posao, mora reci da je zavrsila sa wg.Done(), to smanjuje brojac u WaitGroup za 1
			//defer znaci odlozi izvrsenje funkcije do kraja trenutne funkcije. Kada gorutina zavrsi posao ili se desi greska, wg.Done() ce biti automatski pozvan
			time.Sleep(5 * time.Second)
			data := fetchOMDb(movieID)
			fmt.Println("omdb Sending result to channel:", movieID)
			results <- Result{Source: "OMDb", ID: movieID, Data: data}
			fmt.Println("OMDb done:", movieID)
		}(id)

		//TMdb gorutina
		go func(movieID string) {
			fmt.Println("TMDb start", movieID)
			defer wg.Done()
			time.Sleep(5 * time.Second)
			data := fetchTMDb(movieID)
			fmt.Println("tmdb Sending result to channel:", movieID)
			results <- Result{Source: "TMDb", ID: movieID, Data: data}
			fmt.Println("TMDb done:", movieID)
		}(id)
	}

	//gorutina koja zatvara kanal kada su sve gorutine zavrsile
	go func() {
		wg.Wait()
		fmt.Println("All goroutines finished, closing channel")
		close(results)
	}()

	// //ispis rezultata u konzolu
	// for res := range results {
	// 	fmt.Printf("Result from source %s for movie %s:\n%s\n\n", res.Source, res.ID, res.Data)
	// }

	var allResults []Result
	for res := range results {
		fmt.Printf("Result from source %s for movie %s:\n%s\n\n", res.Source, res.ID, res.Data)
		allResults = append(allResults, res)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(allResults); err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}

func fetchOMDb(id string) string {
	apiKey := "5257cb6d" //skloni ga u .env ili posto ti sale reko da koriste .json, njega koristi
	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", id, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

func fetchTMDb(id string) string {
	apiKey := "081ceea64560bb90d333e9fb727f1927"
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", id, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

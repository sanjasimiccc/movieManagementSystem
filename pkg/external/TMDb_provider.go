package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

type TMDbProvider struct {
	apiKey string
}

func NewTMDbProvider(apiKey string) *TMDbProvider {
	return &TMDbProvider{apiKey: apiKey}
}

func (t *TMDbProvider) FetchMovie(id string) (types.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?api_key=%s", id, t.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return types.Movie{}, err
	}

	defer resp.Body.Close()

	var data types.TMDbFetchData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { //ili Unmarshal
		return types.Movie{}, err
	}

	year := 0
	if len(data.Year) >= 4 {
		year, _ = strconv.Atoi(data.Year[:4])
	}

	var genres []string
	for _, g := range data.Genres {
		genres = append(genres, g.Name)
	}

	source := "TMDb"
	externalID := id

	return types.Movie{
		Title:      data.Title,
		Year:       year,
		Genre:      strings.Join(genres, ", "),
		Director:   "Unknown",
		Plot:       data.Plot,
		IMDBRating: data.Rating,
		Source:     &source,
		ExternalID: &externalID,
	}, nil
}

package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sanjasimiccc/movieManagementSystem/pkg/types"
)

type OMDbProvider struct {
	apiKey string
}

func NewOMDbProvider(apiKey string) *OMDbProvider {
	return &OMDbProvider{apiKey: apiKey}
}

func (o *OMDbProvider) FetchMovie(id string) (types.Movie, error) {
	//apiKey := "5257cb6d" //skloni ga u .env ili posto ti sale reko da koriste .json, njega koristi
	url := fmt.Sprintf("http://www.omdbapi.com/?i=%s&apikey=%s", id, o.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return types.Movie{}, err
	}

	defer resp.Body.Close()

	var data types.OMDbFetchData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { //ili Unmarshal
		return types.Movie{}, err
	}

	year, _ := strconv.Atoi(data.Year)
	rating, _ := strconv.ParseFloat(data.IMDBRating, 64)

	return types.Movie{
		Title:      data.Title,
		Year:       year,
		Genre:      data.Genre,
		Director:   data.Director,
		Plot:       data.Plot,
		IMDBRating: rating,
		Source:     "OMDb",
		ExternalID: id,
	}, nil
}

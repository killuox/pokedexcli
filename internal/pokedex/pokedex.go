package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/killuox/pokedexcli/internal/pokecache"
)

const BaseUrl = "https://pokeapi.co/api/v2"

type PokedexConfig struct {
	Previous string
	Next     string
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Response struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

func Get(url string) (Response, error) {
	var resBody []byte
	fmt.Print(url)
	cache := pokecache.NewCache(5 * time.Second)
	cachedData, ok := cache.Get(url)
	if ok {
		resBody = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return Response{}, fmt.Errorf("Error getting locations: %s", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return Response{}, fmt.Errorf("Error reading Body: %s", err)
		}

		cache.Add(url, body)
		resBody = body
	}

	var response Response
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return Response{}, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	return response, nil
}

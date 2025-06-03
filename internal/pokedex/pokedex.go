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
	PreviousLocation string
	NextLocation     string
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationsResponse struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type LocationResponse struct {
	Id                int         `json:"id"`
	PokemonEncounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonResponse struct {
	BaseExperience int `json:"base_experience"`
	Id             int `json:"id"`
}

func GetLocations(url string) (LocationsResponse, error) {
	var resBody []byte
	cache := pokecache.NewCache(5 * time.Second)
	cachedData, ok := cache.Get(url)
	if ok {
		resBody = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("Error getting locations: %s", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationsResponse{}, fmt.Errorf("Error reading Body: %s", err)
		}

		cache.Add(url, body)
		resBody = body
	}

	var response LocationsResponse
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return LocationsResponse{}, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	return response, nil
}

func GetLocation(name string) (LocationResponse, error) {
	url := fmt.Sprintf("%s/location-area/%s", BaseUrl, name)
	var resBody []byte
	cache := pokecache.NewCache(5 * time.Second)
	cachedData, ok := cache.Get(url)
	if ok {
		resBody = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("Error getting location: %s", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return LocationResponse{}, fmt.Errorf("Error reading Body: %s", err)
		}

		cache.Add(url, body)
		resBody = body
	}

	var response LocationResponse
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	return response, nil
}

func GetPokemon(name string) (PokemonResponse, error) {
	var resBody []byte
	cache := pokecache.NewCache(5 * time.Second)
	url := fmt.Sprintf("%s/pokemon/%s", BaseUrl, name)
	cachedData, ok := cache.Get(url)
	if ok {
		resBody = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("Error getting pokemon: %s", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return PokemonResponse{}, fmt.Errorf("Error reading Body: %s", err)
		}

		cache.Add(url, body)
		resBody = body
	}

	var response PokemonResponse
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return PokemonResponse{}, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	return response, nil
}

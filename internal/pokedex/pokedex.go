package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/killuox/pokedexcli/internal/pokecache"
)

const BaseUrl = "https://pokeapi.co/api/v2"

type PokedexConfig struct {
	PreviousLocation string
	NextLocation     string
	Inventory        *UserInventory
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
	Name           string         `json:"name"`
	Url            string         `json:"url"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}

type PokemonStats struct {
	BaseStat int         `json:"base_stat"`
	Effort   int         `json:"effort"`
	Stat     PokemonStat `json:"stat"`
}

type PokemonStat struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonTypes struct {
	Slot int         `json:"slot"`
	Type PokemonType `json:"type"`
}

type PokemonType struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type UserInventory struct {
	Pokemons map[string]Pokemon
}

const (
	maxPossibleBaseExperience = 300
	baseCatchRate             = 0.50
	minCatchRate              = 0.25
)

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

func GetPokemon(name string) (Pokemon, error) {
	var resBody []byte
	cache := pokecache.NewCache(5 * time.Second)
	url := fmt.Sprintf("%s/pokemon/%s", BaseUrl, name)
	cachedData, ok := cache.Get(url)
	if ok {
		resBody = cachedData
	} else {
		res, err := http.Get(url)
		if err != nil {
			return Pokemon{}, fmt.Errorf("Error getting pokemon: %s", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return Pokemon{}, fmt.Errorf("Error reading Body: %s", err)
		}

		cache.Add(url, body)
		resBody = body
	}

	var response Pokemon
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		return Pokemon{}, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	return response, nil
}

func CalculateCatchChance(p Pokemon) float64 {
	effectiveBaseExp := float64(p.BaseExperience)
	if effectiveBaseExp < 0 {
		effectiveBaseExp = 0
	}
	if effectiveBaseExp > maxPossibleBaseExperience {
		effectiveBaseExp = maxPossibleBaseExperience
	}

	scalingFactor := 1.0 - (effectiveBaseExp / maxPossibleBaseExperience)

	chance := baseCatchRate * scalingFactor

	if chance < minCatchRate {
		return minCatchRate
	}

	return chance
}

func TryToCatch(p Pokemon) bool {
	catchProbability := CalculateCatchChance(p)

	randomRoll := rand.Float64()
	if randomRoll < catchProbability {
		return true
	}

	return false
}

func NewUserInventory() *UserInventory {
	userPokedex := &UserInventory{
		Pokemons: make(map[string]Pokemon),
	}

	return userPokedex
}

func (u *UserInventory) Add(p Pokemon) {
	u.Pokemons[p.Name] = p
}

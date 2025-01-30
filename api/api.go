package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedex/cache"
	"pokedex/utils"
)

func FetchLocationAreas(url *string, cache *cache.Cache) (*utils.Response, error) {
	var locationResponse utils.Response
	res, ok := cache.Get(*url)
	if ok {
		fmt.Println("CACHE HIT!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		if err := json.Unmarshal(res, &locationResponse); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}

		return &locationResponse, nil
	}

	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("error fetching location areas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, &locationResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	cache.Add(*url, body)

	return &locationResponse, nil
}

func FetchPokemon(url *string) (*utils.Pokemon, error) {
	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("error fetching location areas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var PokeResponse utils.Pokemon
	if err := json.Unmarshal(body, &PokeResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &PokeResponse, nil
}

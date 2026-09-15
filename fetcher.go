package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	httpTimeout = 10 * time.Second
	userAgent   = "AstralKernel-Go/2.0"
)

type Post struct {
	Title     string `json:"title"`
	Slug      struct{ Current string `json:"current"` } `json:"slug"`
	MainImage struct{ Asset struct{ URL string `json:"url"` } `json:"asset"` } `json:"mainImage"`
}

func fetchPosts(apiURL, header, secret string) ([]Post, error) {
	client := &http.Client{Timeout: httpTimeout}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creando petición: %w", err)
	}
	req.Header.Add(header, "Bearer "+secret)
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error en la petición: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API devolvió status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("leyendo respuesta: %w", err)
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("parseando JSON: %w", err)
	}

	return posts, nil
}
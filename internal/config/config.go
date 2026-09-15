package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	APIURL  string
	PostURL string
	Header  string
	Secret  string
	BaseURL string
}

func Load() (Config, error) {
	apiURL := os.Getenv("BASE_API_URL")
	postURL := strings.TrimSuffix(os.Getenv("BASE_POST_URL"), "/")
	header := os.Getenv("HEADERS")
	secret := os.Getenv("SECRET_KEY")
	baseURL := os.Getenv("BASE_URL")

	if apiURL == "" || postURL == "" || header == "" || secret == "" || baseURL == "" {
		return Config{}, fmt.Errorf("falta configurar variables de entorno")
	}

	return Config{
		APIURL:  apiURL,
		PostURL: postURL,
		Header:  header,
		Secret:  secret,
		BaseURL: baseURL,
	}, nil
}
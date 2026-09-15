package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	fileName  = "README.md"
	startTag  = "<!-- ARTICLES:START -->"
	endTag    = "<!-- ARTICLES:END -->"
	maxPosts  = 3
	httpTimeout = 10 * time.Second
	userAgent = "AstralKernel-Go/2.0"
)

type Post struct {
	Title     string `json:"title"`
	Slug      struct{ Current string `json:"current"` } `json:"slug"`
	MainImage struct{ Asset struct{ URL string `json:"url"` } `json:"asset"` } `json:"mainImage"`
}

type Config struct {
	APIURL    string
	PostURL   string
	Header    string
	Secret    string
	BaseURL   string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	posts, err := fetchPosts(cfg.APIURL, cfg.Header, cfg.Secret)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if len(posts) < maxPosts {
		fmt.Printf("⚠️ Menos de %d artículos recibidos (got %d)\n", maxPosts, len(posts))
		os.Exit(1)
	}

	table := buildTable(posts[:maxPosts], cfg.PostURL)

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("❌ Error leyendo README: %v\n", err)
		os.Exit(1)
	}

	updated, err := updateReadme(content, table)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(fileName, []byte(updated), 0644); err != nil {
		fmt.Printf("❌ Error escribiendo README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Blog section synchronized.")
}

func loadConfig() (Config, error) {
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

func buildTable(posts []Post, postURL string) string {
	if len(posts) != maxPosts {
		return ""
	}

	var imgParts []string
	var titleParts []string

	for _, p := range posts {
		imgParts = append(imgParts, fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s)",
			p.Title, p.MainImage.Asset.URL, postURL, p.Slug.Current))
		titleParts = append(titleParts, fmt.Sprintf("**[%s](%s/%s)**",
			p.Title, postURL, p.Slug.Current))
	}

	imgRow := strings.Join(imgParts, " | ")
	titleRow := strings.Join(titleParts, " | ")

	return imgRow + "\n--- | --- | ---\n" + titleRow + "\n\n[➡️ More blog posts](" + os.Getenv("BASE_URL") + ")"
}

func updateReadme(content []byte, table string) (string, error) {
	updated := string(content)

	reBlock := regexp.MustCompile("(?s)" + regexp.QuoteMeta(startTag) + ".*?" + regexp.QuoteMeta(endTag))
	updated = reBlock.ReplaceAllString(updated, startTag+"\n"+table+"\n"+endTag)

	timestamp := time.Now().UTC().Format(time.RFC3339)
	tsLine := fmt.Sprintf("<!-- Last updated: %s -->", timestamp)
	reTS := regexp.MustCompile(`<!-- Last updated: .*? -->`)

	if reTS.MatchString(updated) {
		updated = reTS.ReplaceAllString(updated, tsLine)
	} else {
		updated += "\n" + tsLine + "\n"
	}

	return updated, nil
}
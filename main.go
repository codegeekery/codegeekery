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

type Post struct {
	Title     string `json:"title"`
	Slug      struct{ Current string `json:"current"` } `json:"slug"`
	MainImage struct{ Asset struct{ URL string `json:"url"` } `json:"asset"` } `json:"mainImage"`
}

func main() {
	// Configuración
	apiURL := os.Getenv("BASE_API_URL")
	postURL := strings.TrimSuffix(os.Getenv("BASE_POST_URL"), "/")
	header := os.Getenv("HEADERS")
	secret := os.Getenv("SECRET_KEY")
	baseURL := os.Getenv("BASE_URL")

	// Validar variables de entorno
	if apiURL == "" || secret == "" || header == "" {
		fmt.Println("❌ Falta configurar variables de entorno")
		os.Exit(1)
	}

	const (
		fileName = "README.md"
		startTag = "<!-- ARTICLES:START -->"
		endTag   = "<!-- ARTICLES:END -->"
	)

	// Fetch con Headers
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Add(header, "Bearer "+secret)
	req.Header.Add("User-Agent", "AstralKernel-Go/2.0")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error en la petición: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("❌ API devolvió status %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var posts []Post
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n", err)
		os.Exit(1)
	}

	if len(posts) < 3 {
		fmt.Println("⚠️ Menos de 3 artículos recibidos")
		os.Exit(1)
	}

	// Construir tabla
	p := posts[:3]
	imgRow := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s)",
		p[0].Title, p[0].MainImage.Asset.URL, postURL, p[0].Slug.Current,
		p[1].Title, p[1].MainImage.Asset.URL, postURL, p[1].Slug.Current,
		p[2].Title, p[2].MainImage.Asset.URL, postURL, p[2].Slug.Current)

	titleRow := fmt.Sprintf("**[%s](%s/%s)** | **[%s](%s/%s)** | **[%s](%s/%s)**",
		p[0].Title, postURL, p[0].Slug.Current,
		p[1].Title, postURL, p[1].Slug.Current,
		p[2].Title, postURL, p[2].Slug.Current)

	tableMarkdown := imgRow + "\n--- | --- | ---\n" + titleRow + "\n\n[➡️ More blog posts](" + baseURL + ")"

	// Actualizar README
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("❌ Error leyendo README: %v\n", err)
		os.Exit(1)
	}

	updated := string(content)

	// Reemplazar artículos
	reBlock := regexp.MustCompile("(?s)" + regexp.QuoteMeta(startTag) + ".*?" + regexp.QuoteMeta(endTag))
	updated = reBlock.ReplaceAllString(updated, startTag+"\n"+tableMarkdown+"\n"+endTag)

	// Actualizar timestamp
	timestamp := time.Now().UTC().Format(time.RFC3339)
	tsLine := fmt.Sprintf("<!-- Last updated: %s -->", timestamp)
	reTS := regexp.MustCompile(`<!-- Last updated: .*? -->`)

	if reTS.MatchString(updated) {
		updated = reTS.ReplaceAllString(updated, tsLine)
	} else {
		updated += "\n" + tsLine + "\n"
	}

	if err := os.WriteFile(fileName, []byte(updated), 0644); err != nil {
		fmt.Printf("❌ Error escribiendo README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Blog section synchronized.")
}
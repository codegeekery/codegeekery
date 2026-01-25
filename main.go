package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Estructura para mapear la respuesta de Astral Kernel (Sanity)
type Post struct {
	Title     string `json:"title"`
	MainImage struct {
		Asset struct {
			URL string `json:"url"`
		} `json:"asset"`
	} `json:"mainImage"`
	Slug struct {
		Current string `json:"current"`
	} `json:"slug"`
}

func main() {
	// 1. Configuración mediante variables de entorno
	secretKey := os.Getenv("SECRET_KEY")
	apiURL := os.Getenv("BASE_API_URL")
	basePostURL := os.Getenv("BASE_POST_URL")
	baseURL := os.Getenv("BASE_URL")
	userAgent := os.Getenv("USER_AGENT")

	if secretKey == "" || apiURL == "" {
		fmt.Println("Error: Configuración incompleta (SECRET_KEY o BASE_API_URL)")
		os.Exit(1)
	}

	// 2. Obtener los artículos desde el Backend (Go)
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("X-ASTRAL-KERNEL", secretKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error conectando con la API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var posts []Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		fmt.Printf("Error parseando JSON: %v\n", err)
		os.Exit(1)
	}

	// Tomamos solo los 3 últimos
	if len(posts) > 3 {
		posts = posts[:3]
	}

	// 3. Construir la tabla de Markdown
	var imgRow, sepRow, titleRow []string
	for _, post := range posts {
		imgURL := fmt.Sprintf("%s?w=200&h=200", post.MainImage.Asset.URL)
		link := basePostURL + post.Slug.Current

		imgRow = append(imgRow, fmt.Sprintf("[![%s](%s)](%s)", post.Title, imgURL, link))
		sepRow = append(sepRow, "---")
		titleRow = append(titleRow, fmt.Sprintf("**[%s](%s)**", post.Title, link))
	}

	// Unimos las filas con el formato de tabla de Markdown
	tableContent := fmt.Sprintf("\n%s\n%s\n%s\n\n[➡️ More blog posts](%s)\n",
		strings.Join(imgRow, " | "),
		strings.Join(sepRow, " | "),
		strings.Join(titleRow, " | "),
		baseURL)

	// 4. Leer el archivo README.md
	readmePath := "README.md"
	content, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Printf("Error leyendo README: %v\n", err)
		os.Exit(1)
	}
	readmeText := string(content)

	// --- LÓGICA DE IDENTIFICACIÓN ---

	// Identificar y reemplazar el bloque de artículos
	// (?s) permite que el '.' incluya saltos de línea
	reArticles := regexp.MustCompile(`(?s).*?`)
	newArticlesBlock := fmt.Sprintf("%s", tableContent)
	updated := reArticles.ReplaceAllString(readmeText, newArticlesBlock)

	// Identificar y reemplazar el timestamp
	// Formato ISO 8601 (2006-01-02T15:04:05Z)
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	reTime := regexp.MustCompile(``)
	newTimeLine := fmt.Sprintf("%s", timestamp)
	
	updated = reTime.ReplaceAllString(updated, newTimeLine)

	// 5. Guardar los cambios
	err = os.WriteFile(readmePath, []byte(updated), 0644)
	if err != nil {
		fmt.Printf("Error al escribir en README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ README.md actualizado con éxito mediante Go.")
}
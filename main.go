package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Post struct {
	Title    string `json:"title"`
	Slug     string `json:"slug"`
	ImageUrl string `json:"imageUrl"`
}

func main() {
	// Consumimos las variables de ambiente
	basePostURL := os.Getenv("BASE_POST_URL")
	baseAPIURL  := os.Getenv("BASE_API_URL")
	
	// Valores por defecto para tags y nombre de archivo
	fileName := "README.md"
	startTag := "<!-- ARTICLES:START -->"
	endTag   := "<!-- ARTICLES:END -->"

	if baseAPIURL == "" || basePostURL == "" {
		fmt.Println("Error: BASE_API_URL o BASE_POST_URL no están definidas")
		return
	}

	// 1. Obtener datos de tu API
	resp, err := http.Get(baseAPIURL)
	if err != nil {
		fmt.Printf("Error llamando a la API: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Printf("Error parseando JSON: %v\n", err)
		return
	}

	// 2. Construir la tabla Markdown
	markdownTable := generateTable(posts, basePostURL)

	// 3. Leer y actualizar el archivo
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("No se pudo leer el archivo: %v\n", err)
		return
	}

	strContent := string(content)
	startIndex := strings.Index(strContent, startTag)
	endIndex := strings.Index(strContent, endTag)

	if startIndex == -1 || endIndex == -1 {
		fmt.Println("No se encontraron los tags o ")
		return
	}

	// Reensamblar el archivo respetando los tags
	newContent := strContent[:startIndex+len(startTag)] +
		"\n" + markdownTable + "\n" +
		strContent[endIndex:]

	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("Error escribiendo el README: %v\n", err)
		return
	}

	fmt.Println("🚀 README actualizado correctamente con variables de entorno!")
}

func generateTable(posts []Post, basePostURL string) string {
	var row1, row2, row3 []string

	for _, p := range posts {
		// Aseguramos que el slash final no se duplique
		fullURL := strings.TrimSuffix(basePostURL, "/") + "/" + p.Slug
		
		row1 = append(row1, fmt.Sprintf("[![%s](%s?w=200&h=200)](%s)", p.Title, p.ImageUrl, fullURL))
		row2 = append(row2, "---")
		row3 = append(row3, fmt.Sprintf("**[%s](%s)**", p.Title, fullURL))
	}

	return strings.Join(row1, " | ") + "\n" + strings.Join(row2, " | ") + "\n" + strings.Join(row3, " | ")
}
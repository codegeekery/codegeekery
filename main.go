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
	Slug     struct{ Current string `json:"current"` } `json:"slug"`
	ImageUrl string `json:"imageUrl"`
}

func main() {
	// 1. Entorno y Configuración
	basePostURL := os.Getenv("BASE_POST_URL")
	baseAPIURL  := os.Getenv("BASE_API_URL")
	headerKey   := os.Getenv("HEADERS")
	secretValue := os.Getenv("SECRET_KEY")
	
	const fileName = "README.md"
	const startTag = "<!-- ARTICLES:START -->"
	const endTag   = "<!-- ARTICLES:END -->"

	// Still missing vars env
	if baseAPIURL == "" || basePostURL == "" || headerKey == "" || secretValue == "" {
		fmt.Println("❌ Faltan variables de entorno")
		return
	}

	// 2. Petición con Seguridad
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseAPIURL, nil)
	req.Header.Add(headerKey, secretValue)

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("❌ Error en API: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 3. Parsear Datos
	var posts []Post
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &posts); err != nil || len(posts) < 3 {
		fmt.Println("❌ Error en datos recibidos")
		return
	}

	// 4. Construir Tabla Simplificada
	// Nota: Verifica que posts[i].ImageUrl sea una URL pública válida.
	cleanBase := strings.TrimSuffix(basePostURL, "/")
	
	markdownTable := fmt.Sprintf(
		"| [%s](%s/%s) | [%s](%s/%s) | [%s](%s/%s) |\n| :---: | :---: | :---: |\n| [![Img](%s?w=250)](%s/%s) | [![Img](%s?w=250)](%s/%s) | [![Img](%s?w=250)](%s/%s) |",
		posts[0].Title, cleanBase, posts[0].Slug.Current,
		posts[1].Title, cleanBase, posts[1].Slug.Current,
		posts[2].Title, cleanBase, posts[2].Slug.Current,
		posts[0].ImageUrl, cleanBase, posts[0].Slug.Current,
		posts[1].ImageUrl, cleanBase, posts[1].Slug.Current,
		posts[2].ImageUrl, cleanBase, posts[2].Slug.Current,
	)

	// 5. Inyectar en README
	content, _ := os.ReadFile(fileName)
	strContent := string(content)
	startIdx := strings.Index(strContent, startTag)
	endIdx := strings.Index(strContent, endTag)

	if startIdx != -1 && endIdx != -1 {
		newContent := strContent[:startIdx+len(startTag)] + "\n\n" + markdownTable + "\n\n" + strContent[endIdx:]
		os.WriteFile(fileName, []byte(newContent), 0644)
		fmt.Println("🚀 README simplificado y actualizado!")
	}
}
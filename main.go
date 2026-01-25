package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Estructura adaptada a Sanity: Slug es un objeto con un campo "current"
type Post struct {
	Title    string `json:"title"`
	Slug     struct {Current string `json:"current"`} `json:"slug"`
	ImageUrl string `json:"imageUrl"`
}

func main() {
	// 1. Cargar Variables de Entorno
	basePostURL := os.Getenv("BASE_POST_URL")
	baseAPIURL  := os.Getenv("BASE_API_URL")
	headerKey   := os.Getenv("HEADERS")
	secretValue := os.Getenv("SECRET_KEY")
	
	fileName := "README.md"
	startTag := "<!-- ARTICLES:START -->"
	endTag   := "<!-- ARTICLES:END -->"

	if baseAPIURL == "" || basePostURL == "" || headerKey == "" || secretValue == "" {
		fmt.Println("❌ Error: Faltan variables de entorno. Revisa BASE_API_URL, BASE_POST_URL, HEADERS y SECRET_KEY")
		return
	}

	// 2. Configurar la petición HTTP con Headers de Seguridad
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseAPIURL, nil)
	if err != nil {
		fmt.Printf("❌ Error creando la petición: %v\n", err)
		return
	}

	// Añadimos la autenticación que pide tu backend
	req.Header.Add(headerKey, secretValue)

	// 3. Ejecutar Petición
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error de red: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️ Error API (Status %d): %s\n", resp.StatusCode, string(body))
		return
	}

	// 4. Parsear JSON (Ahora coincide con el objeto slug de Sanity)
	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n📥 Body recibido: %s\n", err, string(body))
		return
	}

	if len(posts) < 3 {
		fmt.Printf("⚠️ Se esperaban 3 posts, se recibieron %d\n", len(posts))
		return
	}

	// 5. Construir Tabla Markdown (Accediendo a .Slug.Current)
	cleanBase := strings.TrimSuffix(basePostURL, "/")

	row1 := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s)",
		posts[0].Title, posts[0].ImageUrl, cleanBase, posts[0].Slug.Current,
		posts[1].Title, posts[1].ImageUrl, cleanBase, posts[1].Slug.Current,
		posts[2].Title, posts[2].ImageUrl, cleanBase, posts[2].Slug.Current,
	)

	row2 := "--- | --- | ---"

	row3 := fmt.Sprintf("**[%s](%s/%s)** | **[%s](%s/%s)** | **[%s](%s/%s)**",
		posts[0].Title, cleanBase, posts[0].Slug.Current,
		posts[1].Title, cleanBase, posts[1].Slug.Current,
		posts[2].Title, cleanBase, posts[2].Slug.Current,
	)

	markdownTable := row1 + "\n" + row2 + "\n" + row3

	// 6. Actualizar README.md
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("❌ Error leyendo %s: %v\n", fileName, err)
		return
	}
	strContent := string(content)

	startIdx := strings.Index(strContent, startTag)
	endIdx := strings.Index(strContent, endTag)

	if startIdx == -1 || endIdx == -1 {
		fmt.Println("❌ No se encontraron los marcadores en el archivo")
		return
	}

	// Inyectar el nuevo contenido entre los tags
	newContent := strContent[:startIdx+len(startTag)] + "\n\n" + markdownTable + "\n\n" + strContent[endIdx:]

	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error escribiendo archivo: %v\n", err)
		return
	}

	fmt.Println("🚀 README actualizado con éxito para CodeGeekery!")
}
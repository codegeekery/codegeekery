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
	// 1. Obtener Variables de ambiente
	basePostURL := os.Getenv("BASE_POST_URL")
	baseAPIURL := os.Getenv("BASE_API_URL")
	headerName := os.Getenv("HEADERS")
	secretKey := os.Getenv("SECRET_KEY")
	
	fileName := "README.md"
	startTag := "<!-- START POSTS -->"
	endTag := ""

	if baseAPIURL == "" || basePostURL == "" || headerName == "" || secretKey == "" {
		fmt.Println("❌ Error: Faltan variables de entorno (URLs o Credenciales)")
		return
	}

	// 2. Configurar la petición con Headers
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseAPIURL, nil)
	if err != nil {
		fmt.Printf("❌ Error creando request: %v\n", err)
		return
	}

	// Inyectamos el Header Especial
	req.Header.Add(headerName, secretKey)

	// 3. Ejecutar la petición
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error de red: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️ API Error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n", err)
		return
	}

	if len(posts) < 3 {
		fmt.Println("⚠️ La API devolvió menos de 3 posts.")
		return
	}

	// 4. Construir la tabla (Sin usar FOR)
	cleanBase := strings.TrimSuffix(basePostURL, "/")
	row1 := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s)",
		posts[0].Title, posts[0].ImageUrl, cleanBase, posts[0].Slug,
		posts[1].Title, posts[1].ImageUrl, cleanBase, posts[1].Slug,
		posts[2].Title, posts[2].ImageUrl, cleanBase, posts[2].Slug,
	)
	row2 := "--- | --- | ---"
	row3 := fmt.Sprintf("**[%s](%s/%s)** | **[%s](%s/%s)** | **[%s](%s/%s)**",
		posts[0].Title, cleanBase, posts[0].Slug,
		posts[1].Title, cleanBase, posts[1].Slug,
		posts[2].Title, cleanBase, posts[2].Slug,
	)

	markdownTable := row1 + "\n" + row2 + "\n" + row3

	// 5. Actualizar el README.md
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("❌ No se pudo leer el archivo README.md")
		return
	}
	strContent := string(content)

	startIdx := strings.Index(strContent, startTag)
	endIdx := strings.Index(strContent, endTag)

	if startIdx == -1 || endIdx == -1 {
		fmt.Println("❌ No se encontraron los marcadores en el README")
		return
	}

	newContent := strContent[:startIdx+len(startTag)] + "\n\n" + markdownTable + "\n\n" + strContent[endIdx:]

	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error escribiendo archivo: %v\n", err)
		return
	}

	fmt.Println("🚀 README actualizado con éxito usando autenticación!")
}
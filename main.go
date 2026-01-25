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
	// Variables de ambiente
	basePostURL := os.Getenv("BASE_POST_URL")
	baseAPIURL := os.Getenv("BASE_API_URL")
	
	fileName := "README.md"
	startTag := ""
	endTag := ""

	if baseAPIURL == "" || basePostURL == "" {
		fmt.Println("❌ Error: BASE_API_URL o BASE_POST_URL no están definidas")
		return
	}

	// 1. Obtener datos de Astral Kernel
	resp, err := http.Get(baseAPIURL)
	if err != nil {
		fmt.Printf("❌ Error de red: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Manejo del error del JSON (el famoso carácter 'A')
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("⚠️ API Error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var posts []Post
	if err := json.Unmarshal(body, &posts); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n📥 Recibido: %s\n", err, string(body))
		return
	}

	// Verificamos que existan los 3 posts para evitar pánico al acceder por índice
	if len(posts) < 3 {
		fmt.Println("⚠️ La API devolvió menos de 3 posts.")
		return
	}

	// 2. Construir la tabla (Sin usar FOR)
	cleanBase := strings.TrimSuffix(basePostURL, "/")

	// Fila 1: Imágenes
	row1 := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s)",
		posts[0].Title, posts[0].ImageUrl, cleanBase, posts[0].Slug,
		posts[1].Title, posts[1].ImageUrl, cleanBase, posts[1].Slug,
		posts[2].Title, posts[2].ImageUrl, cleanBase, posts[2].Slug,
	)

	row2 := "--- | --- | ---"

	// Fila 3: Títulos
	row3 := fmt.Sprintf("**[%s](%s/%s)** | **[%s](%s/%s)** | **[%s](%s/%s)**",
		posts[0].Title, cleanBase, posts[0].Slug,
		posts[1].Title, cleanBase, posts[1].Slug,
		posts[2].Title, cleanBase, posts[2].Slug,
	)

	markdownTable := row1 + "\n" + row2 + "\n" + row3

	// 3. Actualizar el README.md
	content, _ := os.ReadFile(fileName)
	strContent := string(content)

	// Buscamos la posición de los tags
	startIdx := strings.Index(strContent, startTag)
	endIdx := strings.Index(strContent, endTag)

	if startIdx == -1 || endIdx == -1 {
		fmt.Println("❌ No se encontraron los marcadores en el README")
		return
	}

	// Reemplazamos el bloque completo manteniendo los marcadores
	// Usamos slicing para inyectar la tabla justo en medio
	newContent := strContent[:startIdx+len(startTag)] + "\n\n" + markdownTable + "\n\n" + strContent[endIdx:]

	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		fmt.Printf("❌ Error escribiendo archivo: %v\n", err)
		return
	}

	fmt.Println("🚀 README actualizado con éxito y sin bucles for!")
}
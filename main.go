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

func main() {
	// 1. Configuración de API y Cliente
	apiURL := os.Getenv("BASE_API_URL")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("X-ASTRAL-KERNEL", os.Getenv("SECRET_KEY"))
	req.Header.Set("User-Agent", os.Getenv("USER_AGENT"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error llamando a la API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()


	var posts []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		fmt.Printf("Error decodificando JSON: %v\n", err)
		os.Exit(1)
	}

	// 3. Construir la tabla Markdown (Top 3)
	var imgs, seps, titles []string
	basePostURL := os.Getenv("BASE_POST_URL")

	for i, p := range posts {
		if i >= 3 {
			break
		}

		slug := p["slug"].(map[string]interface{})["current"].(string)
		mainImg := p["mainImage"].(map[string]interface{})["asset"].(map[string]interface{})["url"].(string)
		title := p["title"].(string)
		link := basePostURL + slug

		imgs = append(imgs, fmt.Sprintf("[![%s](%s?w=200&h=200)](%s)", title, mainImg, link))
		seps = append(seps, "---")
		titles = append(titles, fmt.Sprintf("**[%s](%s)**", title, link))
	}


	readme, _ := os.ReadFile("README.md")

	content := string(readme)

	table := fmt.Sprintf("\n%s\n%s\n%s\n\n[➡️ More blog posts](%s)\n", strings.Join(imgs, " | "), strings.Join(seps, " | "), strings.Join(titles, " | "), os.Getenv("BASE_URL"))

	// Reemplazar bloque de artículos y actualizar timestamp
	re := regexp.MustCompile(`(?s).*?`)
	content = re.ReplaceAllString(content, table)

	ts := fmt.Sprintf("Last updated: %s", time.Now().Format(time.RFC3339))
	reTS := regexp.MustCompile(`Last updated: .+`)

	if reTS.MatchString(content) {
		content = reTS.ReplaceAllString(content, ts)
	} else {
		content += "\n" + ts + "\n"
	}

	os.WriteFile("README.md", []byte(content), 0644)
	fmt.Println("✅ README.md actualizado con éxito.")
}

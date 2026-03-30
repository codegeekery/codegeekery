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
	// Configuración (Variables de entorno)
	apiURL := os.Getenv("BASE_API_URL")
	postURL := strings.TrimSuffix(os.Getenv("BASE_POST_URL"), "/")
	header := os.Getenv("HEADERS") 
	secret := os.Getenv("SECRET_KEY")

	const (
		fileName = "README.md"
		startTag = "<!-- ARTICLES:START -->"
		endTag   = "<!-- ARTICLES:END -->"
	)

	// 2. Fetch con Headers
	client := &http.Client{}
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Add(header, "Bearer "+secret)
	req.Header.Add("User-Agent", "AstralKernel-Go/2.0")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp)
		return
	}
	defer resp.Body.Close()

	var posts []Post
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &posts)

	if len(posts) < 3 {
		fmt.Println("⚠️ Menos de 3 artículos recibidos")
		return
	}

	// 3. Construir Tabla Markdown
	p := posts[:3]
	imgRow := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s) | [![%s](%s?w=200&h=200)](%s/%s)",
		p[0].Title, p[0].MainImage.Asset.URL, postURL, p[0].Slug.Current,
		p[1].Title, p[1].MainImage.Asset.URL, postURL, p[1].Slug.Current,
		p[2].Title, p[2].MainImage.Asset.URL, postURL, p[2].Slug.Current)

	sepRow := "--- | --- | ---"

	titleRow := fmt.Sprintf("**[%s](%s/%s)** | **[%s](%s/%s)** | **[%s](%s/%s)**",
		p[0].Title, postURL, p[0].Slug.Current,
		p[1].Title, postURL, p[1].Slug.Current,
		p[2].Title, postURL, p[2].Slug.Current)

	tableMarkdown := imgRow + "\n" + sepRow + "\n" + titleRow + "\n\n[➡️ More blog posts](" + os.Getenv("BASE_URL") + ")"

	// 4. Actualizar README y Timestamp
	content, _ := os.ReadFile(fileName)
	updated := string(content)

	// Reemplazar bloque de artículos
	reBlock := regexp.MustCompile("(?s)" + startTag + ".*?" + endTag)
	updated = reBlock.ReplaceAllString(updated, startTag+"\n"+tableMarkdown+"\n"+endTag)

	// Reemplazar o añadir Timestamp
	timestamp := time.Now().Format(time.RFC3339)
	tsLine := fmt.Sprintf("<!-- Last updated: %s -->", timestamp)
	reTS := regexp.MustCompile(`<!-- Last updated: .*? -->`)

	if reTS.MatchString(updated) {
		updated = reTS.ReplaceAllString(updated, tsLine)
	} else {
		updated += "\n" + tsLine + "\n"
	}

	os.WriteFile(fileName, []byte(updated), 0644)
	fmt.Println("✅ Blog section synchronized.")
}
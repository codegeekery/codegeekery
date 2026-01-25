package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
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
		fmt.Printf("Error API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var posts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&posts)


	if len(posts) < 3 {
		fmt.Println("No hay suficientes posts.")
		os.Exit(1)
	}


	basePostURL := os.Getenv("BASE_POST_URL")
	
	
	formatCell := func(p map[string]interface{}) (string, string) {
		slug := p["slug"].(map[string]interface{})["current"].(string)
		mainImg := p["mainImage"].(map[string]interface{})["asset"].(map[string]interface{})["url"].(string)
		title := p["title"].(string)
		link := basePostURL + slug
		
		imgMD := fmt.Sprintf("[![%s](%s?w=200&h=200)](%s)", title, mainImg, link)
		titleMD := fmt.Sprintf("**[%s](%s)**", title, link)
		return imgMD, titleMD
	}

	img0, title0 := formatCell(posts[0])
	img1, title1 := formatCell(posts[1])
	img2, title2 := formatCell(posts[2])


	startTag := ""
	endTag := ""
	
	tableMD := fmt.Sprintf("%s | %s | %s\n--- | --- | ---\n%s | %s | %s", 
		img0, img1, img2, title0, title1, title2)

	newContent := fmt.Sprintf("%s\n%s\n\n[➡️ More blog posts](%s)\n%s", 
		startTag, tableMD, os.Getenv("BASE_URL"), endTag)

	// 4. Reemplazo en el archivo
	readme, _ := os.ReadFile("README.md")
	content := string(readme)

	// Reemplazar bloque de artículos
	re := regexp.MustCompile(`(?s).*?`)
	content = re.ReplaceAllString(content, newContent)

	// Actualizar Timestamp correctamente
	tsLine := fmt.Sprintf("%s", time.Now().Format(time.RFC3339))
	reTS := regexp.MustCompile(``)
	
	if reTS.MatchString(content) {
		content = reTS.ReplaceAllString(content, tsLine)
	} else {
		content += "\n" + tsLine
	}

	os.WriteFile("README.md", []byte(content), 0644)
	fmt.Println("🚀 README actualizado sin bucles.")
}
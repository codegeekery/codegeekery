package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	posts, err := fetchPosts(cfg.APIURL, cfg.Header, cfg.Secret)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if len(posts) < maxPosts {
		fmt.Printf("⚠️ Menos de %d artículos recibidos (got %d)\n", maxPosts, len(posts))
		os.Exit(1)
	}

	table := buildTable(posts[:maxPosts], cfg.PostURL)

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("❌ Error leyendo README: %v\n", err)
		os.Exit(1)
	}

	updated, err := updateReadme(content, table)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(fileName, []byte(updated), 0644); err != nil {
		fmt.Printf("❌ Error escribiendo README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Blog section synchronized.")
}
package main

import (
	"fmt"
	"os"

	"codegeekery/updater/internal/config"
	"codegeekery/updater/internal/fetcher"
	"codegeekery/updater/internal/readme"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	posts, err := fetcher.Fetch(cfg.APIURL, cfg.Header, cfg.Secret)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if len(posts) < readme.MaxPosts() {
		fmt.Printf("⚠️ Menos de %d artículos recibidos (got %d)\n", readme.MaxPosts(), len(posts))
		os.Exit(1)
	}

	table := readme.BuildTable(posts[:readme.MaxPosts()], cfg.PostURL)

	content, err := os.ReadFile(readme.FileName())
	if err != nil {
		fmt.Printf("❌ Error leyendo README: %v\n", err)
		os.Exit(1)
	}

	updated, err := readme.Update(content, table)
	if err != nil {
		fmt.Printf("❌ %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(readme.FileName(), []byte(updated), 0644); err != nil {
		fmt.Printf("❌ Error escribiendo README: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Blog section synchronized.")
}
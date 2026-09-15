package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	fileName  = "README.md"
	startTag  = "<!-- ARTICLES:START -->"
	endTag    = "<!-- ARTICLES:END -->"
	maxPosts  = 3
)

func buildTable(posts []Post, postURL string) string {
	if len(posts) != maxPosts {
		return ""
	}

	var imgParts []string
	var titleParts []string

	for _, p := range posts {
		imgParts = append(imgParts, fmt.Sprintf("[![%s](%s?w=200&h=200)](%s/%s)",
			p.Title, p.MainImage.Asset.URL, postURL, p.Slug.Current))
		titleParts = append(titleParts, fmt.Sprintf("**[%s](%s/%s)**",
			p.Title, postURL, p.Slug.Current))
	}

	imgRow := strings.Join(imgParts, " | ")
	titleRow := strings.Join(titleParts, " | ")

	return imgRow + "\n--- | --- | ---\n" + titleRow + "\n\n[➡️ More blog posts](" + os.Getenv("BASE_URL") + ")"
}

func updateReadme(content []byte, table string) (string, error) {
	updated := string(content)

	reBlock := regexp.MustCompile("(?s)" + regexp.QuoteMeta(startTag) + ".*?" + regexp.QuoteMeta(endTag))
	updated = reBlock.ReplaceAllString(updated, startTag+"\n"+table+"\n"+endTag)

	timestamp := time.Now().UTC().Format(time.RFC3339)
	tsLine := fmt.Sprintf("<!-- Last updated: %s -->", timestamp)
	reTS := regexp.MustCompile(`<!-- Last updated: .*? -->`)

	if reTS.MatchString(updated) {
		updated = reTS.ReplaceAllString(updated, tsLine)
	} else {
		updated += "\n" + tsLine + "\n"
	}

	return updated, nil
}
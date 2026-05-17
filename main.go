package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"open-source-issue-finder/handlers"
	"open-source-issue-finder/services"
	"os"
	"strings"
	"time"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed static/*
var staticFS embed.FS

func main() {
	token := os.Getenv("GITHUB_TOKEN")

	githubSvc := services.NewGitHubService(token)

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"timeAgo": func(t time.Time) string {
			since := time.Since(t)
			switch {
			case since < time.Hour:
				mins := int(since.Minutes())
				if mins <= 1 {
					return "just now"
				}
				return fmt.Sprintf("%d minutes ago", mins)
			case since < 24*time.Hour:
				return fmt.Sprintf("%d hours ago", int(since.Hours()))
			case since < 30*24*time.Hour:
				return fmt.Sprintf("%d days ago", int(since.Hours()/24))
			case since < 365*24*time.Hour:
				return fmt.Sprintf("%d months ago", int(since.Hours()/(24*30)))
			default:
				return fmt.Sprintf("%d years ago", int(since.Hours()/(24*365)))
			}
		},
		"difficultyClass": func(d string) string {
			switch d {
			case "Easy":
				return "badge-easy"
			case "Medium":
				return "badge-medium"
			default:
				return "badge-hard"
			}
		},
		"activityClass": func(label string) string {
			switch label {
			case "Very Active":
				return "activity-high"
			case "Active":
				return "activity-medium"
			case "Moderate":
				return "activity-low"
			default:
				return "activity-stale"
			}
		},
		"labelClass": func(name string) string {
			lower := strings.ToLower(name)
			switch {
			case strings.Contains(lower, "good first"):
				return "label-gfi"
			case strings.Contains(lower, "help wanted"):
				return "label-hw"
			case strings.Contains(lower, "bug"):
				return "label-bug"
			case strings.Contains(lower, "enhancement"):
				return "label-enhancement"
			case strings.Contains(lower, "documentation"):
				return "label-docs"
			default:
				return "label-default"
			}
		},
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	issueHandler := handlers.NewIssueHandler(githubSvc, tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("/", issueHandler.Index)
	mux.HandleFunc("/search", issueHandler.SearchIssues)
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Open Source Issue Finder running at http://localhost:%s", port)
	if token == "" {
		log.Println("⚠️  No GITHUB_TOKEN set — API rate limits will be lower (60 req/hr)")
	} else {
		log.Println("✅ GitHub token detected — using authenticated API (5000 req/hr)")
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

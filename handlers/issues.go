package handlers

import (
	"html/template"
	"net/http"
	"open-source-issue-finder/models"
	"open-source-issue-finder/services"
	"strconv"
	"strings"
)

type IssueHandler struct {
	githubSvc *services.GitHubService
	templates *template.Template
}

func NewIssueHandler(svc *services.GitHubService, tmpl *template.Template) *IssueHandler {
	return &IssueHandler{
		githubSvc: svc,
		templates: tmpl,
	}
}

func (h *IssueHandler) Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *IssueHandler) SearchIssues(w http.ResponseWriter, r *http.Request) {
	language := strings.TrimSpace(r.FormValue("language"))
	label := strings.TrimSpace(r.FormValue("label"))
	sort := strings.TrimSpace(r.FormValue("sort"))
	minStarsStr := strings.TrimSpace(r.FormValue("min_stars"))
	pageStr := strings.TrimSpace(r.FormValue("page"))

	minStars := 0
	if minStarsStr != "" {
		if v, err := strconv.Atoi(minStarsStr); err == nil {
			minStars = v
		}
	}

	page := 1
	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	if language == "" {
		language = "go"
	}

	params := models.SearchParams{
		Language: language,
		Label:    label,
		MinStars: minStars,
		Sort:     sort,
		Page:     page,
	}

	result := h.githubSvc.SearchIssues(params)

	data := struct {
		Result models.SearchResult
		Params models.SearchParams
	}{
		Result: result,
		Params: params,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "issues.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

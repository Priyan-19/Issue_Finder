package models

import "time"

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Repository struct {
	FullName       string `json:"full_name"`
	StargazersCount int   `json:"stargazers_count"`
	HTMLURL        string `json:"html_url"`
}

type Issue struct {
	ID          int        `json:"id"`
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	HTMLURL     string     `json:"html_url"`
	Labels      []Label    `json:"labels"`
	Comments    int        `json:"comments"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	RepoName    string
	RepoURL     string
	RepoStars   int

	// Computed
	Difficulty    string
	ActivityScore int
	PriorityScore int
	ActivityLabel string
}

type SearchResponse struct {
	TotalCount int     `json:"total_count"`
	Items      []Issue `json:"items"`
}

type GitHubIssue struct {
	ID        int       `json:"id"`
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	HTMLURL   string    `json:"html_url"`
	Labels    []Label   `json:"labels"`
	Comments  int       `json:"comments"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Repository *Repository `json:"repository"`
}

type GitHubSearchResponse struct {
	TotalCount int           `json:"total_count"`
	Items      []GitHubIssue `json:"items"`
}

type SearchParams struct {
	Language string
	Label    string
	MinStars int
	Sort     string
	Page     int
}

type SearchResult struct {
	Issues     []Issue
	TotalCount int
	Error      string
}

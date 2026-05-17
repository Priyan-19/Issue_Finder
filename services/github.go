package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"open-source-issue-finder/models"
	"strings"
	"sync"
	"time"
)

type CacheEntry struct {
	Result    models.SearchResult
	ExpiresAt time.Time
}

type GitHubService struct {
	client    *http.Client
	cache     map[string]CacheEntry
	cacheMu   sync.RWMutex
	token     string
}

func NewGitHubService(token string) *GitHubService {
	return &GitHubService{
		client: &http.Client{Timeout: 15 * time.Second},
		cache:  make(map[string]CacheEntry),
		token:  token,
	}
}

func (s *GitHubService) cacheKey(params models.SearchParams) string {
	return fmt.Sprintf("%s|%s|%d|%s|%d", params.Language, params.Label, params.MinStars, params.Sort, params.Page)
}

func (s *GitHubService) getFromCache(key string) (models.SearchResult, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()
	entry, ok := s.cache[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return models.SearchResult{}, false
	}
	return entry.Result, true
}

func (s *GitHubService) setCache(key string, result models.SearchResult) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.cache[key] = CacheEntry{
		Result:    result,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func (s *GitHubService) SearchIssues(params models.SearchParams) models.SearchResult {
	cacheKey := s.cacheKey(params)
	if cached, ok := s.getFromCache(cacheKey); ok {
		return cached
	}

	labels := []string{"good first issue", "help wanted"}
	if params.Label != "" && params.Label != "all" {
		labels = []string{params.Label}
	}

	type partialResult struct {
		issues []models.Issue
		total  int
		err    error
	}

	resultCh := make(chan partialResult, len(labels))
	var wg sync.WaitGroup

	for _, label := range labels {
		wg.Add(1)
		go func(lbl string) {
			defer wg.Done()
			issues, total, err := s.fetchIssues(params, lbl)
			resultCh <- partialResult{issues, total, err}
		}(label)
	}

	wg.Wait()
	close(resultCh)

	seen := make(map[int]bool)
	var allIssues []models.Issue
	totalCount := 0

	for partial := range resultCh {
		if partial.err != nil {
			continue
		}
		if partial.total > totalCount {
			totalCount = partial.total
		}
		for _, issue := range partial.issues {
			if !seen[issue.ID] {
				seen[issue.ID] = true
				allIssues = append(allIssues, issue)
			}
		}
	}

	// Sort by priority
	sortIssues(allIssues, params.Sort)

	result := models.SearchResult{
		Issues:     allIssues,
		TotalCount: totalCount,
	}

	s.setCache(cacheKey, result)
	return result
}

func (s *GitHubService) fetchIssues(params models.SearchParams, label string) ([]models.Issue, int, error) {
	query := buildQuery(params, label)

	sortParam := "updated"
	if params.Sort == "comments" {
		sortParam = "comments"
	}

	page := params.Page
	if page < 1 {
		page = 1
	}

	apiURL := fmt.Sprintf(
		"https://api.github.com/search/issues?q=%s&sort=%s&order=desc&per_page=30&page=%d",
		url.QueryEscape(query), sortParam, page,
	)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "OpenSourceIssueFinder/1.0")
	if s.token != "" {
		req.Header.Set("Authorization", "token "+s.token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("github API returned %d", resp.StatusCode)
	}

	var searchResp models.GitHubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, 0, err
	}

	var issues []models.Issue
	for _, item := range searchResp.Items {
		issue := convertAndAnalyze(item)
		issues = append(issues, issue)
	}

	return issues, searchResp.TotalCount, nil
}

func buildQuery(params models.SearchParams, label string) string {
	parts := []string{"is:issue", "is:open", "state:open"}

	if params.Language != "" {
		parts = append(parts, fmt.Sprintf("language:%s", params.Language))
	}

	labelStr := fmt.Sprintf(`label:"%s"`, label)
	parts = append(parts, labelStr)

	if params.MinStars > 0 {
		parts = append(parts, fmt.Sprintf("stars:>=%d", params.MinStars))
	}

	return strings.Join(parts, " ")
}

func convertAndAnalyze(item models.GitHubIssue) models.Issue {
	issue := models.Issue{
		ID:        item.ID,
		Number:    item.Number,
		Title:     item.Title,
		HTMLURL:   item.HTMLURL,
		Labels:    item.Labels,
		Comments:  item.Comments,
		UpdatedAt: item.UpdatedAt,
		CreatedAt: item.CreatedAt,
	}

	// Extract repo name from URL
	if item.Repository != nil {
		issue.RepoName = item.Repository.FullName
		issue.RepoURL = item.Repository.HTMLURL
		issue.RepoStars = item.Repository.StargazersCount
	} else {
		// Parse from issue URL: https://github.com/owner/repo/issues/N
		parts := strings.Split(item.HTMLURL, "/")
		if len(parts) >= 5 {
			issue.RepoName = parts[3] + "/" + parts[4]
			issue.RepoURL = "https://github.com/" + issue.RepoName
		}
	}

	// Difficulty
	issue.Difficulty = computeDifficulty(item)

	// Activity score (0–100 based on recency)
	issue.ActivityScore = computeActivityScore(item.UpdatedAt)
	issue.ActivityLabel = activityLabel(issue.ActivityScore)

	// Priority = comments weight + recency weight
	issue.PriorityScore = computePriority(item.Comments, issue.ActivityScore)

	return issue
}

func computeDifficulty(item models.GitHubIssue) string {
	for _, lbl := range item.Labels {
		if strings.EqualFold(lbl.Name, "good first issue") {
			return "Easy"
		}
	}
	if item.Comments > 10 {
		return "Hard"
	}
	if item.Comments > 4 {
		return "Medium"
	}
	return "Medium"
}

func computeActivityScore(updatedAt time.Time) int {
	since := time.Since(updatedAt)
	switch {
	case since < 24*time.Hour:
		return 100
	case since < 7*24*time.Hour:
		return 80
	case since < 30*24*time.Hour:
		return 60
	case since < 90*24*time.Hour:
		return 40
	case since < 180*24*time.Hour:
		return 20
	default:
		return 5
	}
}

func activityLabel(score int) string {
	switch {
	case score >= 80:
		return "Very Active"
	case score >= 60:
		return "Active"
	case score >= 40:
		return "Moderate"
	case score >= 20:
		return "Quiet"
	default:
		return "Stale"
	}
}

func computePriority(comments, activityScore int) int {
	commentScore := comments * 3
	if commentScore > 60 {
		commentScore = 60
	}
	return commentScore + activityScore
}

func sortIssues(issues []models.Issue, sortBy string) {
	switch sortBy {
	case "easiest":
		stableSort(issues, func(a, b models.Issue) bool {
			da, db := difficultyRank(a.Difficulty), difficultyRank(b.Difficulty)
			if da != db {
				return da < db
			}
			return a.ActivityScore > b.ActivityScore
		})
	case "activity":
		stableSort(issues, func(a, b models.Issue) bool {
			return a.ActivityScore > b.ActivityScore
		})
	case "priority":
		stableSort(issues, func(a, b models.Issue) bool {
			return a.PriorityScore > b.PriorityScore
		})
	default:
		stableSort(issues, func(a, b models.Issue) bool {
			return a.UpdatedAt.After(b.UpdatedAt)
		})
	}
}

func difficultyRank(d string) int {
	switch d {
	case "Easy":
		return 0
	case "Medium":
		return 1
	default:
		return 2
	}
}

func stableSort(issues []models.Issue, less func(a, b models.Issue) bool) {
	n := len(issues)
	for i := 1; i < n; i++ {
		for j := i; j > 0 && less(issues[j], issues[j-1]); j-- {
			issues[j], issues[j-1] = issues[j-1], issues[j]
		}
	}
}

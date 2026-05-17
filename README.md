<div align="center">

# 🔍 Open Source Issue Finder
### High-Performance GitHub Issue Discovery Platform ⚡🧩

[![Go](https://img.shields.io/badge/Go_1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=html5&logoColor=white)](https://htmx.org/)
[![GitHub API](https://img.shields.io/badge/GitHub_API-181717?style=for-the-badge&logo=github)](https://docs.github.com/en/rest)

**Open Source Issue Finder** is a lightning-fast, server-rendered platform for discovering high-quality GitHub issues. Built with Go and HTMX, it eliminates frontend complexity while delivering a premium, real-time user experience.

</div>

---

## 📖 Project Overview

This platform simplifies the process of finding meaningful open-source contributions by combining intelligent filtering, real-time updates, and a high-performance backend.

### Core Value Proposition
- **⚡ Zero-Overhead Frontend**: No frameworks, no build steps
- **🔍 Intelligent Discovery**: Advanced filtering and ranking
- **🚀 High Performance**: Concurrent API fetching using Go
- **🎯 Contribution Focused**: Helps developers find the right issues faster
- **🎨 Premium UX**: Clean, modern dark-mode interface

---

## 🎨 Design Philosophy: Premium Dark Mode

The interface is crafted for developers who spend hours scanning content.

### Visual System
- Background: `#0d1117` (Deep black canvas)
- Cards: `#161b22` (Elevated surfaces)
- Accent: `#58a6ff` (Electric blue highlights)
- Buttons: `#1f6feb` (Action emphasis)

### Typography
- **Outfit** → Branding & headings
- **Plus Jakarta Sans** → Body & UI text
- **JetBrains Mono** → Labels, tags, metrics

### UX Highlights
- Single-viewport landing experience
- Smooth transition into scrollable results
- Sticky header for navigation
- Clean, text-first issue cards

---

## 🏗️ System Architecture

The application follows a **server-rendered, HTMX-enhanced architecture**:

### 🐹 Backend: Go (net/http)
- Native HTTP server (no frameworks)
- Goroutine-based concurrency
- In-memory caching layer

### ⚡ Frontend: HTMX + Templates
- Partial HTML updates via HTMX
- Server-rendered templates
- Zero JavaScript frameworks

### 🔗 Data Layer: GitHub API
- GitHub Search API integration
- Dynamic issue fetching
- Intelligent filtering and scoring

---

## 🚀 Key Features

### 🔍 Advanced Filtering
- Filter by:
  - Programming language
  - Labels
  - Minimum stars
  - Sort order

### ⚡ Live UI Updates
- Instant search & pagination
- No full page reloads (HTMX)

### 🏷️ Difficulty Analysis
- **Easy** → `good first issue`
- **Medium** → Standard issues
- **Hard** → High complexity / high discussion

### 🔥 Activity & Priority Engine
- Scores issues based on:
  - Comment volume
  - Recency of updates

### 💾 Smart Caching
- Thread-safe in-memory cache
- 5-minute TTL
- Reduces API calls significantly

### 🚀 High Performance Engine
- Concurrent GitHub API calls using goroutines
- Fast response times even under load

---

## 📂 Project Structure

```text
open-source-issue-finder/
├── main.go
├── go.mod
├── .gitignore
│
├── handlers/
│   └── issues.go
│
├── services/
│   └── github.go
│
├── models/
│   └── issue.go
│
├── templates/
│   ├── index.html
│   └── issues.html
│
└── static/
    └── style.css
```

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+

---

### 1. Run Locally

```bash
cd open-source-issue-finder

# Optional: Set GitHub Token
# Windows
$env:GITHUB_TOKEN="your_token"

# Linux/macOS
export GITHUB_TOKEN="your_token"

# Start server
go run main.go
```

---

### 2. Access Application

```
http://localhost:8080
```

---

## 🔑 GitHub API Configuration

| Mode | Limit |
|------|------|
| Unauthenticated | 60 requests/hour |
| Authenticated | 5000 requests/hour |

### Recommendation
Use a personal access token:
👉 https://github.com/settings/tokens

✔ No scopes required  
✔ Works with public data  

---

## 📊 Analysis Logic

| Metric | Logic |
|------|------|
| Difficulty | Easy → `good first issue`, Hard → >10 comments, else Medium |
| Activity | Score (0–100) based on last update recency |
| Priority | `(comment_count × 3, max 60) + activity_score` |

---

## 📦 Architecture Highlights

- Zero frontend framework dependency
- Server-rendered UI with HTMX enhancements
- Efficient caching strategy
- Concurrent data fetching engine
- Clean separation of services, handlers, and models

---

## 🎯 Use Cases

- Discover beginner-friendly open source issues
- Find high-impact contribution opportunities
- Explore trending repositories
- Learn from real-world project issues

---

<div align="center">
  <p>Built with 🔍 for Smarter Open Source Contributions</p>
  <p>Developed by <strong>Priyan</strong></p>
  <p>© 2026 Open Source Issue Finder. All Rights Reserved.</p>
</div>
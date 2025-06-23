package server

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"wb_project_0/internal/database"
)

func HomeHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if page < 1 {
			page = 1
		}
		limit := 3
		offset := (page - 1) * limit

		articles, err := db.GetArticles(ctx, limit, offset)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				http.Error(w, "Request timeout", http.StatusGatewayTimeout)
				return
			}
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		renderTemplate(w, "templates/index.html", articles)
	}
}

func ArticleHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		idStr := strings.TrimPrefix(r.URL.Path, "/article/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		article, err := db.GetArticleByID(ctx, id)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				http.Error(w, "Request timeout", http.StatusGatewayTimeout)
				return
			}
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		renderTemplate(w, "templates/article.html", article)
	}
}

func AdminHandler(db *database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		if r.Method == http.MethodPost {
			title := r.FormValue("title")
			content := r.FormValue("content")

			_, err := db.InsertArticle(ctx, title, content)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					http.Error(w, "Request timeout", http.StatusGatewayTimeout)
				}
				http.Error(w, "Failed to save article", http.StatusInternalServerError)
				return
			}
		}
		renderTemplate(w, "templates/admin.html", nil)
	}
}

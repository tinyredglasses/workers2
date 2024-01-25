package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tinyredglasses/workers2/_examples/d1-blog-server/app/model"
	"github.com/tinyredglasses/workers2/cloudflare/d1"
	_ "github.com/tinyredglasses/workers2/cloudflare/d1" // register driver
)

type articleHandler struct{}

var _ http.Handler = (*articleHandler)(nil)

func NewArticleHandler() http.Handler {
	return &articleHandler{}
}

func (h *articleHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// initialize DB.
	// D1 connector requires request's context to initialize DB.
	c, err := d1.OpenConnector(req.Context(), "BlogDB")
	if err != nil {
		h.handleErr(w, http.StatusInternalServerError, fmt.Sprintf("failed to initialize DB: %v", err))
	}
	// use sql.OpenDB instead of sql.Open.
	db := sql.OpenDB(c)

	switch req.Method {
	case http.MethodGet:
		h.listArticles(w, req, db)
		return
	case http.MethodPost:
		h.createArticle(w, req, db)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("not found"))
}

func (h *articleHandler) handleErr(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(msg))
}

func (h *articleHandler) createArticle(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var createArticleReq model.CreateArticleRequest
	if err := json.NewDecoder(req.Body).Decode(&createArticleReq); err != nil {
		h.handleErr(w, http.StatusBadRequest,
			"request format is invalid")
		return
	}

	now := time.Now().Unix()
	article := model.Article{
		Title:     createArticleReq.Title,
		Body:      createArticleReq.Body,
		CreatedAt: uint64(now),
	}

	result, err := db.Exec(`
INSERT INTO articles (title, body, created_at)
VALUES (?, ?, ?)
   `, article.Title, article.Body, article.CreatedAt)
	if err != nil {
		log.Println(err)
		h.handleErr(w, http.StatusInternalServerError,
			"failed to save article")
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		h.handleErr(w, http.StatusInternalServerError,
			"failed to get ID of inserted article")
		return
	}
	article.ID = uint64(id)

	res := model.CreateArticleResponse{
		Article: article,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode response: %w\n", err)
	}
}

func (h *articleHandler) listArticles(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rows, err := db.Query(`
SELECT id, title, body, created_at FROM articles
ORDER BY created_at DESC;
   `)
	if err != nil {
		h.handleErr(w, http.StatusInternalServerError,
			"failed to load article")
		return
	}

	articles := []model.Article{}
	for rows.Next() {
		var a model.Article
		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt)
		if err != nil {
			log.Println(err)
			h.handleErr(w, http.StatusInternalServerError,
				"failed to scan article")
			return
		}
		articles = append(articles, a)
	}
	res := model.ListArticlesResponse{
		Articles: articles,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode response: %w\n", err)
	}
}

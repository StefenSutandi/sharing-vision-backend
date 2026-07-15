package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/handler"
	"github.com/gin-gonic/gin"
)

type MockService struct {
	Articles map[int64]dto.ArticleResponse
	nextID   int64
}

func (m *MockService) CreateArticle(payload dto.ArticlePayload) error {
	m.Articles[m.nextID] = dto.ArticleResponse{
		ID:       m.nextID,
		Title:    payload.Title,
		Content:  payload.Content,
		Category: payload.Category,
		Status:   payload.Status,
	}
	m.nextID++
	return nil
}

func (m *MockService) GetArticles(limit, offset int, status string) (dto.PaginationResponse[dto.ArticleResponse], error) {
	var data []dto.ArticleResponse
	for _, a := range m.Articles {
		if status == "" || a.Status == status {
			data = append(data, a)
		}
	}
	return dto.PaginationResponse[dto.ArticleResponse]{
		Data: data,
		Pagination: dto.PaginationMeta{
			Total:  int64(len(data)),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (m *MockService) GetArticleByID(id int64) (dto.ArticleResponse, error) {
	if a, ok := m.Articles[id]; ok {
		return a, nil
	}
	return dto.ArticleResponse{}, nil
}

func (m *MockService) UpdateArticle(id int64, payload dto.ArticlePayload) error {
	m.Articles[id] = dto.ArticleResponse{
		ID:       id,
		Title:    payload.Title,
		Content:  payload.Content,
		Category: payload.Category,
		Status:   payload.Status,
	}
	return nil
}

func (m *MockService) TrashArticle(id int64) error {
	a := m.Articles[id]
	a.Status = "thrash"
	m.Articles[id] = a
	return nil
}

type MockValidator struct{}
func (m *MockValidator) Validate(i interface{}) error { return nil }

func setupRouter() (*gin.Engine, *MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	svc := &MockService{Articles: make(map[int64]dto.ArticleResponse), nextID: 1}
	val := &MockValidator{}
	h := handler.NewArticleHandler(svc, val)

	router.POST("/article/", h.CreateArticle)
	router.GET("/article/:limit/:offset", h.GetArticles)
	router.GET("/article/:id", h.GetArticleByID)
	router.PUT("/article/:id", h.UpdateArticle)
	router.DELETE("/article/:id", h.TrashArticle)

	return router, svc
}

func TestArticleHandler(t *testing.T) {
	router, svc := setupRouter()

	svc.CreateArticle(dto.ArticlePayload{Title: "Title", Content: "Content", Category: "Tech", Status: "publish"})

	t.Run("Create Article Malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/article/", bytes.NewBufferString("{malformed json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Get Articles Invalid Limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/article/invalid/0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Get Articles Invalid Offset", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/article/10/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})

	t.Run("Get Article Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/article/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})
	
	t.Run("Update Article Malformed JSON", func(t *testing.T) {
		req, _ := http.NewRequest("PUT", "/article/1", bytes.NewBufferString("{malformed json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})
}

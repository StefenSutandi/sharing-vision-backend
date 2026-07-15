package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/handler"
	"github.com/StefenSutandi/sharing-vision-backend/internal/model"
	"github.com/gin-gonic/gin"
)

type MockService struct {
	Articles map[int64]*model.Article
	nextID   int64
}

func (m *MockService) CreateArticle(req *dto.CreateArticleReq) (*model.Article, map[string]string, error) {
	art := &model.Article{
		ID:       m.nextID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}
	m.Articles[m.nextID] = art
	m.nextID++
	return art, nil, nil
}

func (m *MockService) ListArticles(limit, offset int, status string) (*dto.ListArticleRes, error) {
	var data []model.Article
	for _, a := range m.Articles {
		if status == "" || a.Status == status {
			data = append(data, *a)
		}
	}
	return &dto.ListArticleRes{
		Data: data,
		Pagination: dto.PaginationMeta{
			Total:  int64(len(data)),
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (m *MockService) GetArticleByID(id int64) (*model.Article, error) {
	if a, ok := m.Articles[id]; ok {
		return a, nil
	}
	return nil, nil // simple mock
}

func (m *MockService) UpdateArticle(id int64, req *dto.UpdateArticleReq) (*model.Article, map[string]string, error) {
	art := &model.Article{
		ID:       id,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}
	m.Articles[id] = art
	return art, nil, nil
}

func (m *MockService) TrashArticle(id int64) error {
	a := m.Articles[id]
	if a != nil {
		a.Status = "thrash"
	}
	return nil
}

func setupRouter() (*gin.Engine, *MockService) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	svc := &MockService{Articles: make(map[int64]*model.Article), nextID: 1}
	h := handler.NewArticleHandler(svc)

	router.POST("/article/", h.Create)
	router.GET("/article/:limit/:offset", h.List)
	router.GET("/article/:id", h.GetByID)
	router.PUT("/article/:id", h.Update)
	router.DELETE("/article/:id", h.Delete)

	return router, svc
}

func TestArticleHandler(t *testing.T) {
	router, svc := setupRouter()

	svc.CreateArticle(&dto.CreateArticleReq{Title: "Title", Content: "Content", Category: "Tech", Status: "publish"})

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

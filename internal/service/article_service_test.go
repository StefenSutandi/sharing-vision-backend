package service_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/model"
	"github.com/StefenSutandi/sharing-vision-backend/internal/service"
)

type MockRepo struct {
	articles map[int64]*model.Article
	nextID   int64
}

func (m *MockRepo) Create(a *model.Article) error {
	a.ID = m.nextID
	m.articles[m.nextID] = a
	m.nextID++
	return nil
}

func (m *MockRepo) FindAll(limit, offset int, status string) ([]model.Article, int64, error) {
	var result []model.Article
	for _, a := range m.articles {
		if status == "" || a.Status == status {
			result = append(result, *a)
		}
	}
	total := int64(len(result))
	if offset > len(result) {
		offset = len(result)
	}
	result = result[offset:]
	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}
	return result, total, nil
}

func (m *MockRepo) FindByID(id int64) (*model.Article, error) {
	a, ok := m.articles[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return a, nil
}

func (m *MockRepo) Update(a *model.Article) error {
	m.articles[a.ID] = a
	return nil
}

func TestArticleService(t *testing.T) {
	repo := &MockRepo{
		articles: make(map[int64]*model.Article),
		nextID:   1,
	}
	svc := service.NewArticleService(repo)

	validContent := strings.Repeat("a", 200)

	payload := dto.CreateArticleReq{
		Title:    "  Trimmed Title twenty chars  ",
		Content:  validContent,
		Category: " Tech ",
		Status:   "publish",
	}

	createdArt, valErrs, err := svc.CreateArticle(&payload)
	if err != nil {
		t.Errorf("CreateArticle failed: %v", err)
	}
	if len(valErrs) > 0 {
		t.Errorf("Expected 0 validation errors, got %v", valErrs)
	}

	created, err := repo.FindByID(createdArt.ID)
	if err != nil {
		t.Fatalf("Failed to find created article: %v", err)
	}
	if created.Title != "Trimmed Title twenty chars" {
		t.Errorf("Title not trimmed, got %q", created.Title)
	}
	if created.Category != "Tech" {
		t.Errorf("Category not trimmed, got %q", created.Category)
	}

	foundArt, err := svc.GetArticleByID(createdArt.ID)
	if err != nil {
		t.Errorf("GetArticleByID failed: %v", err)
	}
	if foundArt.Title != "Trimmed Title twenty chars" {
		t.Errorf("Expected Trimmed Title, got %q", foundArt.Title)
	}

	_, err = svc.GetArticleByID(99)
	if err == nil {
		t.Errorf("Expected error for missing article")
	}

	updatePayload := dto.UpdateArticleReq{
		Title:    "Updated Title twenty chars",
		Content:  payload.Content,
		Category: "Tech Updated",
		Status:   "draft",
	}
	updatedArt, valErrs2, err := svc.UpdateArticle(createdArt.ID, &updatePayload)
	if err != nil {
		t.Errorf("UpdateArticle failed: %v", err)
	}
	if len(valErrs2) > 0 {
		t.Errorf("Expected 0 validation errors on update, got %v", valErrs2)
	}
	
	if updatedArt.Title != "Updated Title twenty chars" || updatedArt.Status != "draft" {
		t.Errorf("Article not properly updated")
	}

	_, _, err = svc.UpdateArticle(99, &updatePayload)
	if err == nil {
		t.Errorf("Expected error for update missing article")
	}

	err = svc.TrashArticle(createdArt.ID)
	if err != nil {
		t.Errorf("TrashArticle failed: %v", err)
	}
	trashed, _ := repo.FindByID(createdArt.ID)
	if trashed.Status != "thrash" {
		t.Errorf("Expected status to be 'thrash', got %q", trashed.Status)
	}

	err = svc.TrashArticle(99)
	if err == nil {
		t.Errorf("Expected error for trash missing article")
	}

	repo.Create(&model.Article{Title: "A", Status: "publish"})
	repo.Create(&model.Article{Title: "B", Status: "publish"})
	repo.Create(&model.Article{Title: "C", Status: "draft"})

	paginated, err := svc.ListArticles(10, 0, "")
	if err != nil {
		t.Errorf("ListArticles failed: %v", err)
	}
	if paginated.Pagination.Total != 4 {
		t.Errorf("Expected total 4, got %d", paginated.Pagination.Total)
	}

	paginatedDrafts, err := svc.ListArticles(10, 0, "draft")
	if err != nil {
		t.Errorf("ListArticles draft failed: %v", err)
	}
	if paginatedDrafts.Pagination.Total != 1 {
		t.Errorf("Expected total 1 draft, got %d", paginatedDrafts.Pagination.Total)
	}

	paginatedLimit, _ := svc.ListArticles(2, 0, "")
	if len(paginatedLimit.Data) != 2 {
		t.Errorf("Expected limit 2 to return 2 items, got %d", len(paginatedLimit.Data))
	}
	paginatedOffset, _ := svc.ListArticles(10, 2, "")
	if len(paginatedOffset.Data) != 2 {
		t.Errorf("Expected offset 2 out of 4 to return 2 items, got %d", len(paginatedOffset.Data))
	}
}

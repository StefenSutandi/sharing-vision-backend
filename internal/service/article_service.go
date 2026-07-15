package service

import (
	"strings"

	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/model"
	"github.com/StefenSutandi/sharing-vision-backend/internal/repository"
	"github.com/StefenSutandi/sharing-vision-backend/internal/validator"
)

type ArticleService interface {
	CreateArticle(req *dto.CreateArticleReq) (*model.Article, map[string]string, error)
	ListArticles(limit, offset int, status string) (*dto.ListArticleRes, error)
	GetArticleByID(id int64) (*model.Article, error)
	UpdateArticle(id int64, req *dto.UpdateArticleReq) (*model.Article, map[string]string, error)
	TrashArticle(id int64) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (s *articleService) CreateArticle(req *dto.CreateArticleReq) (*model.Article, map[string]string, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.Category = strings.TrimSpace(req.Category)
	req.Status = strings.TrimSpace(req.Status)

	valErrs := validator.ValidateStruct(req)
	if valErrs != nil {
		return nil, valErrs, nil
	}

	article := &model.Article{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   req.Status,
	}

	err := s.repo.Create(article)
	if err != nil {
		return nil, nil, err
	}

	return article, nil, nil
}

func (s *articleService) ListArticles(limit, offset int, status string) (*dto.ListArticleRes, error) {
	articles, total, err := s.repo.FindAll(limit, offset, status)
	if err != nil {
		return nil, err
	}

	return &dto.ListArticleRes{
		Data: articles,
		Pagination: dto.PaginationMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}, nil
}

func (s *articleService) GetArticleByID(id int64) (*model.Article, error) {
	return s.repo.FindByID(id)
}

func (s *articleService) UpdateArticle(id int64, req *dto.UpdateArticleReq) (*model.Article, map[string]string, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	req.Category = strings.TrimSpace(req.Category)
	req.Status = strings.TrimSpace(req.Status)

	valErrs := validator.ValidateStruct(req)
	if valErrs != nil {
		return nil, valErrs, nil
	}

	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Category = req.Category
	article.Status = req.Status

	err = s.repo.Update(article)
	if err != nil {
		return nil, nil, err
	}

	return article, nil, nil
}

func (s *articleService) TrashArticle(id int64) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	article.Status = "thrash"
	return s.repo.Update(article)
}

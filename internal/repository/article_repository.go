package repository

import (
	"github.com/StefenSutandi/sharing-vision-backend/internal/model"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *model.Article) error
	FindAll(limit, offset int, status string) ([]model.Article, int64, error)
	FindByID(id int64) (*model.Article, error)
	Update(article *model.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) FindAll(limit, offset int, status string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_date DESC, id DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, total, err
}

func (r *articleRepository) FindByID(id int64) (*model.Article, error) {
	var article model.Article
	err := r.db.First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

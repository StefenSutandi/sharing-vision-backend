package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/StefenSutandi/sharing-vision-backend/internal/dto"
	"github.com/StefenSutandi/sharing-vision-backend/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "BAD_REQUEST",
				Message: "Malformed JSON request",
			},
		})
		return
	}

	article, valErrs, err := h.service.CreateArticle(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create article",
			},
		})
		return
	}
	if valErrs != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{
				Code:    "VALIDATION_ERROR",
				Message: "The request contains invalid fields.",
				Fields:  valErrs,
			},
		})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) List(c *gin.Context) {
	limitStr := c.Param("param1")
	offsetStr := c.Param("param2")
	status := c.Query("status")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid limit parameter"},
		})
		return
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid offset parameter"},
		})
		return
	}

	if status != "" && status != "publish" && status != "draft" && status != "thrash" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid status parameter"},
		})
		return
	}

	res, err := h.service.ListArticles(limit, offset, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Failed to retrieve articles"},
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ArticleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("param1")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid ID parameter"},
		})
		return
	}

	article, err := h.service.GetArticleByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: "Article not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Failed to retrieve article"},
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid ID parameter"},
		})
		return
	}

	var req dto.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Malformed JSON request"},
		})
		return
	}

	article, valErrs, err := h.service.UpdateArticle(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: "Article not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Failed to update article"},
		})
		return
	}
	if valErrs != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "VALIDATION_ERROR", Message: "The request contains invalid fields.", Fields: valErrs},
		})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "BAD_REQUEST", Message: "Invalid ID parameter"},
		})
		return
	}

	err = h.service.TrashArticle(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorDetail{Code: "NOT_FOUND", Message: "Article not found"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorDetail{Code: "INTERNAL_ERROR", Message: "Failed to trash article"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "article trashed"})
}

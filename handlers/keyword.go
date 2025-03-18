package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/project-box/models"
	"github.com/project-box/services"
)

type KeywordHandler interface {
	GetAllKeywords(c *gin.Context)
	GetKeywords(c *gin.Context)
	GetKeyword(c *gin.Context)
	CreateKeyword(c *gin.Context)
	UpdateKeyword(c *gin.Context)
	DeleteKeyword(c *gin.Context)
}

type keywordHandler struct {
	service services.KeywordService
}

func NewKeywordHandler(service services.KeywordService) KeywordHandler {
	return &keywordHandler{service: service}
}

// GetKeywords godoc
// @Summary Get all keywords
// @Description Get all keywords for a specific program
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Keyword
// @Failure 500 {object} map[string]string
// @Router /v1/keywords/all [get]
func (h *keywordHandler) GetAllKeywords(c *gin.Context) {
	keywords, err := h.service.GetAllKeywords(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, keywords)
}

// GetKeywords godoc
// @Summary Get keywords by program
// @Description Get all keywords for a specific program
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Param program_id query string true "Program ID"
// @Success 200 {array} models.Keyword
// @Failure 500 {object} map[string]string
// @Router /v1/keywords [get]
func (h *keywordHandler) GetKeywords(c *gin.Context) {
	programID := c.Query("program_id")
	keywords, err := h.service.GetKeywords(c.Request.Context(), programID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, keywords)
}

// GetKeyword godoc
// @Summary Get a keyword
// @Description Get a keyword by ID
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Param id path string true "Keyword ID"
// @Success 200 {object} models.Keyword
// @Failure 404 {object} map[string]string
// @Router /v1/keywords/{id} [get]
func (h *keywordHandler) GetKeyword(c *gin.Context) {
	id := c.Param("id")
	keyword, err := h.service.GetKeyword(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keyword not found"})
		return
	}
	c.JSON(http.StatusOK, keyword)
}

// CreateKeyword godoc
// @Summary Create a keyword
// @Description Create a new keyword
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Param keyword body models.Keyword true "Keyword"
// @Success 201 {object} models.Keyword
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/keywords [post]
func (h *keywordHandler) CreateKeyword(c *gin.Context) {
	var keyword models.Keyword
	if err := c.ShouldBindJSON(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateKeyword(c.Request.Context(), &keyword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, keyword)
}

// UpdateKeyword godoc
// @Summary Update a keyword
// @Description Update an existing keyword
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Param keyword body models.Keyword true "Keyword"
// @Success 200 {object} models.Keyword
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/keywords [put]
func (h *keywordHandler) UpdateKeyword(c *gin.Context) {
	var keyword models.Keyword
	if err := c.ShouldBindJSON(&keyword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateKeyword(c.Request.Context(), &keyword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, keyword)
}

// DeleteKeyword godoc
// @Summary Delete a keyword
// @Description Delete a keyword by ID
// @Tags Keywords
// @Accept  json
// @Produce  json
// @Param id path string true "Keyword ID"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /v1/keywords/{id} [delete]
func (h *keywordHandler) DeleteKeyword(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteKeyword(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

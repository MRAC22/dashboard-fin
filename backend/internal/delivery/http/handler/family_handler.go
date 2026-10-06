package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dashboard-fin/internal/domain"
	"dashboard-fin/internal/usecase"
)

type FamilyHandler struct {
	useCase usecase.FamilyUseCase
}

func NewFamilyHandler(uc usecase.FamilyUseCase) *FamilyHandler {
	return &FamilyHandler{useCase: uc}
}

func (h *FamilyHandler) CreateFamily(c *gin.Context) {
	var input usecase.CreateFamilyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	family, err := h.useCase.CreateFamily(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, family)
}

func (h *FamilyHandler) GetFamily(c *gin.Context) {
	id := c.Param("id")
	family, err := h.useCase.GetFamily(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrFamilyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, family)
}

func (h *FamilyHandler) AddMember(c *gin.Context) {
	familyID := c.Param("id")
	var input usecase.AddMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}
	input.FamilyID = familyID

	member, err := h.useCase.AddMember(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *FamilyHandler) ListMembers(c *gin.Context) {
	familyID := c.Param("id")
	members, err := h.useCase.ListMembers(c.Request.Context(), familyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

func (h *FamilyHandler) UpdateMember(c *gin.Context) {
	memberID := c.Param("id")

	var input usecase.UpdateMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}
	input.ID = memberID

	member, err := h.useCase.UpdateMember(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar membro: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

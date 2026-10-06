package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dashboard-fin/internal/domain"
	"dashboard-fin/internal/usecase"
)

type TransactionHandler struct {
	useCase usecase.TransactionUseCase
}

func NewTransactionHandler(uc usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{useCase: uc}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var input usecase.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	familyID := c.GetString("family_id")
	if familyID != "" {
		input.FamilyID = familyID
	}

	tx, err := h.useCase.Create(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) List(c *gin.Context) {
	familyID := c.GetString("family_id")
	if familyID == "" {
		familyID = c.Query("family_id")
	}

	memberID := c.Query("member_id")

	transactions, err := h.useCase.ListByFamily(c.Request.Context(), familyID, memberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	familyID := c.GetString("family_id")
	if familyID == "" {
		familyID = c.Query("family_id")
	}

	err := h.useCase.Delete(c.Request.Context(), id, familyID)
	if err != nil {
		if err == domain.ErrTransactionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

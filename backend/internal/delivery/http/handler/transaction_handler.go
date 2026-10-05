package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

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

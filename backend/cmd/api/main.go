package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"dashboard-fin/internal/delivery/http/handler"
	"dashboard-fin/internal/infra/auth"
	"dashboard-fin/internal/infra/database"
	"dashboard-fin/internal/usecase"
)

func main() {
	_ = godotenv.Load()

	// 1. Conexão com o PostgreSQL
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o PostgreSQL realizada com sucesso!")

	// 2. Injeção de Dependências
	txRepo := database.NewTransactionRepository(db)
	txUseCase := usecase.NewTransactionUseCase(txRepo)
	txHandler := handler.NewTransactionHandler(txUseCase)

	r := gin.Default()

	// 3. Healthcheck público
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 4. Configuração do Keycloak Auth
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakAuth, err := auth.NewKeycloakAuth(keycloakURL)
	if err != nil {
		log.Printf("Aviso: Keycloak não configurado ou indisponível (%v). Rotas protegidas desativadas temporariamente.", err)
	}

	// 5. Rotas da API V1
	api := r.Group("/api/v1")
	if keycloakAuth != nil {
		api.Use(keycloakAuth.Middleware())
	}

	api.POST("/transactions", txHandler.Create)
	api.GET("/transactions", txHandler.List)
	api.DELETE("/transactions/:id", txHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

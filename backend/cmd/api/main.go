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
	// Carrega arquivo .env se existir
	_ = godotenv.Load()

	// 1. Conexão com o PostgreSQL (O GORM cria as tabelas via AutoMigrate)
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o PostgreSQL realizada com sucesso!")

	// 2. Dependências da aplicação
	txRepo := database.NewTransactionRepository(db)
	txUseCase := usecase.NewTransactionUseCase(txRepo)
	txHandler := handler.NewTransactionHandler(txUseCase)

	r := gin.Default()

	// 3. Rota de Healthcheck (Pública)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 4. Configuração opcional do Keycloak JWT Auth
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	keycloakAuth, err := auth.NewKeycloakAuth(keycloakURL)
	if err != nil {
		log.Printf("Aviso: Keycloak não configurado ou indisponível (%v). Rotas protegidas desativadas temporariamente.", err)
	}

	// 5. Rotas da API
	api := r.Group("/api/v1")
	if keycloakAuth != nil {
		api.Use(keycloakAuth.Middleware())
	}

	api.POST("/transactions", txHandler.Create)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

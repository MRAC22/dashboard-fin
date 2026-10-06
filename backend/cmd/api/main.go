package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"dashboard-fin/internal/delivery/http/handler"
	"dashboard-fin/internal/domain"
	"dashboard-fin/internal/infra/auth"
	"dashboard-fin/internal/infra/database"
	"dashboard-fin/internal/usecase"
)

func findIndexHTML() string {
	paths := []string{
		"../index.html",
		"index.html",
		"../../index.html",
	}

	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err == nil {
			if _, err := os.Stat(absPath); err == nil {
				log.Printf("[Static] Arquivo index.html localizado em: %s", absPath)
				return absPath
			}
		}
	}
	return ""
}

func main() {
	_ = godotenv.Load()

	// 1. Conexão com o PostgreSQL
	db, err := database.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	log.Println("Conexão com o PostgreSQL realizada com sucesso!")

	// AutoMigrate
	if err := db.AutoMigrate(&domain.Family{}, &domain.Member{}, &domain.Transaction{}); err != nil {
		log.Fatalf("Erro ao executar AutoMigrate: %v", err)
	}

	// 2. Injeção de Dependências
	txRepo := database.NewTransactionRepository(db)
	txUseCase := usecase.NewTransactionUseCase(txRepo)
	txHandler := handler.NewTransactionHandler(txUseCase)

	familyRepo := database.NewFamilyRepository(db)
	familyUseCase := usecase.NewFamilyUseCase(familyRepo)
	familyHandler := handler.NewFamilyHandler(familyUseCase)

	r := gin.Default()

	// 3. Servir o arquivo estático index.html da raiz
	indexPath := findIndexHTML()
	if indexPath != "" {
		r.StaticFile("/", indexPath)
	} else {
		log.Println("[Aviso] index.html não foi encontrado.")
	}

	// Healthcheck público
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 4. Configuração do Keycloak Auth (Opcional no Dev)
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	var keycloakAuth *auth.KeycloakAuth
	if keycloakURL != "" {
		var err error
		keycloakAuth, err = auth.NewKeycloakAuth(keycloakURL)
		if err != nil {
			log.Printf("Aviso: Keycloak indisponível (%v). Continuando sem auth...", err)
		}
	} else {
		log.Println("Aviso: KEYCLOAK_URL não informada. Rotas da API v1 abertas para desenvolvimento.")
	}

	// 5. Grupo de Rotas da API V1
	api := r.Group("/api/v1")
	if keycloakAuth != nil {
		api.Use(keycloakAuth.Middleware())
	}

	// Rotas de Transações
	api.POST("/transactions", txHandler.Create)
	api.GET("/transactions", txHandler.List)
	api.DELETE("/transactions/:id", txHandler.Delete)

	// Rotas de Famílias e Membros
	api.POST("/families", familyHandler.CreateFamily)
	api.GET("/families/:id", familyHandler.GetFamily)
	api.POST("/families/:id/members", familyHandler.AddMember)
	api.GET("/families/:id/members", familyHandler.ListMembers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Servidor rodando na porta %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

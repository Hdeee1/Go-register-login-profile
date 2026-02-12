package main

import (
	"fmt"
	"log"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http"
	repository "github.com/Hdeee1/go-register-login-profile/internal/repository/mysql"
	"github.com/Hdeee1/go-register-login-profile/internal/usecase"
	"github.com/Hdeee1/go-register-login-profile/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load env")
	}

	db, err := database.ConnectMySQL()
	if err != nil {
		log.Fatalf("Failed to connect database. Error: %s", err.Error())
	}

	repo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatal("Failed to create user repository")
	}

	useCase := usecase.NewUserUsecase(repo)
	h := http.NewUserHandler(useCase)

	r := gin.Default()
	r.POST("/api/user/register", h.Register)
	r.POST("/api/user/login", h.Login)

	fmt.Println("Server started at port :8080")
	r.Run(":8081")
}
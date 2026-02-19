package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUsecase
	tokenBlacklist *jwt.TokenBlacklist
}

type registerResponse struct {
	Id			int `json:"user_id"`
	FullName	string `json:"full_name"`
	Username	string `json:"username"`
	Email		string `json:"email"`
}

type loginResponse struct {
	Username	 string `json:"username"`
	Email		 string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserHandler(u domain.UserUsecase, b *jwt.TokenBlacklist) *UserHandler {
	return &UserHandler{
		userUseCase: u,
		tokenBlacklist: b,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var newUser domain.RegisterRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.Register(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := registerResponse{
		Id: user.Id,
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var newUser domain.LoginRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, accTkn, refTkn, err := h.userUseCase.Login(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res := loginResponse{
		Username: usr.Username,
		Email: usr.Email,
		AccessToken: accTkn,
		RefreshToken: refTkn,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required"})
		return 
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
		return 
	}

	tokenString := parts[1]

	claims, err := jwt.ValidateToken(tokenString, os.Getenv("JWT_ACCESS_SECRET"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return 
	}

	h.tokenBlacklist.AddTokenBlacklist(tokenString, claims.ExpiresAt.Time)
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := value.(int)

	user, err := h.userUseCase.GetProfile(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	res := gin.H{
		"data": gin.H{
			"id": user.Id,
			"full_name": user.FullName,
			"username": user.Username,
			"email": user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, res)
}
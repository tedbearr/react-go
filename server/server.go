package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tedbearr/react-go/config"
	"github.com/tedbearr/react-go/controller"
	"github.com/tedbearr/react-go/repository"
	"github.com/tedbearr/react-go/service"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.DbConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	authService    service.AuthService       = service.NewAuthService(userRepository)
	jwtService     service.JWTService        = service.NewJWTService()
	authController controller.AuthController = controller.NewAuthController(authService, jwtService)
)

func main() {
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("register", authController.Register)
		authRoutes.POST("/login", authController.Login)
	}

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

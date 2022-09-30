package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tedbearr/react-go/dto"
	"github.com/tedbearr/react-go/entity"
	"github.com/tedbearr/react-go/helper"
	"github.com/tedbearr/react-go/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
	}
	authResult := c.authService.VerifyCredential(loginDTO.Username, loginDTO.Password)
	if result, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(result.ID, 10))
		result.Token = generatedToken
		res := helper.BuildResponse(true, "ok!", result)
		ctx.JSON(http.StatusOK, res)
		return
	}
	res := helper.BuildErrorResponse("failed to login", "invalid credential", authResult)
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("failed to proccess request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if !c.authService.IsDuplicateUsername(registerDTO.Username) {
		res := helper.BuildErrorResponse("failed to proccess request", "duplicate username", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, res)
	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		res := helper.BuildResponse(true, "ok!", createdUser)
		ctx.JSON(http.StatusCreated, res)
	}
}

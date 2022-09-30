package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tedbearr/react-go/dto"
	"github.com/tedbearr/react-go/entity"
	"github.com/tedbearr/react-go/helper"
	"github.com/tedbearr/react-go/service"
)

type PostsController interface {
	All(context *gin.Context)
	FindById(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type postController struct {
	postService service.PostService
	jwtService  service.JWTService
}

func NewPostController(postService service.PostService, jwtService service.JWTService) PostsController {
	return &postController{
		postService: postService,
		jwtService:  jwtService,
	}
}

func (controller *postController) All(context *gin.Context) {
	var posts []entity.Posts = controller.postService.All()
	res := helper.BuildResponse(true, "ok!", posts)
	context.JSON(http.StatusOK, res)
}

func (controller *postController) FindById(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("no params was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	}
	var post entity.Posts = controller.postService.FindById(id)
	if (post == entity.Posts{}) {
		res := helper.BuildErrorResponse("data not found", "no data with given id", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "Ok!", post)
		context.JSON(http.StatusOK, res)
	}
}

func (controller *postController) Insert(context *gin.Context) {
	var postCreateDTO dto.PostsCreateDTO
	errDTO := context.ShouldBind(&postCreateDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("failed to proccess request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
	} else {
		autHeader := context.GetHeader("Authorization")
		userID := controller.GetUserIdByToken(autHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err != nil {
			postCreateDTO.UserID = convertedUserID
		}
		result := controller.postService.Insert(postCreateDTO)
		res := helper.BuildResponse(true, "oke", result)
		context.JSON(http.StatusCreated, res)
	}
}

func (controller *postController) Update(context *gin.Context) {

}

func (controller *postController) Delete(context *gin.Context) {

}

func (controller *postController) GetUserIdByToken(token string) string {
	aToken, err := controller.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}

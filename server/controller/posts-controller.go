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
	var postUpdateDTO dto.PostsUpdateDTO
	errDTO := context.ShouldBind((&postUpdateDTO))
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := controller.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	postID, err := strconv.ParseUint(context.Param("id"), 0, 64)
	if err != nil {
		postUpdateDTO.ID = postID
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if controller.postService.IsAllowedToEdit(userID, postUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID != nil {
			postUpdateDTO.UserID = id
		}
		result := controller.postService.Update(postUpdateDTO)
		res := helper.BuildResponse(true, "oke", result)
		context.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("You dont have permission", "you are not owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (controller *postController) Delete(context *gin.Context) {
	var post entity.Posts
	id, err := strconv.ParseUint(context.Param("id"), 10, 64)
	if err != nil {
		res := helper.BuildErrorResponse("failed to get param id", "please insert param id", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	}
	authHeader := context.GetHeader("Authorization")
	token, errToken := controller.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	post.ID = id
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if controller.postService.IsAllowedToEdit(userID, post.ID) {
		controller.postService.Delete(post)
		res := helper.BuildResponse(true, "deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		res := helper.BuildErrorResponse("you dont have permission", "you are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, res)
	}
}

func (controller *postController) GetUserIdByToken(token string) string {
	aToken, err := controller.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["user_id"])
}

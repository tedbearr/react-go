package service

import (
	"fmt"
	"log"

	"github.com/mashingan/smapping"
	"github.com/tedbearr/react-go/dto"
	"github.com/tedbearr/react-go/entity"
	"github.com/tedbearr/react-go/repository"
)

type PostService interface {
	Insert(post dto.PostsCreateDTO) entity.Posts
	Update(post dto.PostsUpdateDTO) entity.Posts
	Delete(post entity.Posts)
	All() []entity.Posts
	FindById(PostID uint64) entity.Posts
	IsAllowedToEdit(userID string, PostID uint64) bool
}

type postService struct {
	postRepository repository.PostsRepository
}

func NewPostService(postRepo repository.PostsRepository) PostService {
	return &postService{
		postRepository: postRepo,
	}
}

func (service *postService) Insert(post dto.PostsCreateDTO) entity.Posts {
	postEntity := entity.Posts{}
	err := smapping.FillStruct(&postEntity, smapping.MapFields(&post))
	if err != nil {
		log.Fatalf("failed map %v", err)
	}
	res := service.postRepository.InsertPosts(postEntity)
	return res
}

func (service *postService) Update(post dto.PostsUpdateDTO) entity.Posts {
	postEntity := entity.Posts{}
	err := smapping.FillStruct(&postEntity, smapping.MapFields(&post))
	if err != nil {
		log.Fatalf("failed map %v", err)
	}
	res := service.postRepository.UpdatePosts(postEntity)
	return res
}

func (service *postService) Delete(post entity.Posts) {
	service.postRepository.DeletePosts(post)
}

func (service *postService) All() []entity.Posts {
	return service.postRepository.AllPosts()
}

func (service *postService) FindById(PostID uint64) entity.Posts {
	return service.postRepository.FindPostsById(PostID)
}

func (service *postService) IsAllowedToEdit(UserID string, postID uint64) bool {
	post := service.postRepository.FindPostsById(postID)
	id := fmt.Sprintf("%v", post.UserID)
	return UserID == id
}

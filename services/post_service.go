package services

import (
	"errors"
	models "moodly/Models"
	"moodly/repositories"
)

type PostService struct {
	repo *repositories.PostRepository
}

func NewPostService(repo *repositories.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) CreatePost(post *models.PostModel) error {
	if post.Title == "" {
		return errors.New("title is required")
	}

	return s.repo.Create(post)
}

func (s *PostService) GetPosts() ([]models.PostModel, error) {
	return s.repo.FindAll()
}

func (s *PostService) GetPostByID(id string) (*models.PostModel, error) {
	return s.repo.FindByID(id)
}

func (s *PostService) UpdatePost(id string, body models.PostModel) (*models.PostModel, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(post, body); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) DeletePost(id string) error {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(post)
}

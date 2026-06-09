package repositories

import (
	models "moodly/Models"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(post *models.PostModel) error {
	return r.db.Create(post).Error
}

func (r *PostRepository) FindAll() ([]models.PostModel, error) {
	var posts []models.PostModel

	err := r.db.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id string) (*models.PostModel, error) {
	var post models.PostModel

	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) Update(post *models.PostModel, body models.PostModel) error {
	return r.db.Model(post).Updates(body).Error
}

func (r *PostRepository) Delete(post *models.PostModel) error {
	return r.db.Delete(post).Error
}

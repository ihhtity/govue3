package services

import (
	"govue3/db"
	"govue3/models"

	"gorm.io/gorm"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService() UserService {
	return &userService{
		db: db.MySQL,
	}
}

// @ID			get-user-by-id
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		id path int true "User ID"
// @Success		200 {object} string "successful operation"
// @Failure		404 {object} string "User not found"
// @Router		/users/{id} [get]
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// @ID			update-user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		user body models.User true "User object"
// @Success		200 {object} string "successful operation"
// @Failure		400 {object} string "Invalid user object"
// @Router		/users/{id} [put]
func (s *userService) UpdateUser(user *models.User) (*models.User, error) {
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// @ID			delete-user
// @Tags		users
// @Accept		json
// @Produce		json
// @Param		id path int true "User ID"
// @Success		200 {object} string "successful operation"
// @Failure		404 {object} string "User not found"
// @Router		/users/{id} [delete]
func (s *userService) DeleteUser(id uint) error {
	var user models.User
	if err := s.db.Delete(&user, id).Error; err != nil {
		return err
	}
	return nil
}

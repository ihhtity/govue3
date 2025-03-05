package services

import (
	"errors"
	"govue3/models"
	"govue3/utils"

	"gorm.io/gorm"
)

type AuthService interface {
	Login(username, password string) (string, *models.User, error) // 修改返回签名
	Register(username, password string) error
}

type authService struct {
	BaseService
}

func NewAuthService() AuthService {
	return &authService{
		BaseService: *NewBaseService(),
	}
}

// ... existing code ...

func (s *authService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	if err := s.MySQL.Where("username = ?", username).First(&user).Error; err != nil {
		return "", nil, errors.New("账号错误")
	}

	// TODO: 密码验证逻辑
	if user.Password != password {
		return "", nil, errors.New("密码错误")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, &user, nil
}

// ... existing code ...

func validateUser(username, password string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}
	if password == "" || len(password) < 2 {
		return errors.New("密码不能为空且长度必须大于等于2")
	}
	return nil
}

func (s *authService) Register(username, password string) error {
	if err := validateUser(username, password); err != nil {
		return err
	}

	if err := s.MySQL.Where("username = ?", username).First(&models.User{}).Error; err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 创建并初始化 user 变量
	user := &models.User{
		Username: username,
		Password: password,
	}

	if err := s.MySQL.Create(user).Error; err != nil {
		return err
	}
	return nil
}

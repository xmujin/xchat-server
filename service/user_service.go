package service

import (
	"xchat-server/model"
	"xchat-server/repository"
	"xchat-server/utils/jwt"
)

type UserService struct {
	// 实际项目中，这里通常会注入 repository 接口
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Register(username string, password string) string {
	// 创建用户
	s.userRepo.CreateUser(&model.User{Name: username, Password: password})

	token := jwt.GenerateToken(username)
	return token
}

package main

import (
	"user_manager/src/repository"
	"user_manager/src/service"
)

type UserManager struct{
	UserRepo repository.UserRepository
}

func NewUserManager() UserManager {
	notifyService := service.NotifyService{}

	return UserManager{
		UserRepo: repository.UserRepository{
			NotifyService: &notifyService,
			Engine:        repository.UserRepositoryEngineInMemory{},
		},
	}
}
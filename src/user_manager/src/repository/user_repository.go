package repository

import (
	"log"
	"user_manager/src/entity"
	"user_manager/src/service"
)

type UserRepositoryEngine interface {
	insert(user *entity.User) error
	update(userId uint, newUser *entity.User) error
	delete(userId uint) error
	findBy(userConditions *entity.User) (error, []entity.User)
}

type UserRepository struct {
	NotifyService *service.NotifyService
	Engine                 UserRepositoryEngine
}

func (repo UserRepository) Insert(user *entity.User) error {
	if err := repo.Engine.insert(user); err != nil {
		return err
	}

	if repo.NotifyService != nil {
		if err := repo.NotifyService.SendEvent(service.EventUserInserted, user); err != nil {
			return err
		}
	}

	log.Printf("User Inserted %v", user)
	return nil
}

func (repo UserRepository) Update(userId uint, newUser *entity.User) error {
	if err := repo.Engine.update(userId, newUser); err != nil {
		return err
	}

	if repo.NotifyService != nil {
		if err := repo.NotifyService.SendEvent(service.EventUserUpdated, newUser); err != nil {
			return err
		}
	}

	log.Printf("User Updated %v", newUser)
	return nil
}

func (repo UserRepository) Delete(userId uint) error {
	tmpUser := entity.User{ID: userId}
	err, users := repo.Engine.findBy(&tmpUser)
	if err != nil {
		return err
	}
	if len(users) <= 0 {
		return ErrNotFoundUser
	}

	user := users[0]

	if err := repo.Engine.delete(user.ID); err != nil {
		return err
	}

	if repo.NotifyService != nil {
		if err := repo.NotifyService.SendEvent(service.EventUserDeleted, &user); err != nil {
			return err
		}
	}

	log.Printf("User Deleted %v", user)
	return nil
}

func (repo UserRepository) FindBy(userConditions *entity.User) (error, []entity.User) {
	err, users := repo.Engine.findBy(userConditions)
	if err != nil {
		return err, []entity.User{}
	}

	log.Printf("Users Found (using following filters %v): %v", userConditions, users)
	return nil, users
}

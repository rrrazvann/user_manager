package repository

import (
	"fmt"
	"sync"
	"user_manager/src/entity"
)

var (
	autoIncrement = uint(1)
	usersLock sync.Mutex

	idsMap = map[uint]*entity.User{}   // helpers for uniqueness check
	emailsMap = map[string]*entity.User{} // helpers for uniqueness check
	nicknamesMap = map[string]*entity.User{} // helpers for uniqueness check

	ErrDuplicateUserEmail    = fmt.Errorf("duplicate user with same email")
	ErrDuplicateUserNickname = fmt.Errorf("duplicate user with same nickname")
	ErrNotFoundUser          = fmt.Errorf("user not found")
)

type UserRepositoryEngineInMemory struct {}

func (repo UserRepositoryEngineInMemory) insert(user *entity.User) error  {
	if user == nil {
		return nil
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	if _, isset := emailsMap[user.Email]; isset {
		return ErrDuplicateUserEmail
	}

	if _, isset := nicknamesMap[user.Nickname]; isset {
		return ErrDuplicateUserNickname
	}

	tmpUser := *user

	user.ID = autoIncrement
	tmpUser.ID = autoIncrement

	idsMap[tmpUser.ID] = &tmpUser
	emailsMap[user.Email] = &tmpUser
	nicknamesMap[user.Nickname] = &tmpUser

	autoIncrement++

	return nil
}

func (repo UserRepositoryEngineInMemory) update(userId uint, newUser *entity.User) error {
	if newUser == nil || userId == 0 {
		return nil
	}

	usersLock.Lock()
	defer usersLock.Unlock()

	foundUser, isset := idsMap[userId];
	if !isset {
		return ErrNotFoundUser
	}

	if foundUser.Email != newUser.Email {
		delete(emailsMap, foundUser.Email)

		foundUser.Email = newUser.Email
		emailsMap[foundUser.Email] = foundUser
	}

	if foundUser.Nickname != newUser.Nickname {
		delete(nicknamesMap, foundUser.Nickname)

		foundUser.Nickname = newUser.Nickname
		nicknamesMap[foundUser.Nickname] = foundUser
	}

	foundUser.FirstName = newUser.FirstName
	foundUser.LastName  = newUser.LastName
	foundUser.Password  = newUser.Password
	foundUser.Country   = newUser.Country

	return nil
}

func (repo UserRepositoryEngineInMemory) delete(userId uint) error {
	usersLock.Lock()
	defer usersLock.Unlock()

	if userId == 0 {
		return nil
	}

	foundUser, isset := idsMap[userId];
	if !isset {
		return ErrNotFoundUser
	}

	delete(idsMap, foundUser.ID)
	delete(nicknamesMap, foundUser.Nickname)
	delete(emailsMap, foundUser.Email)

	return nil
}

func (repo UserRepositoryEngineInMemory) findBy(userConditions *entity.User) (error, []entity.User) {
	if userConditions == nil {
		return nil, []entity.User{}
	}

	if userConditions.ID == 0 &&
		len(userConditions.FirstName) == 0 &&
		len(userConditions.LastName) == 0 &&
		len(userConditions.Email) == 0 &&
		len(userConditions.Nickname) == 0 &&
		len(userConditions.Country) == 0 {
		return nil, []entity.User{}
	}

	// Duplicate idsMap
	usersFoundMap := map[uint]entity.User{}
	for id, user := range idsMap {
		tmpUser := *user

		usersFoundMap[id] = tmpUser
	}

	// Use ID index
	if userConditions.ID > 0 {
		user, isset := idsMap[userConditions.ID]
		if isset {
			tmpUser := *user

			usersFoundMap = map[uint]entity.User{tmpUser.ID: tmpUser}
		} else {
			return nil, []entity.User{}
		}
	}

	// Use Email index
	if len(userConditions.Email) > 0 {
		user, isset := emailsMap[userConditions.Email]
		if isset {
			tmpUser := *user

			usersFoundMap = map[uint]entity.User{tmpUser.ID: tmpUser}
		} else {
			return nil, []entity.User{}
		}
	}

	// Use Nickname index
	if len(userConditions.Nickname) > 0 {
		user, isset := nicknamesMap[userConditions.Nickname]
		if isset {
			tmpUser := *user

			usersFoundMap = map[uint]entity.User{tmpUser.ID: tmpUser}
		} else {
			return nil, []entity.User{}
		}
	}

	// "table scan"
	if len(userConditions.FirstName) > 0 {
		for id, user := range usersFoundMap {
			if user.FirstName != userConditions.FirstName {
				delete(usersFoundMap, id)
			}
		}
	}

	if len(userConditions.LastName) > 0 {
		for id, user := range usersFoundMap {
			if user.LastName != userConditions.LastName {
				delete(usersFoundMap, id)
			}
		}
	}

	if len(userConditions.Country) > 0 {
		for id, user := range usersFoundMap {
			if user.Country != userConditions.Country {
				delete(usersFoundMap, id)
			}
		}
	}

	usersFound := make([]entity.User, 0)
	for _, user := range usersFoundMap {
		usersFound = append(usersFound, user)
	}

	return nil, usersFound
}

package repository

import (
	"math/rand"
	"reflect"
	"testing"
	"user_manager/src/entity"
	"user_manager/src/lib"
	"user_manager/src/repository"
)

var (
	lastNames   = []string{"John", "Billy", "Mike", "Lucas", "Olivia", "William", "Xabi", "Sophia"}
	firstNames  = []string{"Jones", "Garcia", "Davis", "Rodriguez", "Martinez", "Long", "Carter", "Young"}
	countries   = []string{"UK", "RO", "IT", "ES", "DE", "NL", "BG"}
)

func TestInsert(t *testing.T) {
	userRepo := repository.UserRepository{
		NotifyService: nil,
		Engine:        repository.UserRepositoryEngineInMemory{},
	}

	// Add user 1
	user1 := generateRandomUser()
	err := userRepo.Insert(user1)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Find user 1
	user1Conditions := entity.User{Country: user1.Country, Nickname: user1.Nickname}
	err, usersFound := userRepo.FindBy(&user1Conditions)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(usersFound) != 1 {
		t.Errorf("users found count should be 1, not %d", len(usersFound))
	}

	user1Found := usersFound[0]
	if !reflect.DeepEqual(*user1, user1Found) {
		t.Error("user1 is not the same")
		differentUserErrorHandling(t, *user1, user1Found)
	}

	// Add user 2
	user2 := generateRandomUser()
	user2.Country = user1.Country
	err = userRepo.Insert(user2)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Find users: 1 and 2
	usersConditions := entity.User{Country: user1.Country}
	err, usersFound = userRepo.FindBy(&usersConditions)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(usersFound) != 2 {
		t.Errorf("users found count should be 2, not %d", len(usersFound))
	}

	foundCount := 0
	for _, userFound := range usersFound {
		if reflect.DeepEqual(*user1, userFound) {
			foundCount++
		} else if reflect.DeepEqual(*user2, userFound) {
			foundCount++
		}
	}

	if foundCount != 2 {
		t.Error("users found are not the same")
		t.Errorf("expected: %v", []entity.User{*user1, *user2})
		t.Errorf("actual: %v", usersFound)
	}

	// todo: more test cases
}

func generateRandomUser() *entity.User {
	lastName := lastNames[rand.Intn(len(lastNames))]
	firstName := firstNames[rand.Intn(len(firstNames))]
	nickname := lastName + "_" + firstName
	password := lib.HashPassword(nickname)
	country := countries[rand.Intn(len(countries))]
	email := nickname + "@example.com"

	user := entity.User{
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		Password:  password,
		Email:     email,
		Country:   country,
	}
	
	return &user
}

func differentUserErrorHandling(t *testing.T, expectedUser entity.User, actualUser entity.User) {
	t.Errorf("expected: %v", expectedUser)
	t.Errorf("actual: %v", actualUser)
}
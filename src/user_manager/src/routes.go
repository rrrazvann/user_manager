package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"user_manager/src/entity"
	"user_manager/src/form"
	"user_manager/src/repository"
)

var (
	userManager = NewUserManager()
)

type JsonResponse map[string]interface{}

func getAndAddUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		addUser(w, r)
	} else if r.Method == "GET" {
		getUsers(w, r)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[POST] Add User (%s)", r.URL.Path)

	var userForm form.UserForm
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		log.Print(" [400] -> error decoding payload")

		response := map[string]string{
			"error": "error decoding payload",
		}

		sendJsonResponse(w, response, http.StatusBadRequest)
		return
	}

	err, user := form.ValidateUserForm(&userForm);
	if err != nil {
		log.Print(" [400] -> error validation form")

		response := map[string]string{
			"error": err.Error(),
		}

		sendJsonResponse(w, response, http.StatusBadRequest)
		return
	}

	err = userManager.UserRepo.Insert(user)
	if err != nil {
		if err != repository.ErrDuplicateUserNickname && err != repository.ErrDuplicateUserEmail {
			log.Printf(" [500] (%s)", err.Error())

			response := map[string]string{
				"error": "unexpected error",
			}

			sendJsonResponse(w, response, http.StatusInternalServerError)
			return
		}

		log.Printf(" [400] -> error while inserting user (%s)", err.Error())

		response := map[string]string{
			"error": err.Error(),
		}

		sendJsonResponse(w, response, http.StatusBadRequest)
		return
	}

	response := JsonResponse{
		"user": user,
	}

	sendJsonResponse(w, response, http.StatusOK)

	log.Print(" [200]")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET] Users (%s)", r.URL.Path)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(" [400] -> error decoding payload")

		response := map[string]string{
			"error": "error decoding payload",
		}

		sendJsonResponse(w, response, http.StatusBadRequest)
		return
	}

	err, users := userManager.UserRepo.FindBy(&user)
	if err != nil {
		log.Printf(" [500] (%s)", err.Error())

		response := map[string]string{
			"error": "unexpected error",
		}

		sendJsonResponse(w, response, http.StatusInternalServerError)
		return
	}

	response := JsonResponse{
		"users": users,
	}

	sendJsonResponse(w, response, http.StatusOK)
	log.Print(" [200]")
}

func deleteAndUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		deleteUser(w, r)
	} else if r.Method == "PUT" {
		updateUser(w, r)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DELETE] Users (%s)", r.URL.Path)

	// todo: better use another router which has built-in route parameters
	splittedUrl := strings.Split(r.URL.Path, "/")
	if len(splittedUrl) < 3 {
		log.Print(" [404]")
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseUint(splittedUrl[2], 10, 32)
	if err != nil {
		log.Print(" [404]")
		http.NotFound(w, r)
		return
	}

	userId := uint(id)
	err = userManager.UserRepo.Delete(userId)
	if err != nil {
		if err == repository.ErrNotFoundUser {
			log.Print(" [404]")
			http.NotFound(w, r)
			return
		}

		log.Printf(" [500] (%s)", err.Error())

		response := map[string]string{
			"error": "unexpected error",
		}

		sendJsonResponse(w, response, http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, nil, http.StatusNoContent)
	log.Print(" [204]")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("[PUT] Users (%s)", r.URL.Path)

	// todo: better use another router which has built-in route parameters
	splittedUrl := strings.Split(r.URL.Path, "/")
	if len(splittedUrl) < 3 {
		log.Print(" [404]")
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseUint(splittedUrl[2], 10, 32)
	if err != nil {
		log.Print(" [404]")
		http.NotFound(w, r)
		return
	}

	userId := uint(id)

	var user entity.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(" [400] -> error decoding payload")

		response := map[string]string{
			"error": "error decoding payload",
		}

		sendJsonResponse(w, response, http.StatusBadRequest)
		return
	}

	err = userManager.UserRepo.Update(userId, &user)
	if err != nil {
		if err == repository.ErrNotFoundUser {
			log.Print(" [404]")
			http.NotFound(w, r)
			return
		}

		log.Printf(" [500] (%s)", err.Error())

		response := map[string]string{
			"error": "unexpected error",
		}

		sendJsonResponse(w, response, http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, nil, http.StatusNoContent)
	log.Print(" [204]")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("[GET] Health check (%s)", r.URL.Path)

	// todo: maybe check storage system status and add as a key in response
	response := map[string]string{
		"status": "healthy",
	}

	sendJsonResponse(w, response, http.StatusOK)
	log.Print(" [200]")
}

// Utils functions
func sendJsonResponse(w http.ResponseWriter, response interface{}, statusCode int) {
	w.WriteHeader(statusCode)

	if response != nil {
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("  -> couldnt send response (err: %s)", err.Error())
		}
	}
}

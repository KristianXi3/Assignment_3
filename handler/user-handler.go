package handler

import (
	"context"
	"encoding/json"
	"golang-crud-sql/entity"
	"golang-crud-sql/model"
	"golang-crud-sql/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

var UserRepo repository.UserRepoIface

type UserHandlerIface interface {
	UserHandler(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
}

func NewUserHandler() UserHandlerIface {
	return &UserHandler{}
}

func (u *UserHandler) UserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if strings.HasSuffix(r.URL.Path, "users/register") {
		if r.Method == http.MethodPost {
			registerUser(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Method"))
			return
		}

	} else if strings.HasSuffix(r.URL.Path, "users/login") {
		if r.Method == http.MethodPost {
			loginUsers(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid Method"))
			return
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			if id != "" {
				getUserById(w, r, id)
			} else {
				getUsers(w, r)
			}
		case http.MethodPut:
			updateUser(w, r, id)
		case http.MethodDelete:
			deleteUser(w, r, id)
		}
	}
}

func loginUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginUser model.Login

	if err := decoder.Decode(&loginUser); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error decoding json body"))
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	user, errMsg := UserRepo.LoginUser(ctx, loginUser)
	if errMsg != "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(errMsg))
		return
	}

	validToken, err := GenerateJWT(user.Username, user.Age)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to generate token"))
		return
	}

	var token entity.Token
	token.Token = validToken
	json, _ := json.Marshal(token)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getUsers(w http.ResponseWriter, _ *http.Request) {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	users, err := UserRepo.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getUserById(w http.ResponseWriter, _ *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			user, err := UserRepo.GetUserById(ctx, idInt)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			if user.Id != 0 {
				jsonData, _ := json.Marshal(user)
				w.Header().Add("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("User not found"))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid parameter"))
	return
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user entity.User

	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error decoding json body"))
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	result, err := UserRepo.CreateUser(ctx, user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(result))
}

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {
			decoder := json.NewDecoder(r.Body)
			var userSlice model.User
			if err := decoder.Decode(&userSlice); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decoding json body"))
				return
			}

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			result, err := UserRepo.UpdateUser(ctx, idInt, userSlice)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(result))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid parameter"))
	return
}

func deleteUser(w http.ResponseWriter, _ *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			result, err := UserRepo.DeleteUser(ctx, idInt)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(result))
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Invalid parameter"))
	return
}

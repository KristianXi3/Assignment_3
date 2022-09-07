package handler

import (
	"encoding/json"
	"golang-crud-sql/entity"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

type RandomUserHandlerIface interface {
	RandomUserHandler(w http.ResponseWriter, r *http.Request)
}

type RandomUserHandler struct {
}

func NewRandomUserHandler() RandomUserHandlerIface {
	return &RandomUserHandler{}
}

func (u *RandomUserHandler) RandomUserHandler(w http.ResponseWriter, r *http.Request) {

	getRandomUsers(w, r)

}

func getRandomUsers(w http.ResponseWriter, _ *http.Request) {
	res, err := http.Get("https://random-data-api.com/api/users/random_user?size=10")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var users []entity.UserRandom
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Unmarshal(body, &users)

	tpl, err := template.ParseFiles("html/template.html")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl.Execute(w, users)
	return
}

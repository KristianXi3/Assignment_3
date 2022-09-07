package main

import (
	"fmt"
	"golang-crud-sql/context"
	"golang-crud-sql/handler"
	"golang-crud-sql/repository"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const PORT = ":8080"

func main() {

	db := context.Connect()
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	handler.UserRepo = userRepo
	userService := handler.NewUserHandler()

	handler.OrderRepo = orderRepo
	orderService := handler.NewOrderHandler()
	userRandomService := handler.NewRandomUserHandler()

	go handler.GenerateRandomWeather()

	r := mux.NewRouter()
	r.Use(handler.IsAuthorized)

	r.HandleFunc("/weather", handler.WeatherHandler)
	r.HandleFunc("/users", userService.UserHandler)
	r.HandleFunc("/users/{id}", userService.UserHandler)
	r.HandleFunc("/random-users", userRandomService.RandomUserHandler)
	r.HandleFunc("/orders", orderService.OrderHandler)
	r.HandleFunc("/orders/{id}", orderService.OrderHandler)

	fmt.Println("Now listening on port" + PORT)
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0" + PORT,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

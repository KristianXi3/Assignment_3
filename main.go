package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KristianXi3/Assignment_3/context"

	"github.com/KristianXi3/Assignment_3/handler"
	"github.com/KristianXi3/Assignment_3/repository"

	"github.com/gorilla/mux"
)

const PORT = ":8080"

func main() {
	//test
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

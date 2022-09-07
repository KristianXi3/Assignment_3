package handler

import (
	"context"
	"encoding/json"
	"golang-crud-sql/entity"
	"golang-crud-sql/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var OrderRepo repository.OrderRepoIface

type OrderHandlerIface interface {
	OrderHandler(w http.ResponseWriter, r *http.Request)
}

type OrderHandler struct {
}

func NewOrderHandler() OrderHandlerIface {
	return &OrderHandler{}
}

func (o *OrderHandler) OrderHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	switch r.Method {
	case http.MethodGet:
		if id != "" {
			getOrderById(w, r, id)
		} else {
			getOrders(w, r)
		}
	case http.MethodPost:
		createOrder(w, r)
	case http.MethodPut:
		updateOrder(w, r, id)
	case http.MethodDelete:
		deleteOrder(w, r, id)
	}
}

func getOrders(w http.ResponseWriter, _ *http.Request) {

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	orders, err := OrderRepo.GetOrders(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	json, _ := json.Marshal(orders)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func getOrderById(w http.ResponseWriter, _ *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			order, err := OrderRepo.GetOrderById(ctx, idInt)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			if order.OrderId != 0 {
				jsonData, _ := json.Marshal(order)
				w.Header().Add("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}
			w.Write([]byte("Order not found"))
			return
		}
	}
	w.Write([]byte("Invalid parameter"))
	return
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var order entity.Order

	if err := decoder.Decode(&order); err != nil {
		w.Write([]byte("Error decoding json body"))
		return
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()
	result, err := OrderRepo.CreateOrder(ctx, order)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(result))
}

func updateOrder(w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {
			decoder := json.NewDecoder(r.Body)
			var orderSlice entity.Order
			if err := decoder.Decode(&orderSlice); err != nil {
				w.Write([]byte("Error decoding json body"))
				return
			}

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			result, err := OrderRepo.UpdateOrder(ctx, idInt, orderSlice)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(result))
			return
		}
	}
	w.Write([]byte("Invalid parameter"))
	return
}

func deleteOrder(w http.ResponseWriter, _ *http.Request, id string) {
	if id != "" {
		if idInt, err := strconv.Atoi(id); err == nil {

			ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelfunc()
			result, err := OrderRepo.DeleteOrder(ctx, idInt)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(result))
			return
		}
	}
	w.Write([]byte("Invalid parameter"))
	return
}

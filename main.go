package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type order struct {
    Key 				string  `json:"Key"`
    FirstName		string `json:"FirstName"`
    LastName		string `json:"LastName"`
    TableNo			int `json:"TableNo"`
    Order				string `json:"Order"`
    State 			string `json:"State"`
}

type allOrders []order

var orders = allOrders{
	{
    Key: "1",
    FirstName: "John",
    LastName: "Brown",
    TableNo: 32,
    Order: "Rice and jollof",
    State: "Jotted down",
  },
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the waitress app homepage!")
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var newOrder order
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the table number and order")
	}
	
	json.Unmarshal(reqBody, &newOrder)
	orders = append(orders, newOrder)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newOrder)
}

func getOneOrder(w http.ResponseWriter, r *http.Request) {
	orderKey := mux.Vars(r)["key"]

	for _, singleOrder := range orders {
		if singleOrder.Key == orderKey {
			json.NewEncoder(w).Encode(singleOrder)
		}
	}
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(orders)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	orderKey := mux.Vars(r)["id"]
	var updatedOrder order

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedOrder)

	for i, singleOrder := range orders {
		if singleOrder.Key == orderKey {
			singleOrder.FirstName = updatedOrder.FirstName
			singleOrder.LastName = updatedOrder.LastName
			singleOrder.TableNo = updatedOrder.TableNo
			singleOrder.Order = updatedOrder.Order
			singleOrder.State= updatedOrder.State
			orders = append(orders[:i], singleOrder)
			json.NewEncoder(w).Encode(singleOrder)
		}
	}
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	orderKey := mux.Vars(r)["key"]

	for i, singleOrder := range orders {
		if singleOrder.Key == orderKey {
			orders = append(orders[:i], orders[i+1:]...)
			fmt.Fprintf(w, "The order with Key %v has been deleted successfully", orderKey)
		}
	}
}

func main() {
	// initEvents()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/order", createOrder).Methods("POST")
	router.HandleFunc("/order", getAllOrders).Methods("GET")
	router.HandleFunc("/order/{key}", getOneOrder).Methods("GET")
	router.HandleFunc("/order/{key}", updateOrder).Methods("PATCH")
	router.HandleFunc("/order/{key}", deleteOrder).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3030", router))
}

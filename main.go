package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Item structure with exported types
type Item struct {
	Title string
	Body  string
}

// API is a type used to elevate the existing functions to rpc methods
// Needed for rpc:
// both arguments of the method have to be exported types
// the first arg is passed through by the caller
// the second arg represents the result of calling this func/method
// second argument must be a pointer
// the return type of rpc function need to be of error type
type API int

var database []Item

// GetDB method will pass the database to the reply pointer and return it to the client
func (a *API) GetDB(title string, reply *[]Item) error {
	*reply = database
	return nil
}

// GetByName will search db for title and return item
// ToDo: return error if item not found
func (a *API) GetByName(title string, reply *Item) error {
	var getItem Item

	for _, val := range database {
		if val.Title == title {
			getItem = val
		}
	}

	*reply = getItem

	return nil
}

// AddItem to db
func (a *API) AddItem(item Item, reply *Item) error {
	database = append(database, item)
	*reply = item
	return nil
}

// EditItem will take in an item and lookop the db for it and will change only the body
func (a *API) EditItem(edit Item, reply *Item) error {
	var changed Item

	for i, val := range database {
		if val.Title == edit.Title {
			database[i] = Item{edit.Title, edit.Body}
			changed = database[i]
		}
	}

	*reply = changed
	return nil
}

// DeleteItem will delete and item in the db and return old(deleted) item
func (a *API) DeleteItem(item Item, reply *Item) error {
	var del Item

	for i, val := range database {
		if val.Title == item.Title && val.Body == item.Body {
			database = append(database[:i], database[i+1:]...)
			del = item
			break
		}
	}

	*reply = del
	return nil
}

func main() {

	var api = new(API)
	err := rpc.Register(api)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("error setting up listener", err)
	}

	log.Printf("serving rpc on port %d", 4040)
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("error serving", err)
	}

}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct{
	ID uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
}

var users []User
var maxID uint64

func init(){
	users = []User{{
		ID: 1,
		FirstName: "Nahuel",
		LastName: "Costamagna",
		Email: "nahuel@domain.com",
	},{
		ID: 2,
		FirstName: "Eren",
		LastName: "Jaeger",
		Email: "eren@domain.com",
	},{
		ID: 3,
		FirstName: "Paco",
		LastName: "Costa",
		Email: "paco@domain.com",
	}}
	maxID = 3
}

func main() {

	http.HandleFunc("/users", UserServer)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserServer(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetAllUser(w)
	case http.MethodPost:
		decode := json.NewDecoder(r.Body)
		var u User
		if err := decode.Decode(&u); err != nil {
			MsgResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		PostUser(w, u)
	default:
		InvalidMethod(w)		
		
		
	}

}

func GetAllUser(w http.ResponseWriter){
	DataResponse(w, http.StatusOK, users)
}


func PostUser(w http.ResponseWriter, data interface{}){
	user := data.(User)

	if user.FirstName == "" {
		MsgResponse(w, http.StatusBadRequest, "first name is required")
		return
	}

	if user.LastName == "" {
		MsgResponse(w, http.StatusBadRequest, "last name is required")
		return
	}

	if user.Email == "" {
		MsgResponse(w, http.StatusBadRequest, "email is required")
		return
	}

	maxID++
	user.ID = maxID
	users = append(users, user)
	DataResponse(w, http.StatusCreated, user)
}

func InvalidMethod(w http.ResponseWriter){
	status := http.StatusNotFound
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)
}

func MsgResponse(w http.ResponseWriter, status int, message string){
w.WriteHeader(status)
fmt.Fprintf(w,  `{"status": %d, "message":"%s"}`, status, message)
}

func DataResponse(w http.ResponseWriter, status int, users interface{}) {
	value, err := json.Marshal(users)
	if err != nil {
		MsgResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data":%s}`, status, value)
}
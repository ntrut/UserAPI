package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Updated   string `json:"updated"`
}

/*get one user from the database using a query parameter "id"*/
func read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	/*check if the query parameter is correct*/
	keys := mux.Vars(r)

	/*check if we successfully opened the infor.db*/
	database, sqlerr := sql.Open("sqlite3", "infor.db")
	if sqlerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
		return
	}
	defer database.Close()

	output, _ := database.Query("SELECT * FROM users WHERE id = " + keys["id"])
	var user User

	//loop through the output, should be only one user in the output
	defer output.Close()
	for output.Next() {
		scanerr := output.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Updated)
		if scanerr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("500 - Internal Server Error")
			return
		}
	}

	//if id is 0, then that means we didnt get any users from the database with that query parameter
	if user.Id != 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("404 - User Not Found")
	}
}

/*create one user and insert into database"*/
func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decode := json.NewDecoder(r.Body)
	var user User
	err := decode.Decode(&user)
	/*check if there is an error decoding*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
		return
	}

	/*check if any of the json parameters dont exist, we dont need to check timestamp*/
	if user.Id == 0 || user.FirstName == "" || user.LastName == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("400 - Bad Request")
		return
	}

	/*check if id is negative*/
	if user.Id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("400 - Bad Request, Id cant be negative!")
		return
	}

	user.Updated = time.Now().String()

	/*check if we successfully opened the infor.db*/
	database, sqlerr := sql.Open("sqlite3", "infor.db")
	if sqlerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
	} else {
		defer database.Close()

		query := `SELECT 1 FROM users WHERE id = ?`
		output, _ := database.Query(query, user.Id)

		defer output.Close()
		if output.Next() {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode("409 - Conflict, User already Exists with that Id")
		} else {
			fmt.Println(user)
			statement, _ := database.Prepare("INSERT INTO users (id, email, first_name, last_name, updated) VALUES (?,?,?,?,?)")
			statement.Exec(user.Id, user.Email, user.FirstName, user.LastName, user.Updated)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(user)
		}
	}
}

/*update one user*/
func update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")

	decode := json.NewDecoder(r.Body)
	var user User
	err := decode.Decode(&user)
	/*check if there is an error decoding*/
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
		return
	}

	/*check if any of the json parameters dont exist, we dont need to check timestamp*/
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("400 - Bad Request")
		return
	}

	user.Updated = time.Now().String()
	num, _ := strconv.Atoi(vars["id"])
	user.Id = num

	/*check if we successfully opened the infor.db*/
	database, sqlerr := sql.Open("sqlite3", "infor.db")
	if sqlerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
	} else {
		defer database.Close()

		query := `select exists(select 1 from users where id = ?)`
		output := database.QueryRow(query, user.Id)
		var num int
		output.Scan(&num)
		if num != 1 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("404 - Not Found, Id doesn't exists")
		} else {
			statement, _ := database.Prepare("UPDATE users SET first_name = ?, last_name = ?, email = ?, updated = ? WHERE id = ?")
			_, err := statement.Exec(user.FirstName, user.LastName, user.Email, user.Updated, user.Id)
			if err != nil {
				fmt.Println(err)
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		}
	}
}

/*delete a user based on a id query parameter*/
func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	/*check if we successfully opened the infor.db*/
	database, sqlerr := sql.Open("sqlite3", "infor.db")
	if sqlerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
		return
	} else {
		defer database.Close()

		var num int
		query := `select exists(select 1 from users where id = ?)`
		output := database.QueryRow(query, vars["id"])
		output.Scan(&num)

		if num != 1 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("404 - Not Found, Id doesn't exists")
		} else {
			statement, _ := database.Prepare("DELETE FROM users WHERE id = ?")
			_, err := statement.Exec(vars["id"])
			if err != nil {
				fmt.Println(err.Error())
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("User with id " + vars["id"] + " has been deleted")
		}
	}
}

/*list all users from database*/
func list(w http.ResponseWriter, r *http.Request) {

	sort, err := r.URL.Query()["sort"]
	order, err2 := r.URL.Query()["order"]
	//if one of the parameters doesnt exist, then throw a bad request error
	if (err && !err2) || (!err && err2) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("400 - Bad Request")
		return
	}

	/*check if we successfully opened the infor.db*/
	database, sqlerr := sql.Open("sqlite3", "infor.db")
	if sqlerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Internal Server Error")
	}
	defer database.Close()

	//if both are not empty, then continue
	if !err && !err2 {
		query := "SELECT * FROM users"
		var list []User = getList(query, database)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(list)
	} else {
		//parameters need to equal to these strings or else it will fail
		if (sort[0] == "id" || sort[0] == "first_name" || sort[0] == "last_name" || sort[0] == "email" || sort[0] == "updated") && (order[0] == "asc" || order[0] == "desc") {
			query := "SELECT * FROM users ORDER BY " + sort[0] + " " + order[0]
			var list []User = getList(query, database)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(list)
		} else {
			//bad parameters, throw bad request
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("400 - Bad Request")
		}
	}
}

/*return a list of users*/
func getList(query string, database *sql.DB) []User {
	var list []User
	output, _ := database.Query(query)
	var user User
	for output.Next() {
		_ = output.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Updated)
		list = append(list, user)
	}
	return list
}

func initailize() {
	/*initialize the db with 10 users!*/
	log.Println("Creating sqlite database...")
	file, err := os.Create("infor.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	file.Close()
	database, err := sql.Open("sqlite3", "infor.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer database.Close()

	log.Println("Database created!")

	//create 10 users, primary key id auto increments
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, email TEXT, first_name TEXT, last_name TEXT, updated TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO users (email, first_name, last_name, updated) VALUES (?,?,?,?)")
	statement.Exec("nazartrut@gmail.com", "Nazar", "Trut", time.Now())
	statement.Exec("bobmarley@gmail.com", "Bob", "Marley", time.Now())
	statement.Exec("lionelmessi@gmail.com", "Lionel", "Messi", time.Now())
	statement.Exec("sergioaguero@gmail.com", "Sergio", "Aguero", time.Now())
	statement.Exec("davidsilva@gmail.com", "David", "Silva", time.Now())
	statement.Exec("kevindebruyne@gmail.com", "Kevin", "De Bruyne", time.Now())
	statement.Exec("gabrieljesus@gmail.com", "Gabriel", "Jesus", time.Now())
	statement.Exec("bernardosilva@gmail.com", "Bernardo", "Silva", time.Now())
	statement.Exec("philfoden@gmail.com", "Phil", "Foden", time.Now())
	statement.Exec("jackgrealish@gmail.com", "Jack", "Grealish", time.Now())
	statement.Close()
}

func main() {

	initailize()

	//create server on localhost 8080 with gorilla mux
	r := mux.NewRouter()
	r.HandleFunc("/infor/read/{id:[0-9]+}", read).Methods("GET")
	r.HandleFunc("/infor/create", create).Methods("POST")
	r.HandleFunc("/infor/update/{id:[0-9]+}", update).Methods("PUT")
	r.HandleFunc("/infor/delete/{id:[0-9]+}", delete).Methods("DELETE")
	r.HandleFunc("/infor/list", list).Methods("GET")
	http.ListenAndServe(":8080", r)

}

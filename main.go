package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type Info struct {
	UserInfo User
	TaskList TaskList
}

var info Info

type Ping struct {
	Ping string `json:"ping"`
}

type User struct {
	Password string `json:"password"`
	UserName string `json:"username"`
}

type TaskList struct {
	Contents []Task `json:"tasklist"`
}

type Task struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Detail string    `json:"detail"`
}

func handleError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, "")
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	ping, _ := json.Marshal(Ping{
		Ping: "ping",
	})
	fmt.Fprint(w, string(ping))
}

func login(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len) // Content-Length と同じサイズの byte 配列を用意
	r.Body.Read(body)         // byte 配列にリクエストボディを読み込む
	var user User
	_ = json.Unmarshal(body, &user)

	info.UserInfo = user

	fmt.Fprintln(w, "")
}

type PostTaskReq struct {
	User User `json:"user"`
	Task Task `json:"task"`
}

func postTask(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte, len) // Content-Length と同じサイズの byte 配列を用意
	r.Body.Read(body)

	var req PostTaskReq
	err := json.Unmarshal(body, &req)
	if err != nil {
		handleError(w, http.StatusBadRequest)
		return
	}

	user := req.User
	if !reflect.DeepEqual(user, info.UserInfo) {
		handleError(w, http.StatusBadRequest)
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		handleError(w, http.StatusBadRequest)
		return
	}

	task := req.Task
	task.ID = id
	info.TaskList.Contents = append(info.TaskList.Contents, task)

	fmt.Fprint(w, id)
}

func getTask(w http.ResponseWriter, r *http.Request) {

}

func main() {

	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/login", login)
	http.HandleFunc("POST /task", postTask)
	http.HandleFunc("GET /task/id", getTask)

	server.ListenAndServe()
}

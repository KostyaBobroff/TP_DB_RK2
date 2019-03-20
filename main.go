package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/api/forum/create", createForum).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/createBranch", createBranch).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", getDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/threads", getThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", getUsers).Methods("GET")
	r.HandleFunc("/api/post/{id}/details", getDetails).Methods("GET")
	r.HandleFunc("/api/post/{id}/details", postDetails).Methods("POST")
	r.HandleFunc("/api/service/clear", clearService).Methods("POST")
	r.HandleFunc("/api/service/status", getStatus).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/create", createThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/details", getThread).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/details", updateThread).Methods("POST")
	r.HandleFunc("/api/thread/{slug_or_id}/posts", getPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug_or_id}/vote", postVote).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/create", createUser).Methods("POST")
	r.HandleFunc("/api/user/{nickname}/profile", getUser).Methods("GET")
	r.HandleFunc("/api/user/{nickname}/profile", updateUser).Methods("GET")

	//r.PathPrefix("/data/").Handler(http.StripPrefix("/data/", http.FileServer(http.Dir("./static/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Println("Starting server at http://127.0.0.1:5000")
	log.Fatal(srv.ListenAndServe())
}
package main

import (
	"fmt"
	"log"
	"net/http"

	"hello/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/userDetails", controllers.GetUserDetails).Methods("GET")
	router.HandleFunc("/playlistSong", controllers.InsertPlaylistAndSong).Methods("POST")
	router.HandleFunc("/playlist/{playlist_id}", controllers.DeletePlaylist).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

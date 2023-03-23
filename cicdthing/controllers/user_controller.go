package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func InsertPlaylistAndSong(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()

	if err != nil {
		sendErrorResponse(w, "failed")
		return
	}

	playlistName := r.Form.Get("playlistName")
	userId, _ := strconv.Atoi(r.Form.Get("userId"))

	title := r.Form.Get("title")
	duration, _ := strconv.Atoi(r.Form.Get("duration"))
	singer := r.Form.Get("singer")

	_, errQuery := db.Exec("INSERT INTO playlists(playlistName, userId, playlistState)values (?,?,?)",
		playlistName,
		userId,
		1,
	)
	fmt.Println("this one")
	fmt.Println(errQuery)
	fmt.Println("========================")
	_, errQuery = db.Exec("INSERT INTO songs(songTitle, songDuration, songSinger)values (?,?,?)",
		title,
		duration,
		singer,
	)
	fmt.Println(errQuery)
	fmt.Println("========================")
	var playlistId int
	var playlistDateCreated string
	var songId int
	db.QueryRow("select playlistId, playlistDateCreated from playlists where playlistName = ?", playlistName).Scan(&playlistId, &playlistDateCreated)
	db.QueryRow("select songId from songs where songTitle = ?", title).Scan(&songId)
	fmt.Println("playlist created: " + playlistDateCreated)
	_, errQuery = db.Exec("INSERT INTO detailplaylistsong(playlistId, songId, timePlayed)values (?,?,?)",
		playlistId,
		songId,
		0,
	)
	fmt.Println(errQuery)
	fmt.Println("========================")
	var response PlaylistOutputAfterInsert
	var inputSong Song
	inputSong.ID = songId
	inputSong.Title = title
	inputSong.Duration = duration
	inputSong.Singer = singer

	var output PlaylistOutput
	var songs []Song
	output.ID = playlistId
	output.Name = playlistName
	output.DateCreated = playlistDateCreated
	songs = append(songs, inputSong)
	output.Songs = songs
	if errQuery == nil {
		response.Status = 200
		response.Message = "Success"
		response.Data = output
	} else {
		fmt.Println(errQuery)
		response.Status = 400
		response.Message = "Insert Failed"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeletePlaylist(w http.ResponseWriter, r *http.Request) {

	db := connect()
	defer db.Close()

	err := r.ParseForm()

	if err != nil {
		return
	}
	fmt.Println("in here")
	vars := mux.Vars(r)
	playlistId := vars["playlist_id"]

	var count int
	var playlist PlaylistDisplay

	db.QueryRow("select count(*) from detailplaylistsong where playlistId = ?", playlistId).Scan(&count)
	db.QueryRow("select playlists.playlistId, playlists.playlistName, playlists.playlistDateCreated FROM playlists where playlistId = ?", playlistId).Scan(&playlist.ID, &playlist.Name, &playlist.DateCreated)
	var response DeletePlaylistResponse
	if count != 0 {
		response.Status = 400
		response.Message = "Failed, playlist isn't empty"
		response.Data.ID = playlist.ID
		response.Data.Name = playlist.Name
		response.Data.DateCreated = playlist.DateCreated
		response.Data.Status = "false"
	} else {
		_, errQuery := db.Exec("DELETE FROM playlists WHERE playlistId=?",
			playlistId,
		)

		if errQuery == nil {
			response.Status = 200
			response.Message = "Successfully deleted playlist"
			response.Data.ID = playlist.ID
			response.Data.Name = playlist.Name
			response.Data.DateCreated = playlist.DateCreated
			response.Data.Status = "true"
		} else {
			fmt.Println(errQuery)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT playlists.playlistId, playlists.playlistName, playlists.playlistDateCreated, detailplaylistsong.timePlayed FROM playlists INNER JOIN detailplaylistsong ON playlists.playlistID = detailplaylistsong.playlistID"

	id := r.URL.Query()["userId"]
	if id != nil {
		query += " WHERE userId='" + id[0] + "' GROUP BY playlists.playlistId"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}

	var playlistOutput PlaylistOutput
	var playlistOutputs []PlaylistOutput

	for rows.Next() {
		if err := rows.Scan(&playlistOutput.ID, &playlistOutput.Name, &playlistOutput.DateCreated, &playlistOutput.TimePlayed); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			playlistOutputs = append(playlistOutputs, playlistOutput)
		}
	}
	fmt.Println("this one: ")
	fmt.Println(len(playlistOutputs))
	fmt.Println("=================")
	for i := 0; i < len(playlistOutputs); i++ {
		query = "SELECT songs.songId, songs.songTitle, songs.songDuration, songs.songSinger, detailplaylistsong.timePlayed FROM songs INNER JOIN detailplaylistsong ON songs.songId = detailplaylistsong.songId WHERE detailplaylistsong.playlistId ='" + strconv.Itoa(playlistOutputs[i].ID) + "'"
		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		}

		var song Song
		var songs []Song
		var timePlayed int

		for rows.Next() {
			if err := rows.Scan(&song.ID, &song.Title, &song.Duration, &song.Singer, &timePlayed); err != nil {
				log.Println(err)
				sendErrorResponse(w, "Something went wrong, please try again")
				return
			} else {
				songs = append(songs, song)
			}
		}
		playlistOutputs[i].Songs = songs
		playlistOutputs[i].TimePlayed = timePlayed
	}
	query = "SELECT * FROM user WHERE userId = '" + id[0] + "'"
	rows, err = db.Query(query)
	if err != nil {
		log.Println(err)
		sendErrorResponse(w, "Something went wrong, please try again")
		return
	}
	var userdetail UserDetails
	var userdetails []UserDetails
	for rows.Next() {
		if err := rows.Scan(&userdetail.ID, &userdetail.Name); err != nil {
			log.Println(err)
			sendErrorResponse(w, "Something went wrong, please try again")
			return
		} else {
			userdetails = append(userdetails, userdetail)
		}
	}
	fmt.Println(len(playlistOutputs))
	userdetail.Playlists = playlistOutputs

	w.Header().Set("Content-Type", "application/json")

	if len(playlistOutputs) < 100 {
		var response UserDetailsResponse
		response.Data = userdetail
		json.NewEncoder(w).Encode(response)
	} else {
		var response ErrorResponse
		response.Status = 400
		response.Message = "Error Array Size Not Correct"
		json.NewEncoder(w).Encode(response)
	}

}

func sendErrorResponse(w http.ResponseWriter, message string) {
	var response ErrorResponse
	response.Status = 400
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

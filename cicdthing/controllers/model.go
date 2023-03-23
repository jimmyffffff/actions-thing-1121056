package controllers

type User struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

type Playlist struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	DateCreated string `json:"Date Created"`
	State       bool   `json:"Playlist State"`
	UserId      int    `json:"UserId"`
}

type DetailPlaylistSong struct {
	PlaylistId int `json:"PlaylistId"`
	SongId     int `json:"SongId"`
	TimePlayed int `json:"TimePlayed"`
}

type Song struct {
	ID       int    `json:"ID"`
	Title    string `json:"Title"`
	Duration int    `json:"Duration"`
	Singer   string `json:"Singer"`
}

type UserDetailsResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    UserDetails `json:"data"`
}

type UserDetails struct {
	ID        int              `json:"ID"`
	Name      string           `json:"Name"`
	Playlists []PlaylistOutput `json:"Playlists"`
}

type PlaylistOutput struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	DateCreated string `json:"Date Created"`
	Songs       []Song `json:"Songs"`
	TimePlayed  int    `json:"TimePlayed"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type DeletePlaylistResponse struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    PlaylistDisplay `json:"data"`
}

type PlaylistDisplay struct {
	ID          int    `json:"ID"`
	Name        string `json:"Name"`
	DateCreated string `json:"Date Created"`
	Status      string `json:"status"`
}

type PlaylistOutputAfterInsert struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    PlaylistOutput `json:"data"`
}

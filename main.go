package main

import "net/http"

var CurrentGames map[string]GameBoard
var Users map[string]string

func init() {
	CurrentGames = map[string]GameBoard{}
	Users = map[string]string{}
	http.HandleFunc("/play", GameHandler)

}

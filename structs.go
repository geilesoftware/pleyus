package main

import (
	"time"
)

type Game struct {
	Name  string
	Genre string
}

type Lobby struct {
	Title              string
	Created            time.Time
	Updated            time.Time
	Access             string
	Owner              *User
	State              string
	Game               *Game
	SkillLevel         int
	EstimatedPlayTime  int
	EstimatedStartTime time.Time
}

type User struct {
	Id        int64
	Username  string
	Email     string
	Password  string
	Hash      string
	LastLogin time.Time
	Created   time.Time
	Lobby     *Lobby
}

func NewGame(name string, genre string) *Game {
	return &Game{
		Name:  name,
		Genre: genre,
	}
}

func NewLobby(title string, access string, owner *User, state string, game *Game, skillLevel int, estimatedPlayTime int, estimatedStartTime time.Time) *Lobby {
	return &Lobby{
		Title:              title,
		Created:            time.Now(),
		Owner:              owner,
		State:              state,
		Game:               game,
		SkillLevel:         skillLevel,
		EstimatedPlayTime:  estimatedPlayTime,
		EstimatedStartTime: estimatedStartTime,
	}
}

func NewUser(id int64, username string, email string, password string, hash string, lobby *Lobby) *User {
	return &User{
		Id:       id,
		Username: username,
		Email:    email,
		Password: password,
		Hash:     hash,
		Lobby:    lobby,
	}
}

package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
)
import "github.com/joho/godotenv"

func main() {
	// read config from .env
	// loadConfig()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("STUID")
	password := os.Getenv("PASSWORD")
	tokenID := os.Getenv("TOKENID")
	userID := os.Getenv("USERID")
	// base64 decode password
	newPass, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		log.Fatal(err)
	}
	password = string(newPass)
	//fmt.Println(username, password)

	// create new user
	user := User{
		Username: username,
		Password: password,
		Name:     "",
		TokenId:  tokenID,
		UserId:   userID,
	}
	err = cli(&user)
	if err != nil {
		fmt.Println(err)
	}
}

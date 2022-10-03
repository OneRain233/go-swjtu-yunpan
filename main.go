package main

import (
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
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	// create new user
	user := User{
		Username: username,
		Password: password,
		Name:     "",
		TokenId:  "",
		UserId:   "",
	}
	err = cli(&user)
	if err != nil {
		fmt.Println(err)
	}
}

func loadConfig() error {
	// load config from .env
	return nil
}

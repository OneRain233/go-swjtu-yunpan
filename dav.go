package main

// TODO: implement webdav server
import (
	"github.com/gin-gonic/gin"
	"log"
)

var curUser *User

func davMain(user *User) {
	if !isLogin {
		log.Fatal("Please login first")
		return
	}
	curUser = user
	router := gin.Default()
	router.GET("/", index)

	err := router.Run(":8080")
	if err != nil {
		return
	}

}

func index(c *gin.Context) {
	err, entries := curUser.getDocsEntries()
	if err != nil {
		log.Fatal(err)
		return
	}
	c.JSON(200, entries)
}

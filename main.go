package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("assets/templates/*")
	router.GET("/", index)
	router.GET("/login", login)
	router.GET("/signup", signup)
	router.GET("/wj", wj)
	router.Run("127.0.0.1:7000")
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", "")
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", "")
}

func signup(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", "")
}

func wj(c *gin.Context) {
	c.HTML(http.StatusOK, "wj.html", "")
}

package main

import "github.com/gin-gonic/gin"


func html_parse(c *gin.Context) {
	c.HTML(200, "1.html", nil)

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func main(){
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("", html_parse)

	r.Run()
}
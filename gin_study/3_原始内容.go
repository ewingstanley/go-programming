package main

import (
	"bytes"
	"io"

	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("", func(c *gin.Context) {
		byteData, _ := io.ReadAll(c.Request.Body)

		fmt.Println(string(byteData))

		c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

		name := c.PostForm("name")

		fmt.Println(name)

	})

	router.Run(":8080")
}

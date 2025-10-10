package main

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}

func main() {
	// 1.初始化gin
	r := gin.Default()

	// 2.创建路由
	r.GET("/index", Index)

	// 3. 启动服务
	r.Run(":8080")
}

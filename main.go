package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 获取路由对象
	router := gin.Default()
	//fmt.Println(router)

	router.GET("index", func(context *gin.Context) {
		context.String(200, "hello world")
	})

	router.Run(":8080")

}

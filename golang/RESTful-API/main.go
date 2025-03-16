package main

import (
	_ "RESTful-API/bootstrap" // 只执行 bootstrap 的 init()，不直接使用包内函数
)

type User struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Wechat string `json:"wechat"`
}

func main() {

	////配置相关
	//defer cmd.Clean()
	//cmd.Start()
	//
	//db, err := utils.ConnectToDatabase()
	//if err != nil {
	//	fmt.Println("Error connecting to database:", err)
	//	return
	//}
	//fmt.Println(db, err)
	//err = db.AutoMigrate(&mysql_db.CrudList{})
	//err = db.AutoMigrate(&mysql_db.UserList{})

	/*
		//users := []User{{ID: 123, Name: "张三丰"}, {ID: 456, Name: "张无忌"}}
		r := gin.Default()

		//// 获取路由对象
		//router := gin.Default()
		////fmt.Println(router)
		//
		//router.GET("index", func(context *gin.Context) {
		//	context.String(200, "hello world")
		//})

		// get请求
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name": "小红",
				"desc": "你好",
			})
		})

		//get请求单个数据 http://127.0.0.1:8081/users?123?name=老王
		r.GET("/users/:id", func(c *gin.Context) {
			id := c.Param("id")
			name := c.Query("name")
			fmt.Println("id为" + id)
			c.JSON(http.StatusOK, gin.H{ // 页面中输出打印
				name:   name,
				"desc": "你好，我的id" + id,
			})
		})
		// get请求数组 http://127.0.0.1:8081/slice?media=bbc&media=aab
		r.GET("/slice", func(context *gin.Context) {
			context.JSON(http.StatusOK, context.QueryArray("media"))
		})
		// get请求 map键值对 http://127.0.0.1:8081/map?ids[a]=bbc&ids[b]=123
		r.GET("/map", func(context *gin.Context) {
			context.JSON(http.StatusOK, context.QueryMap("ids"))
		})

		// post请求表单数据 https://127.0.0.1:8081/users?name=张三&wechat=123
		r.POST("/users", func(c *gin.Context) {
			//创建用户
			name := c.PostForm("name")
			wechat := c.PostForm("wechat")
			c.JSON(http.StatusOK, gin.H{
				"name":   name,
				"wechat": wechat,
			})

		})
		// post请求json格式数据
		r.POST("/usersjson", func(c *gin.Context) {
			user := &User{}         // 声明并初始化一个指向 User 结构体的指针变量 user。 User 是一个自定义的结构体类型，这里创建了一个空的 User 结构体实例。
			err := c.BindJSON(user) //绑定到user
			if err != nil {         // 如果有错误，抛出错误
				panic(err)
			}
			// 没有错误，打印
			c.JSON(http.StatusOK, gin.H{"user": user})

		})

		r.GET("/users", func(context *gin.Context) {
			//获取全部用户
		})
		r.PUT("/user/:id", func(context *gin.Context) {
			//更新用户id 为 xx 的信息
		})
		r.DELETE("/user/:id", func(context *gin.Context) {
			//删除用户 id为xx
		})

	*/
	//r.Run(":8081")
	//r.Run(viper.GetString("server.addr") + ":" + viper.GetString("server.port"))
}

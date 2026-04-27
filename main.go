package main

import (
	"log"
	"net/http"
	myws "xchat-server/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	hub := myws.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		myws.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	/*db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败")
	}

	// 自动建表
	db.AutoMigrate(&model.User{})

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	userController := controller.NewUserController(service)
	ginServer := gin.Default()

	//ginServer.Use(middleware.JWTAuth())

	ginServer.POST("/user/login", userController.Login)
	ginServer.Run(":8080")*/

}

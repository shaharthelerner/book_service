package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pkg/service/pkg/router"
)

func main() {
	fmt.Println("Hello, World!")
	r := gin.Default()
	router.Routes(r)
	err := r.Run()
	if err != nil {
		return
	}
	//before_refactor.StartService()
	//before_refactor.SaveMockCache()
}

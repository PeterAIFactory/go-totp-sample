package main

import (
	"Go_learning/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	p := gin.Default()
	routes.InitRoutes(p)
	fmt.Println("Start running")
	_ = p.Run(":8080")
}

package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()
	r.POST("/create" createUser)
	r.RUN(":8080")
}

func createUser() {
	
}


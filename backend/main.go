package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func main() {
	// r := gin.Default()
	// r.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "Hello World!",
	// 	})
	// })

	// r.Run()

	fmt.Println(add(1, 2))
}

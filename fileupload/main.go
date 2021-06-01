package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.Default()
	app.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})
	Router(app)
	app.Run(":9999")
}

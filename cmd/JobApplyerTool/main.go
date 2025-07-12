package main

import (
	"github.com/gin-gonic/gin"
)

type Data struct {
	Message string `json:"message"`
	Id      int    `json:"id"`
	BookId  int    `json:"book_id"`
}

func main() {
	//router := gin.Default()
	//
	//router.GET("/someJSON", func(c *gin.Context) {
	//	data := map[string]interface{}{
	//		"lang": "GO语言",
	//		"tag":  "<br>",
	//	}
	//
	//	// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
	//	c.AsciiJSON(http.StatusOK, data)
	//})
	//
	//// Listen and serve on 0.0.0.0:8080
	//router.Run(":8080")

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		data := Data{
			Message: "Hello, World!",
			Id:      1,
			BookId:  101,
		}
		c.AsciiJSON(200, data)
	})

	router.Run(":8080")

}

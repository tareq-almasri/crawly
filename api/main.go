package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type scannedPage struct {
	ID                     int
	HtmlVersion            string
	PageTitle              string
	CountOfH1              uint8
	CountOfH2              uint8
	CountOfH3              uint8
	CountOfH4              uint8
	CountOfH5              uint8
	CountOfH6              uint8
	InternalLinksCount     uint8
	ExternalLinksCount     uint8
	InaccessibleLinksCount uint8
	HasLoginForm           bool
	URL                    string
	Status                 string
}

var example scannedPage = scannedPage{
	ID:                     5555,
	HtmlVersion:            "h",
	PageTitle:              "testzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
	CountOfH1:              1,
	CountOfH2:              1,
	CountOfH3:              1,
	CountOfH4:              1,
	CountOfH5:              1,
	CountOfH6:              1,
	InternalLinksCount:     1,
	ExternalLinksCount:     1,
	InaccessibleLinksCount: 1,
	HasLoginForm:           true,
	URL:                    "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
}

func main() {
	fmt.Println("hello world")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/api")
	{
		api.GET("/results", func(ctx *gin.Context) {
			ctx.JSON(200, "hello world")
		})

	}

	router.Run()
}
